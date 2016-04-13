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
	"time"
)

var (
	filename string
	database string
)

func init() {
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.StringVar(&filename, "f", "none", "dump or load database")
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

	if filename == "none" {
		// Backup database
		savefile := "backup_" + time.Now().String()
		app.DB.Backup(savefile)
	} else {
		app.Info.Printf("Loading database [%s]\n", filename)
		// check for file
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			app.Error.Fatal("File not found [" + filename + "]")
		}
		// open json file
		file, err := ioutil.ReadFile(filename)
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
		app.DB.InsertCustomers(data.Apikeys)
		app.DB.InsertLayers(data.Layers)
	}

	os.Exit(0)

}
