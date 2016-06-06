package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"net/http"
	"strconv"
)

// NewFeatureHandler creates a new feature and adds it to a layer.
// Layer is then saved to database. All active clients viewing layer
// are notified of update via websocket hub.
// @param apikey customer id
// @oaram ds datasource uuid
// @return json
func NewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	network_logger_Info_In.Printf("%v\n", r)

	// Get request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		vars := mux.Vars(r)
		ds := vars["ds"]
		network_logger_Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		network_logger_Warning.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		network_logger_Warning.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		network_logger_Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Unmarshal feature
	feat, err := geojson.UnmarshalFeature(body)
	if err != nil {
		network_logger_Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [400]")
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save feature to database
	err = DB.InsertFeature(ds, feat)
	if err != nil {
		network_logger_Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate message
	data := `{"status":"ok","datasource":"` + ds + `", "message":"feature added"}`
	js, err := json.Marshal(data)
	if err != nil {
		network_logger_Error.Println(r.RemoteAddr, "POST /api/v1/layer [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update websockets
	conn := connection{ds: ds, ip: r.RemoteAddr}
	Hub.broadcast(true, &conn)

	// Return results
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	network_logger_Info.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [200]")
	network_logger_Info_Out.Println(string(js))
	w.Write(js)

}

// ViewFeatureHandler finds feature in layer via array index. Returns feature geojson.
// @param apikey customer id
// @oaram ds datasource uuid
// @return feature geojson
func ViewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	network_logger_Info_In.Printf("%v\n", r)

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	k, err := strconv.Atoi(vars["k"])
	if err != nil {
		network_logger_Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [400]")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		network_logger_Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		network_logger_Warning.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		network_logger_Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Get layer from database
	data, err := DB.GetLayer(ds)
	if err != nil {
		network_logger_Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check for feature
	if k > len(data.Features) {
		network_logger_Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
		err := fmt.Errorf("Not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal feature to json
	js, err := data.Features[k].MarshalJSON()
	if err != nil {
		network_logger_Error.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//
	network_logger_Info.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [200]")
	network_logger_Info_Out.Println(string(js))
	w.Write(js)

}
