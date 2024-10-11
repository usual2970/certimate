FROM golang:1.22-alpine as builder

WORKDIR /app

COPY ../. /app/

RUN go build -o certimate


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/certimate .

ENTRYPOINT ["./certimate", "serve", "--http", "0.0.0.0:8090"]