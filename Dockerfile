FROM golang:1.16

ADD . /bert

WORKDIR /bert

RUN make go-build-linux-amd64
