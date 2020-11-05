#!/bin/env bash

set -e

log() {
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*"
}

trap 'log "Execution finished"' EXIT

log 'Running pkger...'
pkger
log 'Pkger executed correctly'

log 'Running "go run" command...'
go run .
