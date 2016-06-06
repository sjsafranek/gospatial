package app

import (
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://sjsafranek.github.io/gospatial/", 200)
	return
}

// MapHandler returns leaflet map view for customer layers
// @param apikey customer id
// @return map template
func MapHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	// Check for apikey in request
	if apikey == "" {
		network_logger_Error.Println(r.RemoteAddr, "POST /map [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	_, err := DB.GetCustomer(apikey)
	if err != nil {
		network_logger_Warning.Println(r.RemoteAddr, "POST /map [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return results
	htmlFile := "./templates/map.html"
	tmpl, _ := template.ParseFiles(htmlFile)
	network_logger_Info.Println(r.RemoteAddr, "GET /map [200]")
	tmpl.Execute(w, MapData{Apikey: apikey, Version: "1.9.3"})

}

// CustomerManagementHandler returns customer management gui. Allows customers to create and delete both geojson layers and tile baselayers.
func CustomerManagementHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")


	// Check for apikey in request
	if apikey == "" {
		network_logger_Error.Println(r.RemoteAddr, "POST /map [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	_, err := DB.GetCustomer(apikey)
	if err != nil {
		network_logger_Warning.Println(r.RemoteAddr, "POST /map [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}



	// Return results
	htmlFile := "./templates/management.html"
	tmpl, _ := template.ParseFiles(htmlFile)
	network_logger_Info.Println(r.RemoteAddr, "GET /management [200]")
	tmpl.Execute(w, MapData{Apikey: apikey, Version: "1.9.3"})

}
