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
	flag.StringVar(&apikey, "a", "none", "apikey key")
	flag.Parse()
}

func main() {

	// Initiate Database
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()

	if apikey != "none" {
		// Get customer from database
		_, err := app.DB.GetCustomer(apikey)
		if err == nil {
			app.Error.Println("Customer exists")
			os.Exit(1)
		}
	} else {
		apikey = app.NewAPIKey(12)
	}

	// new customer
	customer := app.Customer{Apikey: apikey}
	err := app.DB.InsertCustomer(customer)
	if err != nil {
		app.Error.Println(err)
		os.Exit(1)
	}

	app.Info.Println("New Apikey created", apikey)

	os.Exit(0)

}
