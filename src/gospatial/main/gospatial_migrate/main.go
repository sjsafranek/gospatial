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

func main() {

	// Open db
	db, err := bolt.Open("./"+database+".db", 0600, nil)
	if err != nil {
		app.Error.Fatal(err)
	}
	defer db.Close()

	if option == "dump" {
		app.Info.Println("Dumping database...")

		data := make(map[string]map[string]interface{})
		data["apikeys"] = make(map[string]interface{})
		data["layers"] = make(map[string]interface{})

		// Get all layers
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("layers"))
			b.ForEach(func(k, v []byte) error {
				// fmt.Printf("key=%s, value=%s\n", k, v)
				// data["layers"][string(k)] = string(v)
				geojs := make(map[string]interface{})
				_ = json.Unmarshal(v, &geojs)
				data["apikeys"][string(k)] = geojs
				return nil
			})
			return nil
		})

		// apikey
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("apikeys"))
			b.ForEach(func(k, v []byte) error {
				// fmt.Printf("key=%s, value=%s\n", k, v)
				// data["apikeys"][string(k)] = string(v)
				val := make(map[string]interface{})
				_ = json.Unmarshal(v, &val)
				data["apikeys"][string(k)] = val
				return nil
			})
			return nil
		})

		b, err := json.Marshal(data)
		if err != nil {
			app.Error.Fatal(err)
		}
		// Write to file
		ioutil.WriteFile("dump.json", b, 0644)

	} else if option == "load" {
		app.Info.Println("Loading database...")
	} else {
		app.Error.Fatal("Unknown option:", option)
	}

	os.Exit(0)

}
