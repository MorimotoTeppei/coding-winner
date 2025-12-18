# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot ./cmd/bot

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set timezone to JST
ENV TZ=Asia/Tokyo

# Create app user
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

# Set working directory
WORKDIR /home/app

# Copy binary from builder
COPY --from=builder --chown=app:app /app/bot .
COPY --from=builder --chown=app:app /app/migrations ./migrations

# Switch to app user
USER app

# Run the application
CMD ["./bot"]
