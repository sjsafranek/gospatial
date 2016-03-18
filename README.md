# GoSpatial
[![Go Report Card](https://goreportcard.com/badge/github.com/sjsafranek/gospatial)](https://goreportcard.com/report/github.com/sjsafranek/gospatial)

Documentation: http://sjsafranek.github.io/gospatial/

Golang RESTful api for writing, storing, and serving GeoJSON data. GoSpatial also includes a mapping front end powered by Leaflet.js. Users are able to see eachothers work before submitting.

## Requirements
	github.com/gorilla/websocket
	github.com/gorilla/mux
	github.com/boltdb/bolt

## Install
Get required golang packages by using `go get` or `getRequirements.sh`. Run Makefile.

	./getReuirements.sh
	make install

## Run
Execute the binary file produced in the projects `bin` directory

 	./bin/gospatial


### Command Line Agruments

	Usage of ./bin/gospatial:
	  -d	debug mode
	  -db string
	    	app database (default "bolt")
	  -p int
	    	server port (default 8080)
	  -s string
	    	superuser key (default "su")
	  -v	App Version

 - `-d`: places the server into "debug mode". While the server app is in this mode, logs will be written to a log file.
 - `-db`: Specifies what database file to use. Default database is `bolt.db`.
 - `-p`: Specifies the server port. Default port is `8080`.
 - `-s`: Specifies the superuser key for management routes. Default key is `su`.
 - `-v`: Prints the app version


### Routes

#### Layers:
	[GET]    /api/v1/layers?apikey=<apikey>
	[POST]   /api/v1/layer?apikey=<apikey>
	[GET]    /api/v1/layer/{ds}?apikey=<apikey>
	[DELETE] /api/v1/layer/{ds}?apikey=<apikey>

#### Features:
	[POST]   /api/v1/layer/{ds}/feature?apikey=<apikey>
	[GET]    /api/v1/layer/{ds}/feature/{k}?apikey=<apikey>

#### Web Client
	[GET]    /map/{ds}?apikey=<apikey>
	[WS]     /ws/{ds}
	
#### Management	
	[POST]   /management/customer?apikey=<superuserkey>
	[GET]    /management/mode/{md}
	[GET]    /management/load/{ds}"
	[GET]    /management/unload/{ds}
	[GET]    /management/loaded
	[GET]    /management/profile
