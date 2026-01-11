#!/usr/bin/env bash

BORG_DB="postgres://borg:borg@localhost:5432/borg"
COPY_DB="postgres://borg:borg@localhost:5432/borg_copy"

LISTEN_PORT=8081 DEV_DB_URL=$COPY_DB make dev-backend &
LISTEN_PORT=8080 DEV_DB_URL=$BORG_DB make dev-backend &

cleanup() {
    kill $(jobs -p) 2>/dev/null
    exit
}

trap cleanup SIGINT SIGTERM
wait
