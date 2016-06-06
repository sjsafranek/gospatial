package app

import (
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var routes = routes{
	// Health check
	route{"Ping", "GET", "/ping", PingHandler},

	// Layers
	route{"ViewLayers", "GET", "/api/v1/layers", ViewLayersHandler},
	route{"ViewLayer", "GET", "/api/v1/layer/{ds}", ViewLayerHandler},
	route{"NewLayer", "POST", "/api/v1/layer", NewLayerHandler},
	route{"DeleteLayer", "DELETE", "/api/v1/layer/{ds}", DeleteLayerHandler},
	route{"ShareLayerHandler", "PUT", "/api/v1/layer/{ds}", ShareLayerHandler},

	//
	route{"NewFeature", "POST", "/api/v1/layer/{ds}/feature", NewFeatureHandler},
	route{"ViewFeature", "GET", "/api/v1/layer/{ds}/feature/{k}", ViewFeatureHandler},

	// Superuser routes
	route{"NewCustomerHandler", "POST", "/api/v1/customer", NewCustomerHandler},

	// Web Client routes
	route{"Index", "GET", "/", IndexHandler},
	route{"MapNew", "GET", "/map", MapHandler},
	route{"CustomerManagement", "GET", "/management", CustomerManagementHandler},

	// Web Socket route
	route{"Socket", "GET", "/ws/{ds}", serveWs},

	// Experimental
	route{"UnloadLayer", "GET", "/management/unload/{ds}", UnloadLayer},
	route{"LoadedLayers", "GET", "/management/loaded", LoadedLayers},
	route{"LoadedLayers", "GET", "/management/profile", ServerProfile},
}
