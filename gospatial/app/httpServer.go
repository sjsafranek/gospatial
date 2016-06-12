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
	// Attach Http Hanlders
	router := Router()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Start server
	log.Printf("Magic happens on port %v...\n", self.Port)
	Info.Printf("Magic happens on port %v...\n", self.Port)
	bind := fmt.Sprintf(":%v", self.Port)
	// bind := fmt.Sprintf("0.0.0.0:%v", port)
	err := http.ListenAndServe(bind, router)
	if err != nil {
		panic(err)
	}
}
