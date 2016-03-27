package app

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	// "unsafe"
	// "reflect"
	"time"
)

/*=======================================*/
// Gobals
/*=======================================*/
var DB Database

/*=======================================*/
// Models
/*=======================================*/
type LayerCache struct {
	Geojson *geojson.FeatureCollection
	Time    time.Time
}

type Database struct {
	File    string
	Cache   map[string]*LayerCache
	Apikeys map[string]Customer
}

/*=======================================*/
// Method: Database.connect
// Description:
//		Connects to database
//		Returns open database connection
// @returns *bolt.DB
/*=======================================*/
func (self *Database) connect() *bolt.DB {
	Trace.Printf("Connecting to database: '%s'", self.File)
	conn, err := bolt.Open(self.File, 0644, nil)
	if err != nil {
		conn.Close()
		Error.Fatal(err)
	}
	return conn
}

/*=======================================*/
// Method: Database.Init
// Description:
//		Creates database
//		Creates layers and apikey tables
//		Starts database caching
// @returns Error
/*=======================================*/
func (self *Database) Init() error {
	Trace.Println("Creating database")
	// Start db caching
	m := make(map[string]*LayerCache)
	self.Cache = m
	go self.CacheManager()
	// connect to db
	conn := self.connect()
	// datasources
	Debug.Println("Creating 'layers' bucket if not found")
	err := conn.Update(func(tx *bolt.Tx) error {
		table := []byte("layers")
		_, err := tx.CreateBucketIfNotExists(table)
		return err
	})
	if err != nil {
		Error.Fatal(err)
	}
	// permissions
	Debug.Println("Creating 'apikeys' bucket if not found")
	err = conn.Update(func(tx *bolt.Tx) error {
		table := []byte("apikeys")
		_, err := tx.CreateBucketIfNotExists(table)
		return err
	})
	if err != nil {
		Error.Fatal(err)
	}
	// put apikey into memory
	self.Apikeys = make(map[string]Customer)
	conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("apikeys"))
		b.ForEach(func(k, v []byte) error {
			// val := make(map[string]interface{})
			customer := Customer{}
			err := json.Unmarshal(v, &customer)
			if err != nil {
				conn.Close()
				Error.Fatal(err)
			}
			self.Apikeys[string(k)] = customer
			return nil
		})
		return nil
	})
	// close and return err
	conn.Close()
	return err
}

/*=======================================*/
// Method: Database.insertCustomer
// Description:
//		Inserts customer into apikeys table
// @param customer {Customer}
// @returns string apikey
// @returns Error
/*=======================================*/
func (self *Database) InsertCustomer(customer Customer) error {
	//
	self.Apikeys[customer.Apikey] = customer
	// Connect to database
	conn := self.connect()
	// convert to bytes
	table := []byte("apikeys")
	key := []byte(customer.Apikey)
	value, err := json.Marshal(customer)
	if err != nil {
		Error.Println(err)
	}
	// Insert layer into database
	Debug.Printf("Database insert apikey [%s]", customer.Apikey)
	err = conn.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(table)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		return err
	})
	if err != nil {
		Error.Fatal(err)
	}
	conn.Close()
	return err
}

/*=======================================*/
// Method: Database.getCustomer
// Description:
//		Gets customer from database
// @param apikey {string}
// @returns Customer
// @returns Error
/*=======================================*/
func (self *Database) GetCustomer(apikey string) (Customer, error) {
	// Check apikey cache
	if _, ok := self.Apikeys[apikey]; ok {
		Debug.Printf("Cache read [%s]", apikey)
		return self.Apikeys[apikey], nil
	}
	// If page not found get from database
	Debug.Printf("Database read apikey [%s]", apikey)
	conn := self.connect()
	// Make sure table exists
	table := []byte("apikeys")
	// Get datasrouce from database
	key := []byte(apikey)
	val := []byte{}
	err := conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(table)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}
		val = bucket.Get(key)
		return nil
	})
	if err != nil {
		// bucket doesnt exist
		conn.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(table)
			return err
		})
		conn.Close()
		Error.Println(err)
		return Customer{}, err
	}
	// datasource not found
	if val == nil {
		conn.Close()
		Warning.Printf("Customer not found [%s]", apikey)
		return Customer{}, fmt.Errorf("Apikey not found")
	}
	// Read to struct
	Debug.Printf("Unmarshal customer [%s]", apikey)
	customer := Customer{}
	err = json.Unmarshal(val, &customer)
	if err != nil {
		conn.Close()
		// Error.Printf("Cannot unmarshal customer [%s]", err)
		Error.Println(err)
		return Customer{}, err
	}
	// Close database connection
	conn.Close()
	return customer, nil
}

/*=======================================*/
// Method: Database.insertLayer
// Description:
//		Inserts layer into database
// @param datasource {string}
// @param geojs {Geojson}
// @returns Error
/*=======================================*/
func (self *Database) InsertLayer(datasource string, geojs *geojson.FeatureCollection) error {
	// Caching layer
	Trace.Println("Checking cache")
	if v, ok := self.Cache[datasource]; ok {
		Debug.Printf("Cache update [%s]", datasource)
		v.Geojson = geojs
		v.Time = time.Now()
	} else {
		Debug.Printf("Cache insert [%s]", datasource)
		pgc := &LayerCache{Geojson: geojs, Time: time.Now()}
		self.Cache[datasource] = pgc
	}
	// Connect to database
	conn := self.connect()
	key := []byte(datasource)
	// convert to bytes
	Debug.Printf("Encoding datasource [%s]", datasource)
	value, err := geojs.MarshalJSON()
	if err != nil {
		Error.Println(err)
	}
	// Insert layer into database
	Debug.Printf("Database insert datasource [%s]", datasource)
	err = conn.Update(func(tx *bolt.Tx) error {
		table := []byte("layers")
		bucket, err := tx.CreateBucketIfNotExists(table)
		if err != nil {
			return err
		}
		// err = bucket.Put(key, value)
		err = bucket.Put(key, compressByte(value))
		return err
	})
	if err != nil {
		Error.Fatal(err)
	}
	conn.Close()
	return err
}

