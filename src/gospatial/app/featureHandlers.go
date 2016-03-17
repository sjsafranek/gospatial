package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	// "strings"
)

func NewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check apikey permissions
	if r.FormValue("apikey") == "" {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	customer, err := DB.getCustomer(r.FormValue("apikey"))
	if err != nil {
		Warning.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// get layer from database
	data, err := DB.getLayer(ds)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get geojson
	var geojs Geojson
	err = json.Unmarshal(js, &geojs)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// read request body
	feat := NewFeature()
	err = json.Unmarshal(body, &feat)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// add feature to layer
	geojs.Features = append(geojs.Features, feat)
	lyr := Layer{Datasource: ds, Geojson: geojs}
	lyr.Save()
	// Return results
	js, err = json.Marshal(geojs)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Update websockets
	conn := connection{ds: ds, ip: r.RemoteAddr}
	Hub.broadcast(true, &conn)
	// Report results
	w.Header().Set("Content-Type", "application/json")
	Info.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [200]")
	w.Write(js)
}

func ViewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	// Check apikey permissions
	if r.FormValue("apikey") == "" {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	customer, err := DB.getCustomer(r.FormValue("apikey"))
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	k, err := strconv.Atoi(vars["k"])
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [400]")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Get layer
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
	// Finish
	w.Header().Set("Content-Type", "application/json")
	Info.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [200]")
	w.Write(js)
}
