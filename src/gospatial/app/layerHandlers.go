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
	geojs := NewGeojson()
	ds, _ := NewUUID()
	lyr := Layer{Datasource: ds, Geojson: geojs}
	lyr.Save()
	// Info.Println(r.RemoteAddr, ds, "new layer")
	data := `{"status":"ok","datasource":"` + ds + `"}`
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "| POST | 500 |", "/api/v1/layer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "| POST | 200 |", "/api/v1/layer")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ViewLayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	lyr, err := DB.getLayer(ds)
	if err != nil {
		Info.Println(r.RemoteAddr, "| GET | 500 |", "/api/v1/layer"+ds)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Info.Println(data)
	// http.StatusNotFound
	js, err := json.Marshal(lyr)
	if err != nil {
		Info.Println(r.RemoteAddr, "| GET | 500 |", "/api/v1/layer"+ds)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "| GET | 200 |", "/api/v1/layer"+ds)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DeleteLayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	data := DB.deleteLayer(ds)
	js, err := json.Marshal(data)
	if err != nil {
		Info.Println(r.RemoteAddr, "| DELETE | 500 |", "/api/v1/layer"+ds)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "| DELETE | 200 |", "/api/v1/layer"+ds)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
