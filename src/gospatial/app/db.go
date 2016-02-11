package app

import (
	"encoding/json"
	// "encoding/gob"
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
	Trace.Println("Connecting to database")
	conn, err := bolt.Open(db.File, 0644, nil)
	if err != nil {
		conn.Close()
		Error.Fatal(err)
	}
	// defer db.Close()
	return conn
}

func (db *Database) Init() error {
	Trace.Println("Creating database")
	m := make(map[string]*LayerCache)
	db.Cache = m
	go db.CacheManager()
	conn := db.connect()
	Debug.Println("Creating 'layers' bucket")
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
	Trace.Println("Checking cache")
	if v, ok := db.Cache[datasource]; ok {
		Debug.Printf("Update %s cache", datasource)
		v.Geojson = geojs
		v.Time = time.Now()
	} else {
		Debug.Printf("Insert %s to cache", datasource)
		pgc := &LayerCache{Geojson: geojs, Time: time.Now()}
		db.Cache[datasource] = pgc
	}
	// Connect to database
	conn := db.connect()
	var bucket = []byte("layers")
	key := []byte(datasource)
	// convert to bytes
	Debug.Printf("Encoding %s to []byte", datasource)
	value, err := json.Marshal(geojs)
	if err != nil {
		Error.Println(err)
	}
	// Insert layer into database
	Debug.Printf("Inserting %s to database", datasource)
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
	Trace.Println("Checking cache")
	if v, ok := db.Cache[datasource]; ok {
		Debug.Printf("Retrieve %s from cache", datasource)
		v.Time = time.Now()
		return v.Geojson, nil
	}
	// If page not found get from database
	conn := db.connect()
	Debug.Printf("Retrieve %s from database", datasource)
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
			return fmt.Errorf("Bucket %q not found!", bucket)
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
	Debug.Printf("Unmarshal %s", datasource)
	geojs := Geojson{}
	err = json.Unmarshal(val, &geojs)
	if err != nil {
		// Layer deleted?
		Error.Println(err)
	}
	conn.Close()
	// Store page in memory cache
	Debug.Printf("Insert %s to cache", datasource)
	pgc := &LayerCache{Geojson: geojs, Time: time.Now()}
	db.Cache[datasource] = pgc
	return geojs, nil
}

func (db *Database) deleteLayer(datasource string) error {
	// Connect to database
	Debug.Printf("Connecting to database")
	conn := db.connect()
	var bucket = []byte("layers")
	key := []byte(datasource)
	// Insert layer into database
	Debug.Printf("Deleteing %s", datasource)
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
		Trace.Println("Cache Manager loop")
		for key := range db.Cache {
			if time.Since(db.Cache[key].Time).Seconds() > 90 {
				Debug.Println("Cache remove", key)
				delete(db.Cache, key)
			}
		}
		time.Sleep(15000 * time.Millisecond)
	}
}
