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
	"runtime/pprof"
	"gospatial/utils"
)

import mylogger "gospatial/logs"

var (
	port          int
	database      string
	bind          string
	versionReport bool
	configFile    string
	debugMode bool
)

const (
	VERSION       string = "1.10.4"
	DEFAULT_CONFIG string = "config.json"
)

type serverConfig struct {
	Port    int    `json:"port"`
	Db      string `json:"db"`
	Authkey string `json:"authkey"`
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		app.Error.Fatal(err)
	}
	db := strings.Replace(dir, "bin", "bolt", -1)
	flag.StringVar(&configFile, "c", DEFAULT_CONFIG, "server config file")
	flag.IntVar(&port, "p", 8080, "server port")
	flag.StringVar(&database, "db", db, "app database")
	// flag.StringVar(&app.SuperuserKey, "s", "7q1qcqmsxnvw", "superuser key")
	flag.StringVar(&app.SuperuserKey, "s", "su", "superuser key")
	flag.BoolVar(&versionReport, "V", false, "App Version")
	flag.BoolVar(&app.Verbose, "v", false, "verbose")
	flag.BoolVar(&debugMode, "d", false, "Enable debug mode")
	// TODO:
	// 		log directory
	// 		log level
	// 
	flag.Parse()
	if versionReport {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}
	
	// check if file exists!!!
	if _, err := os.Stat(configFile); err == nil {
		// fmt.Println(configFile)
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
		configuration := serverConfig{}
		err = json.Unmarshal(file, &configuration)
		if err != nil {
			fmt.Println("error:", err)
		}
		// app.Info.Printf("%v\n", configuration)
		port = configuration.Port
		database = configuration.Db
		database = strings.Replace(database, ".db", "", -1)
		app.SuperuserKey = configuration.Authkey
		mylogger.Logger.Info(configuration)
	} else {
		// create config file
		configuration := serverConfig{}
		configuration.Port = port
		configuration.Db = database
		authkey := utils.NewAPIKey(12)
		configuration.Authkey = authkey
		app.SuperuserKey = authkey
		// app.Info.Printf("%v\n", configuration)
		mylogger.Logger.Info(configuration)
	}

}

func main() {

	// mylogger.Logger.Trace("setting to default value")
	// mylogger.Logger.Debug("page request. url + params")
	// mylogger.Logger.Info("Server started")
	// mylogger.Logger.Warn("Cannot talk to database, using backup")
	// mylogger.Logger.Error("Cannot process request!")
	// mylogger.Logger.Critical("Shit BROKE. Shutting down...")

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}


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
			// app.Info.Println("Gracefulling shutting down")
			// app.Info.Println("Waiting for sockets to close...")
			mylogger.Logger.Info("Gracefulling shutting down")
			mylogger.Logger.Info("Waiting for sockets to close...")
			for {
				if len(app.Hub.Sockets) == 0 {
					// app.Info.Println("Shutting down...")
					mylogger.Logger.Info("Shutting down...")
					os.Exit(0)
				}
			}
		}
	}()

	// log.Println("Authkey:", app.SuperuserKey)
	// log.Println("Database:", database)
	mylogger.Logger.Info("Authkey:", app.SuperuserKey)
	mylogger.Logger.Info("Database:", database)

	if debugMode {
		// https://golang.org/pkg/net/http/pprof/
		go func() {
			// log.Printf("Profiling happens on port %v...\n", 6060)
			mylogger.Logger.Info("Profiling happens on port %v...\n", 6060)
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
