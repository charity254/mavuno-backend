FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git — required for go mod download to fetch some dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o mavuno-server ./cmd/server

FROM alpine:latest

# Install ca-certificates — required for HTTPS connections to Neon database
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/mavuno-server .

EXPOSE 8080

CMD ["./mavuno-server"]