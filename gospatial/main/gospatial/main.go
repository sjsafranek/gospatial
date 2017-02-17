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
	tcp_port      int
	database      string
	bind          string
	versionReport bool
	configFile    string
	debugMode     bool
	configuration serverConfig
)

const (
	VERSION             string = "1.11.3"
	DEFAULT_CONFIG_FILE string = "config.json"
	DEFAULT_HTTP_PORT   int    = 8080
	DEFAULT_TCP_PORT    int    = 3333
)

type serverConfig struct {
	HttpPort int    `json:"http_port"`
	TcpPort  int    `json:"tcp_port"`
	Db       string `json:"db"`
	Authkey  string `json:"authkey"`
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		app.ServerLogger.Error(err)
	}
	db := strings.Replace(dir, "bin", "bolt", -1)
	flag.StringVar(&configFile, "c", DEFAULT_CONFIG_FILE, "server config file")
	flag.IntVar(&port, "p", DEFAULT_HTTP_PORT, "http server port")
	flag.IntVar(&tcp_port, "tcp_port", DEFAULT_TCP_PORT, "tcp server port")
	//flag.IntVar(&port, "p", 8888, "server port")
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

	app.Version = VERSION

	app.ResetLogging()

	// check if config file exists!!!
	if _, err := os.Stat(configFile); err != nil {

		// create config object from commandline args
		configuration = serverConfig{}
		configuration.HttpPort = port
		configuration.TcpPort = tcp_port
		configuration.Db = database

		// superuser key
		if "su" == app.SuperuserKey {
			authkey := utils.NewAPIKey(12)
			configuration.Authkey = authkey
			app.SuperuserKey = authkey
		}

		// write to file
		configJson, _ := json.Marshal(configuration)
		err := ioutil.WriteFile(configFile, configJson, 0644)
		if err != nil {
			panic(err)
		}

	} else {
		// read config file
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			panic(err)
		}

		// build config object from file contents
		configuration = serverConfig{}
		err = json.Unmarshal(file, &configuration)
		if err != nil {
			panic(err)
		}

		// apply commandline args as overrides
		if "su" == app.SuperuserKey {
			app.SuperuserKey = configuration.Authkey
		}

		if DEFAULT_HTTP_PORT != port {
			configuration.HttpPort = port
		}
		if DEFAULT_TCP_PORT != tcp_port {
			configuration.TcpPort = tcp_port
		}

		//configuration.Db = strings.Replace(database, ".db", "", -1) //database
		//app.ServerLogger.Info(strings.Replace(database, ".db", "", -1))
	}

}

func main() {

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
			app.ServerLogger.Info("Recieved ", sig)
			app.ServerLogger.Info("Gracefully shutting down")
			app.ServerLogger.Info("Waiting for sockets to close...")
			for {
				if len(app.Hub.Sockets) == 0 {
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

	app.ServerLogger.Info(configuration)

	// start tcp server
	//tcpServer := app.TcpServer{Host: "localhost", Port: "3333"}
	tcpServer := app.TcpServer{Host: "localhost", Port: fmt.Sprintf("%v", configuration.TcpPort)}
	tcpServer.Start()

	// start http server
	httpServer := app.HttpServer{Port: configuration.HttpPort}
	httpServer.Start()

}
