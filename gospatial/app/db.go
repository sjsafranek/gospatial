package app

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/paulmach/go.geojson"
	"gospatial/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// DB application Database
var DB Database
var dbLog io.Writer

// LayerCache keeps track of Database's loaded geojson layers
type LayerCache struct {
	Geojson *geojson.FeatureCollection
	Time    time.Time
}

// Database strust for application.
type Database struct {
	File    string
	Cache   map[string]*LayerCache
	Apikeys map[string]Customer
	Logger  *log.Logger
	guard   sync.RWMutex
}

// Connect to bolt database. Returns open database connection.
// @returns *bolt.DB
func (self *Database) connect() *bolt.DB {
	conn, err := bolt.Open(self.File, 0644, nil)
	if err != nil {
		conn.Close()
		// log.Fatal(err)
		panic(err)
	}
	return conn
}

// Init creates bolt database if existing one not found.
// Creates layers and apikey tables. Starts database caching for layers
// @returns Error
func (self *Database) Init() error {
	// Start db caching
	m := make(map[string]*LayerCache)
	self.Cache = m
	go self.cacheManager()
	self.startLogger()
	// connect to db
	conn := self.connect()
	defer conn.Close()
	// datasources
	err := conn.Update(func(tx *bolt.Tx) error {
		table := []byte("layers")
		_, err := tx.CreateBucketIfNotExists(table)
		return err
	})
	if err != nil {
		return err
	}
	// Add table for datasource owner
	//
	// permissions
	err = conn.Update(func(tx *bolt.Tx) error {
		table := []byte("apikeys")
		_, err := tx.CreateBucketIfNotExists(table)
		return err
	})
	if err != nil {
		return err
	}
	// create apikey/customer cache
	self.Apikeys = make(map[string]Customer)
	// put apikey into memory
	// conn.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("apikeys"))
	// 	b.ForEach(func(k, v []byte) error {
	// 		customer := Customer{}
	// 		err := json.Unmarshal(v, &customer)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		self.Apikeys[string(k)] = customer
	// 		return nil
	// 	})
	// 	return nil
	// })
	// close and return err
	return err
}

// Starts Database logger
func (self *Database) startLogger() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dbLogFile := strings.Replace(dir, "bin", "log/db.log", -1)
	dbLog, err := os.OpenFile(dbLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		dbLog, err = os.OpenFile("db.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}
	self.Logger = log.New(dbLog, "WRITE [DB] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

// TestLogger starts Database logger for db_test.go
func (self *Database) TestLogger() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	f := filepath.Join(dir, "test_db.log")
	dbLog, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	self.Logger = log.New(dbLog, "WRITE [DB] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

// InsertCustomer inserts customer into apikeys table
// @param customer {Customer}
// @returns Error
func (self *Database) InsertCustomer(customer Customer) error {
	self.Apikeys[customer.Apikey] = customer
	// Connect to database
	conn := self.connect()
	defer conn.Close()
	// convert to bytes
	table := []byte("apikeys")
	key := []byte(customer.Apikey)
	value, err := json.Marshal(customer)
	if err != nil {
		return err
	}
	self.Logger.Println(`{"method": "insert_apikey", "data":` + string(value) + `}`)
	// Insert customer into database
	err = conn.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(table)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		return err
	})
	if err != nil {
		panic(err)
	}
	return err
}

// InsertCustomers inserts list of customer into apikeys table
// @param customer {Customer}
// @returns Error
func (self *Database) InsertCustomers(customers map[string]Customer) error {
	// Connect to database
	conn := self.connect()
	defer conn.Close()
	for i := range customers {
		customer := customers[i]
		// convert to bytes
		table := []byte("apikeys")
		key := []byte(customer.Apikey)
		value, err := json.Marshal(customer)
		if err != nil {
			return err
		}
		self.Logger.Println(`{"method": "insert_apikey", "data": ` + string(value) + `}`)
		// Insert layer into database
		err = conn.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(table)
			if err != nil {
				return err
			}
			err = bucket.Put(key, value)
			return err
		})
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// GetCustomer returns customer from database
// @param apikey {string}
// @returns Customer
// @returns Error
func (self *Database) GetCustomer(apikey string) (Customer, error) {
	// Check apikey cache
	if _, ok := self.Apikeys[apikey]; ok {
		return self.Apikeys[apikey], nil
	}
	// If customer not found get from database
	conn := self.connect()
	defer conn.Close()
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
		// bucket doesn't exist
		conn.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(table)
			return err
		})
		return Customer{}, err
	}
	// datasource not found
	if val == nil {
		return Customer{}, fmt.Errorf("Apikey not found")
	}
	// Read to struct
	customer := Customer{}
	err = json.Unmarshal(val, &customer)
	if err != nil {
		// panic(err)
		return Customer{}, err
	}
	// Put apikey into cache
	self.Apikeys[apikey] = customer
	// Close database connection
	return customer, nil
}

