#!/usr/bin/env bash
set -e

GOOS=linux GOARCH=amd64 go build
docker build -t stor.highloadcup.ru/travels/glorious_ibis .
docker push stor.highloadcup.ru/travels/glorious_ibis
