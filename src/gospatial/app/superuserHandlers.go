package app

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

var SuperuserKey string = "stefanrocks"
var AppMode string = "standard"

type Message struct {
	Status  string `json:"status"`
	Message string `json:"Message"`
}

func DebugModeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	md := vars["md"]
	if SuperuserKey != r.FormValue("apikey") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if md == "debug" {
		js, err := json.Marshal(Message{Status: "ok", Message: "debug mode on"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		DebugMode(true)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else if md == "standard" {
		js, err := json.Marshal(Message{Status: "ok", Message: "debug mode off"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		DebugMode(false)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
