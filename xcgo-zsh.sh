#!/usr/bin/env bash

# xcgo-zsh.sh demonstrates interactive usage of neilotoole/xcgo.
# The script starts the container with a zsh in the sqlitr dir.
#
# See README.md for more.
#
# From there you could:
#
# $ go test ./...
# $ golangci-lint run ./...
# $ go install
# $ sqlitr testdata/example.sqlite 'select * from actor'
# $ goreleaser --debug --skip-publish --snapshot --rm-dist
#
# The release packages should then be available in ./dist:
# $ ls ./dist
# build_linux_linux_amd64      checksums.txt                               sqlitr_v0.0.0-snapshot_linux_amd64.tar.gz
# build_macos_darwin_amd64     config.yaml                                 sqlitr_v0.0.0-snapshot_windows_amd64.tar.gz
# build_windows_windows_amd64  sqlitr_v0.0.0-snapshot_darwin_amd64.tar.gz
#
# And try out the newly-built binary:
# ./dist/build_linux_linux_amd64/sqlitr --version
#sqlitr v0.0.0-snapshot  2020-03-11T03:11:42Z  89c75e9322f595608dd007ede3309d475613cab0

docker run -it --privileged  \
-v $(pwd):/go/src/github.com/neilotoole/sqlitr \
-v /var/run/docker.sock:/var/run/docker.sock \
-w /go/src/github.com/neilotoole/sqlitr \
neilotoole/xcgo:latest zsh