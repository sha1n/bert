FROM golang:1.19

ADD . /bert

WORKDIR /bert

RUN make go-build-linux-amd64
