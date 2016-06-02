package app

import (
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://sjsafranek.github.io/gospatial/", 200)
	return
}

/*=======================================*/
// Method: MapHandler
// Description:
//		Returns map client for layer
// @param apikey customer id
// @return map template
/*=======================================*/
func MapHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	/*=======================================*/
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

	/*=======================================*/

	// Return results
	map_tmpl := "./templates/map.html"
	tmpl, _ := template.ParseFiles(map_tmpl)
	network_logger_Info.Println(r.RemoteAddr, "GET /map [200]")
	tmpl.Execute(w, MapData{Apikey: apikey, Version: "1.9.2"})

}

func CustomerManagementHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	/*=======================================*/
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

	/*=======================================*/

	// Return results
	map_tmpl := "./templates/management.html"
	tmpl, _ := template.ParseFiles(map_tmpl)
	network_logger_Info.Println(r.RemoteAddr, "GET /management [200]")
	tmpl.Execute(w, MapData{Apikey: apikey, Version: "1.9.2"})

}
