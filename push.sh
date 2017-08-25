#!/usr/bin/env bash
set -e
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS = -X 'main.version=$(VERSION)'

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}"
docker build -t stor.highloadcup.ru/travels/glorious_ibis .
docker push stor.highloadcup.ru/travels/glorious_ibis
