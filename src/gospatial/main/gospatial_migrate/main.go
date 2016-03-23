/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
//  requires: ogr2ogr
/*=======================================*/

package main

import (
	"flag"
	"gospatial/app"
	"os"
)

var (
	apikey   string
	database string
)

func init() {
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.Parse()
}

func main() {

	// Initiate Database
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()

	os.Exit(0)

}
