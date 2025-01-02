#FROM golang:1.23.1-alpine AS builder
#WORKDIR /app
#RUN printf "go 1.23.1\n\nuse (\n\t./common\n\t./product-service\n)" > go.work
#COPY common ./common
#COPY ./product-service/go.mod ./product-service/go.sum ./product-service/
#RUN go mod download
#COPY ./product-service ./product-service
#RUN go build -o /bin/server ./product-service/cmd
#
#FROM alpine:latest AS server
#WORKDIR /app
#RUN apk add --no-cache ca-certificates
#COPY --from=builder /bin/server ./server
#COPY ./product-service/app.env ./
#CMD ["./server"]

#first way above
#second way easy to understand but have to modify gowork every time build--
#command: docker build -f ./product-service/Dockerfile .
FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY go.work ./
COPY common ./common
COPY ./product-service/go.mod ./product-service/go.sum ./product-service/
RUN go mod download
COPY ./product-service ./product-service
RUN go build -o /bin/server ./product-service/cmd

FROM alpine:latest AS server
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server ./server
COPY ./product-service/app.env ./
CMD ["./server"]