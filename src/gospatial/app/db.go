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

func (db *Database) connect() *bolt.DB {
	Trace.Printf("Database | Connecting to database: '%s'", db.File)
	conn, err := bolt.Open(db.File, 0644, nil)
	if err != nil {
		conn.Close()
		Error.Fatal(err)
	}
	// defer db.Close()
	return conn
}

func (db *Database) Init() error {
	Trace.Println("Database | Creating database")
	m := make(map[string]*LayerCache)
	db.Cache = m
	go db.CacheManager()
	conn := db.connect()
	Debug.Println("Database | Creating 'layers' bucket if not found")
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
	conn.Close()
	return err
}

func (db *Database) insertLayer(datasource string, geojs Geojson) error {
	// Caching layer
	Trace.Println("Database | Checking cache")
	if v, ok := db.Cache[datasource]; ok {
		Debug.Printf("Database | %s | Update cached datasource", datasource)
		v.Geojson = geojs
		v.Time = time.Now()
	} else {
		Debug.Printf("Database | %s | Add datasource to cache", datasource)
		pgc := &LayerCache{Geojson: geojs, Time: time.Now()}
		db.Cache[datasource] = pgc
	}
	// Connect to database
	conn := db.connect()
	var bucket = []byte("layers")
	key := []byte(datasource)
	// convert to bytes
	Debug.Printf("Database | %s | Encoding datasource to []byte", datasource)
	value, err := json.Marshal(geojs)
	if err != nil {
		Error.Println(err)
	}
	// Insert layer into database
	Debug.Printf("Database | %s | Inserting datasource to database", datasource)
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

func (db *Database) getLayer(datasource string) (Geojson, error) {
	// Caching layer
	// Trace.Println("Checking cache")
	if v, ok := db.Cache[datasource]; ok {
		Debug.Printf("Database | %s | Retrieve datasource from cache", datasource)
		v.Time = time.Now()
		return v.Geojson, nil
	}
	// If page not found get from database
	conn := db.connect()
	Debug.Printf("Database | %s | Retrieve datasource from database", datasource)
	var bucket = []byte("layers")
	key := []byte(datasource)
	val := []byte{}
	err := conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		conn.Close()
		DebugMode(true) // this shouldnt happen
		Error.Println(err)
		return Geojson{}, err
	}
	err = conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			DebugMode(true) // this shouldnt happen
			return fmt.Errorf("Database | Bucket %q not found!", bucket)
		}
		val = bucket.Get(key)
		return nil
	})
	if err != nil {
		conn.Close()
		DebugMode(true) // this shouldnt happen
		Error.Println(err)
		return Geojson{}, err
	}
	// Read to struct
	Debug.Printf("Database | %s | Unmarshal datasource", datasource)
	geojs := Geojson{}
	err = json.Unmarshal(val, &geojs)
	if err != nil {
		// Layer deleted?
		Error.Println(err)
	}
	conn.Close()
	// Store page in memory cache
	Debug.Printf("Database | %s | Add datasource to cache", datasource)
	pgc := &LayerCache{Geojson: geojs, Time: time.Now()}
	db.Cache[datasource] = pgc
	return geojs, nil
}

func (db *Database) deleteLayer(datasource string) error {
	// Connect to database
	Debug.Printf("Database | Connecting to database")
	conn := db.connect()
	var bucket = []byte("layers")
	key := []byte(datasource)
	// Insert layer into database
	Debug.Printf("Database | %s | Deleting datasource", datasource)
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
	delete(db.Cache, datasource)
	return err
}

// Manage Cache
func (db *Database) CacheManager() {
	for {
		if len(db.Cache) != 0 {
			Trace.Println("Database | Checking cache...")
			for key := range db.Cache {
				if time.Since(db.Cache[key].Time).Seconds() > 90 {
					Debug.Printf("Database | %s | Uploading datasource", key)
					delete(db.Cache, key)
				}
			}
			time.Sleep(15000 * time.Millisecond)
		} else {
			time.Sleep(60000 * time.Millisecond)
		}
	}
}
