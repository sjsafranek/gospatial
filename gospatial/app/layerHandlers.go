package app

import (
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	"net/http"
	"gospatial/utils"
)

// ViewLayersHandler returns json containing customer layers
// @param apikey customer id
// @return json
func ViewLayersHandler(w http.ResponseWriter, r *http.Request) {
	networkLoggerInfoIn.Printf("%v\n", r)
	// Get params
	apikey := r.FormValue("apikey")

	// Check for apikey in request
	if apikey == "" {
		networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/layers [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		networkLoggerWarning.Println(r.RemoteAddr, "POST /api/v1/layers [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// return results
	js, err := json.Marshal(customer)
	if err != nil {
		networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/layers [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	networkLoggerInfo.Println(r.RemoteAddr, "POST /api/v1/layers [200]")
	networkLoggerInfoOut.Println(string(js))
	w.Write(js)

}

// NewLayerHandler creates a new geojson layer. Saves layer to database and adds layer to customer
// @param apikey
// @return json
func NewLayerHandler(w http.ResponseWriter, r *http.Request) {
	networkLoggerInfoIn.Printf("%v\n", r)

	// Get params
	apikey := r.FormValue("apikey")

	// Check for apikey in request
	if apikey == "" {
		networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		networkLoggerWarning.Println(r.RemoteAddr, "POST /api/v1/layer/ [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create datasource
	ds, err := DB.NewLayer()
	if err != nil {
		networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/layer [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add datasource uuid to customer
	customer.Datasources = append(customer.Datasources, ds)
	DB.InsertCustomer(customer)

	// Generate message
	data := `{"status":"ok","datasource":"` + ds + `"}`
	js, err := json.Marshal(data)
	if err != nil {
		networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/layer [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	networkLoggerInfo.Println(r.RemoteAddr, "POST /api/v1/layer [200]")
	networkLoggerInfoOut.Println(string(js))
	w.Write(js)

}

// ViewLayerHandler returns geojson of requested layer. Apikey/customer is checked for permissions to requested layer.
// @param ds
// @param apikey
// @return geojson
func ViewLayerHandler(w http.ResponseWriter, r *http.Request) {
	networkLoggerInfoIn.Printf("%v\n", r)

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	// Check for apikey in request
	if apikey == "" {
		networkLoggerError.Println(r.RemoteAddr, "GET /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		networkLoggerWarning.Println(r.RemoteAddr, "GET /api/v1/layer/ [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !utils.StringInSlice(ds, customer.Datasources) {
		networkLoggerError.Println(r.RemoteAddr, "GET /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get layer from database
	lyr, err := DB.GetLayer(ds)
	if err != nil {
		networkLoggerWarning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal datasource layer to json
	rawJSON, err := lyr.MarshalJSON()
	if err != nil {
		networkLoggerError.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return layer json
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	networkLoggerInfo.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+" [200]")
	networkLoggerInfoOut.Println(string(rawJSON))
	w.Write(rawJSON)

}

// DeleteLayerHandler deletes layer from database and removes it from customer list.
// @param ds
// @param apikey
// @return json
func DeleteLayerHandler(w http.ResponseWriter, r *http.Request) {
	networkLoggerInfoIn.Printf("%v\n", r)

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	// Check for apikey in request
	if apikey == "" {
		networkLoggerError.Println(r.RemoteAddr, "DELETE /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		networkLoggerWarning.Println(r.RemoteAddr, "DELETE /api/v1/layer/ [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !utils.StringInSlice(ds, customer.Datasources) {
		networkLoggerError.Println(r.RemoteAddr, "DELETE /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// KEEP LAYER IN CASE OF RECOVERY
	// Delete layer from database
	// err = DB.DeleteLayer(ds)
	// if err != nil {
	// 	networkLoggerInfo.Println(r.RemoteAddr, "DELETE /api/v1/layer/"+ds+" [500]")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Delete layer from customer
	i := utils.SliceIndex(ds, customer.Datasources)
	customer.Datasources = append(customer.Datasources[:i], customer.Datasources[i+1:]...)
	DB.InsertCustomer(customer)

	// Generate message
	data := `{"status":"ok","datasource":"` + ds + `", "result":"datasource deleted"}`
	js, err := json.Marshal(data)
	if err != nil {
		networkLoggerInfo.Println(r.RemoteAddr, "DELETE /api/v1/layer/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Returns results
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	networkLoggerInfo.Println(r.RemoteAddr, "DELETE /api/v1/layer/"+ds+" [200]")
	networkLoggerInfoOut.Println(string(js))
	w.Write(js)

}
