# GoSpatial
Golang RESTful api for writing, storing, and serving GeoJSON data 

## Requirements
	github.com/gorilla/websocket
	github.com/gorilla/mux
	github.com/boltdb/bolt

## Install
 - Get required golang packages
 - ./getReuirements.sh
 - make install

## Run
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
