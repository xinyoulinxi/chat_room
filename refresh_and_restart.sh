#!/bin/bash

# find web server process
pid=$(ps -ef | grep 'web_server' | grep -v 'grep' | awk '{print $2}')

# check if pid is not empty
if [ -n "$pid" ]; then
    echo "find web server process，PID: $pid"
    echo "killing process ..."
    kill -9 $pid
    echo "process has killed."
else
    echo "not find web server process."
fi
echo "pull git repo..."
git pull -r
echo "update git repo over."
echo "start compile..."
./build_and_run_backend.sh
echo "compile over。web server running in backend."
