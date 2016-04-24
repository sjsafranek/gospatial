package app

import (
	"github.com/gorilla/mux"
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
// @oaram ds datasource uuid
// @return map template
/*=======================================*/
func MapHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	// Superuser access
	if SuperuserKey == apikey {
		map_tmpl := "./templates/map_admin.html"
		tmpl, _ := template.ParseFiles(map_tmpl)
		network_logger_Info.Println(r.RemoteAddr, "GET /map/"+ds+" [200]")
		tmpl.Execute(w, MapData{Datasource: ds})
	}

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		network_logger_Error.Println(r.RemoteAddr, "POST /map/"+ds+" [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		network_logger_Warning.Println(r.RemoteAddr, "POST /map/"+ds+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		network_logger_Error.Println(r.RemoteAddr, "POST /map/"+ds+" [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Return results
	map_tmpl := "./templates/map_admin.html"
	tmpl, _ := template.ParseFiles(map_tmpl)
	network_logger_Info.Println(r.RemoteAddr, "GET /map/"+ds+" [200]")
	tmpl.Execute(w, MapData{Datasource: ds, Apikey: apikey})

}

func MapHandlerNew(w http.ResponseWriter, r *http.Request) {

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
	map_tmpl := "./templates/new_map.html"
	tmpl, _ := template.ParseFiles(map_tmpl)
	network_logger_Info.Println(r.RemoteAddr, "GET /map [200]")
	tmpl.Execute(w, MapData{Apikey: apikey})

}
