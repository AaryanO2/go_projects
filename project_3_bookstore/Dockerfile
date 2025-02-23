FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:latest

RUN apk update && apk add --no-cache bash

WORKDIR /root/

COPY --from=build /app/main .

COPY wait-for-it.sh .

EXPOSE 9010

CMD ["./wait-for-it.sh", "mysql:3306", "--", "./main"]
