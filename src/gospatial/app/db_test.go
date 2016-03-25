package app

import (
	// "fmt"
	"github.com/paulmach/go.geojson"
	"testing"
)

// go test -bench=.
// go test -bench=. -test.benchmem

const (
	test_db_file         string = "./test.db"
	test_customer_apikey string = "testKey"
	test_datasource      string = "testLayer"
)

/*=======================================*/
// Benchmark Database.InsertCustomer
/*=======================================*/
func BenchmarkDbInsertCustomer(b *testing.B) {
	TestMode()
	test_db := Database{File: test_db_file}
	test_db.Init()
	test_customer := Customer{Apikey: test_customer_apikey}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		test_db.InsertCustomer(test_customer)
	}
}

/*=======================================*/
// Benchmark Database.getCustomer
/*=======================================*/
func BenchmarkDbGetCustomer(b *testing.B) {
	TestMode()
	test_db := Database{File: test_db_file}
	test_db.Init()
	test_customer := Customer{Apikey: test_customer_apikey}
	test_db.InsertCustomer(test_customer)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		test_db.GetCustomer(test_customer_apikey)
	}
}

/*=======================================*/
// Unittest Database.getCustomer
// Unittest Database.InsertCustomer
/*=======================================*/
func TestDbCustomers(t *testing.T) {
	TestMode()
	test_db := Database{File: test_db_file}
	test_db.Init()
	test_customer := Customer{Apikey: test_customer_apikey}
	err := test_db.InsertCustomer(test_customer)
	if err != nil {
		t.Error(err)
	}
	customer, err := test_db.GetCustomer(test_customer_apikey)
	if err != nil {
		t.Error(err)
	}
	if customer.Apikey != test_customer_apikey {
		t.Errorf("Apikey does not match: %s %s", test_customer_apikey, customer.Apikey)
	}
}

/*=======================================*/
// Benchmark Database.InsertLayer
/*=======================================*/
func BenchmarkDbInsertLayer(b *testing.B) {
	TestMode()
	test_db := Database{File: test_db_file}
	test_db.Init()
	data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"},{"geometry":{"coordinates":[[[-87.97851562499999,58.995311187950925],[-87.97851562499999,60.500525410511294],[-84.63867187499997,60.500525410511294],[-84.63867187499997,58.995311187950925],[-87.97851562499999,58.995311187950925]]],"type":"Polygon"},"properties":{"FID":1},"type":"Feature"}],"type":"FeatureCollection"}`)
	geojs, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		test_db.InsertLayer(test_datasource, geojs)
	}
}

/*=======================================*/
// Benchmark Database.GetLayer
/*=======================================*/
func BenchmarkDbGetLayerWithCache(b *testing.B) {
	TestMode()
	// StandardMode()
	test_db := Database{File: test_db_file}
	test_db.Init()
	data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"},{"geometry":{"coordinates":[[[-87.97851562499999,58.995311187950925],[-87.97851562499999,60.500525410511294],[-84.63867187499997,60.500525410511294],[-84.63867187499997,58.995311187950925],[-87.97851562499999,58.995311187950925]]],"type":"Polygon"},"properties":{"FID":1},"type":"Feature"}],"type":"FeatureCollection"}`)
	geojs, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	test_db.InsertLayer(test_datasource, geojs)
	for i := 0; i < b.N; i++ {
		test_db.GetLayer(test_datasource)
	}
}

/*=======================================*/
// Benchmark Database.GetLayer
/*=======================================*/
func BenchmarkDbGetLayerWithoutCache(b *testing.B) {
	TestMode()
	// StandardMode()
	test_db := Database{File: test_db_file}
	test_db.Init()
	data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"},{"geometry":{"coordinates":[[[-87.97851562499999,58.995311187950925],[-87.97851562499999,60.500525410511294],[-84.63867187499997,60.500525410511294],[-84.63867187499997,58.995311187950925],[-87.97851562499999,58.995311187950925]]],"type":"Polygon"},"properties":{"FID":1},"type":"Feature"}],"type":"FeatureCollection"}`)
	geojs, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	test_db.InsertLayer(test_datasource, geojs)
	for i := 0; i < b.N; i++ {
		delete(test_db.Cache, test_datasource)
		test_db.GetLayer(test_datasource)
	}
}

/*=======================================*/
// Unittest: Database.GetLayer
// Unittest: Database.InsertLayer
/*=======================================*/
func TestDbLayers(t *testing.T) {
	TestMode()
	test_db := Database{File: test_db_file}
	test_db.Init()
	data := []byte(`{"crs":{"properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"},"type":"name"},"features":[{"geometry":{"coordinates":[[[-76.64062,50.73645513701065],[-76.64062,65.65827451982659],[-38.67187,65.65827451982659],[-38.67187,50.73645513701065],[-76.64062,50.73645513701065]]],"type":"Polygon"},"properties":{"FID":0},"type":"Feature"},{"geometry":{"coordinates":[[[-87.97851562499999,58.995311187950925],[-87.97851562499999,60.500525410511294],[-84.63867187499997,60.500525410511294],[-84.63867187499997,58.995311187950925],[-87.97851562499999,58.995311187950925]]],"type":"Polygon"},"properties":{"FID":1},"type":"Feature"}],"type":"FeatureCollection"}`)
	geojs, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		t.Error(err)
	}
	err = test_db.InsertLayer(test_datasource, geojs)
	if err != nil {
		t.Error(err)
	}
	// geojs, err = test_db.GetLayer(test_datasource)
	_, err = test_db.GetLayer(test_datasource)
	if err != nil {
		t.Error(err)
	}
	// fmt.Printf("%v\n", geojs)
}
