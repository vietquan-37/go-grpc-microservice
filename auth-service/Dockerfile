# Build stage
FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY ./common ./common
COPY ./auth-service/go.work ./auth-service/go.work.sum ./auth-service/
WORKDIR /app/auth-service
COPY ./auth-service/go.mod ./auth-service/go.sum ./
RUN go mod download
COPY ./auth-service ./
RUN go build -o /bin/server ./cmd
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/latest/download/grpc_health_probe-linux-amd64 \
    /usr/local/bin/grpc_health_probe
RUN chmod +x /usr/local/bin/grpc_health_probe
FROM alpine:latest AS server
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server ./server
COPY --from=builder /usr/local/bin/grpc_health_probe /usr/local/bin/grpc_health_probe
COPY ./auth-service/app.env ./
CMD ["./server"]