// NewLayer creates new datasource layer
// @returns string - datasource id
// @returns Error
func (self *Database) NewLayer() (string, error) {
	// create geojson
	datasource, _ := utils.NewUUID()
	geojs := geojson.NewFeatureCollection()
	// Connect to database
	conn := self.connect()
	defer conn.Close()
	key := []byte(datasource)
	// convert to bytes
	value, err := geojs.MarshalJSON()
	if err != nil {
		return "", nil
	}
	self.Logger.Println(`{"method": "new_layer", "data": { "datasource": "` + datasource + `", "layer": ` + string(value) + `}}`)
	// Insert layer into database
	err = conn.Update(func(tx *bolt.Tx) error {
		table := []byte("layers")
		bucket, err := tx.CreateBucketIfNotExists(table)
		if err != nil {
			return err
		}
		err = bucket.Put(key, self.compressByte(value))
		return err
	})
	if err != nil {
		panic(err)
	}
	return datasource, err
}

// InsertLayer inserts layer into database
// @param datasource {string}
// @param geojs {Geojson}
// @returns Error
func (self *Database) InsertLayer(datasource string, geojs *geojson.FeatureCollection) error {
	// Caching layer
	// self.guard.Lock()
	if v, ok := self.Cache[datasource]; ok {
		self.guard.Lock()
		v.Geojson = geojs
		v.Time = time.Now()
		self.guard.Unlock()
	} else {
		pgc := &LayerCache{Geojson: geojs, Time: time.Now()}
		self.Cache[datasource] = pgc
	}
	// self.guard.Unlock()

	// Connect to database
	conn := self.connect()
	defer conn.Close()
	key := []byte(datasource)
	// convert to bytes
	value, err := geojs.MarshalJSON()
	if err != nil {
		return err
	}
	// self.Logger.Println(`{"method": "insert_layer", "data": { "datasource": ` + datasource + `, "layer": ` + string(value) + `}}`)
	// Insert layer into database
	ServerLogger.Debug("Database insert datasource [%s]", datasource)
	err = conn.Update(func(tx *bolt.Tx) error {
		table := []byte("layers")
		bucket, err := tx.CreateBucketIfNotExists(table)
		if err != nil {
			return err
		}
		err = bucket.Put(key, self.compressByte(value))
		return err
	})
	if err != nil {
		panic(err)
	}
	return err
}

// InsertLayers inserts a map of datasource layers into Database.
// @param datasources map[string]*geojson.FeatureCollection
// @returns Error
func (self *Database) InsertLayers(datsources map[string]*geojson.FeatureCollection) error {
	// Connect to database
	conn := self.connect()
	defer conn.Close()
	for datasource, geojs := range datsources {
		key := []byte(datasource)
		// convert to bytes
		value, err := geojs.MarshalJSON()
		if err != nil {
			return err
		}
		// self.Logger.Println(`{"method": "insert_layer", "data": { "datasource": ` + datasource + `, "layer": ` + string(value) + `}}`)
		// Insert layer into database
		err = conn.Update(func(tx *bolt.Tx) error {
			table := []byte("layers")
			bucket, err := tx.CreateBucketIfNotExists(table)
			if err != nil {
				return err
			}
			err = bucket.Put(key, self.compressByte(value))
			return err
		})
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// GetLayer returns layer from database
// @param datasource {string}
// @returns Geojson
// @returns Error
func (self *Database) GetLayer(datasource string) (*geojson.FeatureCollection, error) {
	// Caching layer
	if v, ok := self.Cache[datasource]; ok {
		self.guard.RLock()
		v.Time = time.Now()
		self.guard.RUnlock()
		return v.Geojson, nil
	}
	// If page not found get from database
	conn := self.connect()
	defer conn.Close()
	table := []byte("layers")
	// Get datasrouce from database
	key := []byte(datasource)
	val := []byte{}
	err := conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(table)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}
		val = self.decompressByte(bucket.Get(key))
		return nil
	})
	if err != nil {
		// bucket doesn't exist
		conn.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(table)
			return err
		})
		return nil, err
	}
	// datasource not found
	if val == nil {
		return nil, fmt.Errorf("Not found")
	}
	// Read to struct
	geojs, err := geojson.UnmarshalFeatureCollection(val)
	if err != nil {
		return geojs, err
	}
	// Store page in memory cache
	pgc := &LayerCache{Geojson: geojs, Time: time.Now()}
	self.Cache[datasource] = pgc
	return geojs, nil
}

