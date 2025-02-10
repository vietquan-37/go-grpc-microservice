# Build stage
FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY ./common ./common
COPY ./gateway/go.work ./gateway/go.work.sum ./gateway/
WORKDIR /app/gateway
COPY ./gateway/go.mod ./gateway/go.sum ./
RUN go mod download
COPY ./gateway ./
RUN go build -o /bin/server ./cmd

FROM alpine:latest AS server
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server ./server
COPY ./gateway/app.env ./
CMD ["./server"]
