FROM golang:1.14 as builder
WORKDIR /go/src/github.com/traefik/whoamiudp
COPY . .
RUN CGO_ENABLED=0 go build ./

# Create a minimal container to run a Golang static binary
FROM scratch
COPY --from=builder /go/src/github.com/traefik/whoamiudp/whoamiudp .
ENTRYPOINT ["/whoamiudp"]
EXPOSE 8080/udp
