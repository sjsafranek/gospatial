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
	port          int
	database      string
	bind          string
	versionReport bool
	configFile    string
	debugMode bool
)

const (
	version       string = "1.10.2"
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
	flag.BoolVar(&versionReport, "v", false, "App Version")
	flag.BoolVar(&debugMode, "d", false, "Enable debug mode")
	flag.Parse()
	if versionReport {
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

	log.Println("Authkey:", app.SuperuserKey)
	log.Println("Database:", database)

	if debugMode {
		// https://golang.org/pkg/net/http/pprof/
		go func() {
			log.Printf("Profiling happens on port %v...\n", 6060)
			app.Info.Println(http.ListenAndServe(":6060", nil))
		}()
	}

	// Initiate Database
	app.DB = app.Database{File: database + ".db"}
	err := app.DB.Init()
	if err != nil {
		panic(err)
	}

	// start tcp server
	tcpServer := app.TcpServer{Host: "localhost", Port: "3333"}
	tcpServer.Start()

	// start http server
	httpServer := app.HttpServer{Port: port}
	httpServer.Start()

}
