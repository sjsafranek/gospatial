package gospatial

import (
	"encoding/json"
	"fmt"
	"net/http"
)

import (
	"./utils"
)

func MarshalJsonFromString(w http.ResponseWriter, r *http.Request, data string) ([]byte, error) {
	js, err := json.Marshal(data)
	if err != nil {
		message := fmt.Sprintf(" %v %v [500]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return js, err
	}
	return js, nil
}

func MarshalJsonFromStruct(w http.ResponseWriter, r *http.Request, data interface{}) ([]byte, error) {
	js, err := json.Marshal(data)
	if err != nil {
		message := fmt.Sprintf(" %v %v [500]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return js, err
	}
	return js, nil
}

// Sends http response
func SendJsonResponse(w http.ResponseWriter, r *http.Request, js []byte) {
	// Log result
	message := fmt.Sprintf(" %v %v [200]", r.Method, r.URL.Path)
	NetworkLogger.Info(r.RemoteAddr, message)
	NetworkLogger.Debug("[Out] ", string(js))
	// set response headers
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// write response content
	w.Write(js)
}

// check request for valid authkey
func CheckAuthKey(w http.ResponseWriter, r *http.Request) bool {
	if SuperuserKey != r.FormValue("authkey") {
		message := fmt.Sprintf(" %v %v [401]", r.Method, r.URL.Path)
		NetworkLogger.Error(r.RemoteAddr, message)
		http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
		return false
	}
	return true
}

// Check for apikey in request
func GetApikeyFromRequest(w http.ResponseWriter, r *http.Request) string {
	// Get params
	apikey := r.FormValue("apikey")
	// Check for apikey in request
	if apikey == "" {
		message := fmt.Sprintf(" %v %v [401]", r.Method, r.URL.Path)
		NetworkLogger.Error(r.RemoteAddr, message)
		http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
	}
	// return apikey
	return apikey
}

// Get customer from database
func GetCustomerFromDatabase(w http.ResponseWriter, r *http.Request, apikey string) (Customer, error) {
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		message := fmt.Sprintf(" %v %v [404]", r.Method, r.URL.Path)
		NetworkLogger.Error(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusNotFound)
		return customer, err
	}
	return customer, err
}

// Check customer datasource list
func CheckCustomerForDatasource(w http.ResponseWriter, r *http.Request, customer Customer, ds string) bool {
	if !utils.StringInSlice(ds, customer.Datasources) {
		message := fmt.Sprintf(" %v %v [401]", r.Method, r.URL.Path)
		NetworkLogger.Error(r.RemoteAddr, message)
		http.Error(w, `{"status": "error", "result": "unauthorized"}`, http.StatusUnauthorized)
		return false
	}
	return true
}
