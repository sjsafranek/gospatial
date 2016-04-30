/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/paulmach/go.geojson"
	"gospatial/app"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	database string
)

type DumpedDatabase struct {
	Apikeys map[string]app.Customer               `json:"apikeys"`
	Layers  map[string]*geojson.FeatureCollection `json:"layers"`
}

func list_datsources() {
	fmt.Println("Datasources:")
	// get datbase
	app.DB = app.Database{File: "./" + database + ".db"}
	conn, err := bolt.Open(app.DB.File, 0644, nil)
	if err != nil {
		conn.Close()
		fmt.Println(err)
		os.Exit(1)
	}
	// Get all layers
	conn.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("layers"))
		b.ForEach(func(k, v []byte) error {
			fmt.Println(string(k))
			return nil
		})
		return nil
	})
	// close database
	conn.Close()
}

func export_datasource(datasource string) {
	fmt.Println("Exporting datasource: ", datasource)
	// setup database
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()
	// get datasource from database
	lyr, err := app.DB.GetLayer(datasource)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// marshal to json
	b, err := lyr.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Write to file
	savename := datasource + ".geojson"
	ioutil.WriteFile(savename, b, 0644)
}

func import_datasource(import_file string) {
	fmt.Println("Importing", import_file)
	// setup database
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()
	// get geojson file
	var geojson_file string
	ext := strings.Split(import_file, ".")[1]
	// convert shapefile
	if ext == "shp" {
		// Convert .shp to .geojson
		// ogr2ogr -f GeoJSON -t_srs crs:84 [name].geojson [name].shp
		geojson_file := strings.Replace(import_file, ".shp", ".geojson", -1)
		// fmt.Println("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojson_file, shapefile)
		out, err := exec.Command("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojson_file, import_file).Output()
		if err != nil {
			fmt.Println(err)
			fmt.Println(string(out))
			os.Exit(1)
		} else {
			fmt.Println(geojson_file, "created")
			fmt.Println(string(out))
		}
	} else if ext == "geojson" {
		geojson_file = import_file
	} else {
		fmt.Println("Unsupported file type", ext)
		os.Exit(1)
	}
	// Read .geojson file
	file, err := ioutil.ReadFile(geojson_file)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	// Unmarshal to geojson struct
	geojs, err := geojson.UnmarshalFeatureCollection(file)
	if err != nil {
		fmt.Printf("Unmarshal GeoJSON error: %v\n", err)
		os.Exit(1)
	}
	// Create datasource
	ds, _ := app.NewUUID()
	app.DB.InsertLayer(ds, geojs)
	fmt.Println("Datasource created:", ds)
	// Cleanup artifacts
	if geojson_file != import_file {
		os.Remove(geojson_file)
	}
}

func init() {
	// flag.Usage = func{
	// 	"stuff"
	// }
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.Parse()
}

func main() {

	required_args := flag.Args()

	if len(required_args) == 0 {
		fmt.Println("Incorrect usage!")
		os.Exit(1)
	}

	method := required_args[0]

	if method == "ls" {
		list_datsources()
	} else if method == "export" {
		if len(required_args) != 2 {
			fmt.Println("Incorrect usage!")
			os.Exit(1)
		}
		datasource := required_args[1]
		export_datasource(datasource)
	} else if method == "import" {
		if len(required_args) != 2 {
			fmt.Println("Incorrect usage!")
			os.Exit(1)
		}
		import_file := required_args[1]
		import_datasource(import_file)
	} else if method == "create" {
		if len(required_args) != 2 {
			fmt.Println("Incorrect usage!")
			os.Exit(1)
		} else if required_args[1] == "datasource" {
			fmt.Println("Creating datasource")
			app.DB = app.Database{File: "./" + database + ".db"}
			app.DB.Init()
			ds, err := app.DB.NewLayer()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Datasource created:", ds)
		} else if required_args[1] == "customer" {
			fmt.Println("Creating customer")
			app.DB = app.Database{File: "./" + database + ".db"}
			app.DB.Init()
			apikey := app.NewAPIKey(12)
			customer := app.Customer{Apikey: apikey}
			err := app.DB.InsertCustomer(customer)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Customer created:", apikey)
		} else {
			fmt.Println("Incorrect usage!")
			os.Exit(1)
		}
	} else if method == "backup" {
		fmt.Println("Backing up database...")
		app.DB = app.Database{File: "./" + database + ".db"}
		app.DB.Init()
		savefile := "backup_" + time.Now().String()
		app.DB.Backup(savefile)
		fmt.Println("Backup created:", savefile)
	} else if method == "load" {
		if len(required_args) != 2 {
			fmt.Println("Incorrect usage!")
			os.Exit(1)
		} else {
			filename := required_args[1]
			fmt.Println("Loading database...")
			app.DB = app.Database{File: "./" + database + ".db"}
			app.DB.Init()
			fmt.Printf("Loading database [%s]\n", filename)
			// check for file
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				fmt.Println("File not found [" + filename + "]")
			}
			// open json file
			file, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Println(err)
			}
			// unmarshal data
			var data DumpedDatabase
			data.Apikeys = make(map[string]app.Customer)
			data.Layers = make(map[string]*geojson.FeatureCollection)
			err = json.Unmarshal(file, &data)
			if err != nil {
				fmt.Println(err)
			}
			app.DB.InsertCustomers(data.Apikeys)
			app.DB.InsertLayers(data.Layers)
		}
	} else {
		fmt.Println("Incorrect usage!")
		os.Exit(1)
	}
	// exit
	os.Exit(0)
}
