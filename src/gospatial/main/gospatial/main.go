package main

import (
	"flag"
	"fmt"
	"os"
	// "github.com/gorilla/mux"
	"gospatial/app"
	"net/http"
)

var (
	port     int
	database string
	bind     string
	debug    bool
	version  bool
)

const VERSION string = "1.2.0 "

func init() {
	flag.IntVar(&port, "p", 8080, "server port")
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.StringVar(&app.SuperuserKey, "s", "su", "superuser key")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.BoolVar(&version, "v", false, "App Version")
	flag.Parse()
	if version {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}
}

func main() {
	app.DebugMode(debug)

	// Initiate Database
	app.DB = app.Database{File: "./" + database + ".db"}
	app.DB.Init()

	router := app.NewRouter()

	// Server static folder
	// router.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Start server
	app.Info.Printf("Magic happens on port %v...\n", port)
	if app.AppMode == "debug" {
		fmt.Printf("Magic happens on port %v...\n", port)
	}

	bind := fmt.Sprintf(":%v", port)
	err := http.ListenAndServe(bind, router)
	if err != nil {
		panic(err)
	}
}
