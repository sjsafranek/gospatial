#!/bin/bash

echo "cleaning files..."

# remove backup files
rm *.json

# remove log files
rm *.log
rm src/gospatial/app/*.log

# remove binaries
cd bin
rm *

# remove test datbase
rm src/gospatial/app/*.db

echo "done"