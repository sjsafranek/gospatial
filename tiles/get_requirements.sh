#!/bin/bash

echo "Setting workspace"
export GOPATH="`pwd`"

# echo "Setting go version"
# gvm use go1.5.2

echo "Installing requirements"

if [ ! -d "`pwd`/src/github.com/mattn/go-sqlite3" ]; then
	echo "installing go-sqlite3..."
	go get github.com/mattn/go-sqlite3
fi

if [ ! -d "`pwd`/src/github.com/fawick/go-mapnik/mapnik" ]; then
	echo "installing go-mapnik..."
	go get -d github.com/fawick/go-mapnik/mapnik
	./src/github.com/fawick/go-mapnik/mapnik/configure.bash
fi

