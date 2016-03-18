/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"flag"
	"fmt"
	"gospatial/app"
	"net/http"
	"os"
	"os/signal"
)

var (
	port     int
	database string
	bind     string
	debug    bool
	version  bool
)

const (
	VERSION string = "1.3.0 "
)

func init() {
	flag.IntVar(&port, "p", 8080, "server port")
	flag.StringVar(&database, "db", "bolt", "app database")
	flag.StringVar(&app.SuperuserKey, "s", "7q1qcqmsxnvw", "superuser key")
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

	// Graceful shut down
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		for sig := range sigs {
			// sig is a ^C, handle it
			fmt.Printf("%s \n", sig)
			app.Info.Println("Gracefulling shutting down")
			app.Info.Println("Waiting for sockets to close...")
			for {
				if len(app.Hub.Sockets) == 0 {
					app.Info.Println("Shutting down...")
					os.Exit(0)
				}
			}
		}
	}()

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
