# This Dockerfile executes https://github.com/neilotoole/sqlitr.
#
# Example:
# $ docker run neilotoole/sqlitr:latest https://github.com/neilotoole/sqlitr/raw/master/testdata/example.sqlite 'SELECT * FROM actor'
# $ docker run neilotoole/sqlitr:latest /example.sqlite 'SELECT * FROM actor'
ARG BRANCH="master"
ARG VERSION="0.1.19"
ARG DEB_URL=https://github.com/neilotoole/sqlitr/releases/download/v${VERSION}/sqlitr_${VERSION}_linux_amd64.deb

FROM golang:1.14.0-buster AS sqlitr_base
ARG VERSION
#ARG DEB_URL=https://github.com/neilotoole/sqlitr/releases/download/v0.1.19/sqlitr_0.1.19_linux_amd64.deb
#ARG DEB_URL=https://github.com/neilotoole/sqlitr/releases/download/${VERSION}/sqlitr_${VERSION}_linux_amd64.deb
ARG DEB_URL
RUN echo "DEB: ${DEB_URL}"

ENV GO111MODULE=on
ENV CGO_ENABLED=1
LABEL maintainer="neilotoole@apache.org"

WORKDIR /var/cache/apt/archives
#mkdir -p /var/cache/apt/archives && cd /var/cache/apt/archives && \
RUN curl -fsSL ${DEB_URL} -O
RUN apt install -y ./sqlitr_${VERSION}_linux_amd64.deb
# curl -fsSL https://github.com/neilotoole/sqlitr/releases/download/v0.1.19/sqlitr_0.1.19_linux_amd64.deb -O

## FIXME: this always builds master
#RUN git clone https://github.com/neilotoole/sqlitr.git /go/src/github.com/neiltoole/sqlitr
#WORKDIR /go/src/github.com/neiltoole/sqlitr
#RUN git checkout "${BRANCH}"
#
#RUN go install
#
#FROM golang:1.14.0-buster AS final
#COPY --from=sqlitr_base /go/bin/sqlitr /usr/local/bin/sqlitr
## Copy the testdata/example.sqlite DB to the final image
## to make testing/examples easy
#COPY --from=sqlitr_base /go/src/github.com/neiltoole/sqlitr/testdata/example.sqlite /example.sqlite
ENTRYPOINT ["/usr/local/bin/sqlitr"]