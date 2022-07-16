FROM golang:1.18.4-alpine3.16 AS builder

WORKDIR /usr/local/go/src/

ADD ./ /usr/local/go/src/

RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/urlapi/main.go

FROM alpine:latest

COPY --from=builder /usr/local/go/src/app /
COPY --from=builder /usr/local/go/src/config.yml /

EXPOSE 8080
CMD ["/app", "-config-source=env", "-store=local"]