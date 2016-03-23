/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"gospatial/app"
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
				data["layers"][string(k)] = v
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
				data["apikeys"][string(k)] = v
				return nil
			})
			return nil
		})

		//
		fmt.Printf("%s\n", data)

	} else if option == "load" {
		app.Info.Println("Loading database...")
	} else {
		app.Error.Fatal("Unknown option:", option)
	}

	os.Exit(0)

}
