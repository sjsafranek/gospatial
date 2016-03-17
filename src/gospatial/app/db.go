package app

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

var DB Database

type LayerCache struct {
	Geojson Geojson
	Time    time.Time
}

type Database struct {
	File  string
	Cache map[string]*LayerCache
}

func (self *Database) connect() *bolt.DB {
	Trace.Printf("Connecting to database: '%s'", self.File)
	conn, err := bolt.Open(self.File, 0644, nil)
	if err != nil {
		conn.Close()
		Error.Fatal(err)
	}
	// defer db.Close()
	return conn
}

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
	var bucket = []byte("layers")
	err := conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		Error.Fatal(err)
	}
	// permissions
	Debug.Println("Creating 'apikeys' bucket if not found")
	bucket = []byte("apikeys")
	err = conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		Error.Fatal(err)
	}
	// close and return err
	conn.Close()
	return err
}

func (self *Database) insertCustomer(customer Customer) (string, error) {
	// Connect to database
	conn := self.connect()
	var bucket = []byte("apikeys")
	// customer := Customer{Apikey: NewAPIKey(12)}
	key := []byte(customer.Apikey)
	// convert to bytes
	value, err := json.Marshal(customer)
	if err != nil {
		Error.Println(err)
	}
	// Insert layer into database
	Debug.Printf("Database insert apikey [%s]", customer.Apikey)
	err = conn.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		Error.Fatal(err)
	}
	conn.Close()
	return customer.Apikey, err
}

func (self *Database) getCustomer(apikey string) (Customer, error) {
	// If page not found get from database
	Debug.Printf("Database read apikey [%s]", apikey)
	conn := self.connect()
	// Make sure table exists
	var bucket = []byte("apikeys")
	err := conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		conn.Close()
		Error.Println(err)
		return Customer{}, err
	}
	// Get datasrouce from database
	key := []byte(apikey)
	val := []byte{}
	err = conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}
		val = bucket.Get(key)
		return nil
	})
	if err != nil {
		conn.Close()
		Error.Println(err)
		return Customer{}, err
	}
	// datasource not found
	if val == nil {
		conn.Close()
		Warning.Printf("Apikey not found [%s]", apikey)
		return Customer{}, fmt.Errorf("Apikey not found")
	}
	// Read to struct
	Debug.Printf("Unmarshal [%s]", apikey)
	customer := Customer{}
	err = json.Unmarshal(val, &customer)
	if err != nil {
		conn.Close()
		Error.Println(err)
		return Customer{}, err
	}
	// Close database connection
	conn.Close()
	return customer, nil
}

func (self *Database) insertLayer(datasource string, geojs Geojson) error {
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
	//
	var bucket = []byte("layers")
	key := []byte(datasource)
	// convert to bytes
	Debug.Printf("Encoding [%s]", datasource)
	value, err := json.Marshal(geojs)
	if err != nil {
		Error.Println(err)
	}
	// Insert layer into database
	Debug.Printf("Database insert [%s]", datasource)
	err = conn.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		Error.Fatal(err)
	}
	conn.Close()
	return err
}

func (self *Database) getLayer(datasource string) (Geojson, error) {
	// Caching layer
	if v, ok := self.Cache[datasource]; ok {
		Debug.Printf("Cache read [%s]", datasource)
		v.Time = time.Now()
		return v.Geojson, nil
	}
	// If page not found get from database
	Debug.Printf("Database read [%s]", datasource)
	conn := self.connect()
	// Make sure table exists
	var bucket = []byte("layers")
	err := conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		conn.Close()
		Error.Println(err)
		return Geojson{}, err
	}
	// Get datasrouce from database
	key := []byte(datasource)
	val := []byte{}
	err = conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}
		val = bucket.Get(key)
		return nil
	})
	if err != nil {
		conn.Close()
		Error.Println(err)
		return Geojson{}, err
	}
	// datasource not found
	if val == nil {
		conn.Close()
		Warning.Printf("Not found [%s]", datasource)
		return Geojson{}, fmt.Errorf("Not found")
	}
	// Read to struct
	Debug.Printf("Unmarshal [%s]", datasource)
	geojs := Geojson{}
	err = json.Unmarshal(val, &geojs)
	if err != nil {
		conn.Close()
		Error.Println(err)
		return Geojson{}, err
	}
	// Close database connection
	conn.Close()
	// Store page in memory cache
	Debug.Printf("Cache insert [%s]", datasource)
	pgc := &LayerCache{Geojson: geojs, Time: time.Now()}
	self.Cache[datasource] = pgc
	return geojs, nil
}

func (self *Database) deleteLayer(datasource string) error {
	// Connect to database
	Debug.Printf("Connecting to database")
	conn := self.connect()
	var bucket = []byte("layers")
	key := []byte(datasource)
	// Insert layer into database
	Debug.Printf("Database delete [%s]", datasource)
	err := conn.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
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

// Manage Cache
func (self *Database) CacheManager() {
	for {
		if len(self.Cache) != 0 {
			Trace.Println("Checking cache...")
			for key := range self.Cache {
				if time.Since(self.Cache[key].Time).Seconds() > 90 {
					Debug.Printf("Cache unload [%s]", key)
					delete(self.Cache, key)
				}
			}
			time.Sleep(15000 * time.Millisecond)
		} else {
			time.Sleep(60000 * time.Millisecond)
		}
	}
}
