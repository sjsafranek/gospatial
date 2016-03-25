#!/bin/bash

echo "Setting workspace"
export GOPATH="`pwd`"

cd src/gospatial/app
go test -bench=. -test.benchmem

