/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
//  requires: ogr2ogr
/*=======================================*/

package main

import (
	// "bytes"
	"flag"
	"fmt"
	"gospatial/app"
	// "io/ioutil"
	"os"
	"os/exec"
	"strings"
	"github.com/paulmach/go.geojson"
)

var (
	shapefile string
	database  string
	version   bool
)

const (
	VERSION string = "1.6.0"
)

func init() {
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.StringVar(&app.SuperuserKey, "s", "7q1qcqmsxnvw", "superuser key")
	flag.StringVar(&shapefile, "shp", "none", "shapefile to upload")
	flag.BoolVar(&version, "v", false, "App Version")
	flag.Parse()
	if version {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}
	if shapefile == "none" {
		fmt.Println("Incorrect usage")
		os.Exit(1)
	}
}

func main() {

	// Initiate Database
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()

	// Convert .shp to .geojson
	// ogr2ogr -f GeoJSON -t_srs crs:84 [name].geojson [name].shp
	geojson_file := strings.Replace(shapefile, ".shp", ".geojson", -1)
	fmt.Println("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojson_file, shapefile)
	out, err := exec.Command("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojson_file, shapefile).Output()
	if err != nil {
		app.Error.Println(err)
		app.Error.Println(string(out))
		os.Exit(1)
	} else {
		app.Info.Println(geojson_file, "created")
		app.Info.Println(string(out))
	}

	// Read .geojson file


}
