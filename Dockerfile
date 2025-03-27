FROM node:20-alpine3.19 AS webui-builder
WORKDIR /app
COPY . /app/
RUN \
  cd /app/ui && \
  npm install && \
  npm run build



FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY ../. /app/
RUN rm -rf /app/ui/dist
COPY --from=webui-builder /app/ui/dist /app/ui/dist
RUN go build -o certimate



FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/certimate .
ENTRYPOINT ["./certimate", "serve", "--http", "0.0.0.0:8090"]
