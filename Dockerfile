FROM golang:1.8-alpine AS builder
WORKDIR /go/src/github.com/djmaze/swarmdns
COPY . .
RUN CGO_ENABLED=0 go build --ldflags "-s"

FROM scratch
EXPOSE 53/udp
COPY --from=builder /go/src/github.com/djmaze/swarmdns/swarmdns /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/swarmdns"]
