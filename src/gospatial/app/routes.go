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

func AttachHttpHandlers() {
	Trace.Println("Attaching HTTP handler for route: [GET] /api/v1/layers")
	Trace.Println("Attaching HTTP handler for route: [GET] /api/v1/layer/{ds}")
	Trace.Println("Attaching HTTP handler for route: [POST] /api/v1/layer")
	Trace.Println("Attaching HTTP handler for route: [DELETE] /api/v1/layer/{ds}")
	Trace.Println("Attaching HTTP handler for route: [PUT] /api/v1/layer/{ds}")
	Trace.Println("Attaching HTTP handler for route: [POST] /api/v1/layer/{ds}/feature")
	Trace.Println("Attaching HTTP handler for route: [GET] /api/v1/layer/{ds}/feature/{k}")
	Trace.Println("Attaching HTTP handler for route: [GET] /management/mode/{md}")
	Trace.Println("Attaching HTTP handler for route: [POST] /management/customer")
	Trace.Println("Attaching HTTP handler for route: [GET] /map/{ds}")
	Trace.Println("Attaching HTTP handler for route: [GET] /")
	Trace.Println("Attaching HTTP handler for route: [GET] /ws/{ds}")
	Trace.Println("Attaching HTTP handler for route: [GET] /management/load/{ds}")
	Trace.Println("Attaching HTTP handler for route: [GET] /management/unload/{ds}")
	Trace.Println("Attaching HTTP handler for route: [GET] /management/loaded")
	Trace.Println("Attaching HTTP handler for route: [GET] /management/profile")
}

var routes = Routes{
	// Layers
	Route{"ViewLayers", "GET", "/api/v1/layers", ViewLayersHandler},
	Route{"ViewLayer", "GET", "/api/v1/layer/{ds}", ViewLayerHandler},
	Route{"NewLayer", "POST", "/api/v1/layer", NewLayerHandler},
	Route{"DeleteLayer", "DELETE", "/api/v1/layer/{ds}", DeleteLayerHandler},
	Route{"ShareLayerHandler", "PUT", "/api/v1/layer/{ds}", ShareLayerHandler},

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
