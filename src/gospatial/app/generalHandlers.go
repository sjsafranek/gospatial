package app

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	// "regexp"
)

// // Register routes
// var validPath = regexp.MustCompile("^/(save|api/v1/layer)/([a-zA-Z0-9]+)$")

// func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		m := validPath.FindStringSubmatch(r.URL.Path)
// 		if m == nil {
// 			geojs := NewGeojson()
// 			ds, _ := NewUUID()
// 			lyr := Layer{Datasource: ds, Geojson: geojs}
// 			lyr.Save()
// 			http.Redirect(w, r, "/api/v1/layer/"+ds, http.StatusFound)
// 			return
// 		}
// 		fn(w, r, m[2])
// 	}
// }

type MapData struct {
	Datasource string
}

func MapHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := vars["ds"]
	if SuperuserKey == r.FormValue("apikey") {
		map_tmpl := "./templates/map_admin.html"
		tmpl, _ := template.ParseFiles(map_tmpl)
		tmpl.Execute(w, MapData{Datasource: ds})
	} else {
		map_tmpl := "./templates/map_standard.html"
		tmpl, _ := template.ParseFiles(map_tmpl)
		tmpl.Execute(w, MapData{Datasource: ds})
	}
}
