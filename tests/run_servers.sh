#!/usr/bin/env bash
PORT1=$1; DB1=$2; PORT2=$3; DB2=$4
cleanup() {
    kill $PID1 $PID2 2>/dev/null
    exit
}
trap cleanup SIGTERM SIGINT
LISTEN_PORT=$PORT1 DEV_DB_URL="postgres://borg:borg@localhost:5432/$DB1" make dev-backend > server1.log 2>&1 &
PID1=$!
echo $PORT
LISTEN_PORT=$PORT2 DEV_DB_URL="postgres://borg:borg@localhost:5432/$DB2" make dev-backend > server2.log 2>&1 &
PID2=$!
echo "$PID1,$PID2"
sleep infinity
