package app

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Layer struct {
	Datasource string  `json:"datasource"`
	Geojson    Geojson `json:"geojson"`
}

func (lyr *Layer) Save() error {
	DB.insertLayer(lyr.Datasource, lyr.Geojson)
	return nil
}

func NewLayerHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("apikey") == "" {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	customer, err := DB.getCustomer(r.FormValue("apikey"))
	if err != nil {
		Warning.Println(r.RemoteAddr, "POST /api/v1/layer/ [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Create datasource
	geojs := NewGeojson()
	ds, _ := NewUUID()
	lyr := Layer{Datasource: ds, Geojson: geojs}
	lyr.Save()
	// add datasource to customer
	customer.Datasources = append(customer.Datasources, ds)
	DB.insertCustomer(customer)
	// return results
	data := `{"status":"ok","datasource":"` + ds + `"}`
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "POST /api/v1/layer [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ViewLayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	// Check apikey permissions
	if r.FormValue("apikey") == "" {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	customer, err := DB.getCustomer(r.FormValue("apikey"))
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/ [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Get layer from database
	lyr, err := DB.getLayer(ds)
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Marshal datasource layer to json
	js, err := json.Marshal(lyr)
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Finish
	Info.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+" [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DeleteLayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	// Check apikey permissions
	if r.FormValue("apikey") == "" {
		Error.Println(r.RemoteAddr, "DELETE /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	customer, err := DB.getCustomer(r.FormValue("apikey"))
	if err != nil {
		Warning.Println(r.RemoteAddr, "DELETE /api/v1/layer/ [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "DELETE /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Delete layer
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
	// Return results
	data := `{"status":"ok","datasource":"` + ds + `", "result":"datasource deleted"}`
	js, err := json.Marshal(data)
	if err != nil {
		Info.Println(r.RemoteAddr, "DELETE /api/v1/layer/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "DELETE /api/v1/layer/"+ds+" [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
