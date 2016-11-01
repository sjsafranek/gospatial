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
	"gospatial/utils"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strings"
)

var (
	port          int
	database      string
	bind          string
	versionReport bool
	configFile    string
	debugMode     bool
)

const (
	VERSION        string = "1.10.5"
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
		app.ServerLogger.Error(err)
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
	flag.StringVar(&app.LogDirectory, "L", "log", "logging directory") // check if directory exists
	flag.StringVar(&app.LogLevel, "l", "trace", "logging level")

	flag.Parse()
	if versionReport {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}

	app.ResetLogging()

	// check if file exists!!!
	if _, err := os.Stat(configFile); err == nil {
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
		configuration := serverConfig{}
		err = json.Unmarshal(file, &configuration)
		if err != nil {
			fmt.Println("error:", err)
		}
		port = configuration.Port
		database = configuration.Db
		database = strings.Replace(database, ".db", "", -1)
		app.SuperuserKey = configuration.Authkey
		app.ServerLogger.Info(configuration)
	} else {
		// create config file
		configuration := serverConfig{}
		configuration.Port = port
		configuration.Db = database

		if "su" == app.SuperuserKey {
			authkey := utils.NewAPIKey(12)
			configuration.Authkey = authkey
			app.SuperuserKey = authkey
		}
		app.ServerLogger.Info(configuration)
	}

}

func main() {

	// app.ServerLogger.Trace("setting to default value")
	// app.ServerLogger.Debug("page request. url + params")
	// app.ServerLogger.Info("Server started")
	// app.ServerLogger.Warn("Cannot talk to database, using backup")
	// app.ServerLogger.Error("Cannot process request!")
	// app.ServerLogger.Critical("Shit BROKE. Shutting down...")

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
			// fmt.Printf("%s \n", sig)
			app.ServerLogger.Info("Recieved ", sig)
			app.ServerLogger.Info("Gracefully shutting down")
			app.ServerLogger.Info("Waiting for sockets to close...")
			for {
				if len(app.Hub.Sockets) == 0 {
					// app.Info.Println("Shutting down...")
					app.ServerLogger.Info("Shutting down...")
					os.Exit(0)
				}
			}
		}
	}()

	app.ServerLogger.Info("Authkey:", app.SuperuserKey)
	app.ServerLogger.Info("Database:", database)

	if debugMode {
		// https://golang.org/pkg/net/http/pprof/
		go func() {
			app.ServerLogger.Info("Profiling happens on port 6060\n")
			app.ServerLogger.Info(http.ListenAndServe(":6060", nil))
		}()
	}

	// Initiate Database
	app.COMMIT_LOG_FILE = database + "_commit.log"
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
