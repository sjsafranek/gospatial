package app

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func LoadLayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	Debug.Printf("Loading [%s]", ds)
	_, err := DB.getLayer(ds)
	// Datasource not found
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /management/load/"+ds+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Response
	data := `{"status":"ok","datasource":"` + ds + `", "result":"datasource loaded"}`
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /management/load/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "GET /management/load/"+ds+" [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func UnloadLayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	Debug.Printf("Unloading [%s]", ds)
	// unload
	delete(DB.Cache, ds)
	// Response
	data := `{"status":"ok","datasource":"` + ds + `", "result":"datasource unloaded"}`
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /management/unload/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "GET /management/unload/"+ds+" [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func LoadedLayers(w http.ResponseWriter, r *http.Request) {
	Debug.Println("Checking loaded datasources...")
	// Response
	js, err := json.Marshal(DB.Cache)
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /management/loaded [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Info.Println(r.RemoteAddr, "GET /management/loaded [200]")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
