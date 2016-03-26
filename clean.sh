#!/bin/bash

echo "cleaning files..."

# remove backup files
rm *.json

# remove binaries
cd bin
rm *

echo "done"