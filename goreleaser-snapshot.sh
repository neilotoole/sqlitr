#!/usr/bin/env bash

docker run --rm --privileged \
-v $(pwd):/go/src/github.com/neilotoole/sqlitr \
-v /var/run/docker.sock:/var/run/docker.sock \
-w /go/src/github.com/neilotoole/sqlitr \
neilotoole/xcgo:latest goreleaser --debug --skip-publish --snapshot --rm-dist