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
	Info.Println(r.RemoteAddr, "| POST\t|", "/api/v1/layer")
	geojs := NewGeojson()
	ds, _ := NewUUID()
	lyr := Layer{Datasource: ds, Geojson: geojs}
	lyr.Save()
	// Info.Println(r.RemoteAddr, ds, "new layer")
	data := `{"status":"ok","datasource":"` + ds + `"}`
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ViewLayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	// Info.Println(r.RemoteAddr, ds, "read layer")
	Info.Println(r.RemoteAddr, "| GET\t|", "/api/v1/layer/"+ds)

	lyr, err := DB.getLayer(ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Info.Println(data)
	// http.StatusNotFound
	js, err := json.Marshal(lyr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DeleteLayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	// Info.Println(r.RemoteAddr, ds, "delete layer")
	Info.Println(r.RemoteAddr, "| DELETE\t|", "/api/v1/layer/"+ds)

	data := DB.deleteLayer(ds)
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
