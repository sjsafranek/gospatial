package app

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func NewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	data, err := DB.getLayer(ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, ds, "new feature")
	// get geojson
	var geojs Geojson
	err = json.Unmarshal(js, &geojs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// read request body
	feat := NewFeature()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println("reading body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &feat)
	if err != nil {
		Error.Println("unmarshal json")
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ViewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	k, err := strconv.Atoi(vars["k"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Info.Println(r.RemoteAddr, ds, "read feature", k)
	data, err := DB.getLayer(ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(data.Features[k])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}