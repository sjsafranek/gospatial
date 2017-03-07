package app

import (
	//"errors"
	// /"github.com/paulmach/go.geojson"
	"encoding/json"
	"log"
	"testing"
)

// go test -bench=.
// go test -bench=. -test.benchmem

// const (
// 	testDbFile         string = "./test.db"
// 	testCustomerApikey string = "testKey"
// 	testDatasource     string = "testLayer"
// )

//var testDb Database
var testTcpServer TcpServer

func init() {
	COMMIT_LOG_FILE = "./test_commit.log"
	DB = Database{File: testDbFile}
	DB.Init()
	enable_test_logging()
	testTcpServer = TcpServer{Host: "localhost", Port: "3333"}
	testTcpServer.Start()
}

/*

{"method":"export_datasources"}
{"method":"export_apikeys"}
{"method": "insert_apikey", "data":{ "apikey": "test"}}

{"method":"create_apikey"}
{"method":"create_datasource"}
{"method":"assign_datasource","apikey":"70P78vbIeKex","datasource":"bf1f964abdab49aea6739bf7f6b32867"}
{"method":"export_apikey","apikey":"70P78vbIeKex"}
{"method":"export_datasource","datasource":"bf1f964abdab49aea6739bf7f6b32867"}

{"method":"import_file","file":"springfield_projects_edit.geojson"}
{"method":"import_file","file":"test.geojson"}

{"method":"edit_feature"}
{"method":"edit_feature","data":{"datasource":"bf1f964abdab49aea6739bf7f6b32867","geo_id":"0"}}
{"method":"insert_feature","data":{"datasource":"bf1f964abdab49aea6739bf7f6b32867","feature":{}}}


*/

func parseRequest(message string) TcpMessage {
	req := TcpMessage{}
	err := json.Unmarshal([]byte(message), &req)
	if err != nil {
		log.Println(err)
	}
	return req
}

func TestTCPCreateApikey(t *testing.T) {
	req := parseRequest(`{"method": "create_apikey"}`)
	resp := testTcpServer.create_apikey(req)
	log.Println(resp)

	if resp != "" {
		t.Error(resp)
	}
	/*
		if customer.Apikey != testCustomerApikey {
			t.Errorf("Apikey does not match: %s %s", testCustomerApikey, customer.Apikey)
		}
	*/
}
