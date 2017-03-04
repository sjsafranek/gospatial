/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"flag"
	"fmt"
	//"github.com/boltdb/bolt"
	"github.com/paulmach/go.geojson"
	"gospatial/app"
	"gospatial/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
		geojsonFile := strings.Replace(importFile, ".shp", ".geojson", -1)
		fmt.Println("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojsonFile, importFile)
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
		fmt.Printf("  import [<filename>.shp || <filename>.geojson]\n\tImports datasource from shapefile or GeoJSON\n")
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

	if method == "import" {
		if len(requiredArgs) != 2 {
			usageError("No file provided")
		}
		importFile := requiredArgs[1]
		importDatasource(importFile)
	} else {
		usageError("Method not found")
	}
	// exit
	os.Exit(0)
}
