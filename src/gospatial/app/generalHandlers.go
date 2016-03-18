package app

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

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
		Info.Println(r.RemoteAddr, "GET /map/"+ds+" [200]")
		tmpl.Execute(w, MapData{Datasource: ds})
	}

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.getCustomer(apikey)
	if err != nil {
		Warning.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Return results
	map_tmpl := "./templates/map_admin.html"
	tmpl, _ := template.ParseFiles(map_tmpl)
	Info.Println(r.RemoteAddr, "GET /map/"+ds+" [200]")
	tmpl.Execute(w, MapData{Datasource: ds, Apikey: apikey})

}
