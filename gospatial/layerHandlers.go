package gospatial

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

import (
	"./utils"
)

// ViewLayersHandler returns json containing customer layers
// @param apikey customer id
// @return json
func ViewLayersHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		return
	}

	customer, err := GetCustomerFromDatabase(w, r, apikey)
	if err != nil {
		return
	}

	js, err := MarshalJsonFromStruct(w, r, customer)
	if err != nil {
		return
	}

	SendJsonResponse(w, r, js)
}

// NewLayerHandler creates a new geojson layer. Saves layer to database and adds layer to customer
// @param apikey
// @return json
func NewLayerHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		return
	}

	customer, err := GetCustomerFromDatabase(w, r, apikey)
	if err != nil {
		return
	}

	// Create datasource
	ds, err := DB.NewLayer()
	if err != nil {
		message := fmt.Sprintf(" %v %v [500]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add datasource uuid to customer
	customer.Datasources = append(customer.Datasources, ds)
	DB.InsertCustomer(customer)

	// Generate message
	data := HttpMessageResponse{Status: "success", Datasource: ds}
	js, err := MarshalJsonFromStruct(w, r, data)
	if err != nil {
		return
	}

	// Return results
	SendJsonResponse(w, r, js)
}

// ViewLayerHandler returns geojson of requested layer. Apikey/customer is checked for permissions to requested layer.
// @param ds
// @param apikey
// @return geojson
func ViewLayerHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

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

	// Get layer from database
	lyr, err := DB.GetLayer(ds)
	if err != nil {
		message := fmt.Sprintf(" %v %v [404]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal datasource layer to json
	js, err := lyr.MarshalJSON()
	if err != nil {
		message := fmt.Sprintf(" %v %v [500]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return layer json
	SendJsonResponse(w, r, js)
}

// DeleteLayerHandler deletes layer from database and removes it from customer list.
// @param ds
// @param apikey
// @return json
func DeleteLayerHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

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

	// Delete layer from customer
	i := utils.SliceIndex(ds, customer.Datasources)
	customer.Datasources = append(customer.Datasources[:i], customer.Datasources[i+1:]...)
	DB.InsertCustomer(customer)

	// Delete layer from database
	err = DB.DeleteLayer(ds)
	if err != nil {
		message := fmt.Sprintf(" %v %v [500]", r.Method, r.URL.Path)
		NetworkLogger.Critical(r.RemoteAddr, message)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate message
	data := HttpMessageResponse{Status: "success", Datasource: ds, Data: "datasource deleted"}
	js, err := MarshalJsonFromStruct(w, r, data)
	if err != nil {
		return
	}

	// Returns results
	SendJsonResponse(w, r, js)
}
