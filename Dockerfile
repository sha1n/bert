# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install make and git (required for Makefile and versioning)
RUN apk add --no-cache make git

WORKDIR /bert
COPY . .

# Build the linux amd64 binary
RUN make go-build-linux-amd64

# Stage 2: Runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /bert
COPY --from=builder /bert/bin/bert-linux-amd64 /bert/bin/bert-linux-amd64
COPY --from=builder /bert/test /bert/test

# Set the binary as the entrypoint
ENTRYPOINT ["/bert/bin/bert-linux-amd64"]
