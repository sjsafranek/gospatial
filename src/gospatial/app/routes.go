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
	Info.Println("Attaching HTTP handler for route: [GET] /ping")
	Info.Println("Attaching HTTP handler for route: [GET] /api/v1/layers")
	Info.Println("Attaching HTTP handler for route: [GET] /api/v1/layer/{ds}")
	Info.Println("Attaching HTTP handler for route: [POST] /api/v1/layer")
	Info.Println("Attaching HTTP handler for route: [DELETE] /api/v1/layer/{ds}")
	Info.Println("Attaching HTTP handler for route: [PUT] /api/v1/layer/{ds}")
	Info.Println("Attaching HTTP handler for route: [POST] /api/v1/layer/{ds}/feature")
	Info.Println("Attaching HTTP handler for route: [GET] /api/v1/layer/{ds}/feature/{k}")
	Info.Println("Attaching HTTP handler for route: [POST] /management/customer")
	Info.Println("Attaching HTTP handler for route: [GET] /map/{ds}")
	Info.Println("Attaching HTTP handler for route: [GET] /")
	Info.Println("Attaching HTTP handler for route: [GET] /ws/{ds}")
	Info.Println("Attaching HTTP handler for route: [GET] /management/unload/{ds}")
	Info.Println("Attaching HTTP handler for route: [GET] /management/loaded")
	Info.Println("Attaching HTTP handler for route: [GET] /management/profile")
}

var routes = Routes{
	// General
	Route{"Ping", "GET", "/ping", PingHandler},

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
	Route{"NewCustomerHandler", "POST", "/management/customer", NewCustomerHandler},

	// Web Client Routes
	Route{"Map", "GET", "/map/{ds}", MapHandler},
	Route{"MapNew", "GET", "/map", MapHandlerNew},
	Route{"Index", "GET", "/", IndexHandler},

	// Web Socket Route
	Route{"Socket", "GET", "/ws/{ds}", serveWs},

	// Experimental
	Route{"UnloadLayer", "GET", "/management/unload/{ds}", UnloadLayer},
	Route{"LoadedLayers", "GET", "/management/loaded", LoadedLayers},
	Route{"LoadedLayers", "GET", "/management/profile", server_profile},
}
