package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type MapData struct {
	Datasource string
}

func MapHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	if SuperuserKey == r.FormValue("apikey") {
		map_tmpl := "./templates/map_admin.html"
		tmpl, _ := template.ParseFiles(map_tmpl)
		Info.Println(r.RemoteAddr, "GET /map/"+ds+" [200]")
		tmpl.Execute(w, MapData{Datasource: ds})
	} else if r.FormValue("apikey") != "" {
		Warning.Println(r.RemoteAddr, "GET /map/"+ds+" [401]")
		err := fmt.Errorf("Unauthorized")
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		map_tmpl := "./templates/map_standard.html"
		tmpl, _ := template.ParseFiles(map_tmpl)
		Info.Println(r.RemoteAddr, "GET /map/"+ds+" [200]")
		tmpl.Execute(w, MapData{Datasource: ds})
	}
}
