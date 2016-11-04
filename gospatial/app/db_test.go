package app

import (
	"errors"
	"github.com/paulmach/go.geojson"
	//"log"
	"testing"
)

// go test -bench=.
// go test -bench=. -test.benchmem

const (
	testDbFile         string = "./test.db"
	testCustomerApikey string = "testKey"
	testDatasource     string = "testLayer"
)

// Benchmark Database.InsertCustomer
func BenchmarkDbInsertCustomer(b *testing.B) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()
	testCustomer := Customer{Apikey: testCustomerApikey}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testDb.InsertCustomer(testCustomer)
	}
}

// Benchmark Database.getCustomer
func BenchmarkDbGetCustomerWithCache(b *testing.B) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()
	testCustomer := Customer{Apikey: testCustomerApikey}
	testDb.InsertCustomer(testCustomer)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testDb.GetCustomer(testCustomerApikey)
	}
}

func BenchmarkDbGetCustomerWithOutCache(b *testing.B) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()
	testCustomer := Customer{Apikey: testCustomerApikey}
	testDb.InsertCustomer(testCustomer)
	b.ResetTimer()
	testDb.Apikeys = make(map[string]Customer)
	for i := 0; i < b.N; i++ {
		delete(testDb.Apikeys, testCustomerApikey)
		testDb.GetCustomer(testCustomerApikey)
	}
}

// Unittest Database.GetCustomer
// Unittest Database.InsertCustomer
func TestDbCustomers(t *testing.T) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()
	testCustomer := Customer{Apikey: testCustomerApikey}
	err := testDb.InsertCustomer(testCustomer)
	if err != nil {
		t.Error(err)
	}
	customer, err := testDb.GetCustomer(testCustomerApikey)
	if err != nil {
		t.Error(err)
	}
	if customer.Apikey != testCustomerApikey {
		t.Errorf("Apikey does not match: %s %s", testCustomerApikey, customer.Apikey)
	}
}

// Benchmark Database.NewLayer
func BenchmarkDbNewLayer(b *testing.B) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testDb.NewLayer()
	}
}

// Benchmark Database.InsertLayer
func BenchmarkDbInsertLayer(b *testing.B) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()
	data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"},{"geometry":{"coordinates":[[[-87.97851562499999,58.995311187950925],[-87.97851562499999,60.500525410511294],[-84.63867187499997,60.500525410511294],[-84.63867187499997,58.995311187950925],[-87.97851562499999,58.995311187950925]]],"type":"Polygon"},"properties":{"FID":1},"type":"Feature"}],"type":"FeatureCollection"}`)
	geojs, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testDb.InsertLayer(testDatasource, geojs)
	}
}

// Benchmark Database.GetLayer
func BenchmarkDbGetLayerWithCache(b *testing.B) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()
	data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"},{"geometry":{"coordinates":[[[-87.97851562499999,58.995311187950925],[-87.97851562499999,60.500525410511294],[-84.63867187499997,60.500525410511294],[-84.63867187499997,58.995311187950925],[-87.97851562499999,58.995311187950925]]],"type":"Polygon"},"properties":{"FID":1},"type":"Feature"}],"type":"FeatureCollection"}`)
	geojs, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	testDb.InsertLayer(testDatasource, geojs)
	for i := 0; i < b.N; i++ {
		testDb.GetLayer(testDatasource)
	}
}

// Benchmark Database.GetLayer
func BenchmarkDbGetLayerWithoutCache(b *testing.B) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()
	data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"},{"geometry":{"coordinates":[[[-87.97851562499999,58.995311187950925],[-87.97851562499999,60.500525410511294],[-84.63867187499997,60.500525410511294],[-84.63867187499997,58.995311187950925],[-87.97851562499999,58.995311187950925]]],"type":"Polygon"},"properties":{"FID":1},"type":"Feature"}],"type":"FeatureCollection"}`)
	geojs, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	testDb.InsertLayer(testDatasource, geojs)
	for i := 0; i < b.N; i++ {
		delete(testDb.Cache, testDatasource)
		testDb.GetLayer(testDatasource)
	}
}

// Benchmark Database.InsertFeature
func BenchmarkDbInsertFeatureWithCache(b *testing.B) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()

	lyr_data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[],"type":"FeatureCollection"}`)
	layer, err := geojson.UnmarshalFeatureCollection(lyr_data)
	if err != nil {
		b.Error(err)
	}
	testDb.InsertLayer(testDatasource, layer)

	feat_data := []byte(`{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"}`)
	feature, err := geojson.UnmarshalFeature(feat_data)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = testDb.InsertFeature(testDatasource, feature)
		if err != nil {
			b.Error(err)
		}
	}
}

// Benchmark Database.InsertFeature
func BenchmarkDbInsertFeatureWithOutCache(b *testing.B) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()

	lyr_data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[],"type":"FeatureCollection"}`)
	layer, err := geojson.UnmarshalFeatureCollection(lyr_data)
	if err != nil {
		b.Error(err)
	}
	testDb.InsertLayer(testDatasource, layer)

	feat_data := []byte(`{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"}`)
	feature, err := geojson.UnmarshalFeature(feat_data)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delete(testDb.Cache, testDatasource)
		err = testDb.InsertFeature(testDatasource, feature)
		if err != nil {
			b.Error(err)
		}
	}
}

// Unittest: Database.GetLayer
// Unittest: Database.InsertLayer
func TestDbLayers(t *testing.T) {
	COMMIT_LOG_FILE = "./test_commit.log"
	testDb := Database{File: testDbFile}
	testDb.Init()
	enable_test_logging()
	data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"},{"geometry":{"coordinates":[[[-87.97851562499999,58.995311187950925],[-87.97851562499999,60.500525410511294],[-84.63867187499997,60.500525410511294],[-84.63867187499997,58.995311187950925],[-87.97851562499999,58.995311187950925]]],"type":"Polygon"},"properties":{"FID":1},"type":"Feature"}],"type":"FeatureCollection"}`)
	geojs, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		t.Error(err)
	}
	err = testDb.InsertLayer(testDatasource, geojs)
	if err != nil {
		t.Error(err)
	}
	lyr, err := testDb.GetLayer(testDatasource)
	if err != nil {
		t.Error(err)
	}
	if 2 != len(lyr.Features) {
		t.Error(errors.New("missing features!"))
	}
}

// Test NewLayer
// Test InsertFeature
