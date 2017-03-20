#!/bin/bash

export GOPATH="`pwd`"

echo "creating workspace..."

# Setup working directory
echo "creating directories..."
if [ ! -d "`pwd`/bin" ]; then
    mkdir bin
fi
if [ ! -d "`pwd`/pkg" ]; then
    mkdir pkg
fi
if [ ! -d "`pwd`/log" ]; then
    mkdir log
fi
if [ ! -d "`pwd`/src" ]; then
    mkdir src
    #mkdir src/gospatial
fi

# Move source files
#echo "copying source files..."
#cp -R gospatial/* src/gospatial/

# Download required libraries
echo "checking requirements..."
if [ ! -d "`pwd`/src/github.com/gorilla/mux" ]; then
    echo "downloading gorilla mux..."
    go get github.com/gorilla/mux
fi

if [ ! -d "`pwd`/src/github.com/boltdb/bolt" ]; then
    echo "downloading bolt..."
    go get github.com/boltdb/bolt
fi

if [ ! -d "`pwd`/src/github.com/gorilla/websocket" ]; then
    echo "downloading gorilla websocket..."
    go get github.com/gorilla/websocket
fi

if [ ! -d "`pwd`/src/github.com/paulmach/go.geojson" ]; then
    echo "downloading go.geojson..."
    go get github.com/paulmach/go.geojson
fi
if [ ! -d "`pwd`/src/github.com/cihub/seelog" ]; then
    echo "downloading seelog..."
    go get github.com/cihub/seelog
fi

if [ ! -d "`pwd`/src/github.com/mattn/go-sqlite3" ]; then
    echo "downloading go-sqlite3..."
    go get github.com/mattn/go-sqlite3
fi

if [ ! -d "`pwd`/src/github.com/lib/pq" ]; then
    echo "downloading pg..."
    go get github.com/lib/pq
fi

# format
./manage.sh -f

# run make install
make install

echo "done!"
