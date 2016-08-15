package app

import (
	"fmt"
	// "log"
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
	ServerLogger.Info("Magic happens on port ", self.Port)

	bind := fmt.Sprintf(":%v", self.Port)
	// bind := fmt.Sprintf("0.0.0.0:%v", port)

	err := http.ListenAndServe(bind, router)
	if err != nil {
		panic(err)
	}
	
}
