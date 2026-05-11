FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY . .

RUN go build -mod=vendor -o mavuno-server ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/mavuno-server .

EXPOSE 8080

CMD ["./mavuno-server"]