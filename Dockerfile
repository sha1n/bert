FROM golang:1.21

ADD . /bert

WORKDIR /bert

RUN make go-build-linux-amd64
