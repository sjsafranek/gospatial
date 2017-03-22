package gospatial

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Router for http api calls
func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		// log.Println("Attaching HTTP handler for route:", route.Method, route.Pattern)
		ServerLogger.Info("Attaching HTTP handler for route: ", route.Method, " ", route.Pattern)
		handler = route.HandlerFunc
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
