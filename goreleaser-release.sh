#!/usr/bin/env bash

docker run --rm --privileged \
-v $(pwd):/go/src/github.com/neilotoole/sqlitr \
-v /var/run/docker.sock:/var/run/docker.sock \
-e "GITHUB_TOKEN=$GITHUB_TOKEN" \
-e "DOCKER_USERNAME=$DOCKER_USERNAME" -e "DOCKER_PASSWORD=$DOCKER_PASSWORD" -e "DOCKER_REGISTRY=$DOCKER_REGISTRY" \
-v "${HOME}/.snapcraft.login":/.snapcraft.login \
-w /go/src/github.com/neilotoole/sqlitr \
neilotoole/xcgo:latest goreleaser --rm-dist