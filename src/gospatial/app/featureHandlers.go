package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*=======================================*/
// Method: NewFeatureHandler
// Description:
//		Adds a new feature to a layer
//		Saves layer to database
// @param apikey customer id
// @oaram ds datasource uuid
// @return json
/*=======================================*/
func NewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	// Get request body
	// If this id done later in this function an EOF error occurs
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	/*
		t := NewFeature()
		decoder := json.NewDecoder(bd)
		err := decoder.Decode(&t)
		if err != nil {
			Error.Println(err)
		}
		Info.Println(t)
	*/

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.getCustomer(apikey)
	if err != nil {
		Warning.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Get layer from database
	geojs, err := DB.getLayer(ds)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read request body and marshal to geojson feature
	feat := NewFeature()
	err = json.Unmarshal(body, &feat)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add new feature to layer
	geojs.Features = append(geojs.Features, feat)
	lyr := Layer{Datasource: ds, Geojson: geojs}
	lyr.Save()

	// Generate message
	data := `{"status":"ok","datasource":"` + ds + `", "message":"feature added"}`
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update websockets
	conn := connection{ds: ds, ip: r.RemoteAddr}
	Hub.broadcast(true, &conn)

	// Return results
	w.Header().Set("Content-Type", "application/json")
	Info.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [200]")
	w.Write(js)

}

/*=======================================*/
// Method: ViewFeatureHandler
// Description:
//		Finds feature from layer
// @param apikey customer id
// @oaram ds datasource uuid
// @return feature geojson
/*=======================================*/
func ViewFeatureHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	k, err := strconv.Atoi(vars["k"])
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [400]")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.getCustomer(apikey)
	if err != nil {
		Warning.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Get layer from database
	data, err := DB.getLayer(ds)
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check for feature
	if k > len(data.Features) {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
		err := fmt.Errorf("Not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal feature to json
	js, err := json.Marshal(data.Features[k])
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	Info.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [200]")
	w.Write(js)

}
