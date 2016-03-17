package app

import (
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type MapData struct {
	Datasource string
	Apikey     string
}

func MapHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	if SuperuserKey == r.FormValue("apikey") {
		map_tmpl := "./templates/map_admin.html"
		tmpl, _ := template.ParseFiles(map_tmpl)
		Info.Println(r.RemoteAddr, "GET /map/"+ds+" [200]")
		tmpl.Execute(w, MapData{Datasource: ds})
	}
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
	map_tmpl := "./templates/map_admin.html"
	tmpl, _ := template.ParseFiles(map_tmpl)
	Info.Println(r.RemoteAddr, "GET /map/"+ds+" [200]")
	tmpl.Execute(w, MapData{Datasource: ds, Apikey: r.FormValue("apikey")})
	// else if r.FormValue("apikey") != "" {
	// 	Warning.Println(r.RemoteAddr, "GET /map/"+ds+" [401]")
	// 	err := fmt.Errorf("Unauthorized")
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// }
	// else {
	// 	map_tmpl := "./templates/map_standard.html"
	// 	tmpl, _ := template.ParseFiles(map_tmpl)
	// 	Info.Println(r.RemoteAddr, "GET /map/"+ds+" [200]")
	// 	tmpl.Execute(w, MapData{Datasource: ds})
	// }
}

// SUPERUSER
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	if SuperuserKey != r.FormValue("apikey") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	} else {
		// new customer
		customer := Customer{Apikey: NewAPIKey(12)}
		apikey, err := DB.insertCustomer(customer)
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
		w.Write(js)
	}
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer from db
	vars := mux.Vars(r)
	key := vars["key"]
	customer, err := DB.getCustomer(key)
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /management/customer/"+key+" [404]")
		// Warning.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// return results
	js, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	DebugMode(false)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
