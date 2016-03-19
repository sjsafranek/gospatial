package app

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	// Layers
	Route{"ViewLayers", "GET", "/api/v1/layers", ViewLayersHandler},
	Route{"ViewLayer", "GET", "/api/v1/layer/{ds}", ViewLayerHandler},
	Route{"NewLayer", "POST", "/api/v1/layer", NewLayerHandler},
	Route{"DeleteLayer", "DELETE", "/api/v1/layer/{ds}", DeleteLayerHandler},

	// Features
	Route{"NewFeature", "POST", "/api/v1/layer/{ds}/feature", NewFeatureHandler},
	Route{"ViewFeature", "GET", "/api/v1/layer/{ds}/feature/{k}", ViewFeatureHandler},

	// Superuser Routes
	Route{"DebugMode", "GET", "/management/mode/{md}", DebugModeHandler},
	Route{"NewCustomerHandler", "POST", "/management/customer", NewCustomerHandler},

	// Web Client Routes
	Route{"Map", "GET", "/map/{ds}", MapHandler},
	Route{"Index", "GET", "/", IndexHandler},

	// Web Socket Route
	Route{"Socket", "GET", "/ws/{ds}", serveWs},

	// Experimental
	Route{"LoadLayer", "GET", "/management/load/{ds}", LoadLayer},
	Route{"UnloadLayer", "GET", "/management/unload/{ds}", UnloadLayer},
	Route{"LoadedLayers", "GET", "/management/loaded", LoadedLayers},
	Route{"LoadedLayers", "GET", "/management/profile", server_profile},
}
