FROM golang:1.21-alpine AS builder

WORKDIR /rcs_crawler_proxy_server

COPY src/ ./

RUN CGO_ENABLED=0 go mod tidy && CGO_ENABLED=0 go build -o MYAPP ./main.go

FROM alpine:latest

WORKDIR /rcs_crawler_proxy_server

COPY --from=builder /rcs_crawler_proxy_server/MYAPP /rcs_crawler_proxy_server/MYAPP

CMD ["./MYAPP"]