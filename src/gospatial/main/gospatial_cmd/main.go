/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"encoding/json"
	"flag"
	"github.com/boltdb/bolt"
	// "github.com/paulmach/go.geojson"
	"fmt"
	"gospatial/app"
	"io/ioutil"
	"os"
)

var (
	datasource string
	database   string
	list       bool
	imp bool
	exp bool
)

func list_datsources() {

	fmt.Println("Datasources:")

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

	conn.Close()

}

func export_datasource() {
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()
	lyr, err := app.DB.GetLayer(datasource)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b, err := json.Marshal(lyr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Write to file
	savename := datasource + ".geojson"
	ioutil.WriteFile(savename, b, 0644)

}

func init() {
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.StringVar(&datasource, "ds", "", "datasource")
	flag.BoolVar(&list, "ls", false, "list datasources")
	// 
	flag.BoolVar(&imp, "i", false, "import")
	flag.BoolVar(&exp, "e", false, "export")
	// 
	flag.Parse()
}

func main() {

	if list {
		list_datsources()
		os.Exit(0)
	}

	if datasource == "" {
		fmt.Println("Incorrect usage!")
		os.Exit(1)
	}

	export_datasource()

	os.Exit(0)

}
