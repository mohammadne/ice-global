# STEP 1 build executable binary
FROM golang:1.23 AS builder

WORKDIR /src

COPY . ./

RUN cd cmd/web-api && CGO_ENABLED=0 go build -o ice-global && mv ice-global /usr/bin

# STEP 2 build a small image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/bin/ice-global /usr/bin/ice-global

ENV USER=root

ENTRYPOINT ["/usr/bin/ice-global"]
