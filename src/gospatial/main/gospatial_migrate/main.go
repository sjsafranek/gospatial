/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"encoding/json"
	"flag"
	// "github.com/boltdb/bolt"
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

type apikey struct {
	Apikey      string   `json:"apikey"`
	Datasources []string `json:"datasources"`
}

type DumpedDatabase struct {
	Apikeys map[string]apikey                     `json:"apikeys"`
	Layers  map[string]*geojson.FeatureCollection `json:"layers"`
}

func main() {

	// Open db
	// db, err := bolt.Open("./"+database+".db", 0600, nil)
	// if err != nil {
	// 	app.Error.Fatal(err)
	// }
	// defer db.Close()

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
		data.Apikeys = make(map[string]apikey)
		data.Layers = make(map[string]*geojson.FeatureCollection)
		err = json.Unmarshal(file, &data)
		if err != nil {
			app.Error.Fatal(err)
		}
		app.Info.Printf("%v\n", data.Layers)
	} else {
		app.Error.Fatal("Unknown option:", option)
	}

	os.Exit(0)

}
