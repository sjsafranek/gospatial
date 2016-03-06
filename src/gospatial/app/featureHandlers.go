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
		// Info.Println("NewFeatureHandler |", r.RemoteAddr, "| DELETE | 500 |", "/api/v1/layer/"+ds+"/feature")
		Error.Println(r.RemoteAddr, "| POST | 500 |", "/api/v1/layer/"+ds+"/feature")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "| POST | 500 |", "/api/v1/layer/"+ds+"/feature")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Info.Println(r.RemoteAddr, ds, "new feature")
	// get geojson
	var geojs Geojson
	err = json.Unmarshal(js, &geojs)
	if err != nil {
		Error.Println(r.RemoteAddr, "| POST | 500 |", "/api/v1/layer/"+ds+"/feature")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// read request body
	feat := NewFeature()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println(r.RemoteAddr, "| POST | 500 |", "/api/v1/layer/"+ds+"/feature")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &feat)
	if err != nil {
		Error.Println(r.RemoteAddr, "| POST | 500 |", "/api/v1/layer/"+ds+"/feature")
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
		Error.Println(r.RemoteAddr, "| POST | 500 |", "/api/v1/layer/"+ds+"/feature")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Update websockets
	conn := connection{ds: ds, ip: r.RemoteAddr}
	Hub.broadcast(true, &conn)
	// Report results
	w.Header().Set("Content-Type", "application/json")
	Info.Println(r.RemoteAddr, "| POST | 200 |", "/api/v1/layer/"+ds+"/feature")
	w.Write(js)
}

func ViewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	k, err := strconv.Atoi(vars["k"])
	if err != nil {
		Error.Println(r.RemoteAddr, "| GET | 500 |", "/api/v1/layer/"+ds+"/feature/"+vars["k"])
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := DB.getLayer(ds)
	if err != nil {
		Error.Println(r.RemoteAddr, "| GET | 500 |", "/api/v1/layer/"+ds+"/feature/"+vars["k"])
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(data.Features[k])
	if err != nil {
		Error.Println(r.RemoteAddr, "| GET | 500 |", "/api/v1/layer/"+ds+"/feature/"+vars["k"])
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	Info.Println(r.RemoteAddr, "| GET | 200 |", "/api/v1/layer/"+ds+"/feature/"+vars["k"])
	w.Write(js)
}
