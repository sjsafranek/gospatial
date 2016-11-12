package app

import (
	"fmt"
	"net/http"
)

// Sends http response
func SendJsonResponse(w http.ResponseWriter, js []byte) {
	// set response headers
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// write response content
	w.Write(js)
}

// check request for valid authkey
func CheckAuthKey(w http.ResponseWriter, r *http.Request) bool {
	if SuperuserKey != r.FormValue("authkey") {
		message := fmt.Sprintf(" %v %v [401]", r.Method, r.URL.Path)
		NetworkLogger.Error(r.RemoteAddr, message)
		// NetworkLogger.Error(r.RemoteAddr, +" "+r.Method+" "+r.URL.Path+" [401]")
		http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
		return false
	}
	return true
}

func GetApikeyFromRequest(w http.ResponseWriter, r *http.Request) string {
	// Get params
	apikey := r.FormValue("apikey")

	// Check for apikey in request
	if apikey == "" {
		http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
		ServerLogger.Info(r.URL.Path)
	}

	return apikey
}
