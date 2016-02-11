package app

import (
	"encoding/json"
	// "github.com/gorilla/mux"
	"net/http"
)

var SuperuserKey string = "stefanrocks"
var AppMode string = "standard"

type Message struct {
	Status  string `json:"status"`
	Message string `json:"Message"`
}

func DebugModeHandler(w http.ResponseWriter, r *http.Request) {
	if SuperuserKey != r.FormValue("apikey") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	js, err := json.Marshal(Message{Status: "ok", Message: "debug mode on"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	DebugMode(true)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
