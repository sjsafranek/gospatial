##=======================================================================##
## Makefile
## Created: Wed Aug 05 14:35:14 PDT 2015 @941 /Internet Time/
# :mode=makefile:tabSize=3:indentSize=3:
## Purpose: 
##======================================================================##

SHELL=/bin/bash
PROJECT_NAME = gospatial
GPATH = $(shell pwd)

.PHONY: fmt install get scrape build clean 

install: fmt
	@GOPATH=${GPATH} go install ${PROJECT_NAME}/main/${PROJECT_NAME}
	@GOPATH=${GPATH} go install ${PROJECT_NAME}/main/gospatial_loader
	@GOPATH=${GPATH} go install ${PROJECT_NAME}/main/gospatial_apikey
	@GOPATH=${GPATH} go install ${PROJECT_NAME}/main/gospatial_backup
	g++ -o setup src/settings/setup.cpp

fmt:
	@GOPATH=${GPATH} gofmt -s -w src/${PROJECT_NAME}

get:
	@GOPATH=${GPATH} go get ${OPTS} ${ARGS}

requirements:
	if [ ! -d "`pwd`/src/github.com/gorilla/mux" ]; then
		echo "installing mux..."
		@GOPATH=${GPATH} go get github.com/gorilla/mux
	fi

	if [ ! -d "`pwd`/src/github.com/boltdb/bolt" ]; then
		echo "installing bolt..."
		@GOPATH=${GPATH} go get github.com/boltdb/bolt
	fi

	if [ ! -d "`pwd`/src/github.com/gorilla/websocket" ]; then
		echo "installing websocket..."
		@GOPATH=${GPATH} go get github.com/gorilla/websocket
	fi

	if [ ! -d "`pwd`/src/github.com/paulmach/go.geojson" ]; then
		echo "installing go.geojson..."
		@GOPATH=${GPATH} go get github.com/paulmach/go.geojson
	fi

scrape:
	@find src -type d -name '.hg' -or -type d -name '.git' | xargs rm -rf

clean:
	@GOPATH=${GPATH} go clean
	rm bin/* && rm setup
	rm *.log && rm *.json
	rm src/gospatial/app/*.log && rm src/gospatial/app/*.db
