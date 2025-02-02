# STEP 1 build executable binary
FROM golang:1.23 AS builder
WORKDIR /src
COPY . ./
RUN cd cmd/web-api && CGO_ENABLED=0 go build -o ice-global && mv ice-global /usr/bin

# STEP 2 build a small image
FROM alpine:3.20.5
LABEL maintainer="Mohammad Nasr <mohammadne.dev@gmail.com>"
ENV USER=root
ENTRYPOINT ["/usr/bin/ice-global"]
