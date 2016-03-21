package app

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/paulmach/go.geojson"
	"net/http"
)

/*=======================================*/
// Method: ViewLayersHandler
// Description:
//		Returns customer layers
// @param apikey customer id
// @return json
/*=======================================*/
func ViewLayersHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.getCustomer(apikey)
	if err != nil {
		Warning.Println(r.RemoteAddr, "POST /api/v1/layer/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	/*=======================================*/

	// return results
	js, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)

}

/*=======================================*/
// Method: NewLayerHandler
// Description:
//		Creates a new layer
// 		Layer is saved to database
//		Layer uuid is added to customer list
// @param apikey
// @return json
/*=======================================*/
func NewLayerHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	// Check for apikey in request
	if apikey == "" {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.getCustomer(apikey)
	if err != nil {
		Warning.Println(r.RemoteAddr, "POST /api/v1/layer/ [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create datasource
	// geojs := NewGeojson()
	ds, _ := NewUUID()
	// lyr := Layer{Datasource: ds, Geojson: geojs}
	// lyr.Save()
	featCollection := geojson.NewFeatureCollection()
	DB.insertLayer(ds, featCollection)

	// Add datasource uuid to customer
	customer.Datasources = append(customer.Datasources, ds)
	DB.insertCustomer(customer)

	// Generate message
	data := `{"status":"ok","datasource":"` + ds + `"}`
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return results
	Info.Println(r.RemoteAddr, "POST /api/v1/layer [200]")
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)

}

/*=======================================*/
// Method: ViewLayerHandler
// Description:
//		Gets requested layer from database
//		Checks apikey/customer permissions
//		Returns layer geojson
// @param ds
// @param apikey
// @return geojson
/*=======================================*/
func ViewLayerHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.getCustomer(apikey)
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/ [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Get layer from database
	lyr, err := DB.getLayer(ds)
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal datasource layer to json
	// js, err := json.Marshal(lyr)
	rawJSON, err := lyr.MarshalJSON()
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return layer json
	Info.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+" [200]")
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(rawJSON)

}

/*=======================================*/
// Method: DeleteLayerHandler
// Description:
//		Checks apikey/customer permissions
//		Deletes layer from database
//		Deletes layer from customer list
// @param ds
// @param apikey
// @return json
/*=======================================*/
func DeleteLayerHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		Error.Println(r.RemoteAddr, "DELETE /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.getCustomer(apikey)
	if err != nil {
		Warning.Println(r.RemoteAddr, "DELETE /api/v1/layer/ [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "DELETE /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Delete layer from database
	err = DB.deleteLayer(ds)
	if err != nil {
		Info.Println(r.RemoteAddr, "DELETE /api/v1/layer/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete layer from customer
	i := sliceIndex(ds, customer.Datasources)
	customer.Datasources = append(customer.Datasources[:i], customer.Datasources[i+1:]...)
	DB.insertCustomer(customer)

	// Generate message
	data := `{"status":"ok","datasource":"` + ds + `", "result":"datasource deleted"}`
	js, err := json.Marshal(data)
	if err != nil {
		Info.Println(r.RemoteAddr, "DELETE /api/v1/layer/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Returns results
	Info.Println(r.RemoteAddr, "DELETE /api/v1/layer/"+ds+" [200]")
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)

}
