package app

import (
	//"errors"
	//"github.com/paulmach/go.geojson"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
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

{"method":"assign_datasource","apikey":"70P78vbIeKex","datasource":"bf1f964abdab49aea6739bf7f6b32867"}

{"method":"import_file","file":"springfield_projects_edit.geojson"}
{"method":"import_file","file":"test.geojson"}

{"method":"edit_feature"}
{"method":"edit_feature","data":{"datasource":"bf1f964abdab49aea6739bf7f6b32867","geo_id":"0"}}
{"method":"insert_feature","data":{"datasource":"bf1f964abdab49aea6739bf7f6b32867","feature":{}}}
{"method":"insert_feature","data":{"datasource":"bf1f964abdab49aea6739bf7f6b32867","feature":{"type":"Feature","geometry":{"type":"Point","coordinates":[47.279229,47.27922900257082]},"properties":{"date_created":1.487653451e+09,"date_modified":1.487653451e+09,"geo_id":"1487653451","is_active":true,"is_deleted":false,"name":"Dot"}}}}

*/

func parseRequest(message string) TcpMessage {
	req := TcpMessage{}
	err := json.Unmarshal([]byte(message), &req)
	if err != nil {
		log.Println(err)
	}
	return req
}

func parseResponse(message string) map[string]interface{} {
	//data := map[string]interface{}
	var data map[string]interface{}
	err := json.Unmarshal(message & data)
	if err != nil {
		log.Println(err)
	}
	return data
}

func TestTCPCreateApikeySuccess(t *testing.T) {
	req := parseRequest(`{"method": "create_apikey"}`)
	resp := testTcpServer.create_apikey(req)
	// check for error in response
	if strings.Contains(resp, `"status": "error"`) {
		t.Error(resp)
	}
	// check if "apikey" in response
	if !strings.Contains(resp, `"apikey":`) {
		t.Error(resp)
	}
}

/*
func TestTCPMethodError(t *testing.T) {
	req := parseRequest(`{"method": "this_is_not_a_supported_method"}`)
	resp := testTcpServer.create_apikey(req)

	log.Println(resp)

	// check for error in response
	if strings.Contains(resp, `"status": "error"`) {
		t.Error(resp)
	}
}
*/

func TestTCPInsertExportApikeySuccess(t *testing.T) {
	// create and insert apikey
	now := time.Now().Second()
	test_apikey := fmt.Sprintf("test_apikey_%s", now)
	req := parseRequest(`{"method": "insert_apikey", "data": { "apikey": "` + test_apikey + `" } }`)
	resp := testTcpServer.insert_apikey(req)
	// check for error in response
	if !strings.Contains(resp, `"status": "ok"`) {
		t.Error(resp)
	}
	// check if "test_apikey" in response
	if !strings.Contains(resp, test_apikey) {
		t.Error(resp)
	}
	//
	req = parseRequest(`{"method": "export_apikey", "apikey": "` + test_apikey + `" }`)
	resp = testTcpServer.export_apikey(req)
	// check if "test_apikey" in response
	if !strings.Contains(resp, test_apikey) {
		t.Error(resp)
	}
	//
	req = parseRequest(`{"method": "export_apikeys"}`)
	resp = testTcpServer.export_apikeys(req)
	// check if "test_apikey" in response
	if !strings.Contains(resp, test_apikey) {
		t.Error(resp)
	}
}

func TestTCPCreateExportDatasources(t *testing.T) {
	req := parseRequest(`{"method":"create_datasource"}`)
	resp := testTcpServer.create_datasource(req)
	// check for error in response
	if strings.Contains(resp, `"status": "error"`) {
		t.Error(resp)
	}
	// // check if "apikey" in response
	// if !strings.Contains(resp, `"apikey":`) {
	// 	t.Error(resp)
	// }
	log.Println(parseResponse(resp))
}

// {"method":"export_datasource","datasource":"bf1f964abdab49aea6739bf7f6b32867"}
// {"method":"export_datasources"}
