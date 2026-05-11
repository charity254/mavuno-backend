# Use the official Go image as the base
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN go build -o mavuno-server ./cmd/server

# Use a minimal image for the final container
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/mavuno-server .

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./mavuno-server"]