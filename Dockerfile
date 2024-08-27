FROM node:20-alpine

WORKDIR /app

COPY ../. /app/

RUN npm --prefix=./ui ci && npm --prefix=./ui run build

FROM golang:1.22-alpine

WORKDIR /app

COPY ../. /app/

RUN go build -o certimate

ENTRYPOINT ["./certimate", "serve", "--http", "0.0.0.0:8090"]