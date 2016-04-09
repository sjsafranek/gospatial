#!/bin/bash

echo "cleaning files..."

# remove backup files
rm *.json

# remove log files
rm *.log

# remove binaries
cd bin
rm *

echo "done"