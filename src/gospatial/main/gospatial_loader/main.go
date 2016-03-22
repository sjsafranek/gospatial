/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"flag"
	"fmt"
	"gospatial/app"
	"os"
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
}

func main() {

	// Initiate Database
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()

}
