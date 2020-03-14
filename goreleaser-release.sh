#!/usr/bin/env bash

docker run --rm --privileged \
-v $(pwd):/go/src/github.com/neilotoole/sqlitr \
-v /var/run/docker.sock:/var/run/docker.sock \
-e "GITHUB_TOKEN=$GITHUB_TOKEN" \
-e "SNAPCRAFT_LOGIN_FILE=/.snapcraft.login" \
-v "${HOME}/.snapcraft.login":/.snapcraft.login \
-w /go/src/github.com/neilotoole/sqlitr \
neilotoole/xcgo:latest goreleaser --rm-dist