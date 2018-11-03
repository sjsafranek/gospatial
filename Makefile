##=======================================================================##
## Makefile
## Created: Wed Aug 05 14:35:14 PDT 2015 @941 /Internet Time/
# :mode=makefile:tabSize=3:indentSize=3:
## Purpose: 
##======================================================================##

SHELL=/bin/bash
PROJECT_NAME = gospatial
GPATH = $(shell pwd)

.PHONY: fmt deps test install build scrape clean

install: fmt deps
	@GOPATH=${GPATH} go build -o gospatial_server server.go
	@GOPATH=${GPATH} go build -o gospatial_importer importer.go
	@GOPATH=${GPATH} go build -o gospatial_ts timeseries_tool.go

build: fmt deps
	@GOPATH=${GPATH} go build -o gospatial_server server.go
	@GOPATH=${GPATH} go build -o gospatial_importer importer.go
	@GOPATH=${GPATH} go build -o gospatial_ts timeseries_tool.go

deps:
	mkdir -p "src"
	mkdir -p "pkg"
	mkdir -p "log"
	@GOPATH=${GPATH} go get github.com/boltdb/bolt
	@GOPATH=${GPATH} go get github.com/cihub/seelog
	@GOPATH=${GPATH} go get github.com/gorilla/mux
	@GOPATH=${GPATH} go get github.com/gorilla/websocket
	@GOPATH=${GPATH} go get github.com/paulmach/go.geojson
	@GOPATH=${GPATH} go get github.com/sjsafranek/DiffDB/diff_store
	@GOPATH=${GPATH} go get github.com/sjsafranek/DiffDB/diff_db

fmt:
	@GOPATH=${GPATH} gofmt -s -w ${PROJECT_NAME}
	@GOPATH=${GPATH} gofmt -s -w server.go
	@GOPATH=${GPATH} gofmt -s -w importer.go
	@GOPATH=${GPATH} gofmt -s -w timeseries_tool.go

test:
	##./tcp_test.sh
	./benchmark.sh

scrape:
	@find src -type d -name '.hg' -or -type d -name '.git' | xargs rm -rf

clean:
	@GOPATH=${GPATH} go clean
