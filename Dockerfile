FROM node:20-alpine as build-ui

WORKDIR /app

COPY ../. /app/

RUN npm --prefix=./ui ci && npm --prefix=./ui run build

FROM golang:1.22-alpine

WORKDIR /app

COPY ../. /app/
COPY --from=build-ui /app/ui/dist /app/ui/

RUN go build -o certimate

ENTRYPOINT ["./certimate", "serve", "--http", "0.0.0.0:8090"]