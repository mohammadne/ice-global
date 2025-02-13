# STEP 1 build executable binary
FROM golang:1.23 AS builder
WORKDIR /src
COPY . ./
RUN cd cmd/web-api && CGO_ENABLED=0 go build -o shopping-cart-manager && mv shopping-cart-manager /usr/bin

# STEP 2 build a small image
FROM alpine:3.20
LABEL maintainer="Mohammad Nasr <mohammadne.dev@gmail.com>"
RUN apk add --no-cache bind-tools busybox-extras
COPY --from=builder /usr/bin/shopping-cart-manager /usr/bin/shopping-cart-manager
ENV USER=root
ENTRYPOINT ["/usr/bin/shopping-cart-manager"]
