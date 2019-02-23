#!/bin/bash

export TTROOT="task_tracker"

chmod 755 -R "${TTROOT}/webroot"

screen -S tt_api -X quit
echo "starting client"
screen -S tt_api -d -m bash -c "cd ${TTROOT} && ./tt_api"
sleep 1
screen -list
