# Build stage
FROM golang:1.25-bookworm AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o web ./cmd/web

# Final stage
FROM alpine:3.21

WORKDIR /app

# Install CA certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -s /sbin/nologin appuser

# Copy the binary from builder
COPY --from=builder /app/web .

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./web"]
