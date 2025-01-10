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

FROM alpine:latest AS server
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server ./server
COPY ./order-service/app.env ./
CMD ["./server"]
