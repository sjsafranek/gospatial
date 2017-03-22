#!/bin/bash

export GOPATH="`pwd`"

cd gospatial
go test -bench=. -test.benchmem
