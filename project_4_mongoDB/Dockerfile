FROM golang:1.22-alpine as build 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest 

RUN apk update && apk add --no-cache bash
COPY wait-for-it.sh .
RUN chmod +x wait-for-it.sh

COPY --from=build /app/main .

EXPOSE 9010

CMD ["./wait-for-it.sh", "mongodb:27017", "--", "./main"]