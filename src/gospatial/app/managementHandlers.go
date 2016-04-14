package app

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"runtime"
	"time"
)

var START_TIME = time.Now()

var SuperuserKey string = "su"
var AppMode string = "standard"

func PingHandler(w http.ResponseWriter, r *http.Request) {
	data := `{"status": "ok", "message": "pong"}`
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /ping [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "GET /ping [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func UnloadLayer(w http.ResponseWriter, r *http.Request) {
	// Check auth key
	if SuperuserKey != r.FormValue("authkey") {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Parse url params
	vars := mux.Vars(r)
	ds := vars["ds"]
	Debug.Printf("Unloading [%s]", ds)
	// unload
	delete(DB.Cache, ds)
	// Response
	data := `{"status":"ok","datasource":"` + ds + `", "result":"datasource unloaded"}`
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /management/unload/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "GET /management/unload/"+ds+" [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func LoadedLayers(w http.ResponseWriter, r *http.Request) {
	// Check auth key
	if SuperuserKey != r.FormValue("authkey") {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	Debug.Println("Checking loaded datasources...")
	// collect datasource ids
	i := 0
	keys := make([]string, len(DB.Cache))
	for k := range DB.Cache {
		keys[i] = k
		i++
	}
	data := make(map[string]interface{})
	data["datasources"] = keys
	// marshal and send response
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /management/loaded [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "GET /management/loaded [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func server_profile(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data = make(map[string]interface{})
	data["registered"] = START_TIME.UTC()
	data["uptime"] = time.Since(START_TIME).Seconds()
	data["status"] = AppMode // debug, static, standard
	data["num_cores"] = runtime.NumCPU()
	// data["free_mem"] = runtime.MemStats()
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /management/profile [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "GET /management/profile/ [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DebugModeHandler(w http.ResponseWriter, r *http.Request) {
	// Check authkey
	if SuperuserKey != r.FormValue("authkey") {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Get url params
	vars := mux.Vars(r)
	md := vars["md"]
	if md == "debug" {
		js, err := json.Marshal(`{"status": "ok", "message": "debug mode on"}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		DebugMode()
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else if md == "standard" {
		js, err := json.Marshal(`{"status": "ok", "message": "debug mode off"}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		StandardMode()
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// SUPERUSER
func NewCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Check auth key
	if SuperuserKey != r.FormValue("authkey") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	} else {
		// new customer
		apikey := NewAPIKey(12)
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
			Error.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		// allow cross domain AJAX requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	}
}

/*=======================================*/
// Method: ShareLayerHandler
// Description:
//		Gives customer access to datasource
// @param apikey - customer to give access
// @param authkey
// @return json
/*=======================================*/
func ShareLayerHandler(w http.ResponseWriter, r *http.Request) {

	// Get url params
	apikey := r.FormValue("apikey")
	authkey := r.FormValue("authkey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	// superuser access
	if SuperuserKey != authkey {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	} else {

		if apikey == "" {
			Error.Println(r.RemoteAddr, "PUT /api/v1/layer/{ds} [401]")
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Get customer from database
		customer, err := DB.GetCustomer(apikey)
		if err != nil {
			Warning.Println(r.RemoteAddr, "PUT /api/v1/layer/{ds} [404]")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Add datasource uuid to customer
		customer.Datasources = append(customer.Datasources, ds)
		DB.InsertCustomer(customer)

		// Generate message
		data := `{"status":"ok","datasource":"` + ds + `"}`
		js, err := json.Marshal(data)
		if err != nil {
			Error.Println(r.RemoteAddr, "PUT /api/v1/layer [500]")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return results
		Info.Println(r.RemoteAddr, "PUT /api/v1/layer [200]")
		w.Header().Set("Content-Type", "application/json")
		// allow cross domain AJAX requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)

	}

}
