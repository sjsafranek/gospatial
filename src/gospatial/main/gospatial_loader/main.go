/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
//  requires: ogr2ogr
/*=======================================*/

package main

import (
	"flag"
	"fmt"
	"github.com/paulmach/go.geojson"
	"gospatial/app"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	upload_file string
	database    string
	apikey      string
)

func init() {
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.StringVar(&apikey, "a", "7q1qcqmsxnvw", "apikey key")
	flag.StringVar(&upload_file, "f", "", "shapefile or geojson to upload")
	flag.Parse()
	if upload_file == "" {
		fmt.Println("Incorrect usage")
		os.Exit(1)
	}
}

func main() {

	// Initiate Database
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()

	// Get customer from database
	customer, err := app.DB.GetCustomer(apikey)
	if err != nil {
		app.Error.Println(err)
		os.Exit(1)
	}

	var geojson_file string
	ext := strings.Split(upload_file, ".")[1]
	if ext == "shp" {
		// Convert .shp to .geojson
		// ogr2ogr -f GeoJSON -t_srs crs:84 [name].geojson [name].shp
		geojson_file := strings.Replace(upload_file, ".shp", ".geojson", -1)
		// fmt.Println("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojson_file, shapefile)
		out, err := exec.Command("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojson_file, upload_file).Output()
		if err != nil {
			app.Error.Println(err)
			app.Error.Println(string(out))
			os.Exit(1)
		} else {
			app.Info.Println(geojson_file, "created")
			app.Info.Println(string(out))
		}
	} else if ext == "geojson" {
		geojson_file = upload_file
	} else {
		app.Error.Println("Unsupported file type", ext)
		os.Exit(1)
	}

	// Read .geojson file
	file, err := ioutil.ReadFile(geojson_file)
	if err != nil {
		app.Error.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal to geojson struct
	geojs, err := geojson.UnmarshalFeatureCollection(file)
	if err != nil {
		app.Error.Printf("Unmarshal GeoJSON error: %v\n", err)
		os.Exit(1)
	}

	// Create datasource
	ds, _ := app.NewUUID()
	app.DB.InsertLayer(ds, geojs)
	app.Info.Println(ds, "created")

	// Add datasource uuid to customer
	customer.Datasources = append(customer.Datasources, ds)
	app.DB.InsertCustomer(customer)
	app.Info.Println(ds, "added to", apikey)

	// Cleanup
	if geojson_file != upload_file {
		os.Remove(geojson_file)
	}

	os.Exit(0)

}
