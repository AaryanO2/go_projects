FROM golang:alpine as Build 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a statically-linked binary
RUN CGO_ENABLED=0 go build -o main main.go

FROM alpine

RUN adduser -S -D -h /app appuser

USER appuser

WORKDIR /app

COPY --from=Build /app/main .

EXPOSE 3000

CMD [ "/app/main" ]