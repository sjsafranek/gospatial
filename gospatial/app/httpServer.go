package app

import (
	"fmt"
	"log"
	"net/http"
)

const HTTP_DEFAULT_PORT = 8080

type HttpServer struct {
	Port int
}

func (self HttpServer) Start() {

	// Report available routes
	log.Println("Attaching HTTP handler for route: [GET]    /")
	log.Println("Attaching HTTP handler for route: [GET]    /ping")
	// log.Println("Attaching HTTP handler for route: [GET]    /api/v1/layers")
	log.Println("Attaching HTTP handler for route: [POST]   /api/v1/layer")
	log.Println("Attaching HTTP handler for route: [GET]    /api/v1/layer/{ds}")
	log.Println("Attaching HTTP handler for route: [DELETE] /api/v1/layer/{ds}")
	// log.Println("Attaching HTTP handler for route: [PUT]    /api/v1/layer/{ds}")
	log.Println("Attaching HTTP handler for route: [POST]   /api/v1/layer/{ds}/feature")
	log.Println("Attaching HTTP handler for route: [GET]    /api/v1/layer/{ds}/feature/{k}")
	log.Println("Attaching HTTP handler for route: [POST]   /api/v1/tilelayer")
	log.Println("Attaching HTTP handler for route: [GET]    /api/v1/customer")
	log.Println("Attaching HTTP handler for route: [POST]   /api/v1/customer")
	log.Println("Attaching HTTP handler for route: [GET]    /management")
	log.Println("Attaching HTTP handler for route: [GET]    /map/{ds}")
	log.Println("Attaching HTTP handler for route: [GET]    /ws/{ds}")
	log.Println("Attaching HTTP handler for route: [GET]    /management/unload/{ds}")
	log.Println("Attaching HTTP handler for route: [GET]    /management/loaded")
	log.Println("Attaching HTTP handler for route: [GET]    /management/profile")

	// Attach Http Hanlders
	router := Router()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Start server
	log.Printf("Magic happens on port %v...\n", port)
	Info.Printf("Magic happens on port %v...\n", self.Port)
	bind := fmt.Sprintf(":%v", self.Port)
	// bind := fmt.Sprintf("0.0.0.0:%v", port)
	err := http.ListenAndServe(bind, router)
	if err != nil {
		panic(err)
	}
}