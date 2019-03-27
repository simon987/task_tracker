#!/bin/bash

export TTROOT="task_tracker"

chmod 755 -R "${TTROOT}/webroot"

screen -S tt_api -X quit
echo "starting client"
screen -S tt_api -d -m bash -c "cd ${TTROOT} && chmod +x tt_api && ./tt_api 2> stderr.txt"
sleep 1
screen -list
