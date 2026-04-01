# Stage 1: Build stage
FROM golang:1.25-alpine AS builder

# Install required build tools for SQLite (CGO)
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o findit-server ./cmd/server

# Stage 2: Runtime stage
FROM alpine:latest

# Install runtime dependencies for SQLite
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/findit-server .

# Copy frontend files
COPY --from=builder /app/frontend ./frontend

# Copy database directory (will be created at runtime if needed)
RUN mkdir -p /app/data

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./findit-server"]
