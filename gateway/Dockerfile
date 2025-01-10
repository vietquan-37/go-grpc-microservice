FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/server ./cmd

FROM alpine:latest AS server
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server ./server
COPY ./app.env ./
CMD ["./server"]