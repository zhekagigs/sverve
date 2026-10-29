# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files if they exist
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy source code
COPY . .

# Build the server
RUN go build -o server .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 9000

# Run the server
CMD ["./server"]
