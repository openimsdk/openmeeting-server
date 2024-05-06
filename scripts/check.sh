#!/usr/bin/env bash

PORTS="7880"

for PORT in $PORTS; do
    PID=$(lsof -i :"$PORT" -sTCP:LISTEN -t)
    if [ -z "$PID" ]; then
        echo "port $PORT is not open"
    else
        echo "port $PORT is open with process PID $PID"
    fi
done

