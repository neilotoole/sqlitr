#!/usr/bin/env bash

# xcgo-zsh.sh runs neilotoole/xcgo, starting the container
# with a zsh in the sqlitr dir. From there you could
#
# $ go install
# $ sqlitr testdata/example.sqlite 'select * from actor'

docker run -it --privileged  \
-v $(pwd):/go/src/github.com/neilotoole/sqlitr \
-v /var/run/docker.sock:/var/run/docker.sock \
-w /go/src/github.com/neilotoole/sqlitr \
neilotoole/xcgo:latest zsh