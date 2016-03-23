/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"encoding/json"
	"flag"
	"github.com/paulmach/go.geojson"
	"gospatial/app"
	"io/ioutil"
	"os"
)

var (
	option   string
	database string
)

func init() {
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.StringVar(&option, "o", "dump", "dump or load database")
	flag.Parse()
}

type DumpedDatabase struct {
	Apikeys map[string]app.Customer               `json:"apikeys"`
	Layers  map[string]*geojson.FeatureCollection `json:"layers"`
}

func main() {

	// Initiate Database
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()

	if option == "dump" {
		app.Info.Println("Dumping database...")
		// Create struct to store db data
		data := app.DB.Dump()
		// marshal to json
		b, err := json.Marshal(data)
		if err != nil {
			app.Error.Fatal(err)
		}
		// Write to file
		ioutil.WriteFile("dump.json", b, 0644)

	} else if option == "load" {
		app.Info.Println("Loading database...")
		// open json file
		file, err := ioutil.ReadFile("dump.json")
		if err != nil {
			app.Error.Fatal(err)
		}
		// unmarshal data
		var data DumpedDatabase
		data.Apikeys = make(map[string]app.Customer)
		data.Layers = make(map[string]*geojson.FeatureCollection)
		err = json.Unmarshal(file, &data)
		if err != nil {
			app.Error.Fatal(err)
		}
		// load api keus
		for k := range data.Apikeys {
			app.DB.InsertCustomer(data.Apikeys[k])
		}
		// load lauers
		for k := range data.Layers {
			app.DB.InsertLayer(k, data.Layers[k])
		}

	} else {
		app.Error.Fatal("Unknown option:", option)
	}

	os.Exit(0)

}
