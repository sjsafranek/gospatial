package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/paulmach/go.geojson"
	"gospatial/utils"
	"io/ioutil"
	"net/http"
	// "strconv"
	"time"
)

// NewFeatureHandler creates a new feature and adds it to a layer.
// Layer is then saved to database. All active clients viewing layer
// are notified of update via websocket hub.
// @param apikey customer id
// @oaram ds datasource uuid
// @return json
func NewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Get request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		vars := mux.Vars(r)
		ds := vars["ds"]
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	// Get params
	//apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	/*=======================================*/
	// Check for apikey in request
	// if apikey == "" {
	// 	NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer/"+ds+"/feature [401]")
	// 	http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
	// 	return
	// }
	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer/"+ds+"/feature [401]")
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer/"+ds+"/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !utils.StringInSlice(ds, customer.Datasources) {
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
		return
	}

	// Unmarshal feature
	feat, err := geojson.UnmarshalFeature(body)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer/"+ds+"/feature [400]")
		// Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fill required attributes
	now := time.Now().Unix()
	feat.Properties["is_active"] = true
	feat.Properties["is_deleted"] = false
	feat.Properties["date_created"] = now
	feat.Properties["date_modified"] = now
	feat.Properties["geo_id"] = fmt.Sprintf("%v", now)
	ServerLogger.Info(feat.Properties)

	// Save feature to database
	err = DB.InsertFeature(ds, feat)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer/"+ds+"/feature [500]")
		// Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate message
	data := `{"status":"success","datasource":"` + ds + `", "message":"feature added"}`
	js, err := json.Marshal(data)
	if err != nil {
		NetworkLogger.Critical(r.RemoteAddr, " POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update websockets
	conn := connection{ds: ds, ip: r.RemoteAddr}
	Hub.broadcast(true, &conn)

	// Return results
	NetworkLogger.Info(r.RemoteAddr, " POST /api/v1/layer/"+ds+"/feature [200]")
	NetworkLogger.Debug("[Out] ", string(js))
	SendJsonResponse(w, js)

}

// ViewFeatureHandler finds feature in layer via array index. Returns feature geojson.
// @param apikey customer id
// @oaram ds datasource uuid
// @return feature geojson
func ViewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Get params
	//apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	/*=======================================*/
	// Check for apikey in request
	// if apikey == "" {
	// 	NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [401]")
	// 	http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
	// 	return
	// }
	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [401]")
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !utils.StringInSlice(ds, customer.Datasources) {
		NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [401]")
		http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Get layer from database
	data, err := DB.GetLayer(ds)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check for feature
	var js []byte
	for _, v := range data.Features {
		geo_id := fmt.Sprintf("%v", v.Properties["geo_id"])
		if geo_id == vars["k"] {
			js, err = v.MarshalJSON()
			if err != nil {
				NetworkLogger.Critical(r.RemoteAddr, " GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [500]")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Return results
			NetworkLogger.Info(r.RemoteAddr, " GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [200]")
			NetworkLogger.Debug("[Out] ", string(js))
			SendJsonResponse(w, js)
			return
		}
	}

	// Feature not found
	NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
	err = fmt.Errorf("Not found")
	http.Error(w, err.Error(), http.StatusNotFound)

}
