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
	"gospatial/utils"
)

var (
	database string
)

type dumpedDatabase struct {
	Apikeys map[string]app.Customer               `json:"apikeys"`
	Layers  map[string]*geojson.FeatureCollection `json:"layers"`
}

func usageError(message string) {
	fmt.Println("Incorrect usage!")
	fmt.Println(message)
	os.Exit(1)
}

func setupDb() {
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()
}

func listDatasources() {
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

func exportDatasource(datasource string) {
	fmt.Println("Exporting datasource: ", datasource)
	// setup database
	setupDb()
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

func importDatasource(importFile string) {
	fmt.Println("Importing", importFile)
	// setup database
	setupDb()
	// get geojson file
	var geojsonFile string
	ext := strings.Split(importFile, ".")[1]
	// convert shapefile
	if ext == "shp" {
		// Convert .shp to .geojson
		// ogr2ogr -f GeoJSON -t_srs crs:84 [name].geojson [name].shp
		geojsonFile := strings.Replace(importFile, ".shp", ".geojson", -1)
		// fmt.Println("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojsonFile, shapefile)
		out, err := exec.Command("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojsonFile, importFile).Output()
		if err != nil {
			fmt.Println(err)
			fmt.Println(string(out))
			os.Exit(1)
		} else {
			fmt.Println(geojsonFile, "created")
			fmt.Println(string(out))
		}
	} else if ext == "geojson" {
		geojsonFile = importFile
	} else {
		fmt.Println("Unsupported file type", ext)
		os.Exit(1)
	}
	// Read .geojson file
	file, err := ioutil.ReadFile(geojsonFile)
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
	ds, _ := utils.NewUUID()
	app.DB.InsertLayer(ds, geojs)
	fmt.Println("Datasource created:", ds)
	// Cleanup artifacts
	if geojsonFile != importFile {
		os.Remove(geojsonFile)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: gospatial_cmd [method] [option]")
		fmt.Printf("Methods:\n")
		// CHANGE [ls] to [datasource ls] and [customer ls]
		fmt.Printf("  ls\n\tList all datasources from database\n")
		fmt.Printf("  export [datasource]\n\tExports datasource to GeoJSON file\n")
		fmt.Printf("  import [<filename>.shp || <filename>.geojson]\n\tImports datasource from shapefile or GeoJSON\n")
		fmt.Printf("  create [datasource || customer]\n\tCreates new datasource or customer\n")
		fmt.Printf("  assign [datasource] [customer]\n\tAssigns datasource to customer\n")
		fmt.Printf("\n")
		fmt.Printf("Defaults:\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.Parse()
}

func main() {

	requiredArgs := flag.Args()

	if len(requiredArgs) == 0 {
		usageError("No method provided")
	}

	method := requiredArgs[0]

	if method == "ls" {
		listDatasources()
	} else if method == "export" {
		if len(requiredArgs) != 2 {
			usageError("No datasource provided")
		}
		datasource := requiredArgs[1]
		exportDatasource(datasource)
	} else if method == "import" {
		if len(requiredArgs) != 2 {
			usageError("No file provided")
		}
		importFile := requiredArgs[1]
		importDatasource(importFile)
	} else if method == "create" {
		if len(requiredArgs) != 2 {
			usageError("Please specify either 'datasource' or 'customer' to create")
		} else if requiredArgs[1] == "datasource" {
			fmt.Println("Creating datasource")
			setupDb()
			ds, err := app.DB.NewLayer()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Datasource created:", ds)
		} else if requiredArgs[1] == "customer" {
			fmt.Println("Creating customer")
			setupDb()
			apikey := utils.NewAPIKey(12)
			customer := app.Customer{Apikey: apikey}
			err := app.DB.InsertCustomer(customer)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Customer created:", apikey)
		} else {
			usageError("Cannot create '" + requiredArgs[1] + "'")
		}
	} else if method == "backup" {
		fmt.Println("Backing up database...")
		setupDb()
		savefile := "backup_" + time.Now().String()
		app.DB.Backup(savefile)
		fmt.Println("Backup created:", savefile)
	} else if method == "load" {
		if len(requiredArgs) != 2 {
			usageError("Please provide a database to load")
		} else {
			filename := requiredArgs[1]
			fmt.Println("Loading database...")
			setupDb()
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
			var data dumpedDatabase
			data.Apikeys = make(map[string]app.Customer)
			data.Layers = make(map[string]*geojson.FeatureCollection)
			err = json.Unmarshal(file, &data)
			if err != nil {
				fmt.Println(err)
			}
			app.DB.InsertCustomers(data.Apikeys)
			app.DB.InsertLayers(data.Layers)
		}
	} else if method == "assign" {
		if len(requiredArgs) != 3 {
			usageError("Please datasource and customer key")
		} else {
			setupDb()
			customer, err := app.DB.GetCustomer(requiredArgs[2])
			if err != nil {
				fmt.Println("Customer key not found!")
				os.Exit(1)
			}
			// CHECK IF DATASOURCE EXISTS
			// CHECK IF DATASOURCE ALREADY ADDED TO CUSTOMER
			// Add datasource uuid to customer
			customer.Datasources = append(customer.Datasources, requiredArgs[1])
			app.DB.InsertCustomer(customer)
		}
	} else {
		usageError("Method not found")
	}
	// exit
	os.Exit(0)
}
