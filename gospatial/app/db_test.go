package app

import (
	"github.com/paulmach/go.geojson"
	"testing"
)

// go test -bench=.
// go test -bench=. -test.benchmem

// Test NewLayer
// Benchmark InsertFeature
// Test InsertFeature

const (
	testDbFile         string = "./test.db"
	testCustomerApikey string = "testKey"
	testDatasource     string = "testLayer"
)

// Benchmark Database.InsertCustomer
func BenchmarkDbInsertCustomer(b *testing.B) {
	test_logger_init()
	testDb := Database{File: testDbFile}
	testDb.Init()
	testDb.TestLogger()
	testCustomer := Customer{Apikey: testCustomerApikey}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testDb.InsertCustomer(testCustomer)
	}
}

// Benchmark Database.getCustomer
func BenchmarkDbGetCustomerWithCache(b *testing.B) {
	test_logger_init()
	testDb := Database{File: testDbFile}
	testDb.Init()
	testDb.TestLogger()
	testCustomer := Customer{Apikey: testCustomerApikey}
	testDb.InsertCustomer(testCustomer)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testDb.GetCustomer(testCustomerApikey)
	}
}

func BenchmarkDbGetCustomerWithOutCache(b *testing.B) {
	test_logger_init()
	testDb := Database{File: testDbFile}
	testDb.Init()
	testDb.TestLogger()
	testCustomer := Customer{Apikey: testCustomerApikey}
	testDb.InsertCustomer(testCustomer)
	b.ResetTimer()
	testDb.Apikeys = make(map[string]Customer)
	for i := 0; i < b.N; i++ {
		testDb.GetCustomer(testCustomerApikey)
	}
}

// Unittest Database.GetCustomer
// Unittest Database.InsertCustomer
func TestDbCustomers(t *testing.T) {
	test_logger_init()
	testDb := Database{File: testDbFile}
	testDb.Init()
	testDb.TestLogger()
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
	test_logger_init()
	testDb := Database{File: testDbFile}
	testDb.Init()
	testDb.TestLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testDb.NewLayer()
	}
}

// Benchmark Database.InsertLayer
func BenchmarkDbInsertLayer(b *testing.B) {
	test_logger_init()
	testDb := Database{File: testDbFile}
	testDb.Init()
	testDb.TestLogger()
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
	test_logger_init()
	testDb := Database{File: testDbFile}
	testDb.Init()
	testDb.TestLogger()
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
	test_logger_init()
	testDb := Database{File: testDbFile}
	testDb.Init()
	testDb.TestLogger()
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

// Unittest: Database.GetLayer
// Unittest: Database.InsertLayer
func TestDbLayers(t *testing.T) {
	test_logger_init()
	testDb := Database{File: testDbFile}
	testDb.Init()
	testDb.TestLogger()
	data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"},{"geometry":{"coordinates":[[[-87.97851562499999,58.995311187950925],[-87.97851562499999,60.500525410511294],[-84.63867187499997,60.500525410511294],[-84.63867187499997,58.995311187950925],[-87.97851562499999,58.995311187950925]]],"type":"Polygon"},"properties":{"FID":1},"type":"Feature"}],"type":"FeatureCollection"}`)
	geojs, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		t.Error(err)
	}
	err = testDb.InsertLayer(testDatasource, geojs)
	if err != nil {
		t.Error(err)
	}
	_, err = testDb.GetLayer(testDatasource)
	if err != nil {
		t.Error(err)
	}
}
