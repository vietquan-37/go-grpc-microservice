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

FROM alpine:latest AS server
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server ./server
COPY ./auth-service/app.env ./
CMD ["./server"]
