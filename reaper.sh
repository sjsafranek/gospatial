#!/bin/bash
while :
do
        PROCESS_FOUND=`ps -ef|grep ml.internalpositioning|grep -v grep`
        if [ "$PROCESS_FOUND" = "" ]
        then
                echo "$(date)" > reaper.log
                /etc/init.d/ml.internalpositioning.init restart
        fi
        sleep 60
done
