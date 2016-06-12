package app

import (
	"encoding/json"
	// "github.com/gorilla/mux"
	"gospatial/utils"
	"net/http"
	"runtime"
	"time"
)

var startTime = time.Now()

// SuperuserKey api servers superuser key
var SuperuserKey string = "su"

// PingHandler provides an api route for server health check
func PingHandler(w http.ResponseWriter, r *http.Request) {
	data := `{"status": "ok", "message": "pong"}`
	js, err := json.Marshal(data)
	if err != nil {
		networkLoggerError.Println(r.RemoteAddr, "GET /ping [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	networkLoggerInfo.Println(r.RemoteAddr, "GET /ping [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// // UnloadLayer unloads layer from memory cache
// func UnloadLayer(w http.ResponseWriter, r *http.Request) {
// 	// Check auth key
// 	if SuperuserKey != r.FormValue("authkey") {
// 		networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/layer [401]")
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	// Parse url params
// 	vars := mux.Vars(r)
// 	ds := vars["ds"]
// 	Debug.Printf("Unloading [%s]", ds)
// 	// unload
// 	delete(DB.Cache, ds)
// 	// Response
// 	data := `{"status":"ok","datasource":"` + ds + `", "result":"datasource unloaded"}`
// 	js, err := json.Marshal(data)
// 	if err != nil {
// 		networkLoggerError.Println(r.RemoteAddr, "GET /management/unload/"+ds+" [500]")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	networkLoggerInfo.Println(r.RemoteAddr, "GET /management/unload/"+ds+" [200]")
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }

// // LoadedLayers returns list of layers loaded in memory
// func LoadedLayers(w http.ResponseWriter, r *http.Request) {
// 	// Check auth key
// 	if SuperuserKey != r.FormValue("authkey") {
// 		networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/layer [401]")
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	Debug.Println("Checking loaded datasources...")
// 	// collect datasource ids
// 	i := 0
// 	keys := make([]string, len(DB.Cache))
// 	for k := range DB.Cache {
// 		keys[i] = k
// 		i++
// 	}
// 	data := make(map[string]interface{})
// 	data["datasources"] = keys
// 	// marshal and send response
// 	js, err := json.Marshal(data)
// 	if err != nil {
// 		networkLoggerError.Println(r.RemoteAddr, "GET /management/loaded [500]")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	networkLoggerInfo.Println(r.RemoteAddr, "GET /management/loaded [200]")
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }

// ServerProfile returns basic server stats
func ServerProfile(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data = make(map[string]interface{})
	data["registered"] = startTime.UTC()
	data["uptime"] = time.Since(startTime).Seconds()
	// data["status"] = AppMode // debug, static, standard
	data["num_cores"] = runtime.NumCPU()
	// data["free_mem"] = runtime.MemStats()
	js, err := json.Marshal(data)
	if err != nil {
		networkLoggerError.Println(r.RemoteAddr, "GET /management/profile [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	networkLoggerInfo.Println(r.RemoteAddr, "GET /management/profile/ [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// NewCustomerHandler superuser route to create new api customers/apikeys
func NewCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Check auth key
	if SuperuserKey != r.FormValue("authkey") {
		networkLoggerError.Println(r.RemoteAddr, "POST /management/customer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// new customer
	apikey := utils.NewAPIKey(12)
	customer := Customer{Apikey: apikey}
	err := DB.InsertCustomer(customer)
	if err != nil {
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return results
	data := `{"status":"ok","apikey":"` + apikey + `", "result":"customer created"}`
	js, err := json.Marshal(data)
	if err != nil {
		networkLoggerError.Println(r.RemoteAddr, "POST /management/customer [500]")
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	networkLoggerInfo.Println(r.RemoteAddr, "POST /management/customer [200]")
	w.Write(js)
}

// ShareLayerHandler gives customer access to an existing datasource.
// @param apikey - customer to give access
// @param authkey
// @return json
// func ShareLayerHandler(w http.ResponseWriter, r *http.Request) {

// 	// Get url params
// 	apikey := r.FormValue("apikey")
// 	authkey := r.FormValue("authkey")

// 	// Get ds from url path
// 	vars := mux.Vars(r)
// 	ds := vars["ds"]

// 	// superuser access
// 	if SuperuserKey != authkey {
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	if apikey == "" {
// 		networkLoggerError.Println(r.RemoteAddr, "PUT /api/v1/layer/{ds} [401]")
// 		http.Error(w, "bad request", http.StatusBadRequest)
// 		return
// 	}

// 	// Get customer from database
// 	customer, err := DB.GetCustomer(apikey)
// 	if err != nil {
// 		networkLoggerWarning.Println(r.RemoteAddr, "PUT /api/v1/layer/{ds} [404]")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	// Add datasource uuid to customer
// 	customer.Datasources = append(customer.Datasources, ds)
// 	DB.InsertCustomer(customer)

// 	// Generate message
// 	data := `{"status":"ok","datasource":"` + ds + `"}`
// 	js, err := json.Marshal(data)
// 	if err != nil {
// 		networkLoggerError.Println(r.RemoteAddr, "PUT /api/v1/layer [500]")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Return results
// 	networkLoggerInfo.Println(r.RemoteAddr, "PUT /api/v1/layer [200]")
// 	w.Header().Set("Content-Type", "application/json")
// 	// allow cross domain AJAX requests
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Write(js)

// }
