#!/bin/bash

export GOPATH=`pwd`
go get github.com/gorilla/mux
go get github.com/boltdb/bolt
go get github.com/gorilla/websocket