package app

import (
	"fmt"
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

	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		return
	}

	_, err := GetCustomerFromDatabase(w, r, apikey)
	if err != nil {
		return
	}

	// Return results
	htmlFile := "./templates/map.html"
	tmpl, _ := template.ParseFiles(htmlFile)
	message := fmt.Sprintf(" %v %v [200]", r.Method, r.URL.Path)
	NetworkLogger.Info(r.RemoteAddr, message)
	tmpl.Execute(w, PageViewData{Apikey: apikey, Version: VERSION})

}

// DashboardHandler returns customer management gui.
// Allows customers to create and delete both geojson layers and tile baselayers.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		return
	}

	_, err := GetCustomerFromDatabase(w, r, apikey)
	if err != nil {
		return
	}

	// Return results
	htmlFile := "./templates/management.html"
	tmpl, _ := template.ParseFiles(htmlFile)
	message := fmt.Sprintf(" %v %v [200]", r.Method, r.URL.Path)
	NetworkLogger.Info(r.RemoteAddr, message)
	tmpl.Execute(w, PageViewData{Apikey: apikey, Version: VERSION})

}
