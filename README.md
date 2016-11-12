# GoSpatial
[![Go Report Card](https://goreportcard.com/badge/github.com/sjsafranek/gospatial)](https://goreportcard.com/report/github.com/sjsafranek/gospatial)
[![Version 1.10.4](https://img.shields.io/badge/version-1.10-brightgreen.svg)](http://sjsafranek.github.io/gospatial/)
[![License MIT License](https://img.shields.io/github/license/mashape/apistatus.svg)](http://sjsafranek.github.io/gospatial/)

Full documentation: http://sjsafranek.github.io/gospatial/

Golang RESTful api for writing, storing, and serving GeoJSON data. GoSpatial also includes a mapping front end powered by Leaflet.js. Users are able to see eachothers work before submitting.


## Install
``./install.sh`` will install the following packages and setup the workspace:

	github.com/gorilla/websocket
	github.com/gorilla/mux
	github.com/boltdb/bolt

Run ``make install`` to build the binary for the application


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

### Service File

	vim /lib/systemd/system/gospatial.service
	systemctl daemon-reload

### Restore database from commit log

	`nc localhost 3333 < test_commit.log 2>&1 | tee -a file_load.log`

