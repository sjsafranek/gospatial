#!/bin/bash

export GOPATH="`pwd`"

key="$1"

case $key in
    -c|--clean)
        echo "cleaning working directory..."
        rm bin/* && rm setup || echo "no binaries found"
        rm log/*.log || echo "no log files found"
        # rm *.json || echo "no json files found"
        rm src/gospatial/app/*.log || echo "no testing log found"
        rm src/gospatial/app/*.db || echo "no testing db found"
        echo "done"
    ;;
    -d|--destroy)
        echo "destroy woring directory"
        rm -rf bin
        rm -rf src
        rm -rf log
        rm -rf pkg
        echo "done"
    ;;
    -t| --test)
        cd src/gospatial/app
        go test -bench=. -test.benchmem
        # cd ../../..
        # cd tests
        # python3 api_tester.py        
    ;;
    -h| --help)
        echo "-c -d -t -h"        
    ;;
    *)
        echo "unknown option"
    ;;
esac









