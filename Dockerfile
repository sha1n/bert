FROM golang:1.16

ADD . /benchy

WORKDIR /benchy

RUN make go-build-linux-amd64
