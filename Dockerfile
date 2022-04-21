FROM golang:1.18

ADD . /bert

WORKDIR /bert

RUN make go-build-linux-amd64
