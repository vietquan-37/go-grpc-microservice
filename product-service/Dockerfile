# Build stage
FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY ./common ./common
COPY ./product-service/go.work ./product-service/go.work.sum ./product-service/
WORKDIR /app/product-service
COPY ./product-service/go.mod ./product-service/go.sum ./
RUN go mod download
COPY ./product-service ./
RUN go build -o /bin/server ./cmd

FROM alpine:latest AS server
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server ./server
COPY ./product-service/app.env ./
CMD ["./server"]
