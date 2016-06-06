/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gospatial/app"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

var (
	port       int
	database   string
	bind       string
	version    bool
	configFile string
)

const (
	version       string = "1.9.3"
	configDefault string = ""
)

type serverConfig struct {
	Port    int    `json:"port"`
	Db      string `json:"db"`
	Authkey string `json:"authkey"`
}

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		app.Error.Fatal(err)
	}
	db := strings.Replace(dir, "bin", "bolt", -1)
	flag.StringVar(&configFile, "c", configDefault, "server config file")
	flag.IntVar(&port, "p", 8080, "server port")
	flag.StringVar(&database, "db", db, "app database")
	// flag.StringVar(&app.SuperuserKey, "s", "7q1qcqmsxnvw", "superuser key")
	flag.StringVar(&app.SuperuserKey, "s", "su", "superuser key")
	flag.BoolVar(&version, "v", false, "App Version")
	flag.Parse()
	if version {
		fmt.Println("Version:", version)
		os.Exit(0)
	}
	if configFile != "" {
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
		configuration := serverConfig{}
		err = json.Unmarshal(file, &configuration)
		if err != nil {
			fmt.Println("error:", err)
		}
		app.Info.Printf("%v\n", configuration)
		port = configuration.Port
		database = configuration.Db
		database = strings.Replace(database, ".db", "", -1)
		app.SuperuserKey = configuration.Authkey
	}
}

func main() {

	// source: http://patorjk.com/software/taag/#p=display&f=Slant&t=Gospatial
	// HyperCube Platforms
	fmt.Println(`
   ______                       __  _       __
  / ____/___  _________  ____ _/ /_(_)___ _/ /
 / / __/ __ \/ ___/ __ \/ __ '/ __/ / __ '/ /
/ /_/ / /_/ (__  ) /_/ / /_/ / /_/ / /_/ / /
\____/\____/____/ .___/\__,_/\__/_/\__,_/_/
               /_/
	`)

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

	app.Network_logger_init()
	// app.StdOutMode()

	log.Println("Authkey:", app.SuperuserKey)
	log.Println("Database:", database)
	log.Printf("Profiling happens on port %v...\n", 6060)
	log.Printf("Magic happens on port %v...\n", port)
	// https://golang.org/pkg/net/http/pprof/
	go func() {
		app.Info.Println(http.ListenAndServe(":6060", nil))
	}()

	// Initiate Database
	app.DB = app.Database{File: database + ".db"}
	app.DB.Init()

	// Attach Http Hanlders
	router := app.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Report available routes
	Info.Println("Attaching HTTP handler for route: [GET] /ping")
	Info.Println("Attaching HTTP handler for route: [GET] /api/v1/layers")
	Info.Println("Attaching HTTP handler for route: [GET] /api/v1/layer/{ds}")
	Info.Println("Attaching HTTP handler for route: [POST] /api/v1/layer")
	Info.Println("Attaching HTTP handler for route: [DELETE] /api/v1/layer/{ds}")
	Info.Println("Attaching HTTP handler for route: [PUT] /api/v1/layer/{ds}")
	Info.Println("Attaching HTTP handler for route: [POST] /api/v1/layer/{ds}/feature")
	Info.Println("Attaching HTTP handler for route: [GET] /api/v1/layer/{ds}/feature/{k}")
	Info.Println("Attaching HTTP handler for route: [POST] /api/v1/customer")
	Info.Println("Attaching HTTP handler for route: [GET] /")
	Info.Println("Attaching HTTP handler for route: [GET] /map/{ds}")
	Info.Println("Attaching HTTP handler for route: [GET] /management")
	Info.Println("Attaching HTTP handler for route: [GET] /ws/{ds}")
	Info.Println("Attaching HTTP handler for route: [GET] /management/unload/{ds}")
	Info.Println("Attaching HTTP handler for route: [GET] /management/loaded")
	Info.Println("Attaching HTTP handler for route: [GET] /management/profile")

	// Start server
	app.Info.Printf("Magic happens on port %v...\n", port)
	bind := fmt.Sprintf(":%v", port)
	// bind := fmt.Sprintf("0.0.0.0:%v", port)
	err := http.ListenAndServe(bind, router)
	if err != nil {
		panic(err)
	}
}
