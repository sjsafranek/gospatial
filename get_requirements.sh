#!/bin/bash

echo "Setting workspace"
export GOPATH="`pwd`"

# echo "Setting go version"
# gvm use go1.5.2

echo "Installing requirements"
if [ ! -d "`pwd`/src/github.com/gorilla/mux" ]; then
	go get github.com/gorilla/mux
fi

if [ ! -d "`pwd`/src/github.com/boltdb/bolt" ]; then
	go get github.com/boltdb/bolt
fi

if [ ! -d "`pwd`/src/github.com/gorilla/websocket" ]; then
	go get github.com/gorilla/websocket
fi



if [ ! -d "`pwd`/src/github.com/paulmach/go.geojson" ]; then
	go get github.com/paulmach/go.geojson
fi

