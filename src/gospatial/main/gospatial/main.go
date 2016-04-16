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
	port        int
	database    string
	bind        string
	version     bool
	config_file string
)

const (
	VERSION        string = "1.8.0"
	default_config string = ""
)

type Configuration struct {
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
	flag.StringVar(&config_file, "c", default_config, "server config file")
	flag.IntVar(&port, "p", 8080, "server port")
	flag.StringVar(&database, "db", db, "app database")
	// flag.StringVar(&app.SuperuserKey, "s", "7q1qcqmsxnvw", "superuser key")
	flag.StringVar(&app.SuperuserKey, "s", "su", "superuser key")
	flag.BoolVar(&version, "v", false, "App Version")
	flag.Parse()
	if version {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}
	if config_file != "" {
		file, err := ioutil.ReadFile(config_file)
		if err != nil {
			panic(err)
		}
		configuration := Configuration{}
		err = json.Unmarshal(file, &configuration)
		if err != nil {
			fmt.Println("error:", err)
		}
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
					// app.DB.Backup("backup")
					app.Info.Println("Shutting down...")
					os.Exit(0)
				}
			}
		}
	}()

	app.DebugMode()
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
	// app.DB.Backup("backup")

	// Attach Http Hanlders
	app.AttachHttpHandlers()
	router := app.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Start server
	app.Info.Printf("Magic happens on port %v...\n", port)
	bind := fmt.Sprintf(":%v", port)
	// bind := fmt.Sprintf("0.0.0.0:%v", port)
	err := http.ListenAndServe(bind, router)
	if err != nil {
		panic(err)
	}
}
