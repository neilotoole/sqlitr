# This Dockerfile builds an image whose entrypoint
# executes https://github.com/neilotoole/sqlitr.
# A sample database is included at /example.sqlite
#
# Example:
# $ docker run neilotoole/sqlitr:latest https://github.com/neilotoole/sqlitr/raw/master/testdata/example.sqlite 'SELECT * FROM actor'
# $ docker run neilotoole/sqlitr:latest /example.sqlite 'SELECT * FROM actor'
FROM scratch
LABEL maintainer="neilotoole@apache.org"

# This Dockerfile is intended to be built by goreleaser, which
# sets up a temporary folder that can be copied from.
COPY sqlitr /usr/local/bin/
COPY LICENSE /
COPY README.md /
COPY testdata/example.sqlite /example.sqlite

WORKDIR /
ENTRYPOINT ["/usr/local/bin/sqlitr"]