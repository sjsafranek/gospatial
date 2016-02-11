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
	Route{
		"ViewLayer",
		"GET",
		"/api/v1/layer/{ds}",
		ViewLayerHandler,
	},
	Route{
		"NewLayer",
		"POST",
		"/api/v1/layer",
		NewLayerHandler,
	},
	Route{
		"DeleteLayer",
		"DELETE",
		"/api/v1/layer/{ds}",
		DeleteLayerHandler,
	},
	// Features
	Route{
		"NewFeature",
		"POST",
		"/api/v1/layer/{ds}/feature",
		NewFeatureHandler,
	},
	Route{
		"ViewFeature",
		"GET",
		"/api/v1/layer/{ds}/feature/{k}",
		ViewFeatureHandler,
	},
	// Superuser Routes
	Route{
		"DebugMode",
		"GET",
		"/management/mode/debug",
		DebugModeHandler,
	},
	// Web Client Routes
	Route{
		"ClientLogging",
		"POST",
		"/log",
		WebClientLogHandler,
	},
	Route{
		"Map",
		"GET",
		"/map/{ds}",
		MapHandler,
	},
}