// DeleteLayer deletes layer from database
// @param datasource {string}
// @returns Error
func (self *Database) DeleteLayer(datasource string) error {
	// Connect to database
	conn := self.connect()
	defer conn.Close()
	key := []byte(datasource)
	//
	self.Logger.Println(string(`{"method": "delete_layer", "data": { "datasource": "` + datasource + `"}}`))
	// Insert layer into database
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
		panic(err)
	}
	self.guard.Lock()
	delete(self.Cache, datasource)
	self.guard.Unlock()
	return err
}

// InsertFeature adds feature to layer. Updates layer in Database
// @param datasource {string}
// @param feat {Geojson Feature}
// @returns Error
func (self *Database) InsertFeature(datasource string, feat *geojson.Feature) error {
	// Get layer from database
	featCollection, err := self.GetLayer(datasource)
	if err != nil {
		return err
	}
	// Add new feature to layer
	value, err := feat.MarshalJSON()
	if err != nil {
		return err
	}
	self.Logger.Println(`{"method": "insert_feature", "data": { "datasource": "` + datasource + `", "feature": ` + string(value) + `}}`)
	featCollection.AddFeature(feat)
	err = self.InsertLayer(datasource, featCollection)
	if err != nil {
		panic(err)
	}
	return err

}

// Backup dumps database contents to json file
func (self *Database) Backup(filename ...string) {

	// Connect to database
	conn := self.connect()
	defer conn.Close()

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
			err := json.Unmarshal(self.decompressByte(v), &geojs)
			if err != nil {
				ServerLogger.Error(err)
			}
			data["layers"][string(k)] = geojs
			return nil
		})
		return nil
	})

	// Get all apikeyss
	conn.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("apikeys"))
		b.ForEach(func(k, v []byte) error {
			val := make(map[string]interface{})
			err := json.Unmarshal(v, &val)
			if err != nil {
				ServerLogger.Error(err)
			}
			data["apikeys"][string(k)] = val
			return nil
		})
		return nil
	})

	// Write data to json file
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	// Write to file
	savename := "backup_" + time.Now().String() + ".json"
	if len(filename) > 0 {
		savename = filename[0] + ".json"
	}
	ioutil.WriteFile(savename, b, 0644)
}

// cacheManager for Database. Stores layers in memory.
//		Unloads layers older than 90 sec
//		When empty --> 60 sec timer
//		When items in cache --> 15 sec timer
func (self *Database) cacheManager() {
	for {
		n := float64(len(self.Cache))
		if n != 0 {
			for key := range self.Cache {
				// CHECK AVAILABLE SYSTEM MEMORY
				f := float64(len(self.Cache[key].Geojson.Features))
				limit := (300.0 - (f * (f * 0.25))) - (n * 2.0)
				if limit < 0.0 {
					limit = 10.0
				}
				if time.Since(self.Cache[key].Time).Seconds() > limit {
					self.guard.Lock()
					delete(self.Cache, key)
					self.guard.Unlock()
				}
			}
		}
		time.Sleep(10000 * time.Millisecond)
	}
}

// Methods: Compression
// Source: https://github.com/schollz/gofind/blob/master/utils.go#L146-L169
//         https://github.com/schollz/gofind/blob/master/fingerprint.go#L43-L54
// Description:
//		Compress and Decompress bytes
func (self *Database) compressByte(src []byte) []byte {
	compressedData := new(bytes.Buffer)
	self.compress(src, compressedData, 9)
	return compressedData.Bytes()
}

func (self *Database) decompressByte(src []byte) []byte {
	compressedData := bytes.NewBuffer(src)
	deCompressedData := new(bytes.Buffer)
	self.decompress(compressedData, deCompressedData)
	return deCompressedData.Bytes()
}

func (self *Database) compress(src []byte, dest io.Writer, level int) {
	compressor, _ := flate.NewWriter(dest, level)
	compressor.Write(src)
	compressor.Close()
}

func (self *Database) decompress(src io.Reader, dest io.Writer) {
	decompressor := flate.NewReader(src)
	io.Copy(dest, decompressor)
	decompressor.Close()
}
