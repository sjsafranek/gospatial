package app

import (
	"html/template"
	"net/http"
)

// IndexHandler returns html page containing api docs
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
		networkLoggerError.Println(r.RemoteAddr, "POST /map [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	_, err := DB.GetCustomer(apikey)
	if err != nil {
		networkLoggerWarning.Println(r.RemoteAddr, "POST /map [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return results
	htmlFile := "./templates/map.html"
	tmpl, _ := template.ParseFiles(htmlFile)
	networkLoggerInfo.Println(r.RemoteAddr, "GET /map [200]")
	tmpl.Execute(w, PageViewData{Apikey: apikey, Version: "1.9.3"})

}

// DashboardHandler returns customer management gui.
// Allows customers to create and delete both geojson layers and tile baselayers.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	// Check for apikey in request
	if apikey == "" {
		networkLoggerError.Println(r.RemoteAddr, "POST /map [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	_, err := DB.GetCustomer(apikey)
	if err != nil {
		networkLoggerWarning.Println(r.RemoteAddr, "POST /map [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return results
	htmlFile := "./templates/management.html"
	tmpl, _ := template.ParseFiles(htmlFile)
	networkLoggerInfo.Println(r.RemoteAddr, "GET /management [200]")
	tmpl.Execute(w, PageViewData{Apikey: apikey, Version: "1.9.3"})

}
