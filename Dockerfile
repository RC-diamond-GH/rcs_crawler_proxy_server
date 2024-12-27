FROM golang:1.21-alpine AS builder

WORKDIR /rcs_crawler_proxy_server

COPY src/ ./

RUN go mod tidy

RUN go build -o proxy ./main.go

FROM alpine:latest

COPY --from=builder /rcs_crawler_proxy_server/proxy /rcs_crawler_proxy_server

CMD ["/proxy"]
