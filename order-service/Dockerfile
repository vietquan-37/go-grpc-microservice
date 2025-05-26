# Build stage
FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY ./common ./common
COPY ./order-service/go.work ./order-service/go.work.sum ./order-service/
WORKDIR /app/order-service
COPY ./order-service/go.mod ./order-service/go.sum ./
RUN go mod download
COPY ./order-service ./
RUN go build -o /bin/server ./cmd
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/latest/download/grpc_health_probe-linux-amd64 \
    /usr/local/bin/grpc_health_probe
RUN chmod +x /usr/local/bin/grpc_health_probe
FROM alpine:latest AS server
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server ./server
COPY --from=builder /usr/local/bin/grpc_health_probe /usr/local/bin/grpc_health_probe
COPY ./order-service/app.env ./
CMD ["./server"]
