package app

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

var SuperuserKey string = "su"
var AppMode string = "standard"

func DebugModeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	md := vars["md"]
	if SuperuserKey != r.FormValue("authkey") {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if md == "debug" {
		js, err := json.Marshal(`{"status": "ok", "message": "debug mode on"}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		DebugMode(true)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else if md == "standard" {
		js, err := json.Marshal(`{"status": "ok", "message": "debug mode off"}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		DebugMode(false)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// SUPERUSER
func NewCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if SuperuserKey != r.FormValue("authkey") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	} else {
		// new customer
		apikey := NewAPIKey(12)
		customer := Customer{Apikey: apikey}
		err := DB.insertCustomer(customer)
		if err != nil {
			Error.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// return results
		data := `{"status":"ok","apikey":"` + apikey + `", "result":"customer created"}`
		js, err := json.Marshal(data)
		if err != nil {
			Error.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		DebugMode(false)
		w.Header().Set("Content-Type", "application/json")
		// allow cross domain AJAX requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	}
}
