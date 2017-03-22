#!/bin/bash

echo '{"method": "ping"}' | nc localhost 3333
exit 0