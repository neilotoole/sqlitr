# This Dockerfile executes https://github.com/neilotoole/sqlitr.
#
# Example:
# $ docker run neilotoole/sqlitr:latest https://github.com/neilotoole/sqlitr/raw/master/testdata/example.sqlite 'SELECT * FROM actor'
# $ docker run neilotoole/sqlitr:latest /example.sqlite 'SELECT * FROM actor'

FROM scratch
COPY sqlitr /
ENTRYPOINT ["/sqlitr"]

#FROM golang:1.14.0-buster AS builder
#FROM debian:buster-slim AS builder
##FROM debian:buster-slim
#LABEL maintainer="neilotoole@apache.org"
#ARG VERSION="0.1.19"
#ARG DEB_URL=https://github.com/neilotoole/sqlitr/releases/download/v${VERSION}/sqlitr_${VERSION}_linux_amd64.deb
##ARG DEB_URL=https://github.com/neilotoole/sqlitr/releases/download/v0.1.19/sqlitr_0.1.19_linux_amd64.deb
#ARG EXAMPLE_DB_URL=https://github.com/neilotoole/sqlitr/raw/v${VERSION}/testdata/example.sqlite
#
## https://github.com/neilotoole/sqlitr/releases/download/v0.1.19/sqlitr_0.1.19_linux_amd64.tar.gz
#WORKDIR /root
#ADD https://github.com/neilotoole/sqlitr/releases/download/v0.1.19/sqlitr_0.1.19_linux_amd64.tar.gz /root/
#RUN tar xvzf *.tar.gz && cp sqlitr /usr/local/bin/
##ADD sqlitr /usr/local/bin
##ADD ./dist/build_linux_linux_amd64/sqlitr /usr/local/bin
##RUN
#
#ADD ${EXAMPLE_DB_URL} /example.sqlite
##WORKDIR /var/cache/apt/archives
##RUN curl -fsSL ${DEB_URL} -O
##RUN apt install -y ./sqlitr_${VERSION}_linux_amd64.deb
#
#WORKDIR /
##ENTRYPOINT ["/usr/local/bin/sqlitr"]
#
#FROM scratch AS final
##FROM scratch AS final
#COPY --from=builder /usr/local/bin/sqlitr /usr/local/bin/sqlitr
#COPY --from=builder /example.sqlite /example.sqlite
#ENTRYPOINT ["/usr/local/bin/sqlitr"]






##FROM golang:1.14.0-buster AS builder
#FROM debian:buster-slim AS builder
##FROM debian:buster-slim
#LABEL maintainer="neilotoole@apache.org"
#ARG VERSION="0.1.19"
#ARG DEB_URL=https://github.com/neilotoole/sqlitr/releases/download/v${VERSION}/sqlitr_${VERSION}_linux_amd64.deb
##ARG DEB_URL=https://github.com/neilotoole/sqlitr/releases/download/v0.1.19/sqlitr_0.1.19_linux_amd64.deb
#ARG EXAMPLE_DB_URL=https://github.com/neilotoole/sqlitr/raw/v${VERSION}/testdata/example.sqlite
#
## https://github.com/neilotoole/sqlitr/releases/download/v0.1.19/sqlitr_0.1.19_linux_amd64.tar.gz
#WORKDIR /root
#ADD https://github.com/neilotoole/sqlitr/releases/download/v0.1.19/sqlitr_0.1.19_linux_amd64.tar.gz /root/
#RUN tar xvzf *.tar.gz && cp sqlitr /usr/local/bin/
##ADD sqlitr /usr/local/bin
##ADD ./dist/build_linux_linux_amd64/sqlitr /usr/local/bin
##RUN
#
#ADD ${EXAMPLE_DB_URL} /example.sqlite
##WORKDIR /var/cache/apt/archives
##RUN curl -fsSL ${DEB_URL} -O
##RUN apt install -y ./sqlitr_${VERSION}_linux_amd64.deb
#
#WORKDIR /
##ENTRYPOINT ["/usr/local/bin/sqlitr"]
#
#FROM scratch AS final
##FROM scratch AS final
#COPY --from=builder /usr/local/bin/sqlitr /usr/local/bin/sqlitr
#COPY --from=builder /example.sqlite /example.sqlite
#ENTRYPOINT ["/usr/local/bin/sqlitr"]