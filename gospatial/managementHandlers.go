package gospatial

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

import (
	"./utils"
)

// PingHandler provides an api route for server health check
func PingHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	var data map[string]interface{}
	data = make(map[string]interface{})
	data["status"] = "success"
	result := make(map[string]interface{})
	result["result"] = "pong"
	result["registered"] = startTime.UTC()
	result["uptime"] = time.Since(startTime).Seconds()
	result["num_cores"] = runtime.NumCPU()
	data["data"] = result

	js, err := MarshalJsonFromStruct(w, r, data)
	if err != nil {
		return
	}

	SendJsonResponse(w, r, js)
}

// NewCustomerHandler superuser route to create new api customers/apikeys
func NewCustomerHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Check auth key
	if !CheckAuthKey(w, r) {
		return
	}

	// new customer
	apikey := utils.NewAPIKey(12)
	customer := Customer{Apikey: apikey}
	err := DB.InsertCustomer(customer)
	if err != nil {
		ServerLogger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := HttpMessageResponse{Status: "success", Apikey: apikey, Data: "customer created"}
	js, err := MarshalJsonFromStruct(w, r, data)
	if err != nil {
		return
	}

	SendJsonResponse(w, r, js)
}

// Pull all customer datasource pairs
// Distributed System
func AllCustomerDatasources(w http.ResponseWriter, r *http.Request) {
	results := []Customer{}

	if !CheckAuthKey(w, r) {
		return
	}

	customers, err := DB.SelectAll("apikeys")
	if err != nil {
		ServerLogger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range customers {
		val, err := DB.Select("apikeys", v)
		if err != nil {
			ServerLogger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customer := Customer{}
		err = json.Unmarshal(val, &customer)
		if err != nil {
			panic(err)
		}
		results = append(results, customer)
	}

	js, err := MarshalJsonFromStruct(w, r, results)
	if err != nil {
		return
	}

	SendJsonResponse(w, r, js)
}
