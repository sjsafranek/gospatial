package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"net/http"
	//"time"
)

// NewFeatureHandler creates a new feature and adds it to a layer.
// Layer is then saved to database. All active clients viewing layer
// are notified of update via websocket hub.
// @param apikey customer id
// @oaram ds datasource uuid
// @return json
func NewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Get request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf(" %v %v [500]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	/*=======================================*/
	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		return
	}

	customer, err := GetCustomerFromDatabase(w, r, apikey)
	if err != nil {
		return
	}

	if !CheckCustomerForDatasource(w, r, customer, ds) {
		return
	}

	// Unmarshal feature
	feat, err := geojson.UnmarshalFeature(body)
	if err != nil {
		message := fmt.Sprintf(" %v %v [400]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save feature to database
	err = DB.InsertFeature(ds, feat)
	if err != nil {
		message := fmt.Sprintf(" %v %v [500]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate message
	data := HttpMessageResponse{Status: "success", Datasource: ds, Data: "feature added"}
	js, err := MarshalJsonFromStruct(w, r, data)
	if err != nil {
		return
	}

	// Update websockets
	conn := connection{ds: ds, ip: r.RemoteAddr}
	Hub.broadcast(true, &conn)

	// Return results
	SendJsonResponse(w, r, js)
}

// ViewFeatureHandler finds feature in layer via array index. Returns feature geojson.
// @param apikey customer id
// @oaram ds datasource uuid
// @return feature geojson
func ViewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	/*=======================================*/
	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		return
	}

	customer, err := GetCustomerFromDatabase(w, r, apikey)
	if err != nil {
		return
	}

	if !CheckCustomerForDatasource(w, r, customer, ds) {
		return
	}
	/*=======================================*/

	// Get layer from database
	data, err := DB.GetLayer(ds)
	if err != nil {
		message := fmt.Sprintf(" %v %v [404]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check for feature
	var js []byte
	for _, v := range data.Features {
		geo_id := fmt.Sprintf("%v", v.Properties["geo_id"])
		if geo_id == vars["k"] {
			js, err = v.MarshalJSON()
			if err != nil {
				message := fmt.Sprintf(" %v %v [500]", r.Method, r.URL.Path)
				NetworkLogger.Critical(r.RemoteAddr, message)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Return results
			SendJsonResponse(w, r, js)
			return
		}
	}

	// Feature not found
	message := fmt.Sprintf(" %v %v [404]", r.Method, r.URL.Path)
	NetworkLogger.Critical(r.RemoteAddr, message)
	err = fmt.Errorf("Not found")
	http.Error(w, err.Error(), http.StatusNotFound)
}

// EditFeatureHandler finds feature in layer via array index. Edits feature.
// @param apikey customer id
// @oaram ds datasource uuid
func EditFeatureHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Get request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf(" %v %v [500]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	//TESTING
	fmt.Printf(string(body))

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]
	geo_id := vars["k"]

	/*=======================================*/
	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		return
	}

	customer, err := GetCustomerFromDatabase(w, r, apikey)
	if err != nil {
		return
	}

	if !CheckCustomerForDatasource(w, r, customer, ds) {
		return
	}
	/*=======================================*/

	// Unmarshal feature
	feat, err := geojson.UnmarshalFeature(body)
	if err != nil {
		message := fmt.Sprintf(" %v %v [400]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = DB.EditFeature(ds, geo_id, feat)
	if err != nil {
		// Feature not found
		message := fmt.Sprintf(" %v %v [404]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		err = fmt.Errorf("Not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	// Generate message
	data := HttpMessageResponse{Status: "success", Datasource: ds, Data: "feature edited"}
	js, err := MarshalJsonFromStruct(w, r, data)
	if err != nil {
		return
	}

	// Update websockets
	conn := connection{ds: ds, ip: r.RemoteAddr}
	Hub.broadcast(true, &conn)

	// Feature not found
	SendJsonResponse(w, r, js)
}
