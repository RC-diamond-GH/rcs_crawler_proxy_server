FROM golang:1.20-alpine AS builder

WORKDIR /rcsproxy

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ ./src/

RUN go build -o proxy ./src/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/proxy .

COPY src/config.json .

EXPOSE 8080

CMD ["./rcs_crawler_proxy_server"]