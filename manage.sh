#!/bin/bash

export GOPATH="`pwd`"

key="$1"

case $key in
    # -i|--install)

    #     # Setup working directory
    #     echo "creating directories..."
    #     if [ ! -d "`pwd`/bin" ]; then
    #         mkdir bin
    #     fi
    #     if [ ! -d "`pwd`/pkg" ]; then
    #         mkdir pkg
    #     fi
    #     if [ ! -d "`pwd`/log" ]; then
    #         mkdir log
    #     fi
    #     if [ ! -d "`pwd`/src" ]; then
    #         mkdir src
    #         mkdir src/gospatial
    #     fi

    #     # Move source files
    #     cp -R gospatial/* src/gospatial/

    #     # Download required libraries
    #     echo "checking requirements..."
    #     if [ ! -d "`pwd`/src/github.com/gorilla/mux" ]; then
    #         echo "installing mux..."
    #         go get github.com/gorilla/mux
    #     fi

    #     if [ ! -d "`pwd`/src/github.com/boltdb/bolt" ]; then
    #         echo "installing bolt..."
    #         go get github.com/boltdb/bolt
    #     fi

    #     if [ ! -d "`pwd`/src/github.com/gorilla/websocket" ]; then
    #         echo "installing websocket..."
    #         go get github.com/gorilla/websocket
    #     fi

    #     if [ ! -d "`pwd`/src/github.com/paulmach/go.geojson" ]; then
    #         echo "installing go.geojson..."
    #         go get github.com/paulmach/go.geojson
    #     fi
    #     echo "done!"
    # ;;
    -c|--clean)
        echo "cleaning working directory..."
        rm bin/* && rm setup || echo "no binaries found"
        rm log/*.log || echo "no log files found"
        # rm *.json || echo "no json files found"
        rm src/gospatial/app/*.log || echo "no testing log found"
        rm src/gospatial/app/*.db || echo "no testing db found"
        echo "done"
    ;;
    # -d|--destroy)
    #     echo "destroy woring directory"
    #     rm -rf bin
    #     rm -rf src
    #     rm -rf log
    #     rm -rf pkg
    #     echo "done"
    # ;;
    -t| --test)
        cd src/gospatial/app
        go test -bench=. -test.benchmem
        # cd ../../..
        # cd tests
        # python3 api_tester.py        
    ;;
    -h| --help)
# HELP STUFF
        echo "-c -t -h -r"
        # cd ../../..
        # cd tests
        # python3 api_tester.py        
    ;;
    *)
        echo "unknown option"
    ;;
esac









