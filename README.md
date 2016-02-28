# GoSpatial
Golang RESTful api for writing, storing, and serving GeoJSON data. GoSpatial also includes a mapping front end powered by Leaflet.js. Users are able to see what other 

## Requirements
	github.com/gorilla/websocket
	github.com/gorilla/mux
	github.com/boltdb/bolt

## Install
Get required golang packages by using `go get` or `getRequirements.sh`. Run Makefile.

Example:

	./getReuirements.sh
	make install

## Run
Execute the binary file produced in the projects `bin` directory

Example:

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

 - `-d`

