FROM golang:1.22-alpine as Build 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=Build /app/main .

RUN apk update && apk add --no-cache bash 

COPY wait-for-it.sh .

RUN chmod +x wait-for-it.sh

CMD [ "./wait-for-it.sh", "postgresql_db:5432","--","./main" ]