/*=======================================*/
// Method: Database.getLayer
// Description:
//		Gets layer from database
// @param datasource {string}
// @returns Geojson
// @returns Error
/*=======================================*/
// func (self *Database) getLayer(datasource string) (Geojson, error) {
func (self *Database) GetLayer(datasource string) (*geojson.FeatureCollection, error) {
	// Caching layer
	if v, ok := self.Cache[datasource]; ok {
		Debug.Printf("Cache read [%s]", datasource)
		v.Time = time.Now()
		return v.Geojson, nil
	}
	// If page not found get from database
	Debug.Printf("Database read [%s]", datasource)
	conn := self.connect()
	table := []byte("layers")
	// Get datasrouce from database
	key := []byte(datasource)
	val := []byte{}
	err := conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(table)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}
		val = decompressByte(bucket.Get(key))
		return nil
	})
	if err != nil {
		// bucket doesnt exist
		conn.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(table)
			return err
		})
		conn.Close()
		Error.Println(err)
		return nil, err
	}
	// datasource not found
	if val == nil {
		conn.Close()
		Warning.Printf("Datasource not found [%s]", datasource)
		return nil, fmt.Errorf("Not found")
	}
	// Read to struct
	Debug.Printf("Unmarshal datasource [%s]", datasource)

	geojs, err := geojson.UnmarshalFeatureCollection(val)
	if err != nil {
		conn.Close()
		Error.Println(err)
		return geojs, err
	}
	// Close database connection
	conn.Close()
	// Store page in memory cache
	Debug.Printf("Cache insert [%s]", datasource)
	pgc := &LayerCache{Geojson: geojs, Time: time.Now()}
	self.Cache[datasource] = pgc
	return geojs, nil
}

/*=======================================*/
// Method: Database.deleteLayer
// Description:
//		Deletes layer from database
// @param datasource {string}
// @returns Error
/*=======================================*/
func (self *Database) DeleteLayer(datasource string) error {
	// Connect to database
	conn := self.connect()
	key := []byte(datasource)
	// Insert layer into database
	Debug.Printf("Database delete [%s]", datasource)
	err := conn.Update(func(tx *bolt.Tx) error {
		table := []byte("layers")
		bucket, err := tx.CreateBucketIfNotExists(table)
		if err != nil {
			return err
		}
		err = bucket.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		Error.Println(err)
	}
	conn.Close()
	delete(self.Cache, datasource)
	return err
}

/*=======================================*/
// Method: Dumps database
/*=======================================*/
func (self *Database) Dump() map[string]map[string]interface{} {
	Info.Println("Extract all data from database...")
	conn := self.connect()
	// Create struct to store db data
	data := make(map[string]map[string]interface{})
	data["apikeys"] = make(map[string]interface{})
	data["layers"] = make(map[string]interface{})
	// Get all layers
	conn.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("layers"))
		b.ForEach(func(k, v []byte) error {
			geojs := make(map[string]interface{})
			err := json.Unmarshal(decompressByte(v), &geojs)
			if err != nil {
				conn.Close()
				Error.Fatal(err)
			}
			data["layers"][string(k)] = geojs
			return nil
		})
		return nil
	})
	// apikey
	conn.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("apikeys"))
		b.ForEach(func(k, v []byte) error {
			val := make(map[string]interface{})
			err := json.Unmarshal(v, &val)
			if err != nil {
				conn.Close()
				Error.Fatal(err)
			}
			data["apikeys"][string(k)] = val
			return nil
		})
		return nil
	})
	//
	conn.Close()
	return data
}

func (self *Database) Backup(filename ...string) {
	Info.Println("Backing up database...")
	// Create struct to store db data
	data := self.Dump()
	// marshal to json
	b, err := json.Marshal(data)
	if err != nil {
		Error.Fatal(err)
	}
	// Write to file
	savename := "backup_" + time.Now().String() + ".json"
	if len(filename) > 0 {
		savename = filename[0] + ".json"
	}
	ioutil.WriteFile(savename, b, 0644)
}

/*=======================================*/
// Method: Database.CacheManager
// Description:
//		Database caching layer
//		Unloads layers older than 90 sec
//		When empty --> 60 sec timer
//		When items in cache --> 15 sec timer
/*=======================================*/
func (self *Database) CacheManager() {
	for {
		n := len(self.Cache)
		// timout := 90000
		if n != 0 {
			// Trace.Println("Checking cache...")
			limit := 90.0
			switch {
			case n > 5:
				limit = 60.0
			case n > 10:
				limit = 45.0
			case n > 15:
				limit = 30.0
			case n > 20:
				limit = 15.0
			}
			for key := range self.Cache {
				// s := unsafe.Sizeof(self.Cache[key])
				// s := reflect.TypeOf(self.Cache[key]).Size()
				// Info.Println(s)
				if time.Since(self.Cache[key].Time).Seconds() > limit {
					Debug.Printf("Cache unload [%s]", key)
					delete(self.Cache, key)
				}
			}
			time.Sleep(time.Duration(15000) * time.Millisecond)
		} else {
			time.Sleep(time.Duration(15000) * time.Millisecond)
		}
	}
}
