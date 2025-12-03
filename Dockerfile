FROM golang:1.24

ADD . /bert

WORKDIR /bert

RUN make go-build-linux-amd64
