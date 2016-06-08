package app

import (
	"encoding/json"
	"net/http"
)

// NewLayerHandler creates a new geojson layer. Saves layer to database and adds layer to customer
// @param apikey
// @return json
func NewTileLayerHandler(w http.ResponseWriter, r *http.Request) {
	networkLoggerInfoIn.Printf("%v\n", r)

	// Get params
	apikey := r.FormValue("apikey")

	tilelayer_url := r.FormValue("tilelayer_url")
	tilelayer_name := r.FormValue("tilelayer_name")
	tilelayer := TileLayer{Url: tilelayer_url, Name: tilelayer_name}

	// if isUrl(tilelayer.Url) != true {
	// 	networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/tilelayer [400]")
	// 	http.Error(w, "not a valid url", http.StatusBadRequest)
	// 	return
	// }

	// Check for apikey in request
	if apikey == "" {
		networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/tilelayer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		networkLoggerWarning.Println(r.RemoteAddr, "POST /api/v1/tilelayer [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Add datasource uuid to customer
	customer.TileLayers = append(customer.TileLayers, tilelayer)
	// customer.TileLayers[tilelayer_name] = tilelayer_url
	DB.InsertCustomer(customer)

	// Generate message
	data := `{"status": "success", "data": {"tilelayer": {"url": "` + tilelayer_url + `", "name": "` + tilelayer_name + `"}}}`
	js, err := json.Marshal(data)
	if err != nil {
		networkLoggerError.Println(r.RemoteAddr, "POST /api/v1/tilelayer [500]")
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

// // DeleteLayerHandler deletes layer from database and removes it from customer list.
// // @param ds
// // @param apikey
// // @return json
// func DeleteBasemapHandler(w http.ResponseWriter, r *http.Request) {
// 	networkLoggerInfoIn.Printf("%v\n", r)

// 	// Get params
// 	apikey := r.FormValue("apikey")

// 	// Get ds from url path
// 	vars := mux.Vars(r)
// 	ds := vars["ds"]

// 	// Check for apikey in request
// 	if apikey == "" {
// 		networkLoggerError.Println(r.RemoteAddr, "DELETE /api/v1/tilelayer [401]")
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Get customer from database
// 	customer, err := DB.GetCustomer(apikey)
// 	if err != nil {
// 		networkLoggerWarning.Println(r.RemoteAddr, "DELETE /api/v1/tilelayer/ [404]")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	// Check customer basemap list
// 	if !stringInSlice(ds, customer.TileLayers) {
// 		networkLoggerError.Println(r.RemoteAddr, "DELETE /api/v1/tilelayer [401]")
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Delete layer from customer
// 	i := sliceIndex(ds, customer.TileLayers)
// 	customer.TileLayers = append(customer.TileLayers[:i], customer.TileLayers[i+1:]...)
// 	DB.InsertCustomer(customer)

// 	// Generate message
// 	data := `{"status":"ok","tilelayer":"` + ds + `", "result":"tilelayer deleted"}`
// 	js, err := json.Marshal(data)
// 	if err != nil {
// 		networkLoggerInfo.Println(r.RemoteAddr, "DELETE /api/v1/tilelayer/"+ds+" [500]")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Returns results
// 	w.Header().Set("Content-Type", "application/json")
// 	// allow cross domain AJAX requests
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	networkLoggerInfo.Println(r.RemoteAddr, "DELETE /api/v1/tilelayer/"+ds+" [200]")
// 	networkLoggerInfoOut.Println(string(js))
// 	w.Write(js)

// }
