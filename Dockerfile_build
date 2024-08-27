FROM golang:1.22-alpine

ENV GOPROXY=https://goproxy.cn

WORKDIR /app

COPY ../. /app/

RUN go build -o certimate

ENTRYPOINT ["./certimate", "serve", "--http", "0.0.0.0:8090"]