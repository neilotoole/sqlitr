# This Dockerfile executes https://github.com/neilotoole/sqlitr.
#
# Example:
# $ docker run neilotoole/sqlitr:latest https://github.com/neilotoole/sqlitr/raw/master/testdata/example.sqlite 'SELECT * FROM actor'
# $ docker run neilotoole/sqlitr:latest /example.sqlite 'SELECT * FROM actor'

FROM golang:1.14.0-buster AS sqlitr_base
LABEL maintainer="neilotoole@apache.org"
WORKDIR /go/src/github.com/neiltoole
ENV GO111MODULE=on
ENV CGO_ENABLED=1
RUN git clone https://github.com/neilotoole/sqlitr.git
WORKDIR /go/src/github.com/neiltoole/sqlitr
RUN go install

FROM golang:1.14.0-buster AS final
COPY --from=sqlitr_base /go/bin/sqlitr /usr/local/bin/sqlitr
# Copy the testdata/example.sqlite DB to the final image
# to make testing/examples easy
COPY --from=sqlitr_base /go/src/github.com/neiltoole/sqlitr/testdata/example.sqlite /example.sqlite
ENTRYPOINT ["/usr/local/bin/sqlitr"]