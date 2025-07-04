# Multi-stage build for optimized image size
# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy source code
COPY mandelbrot_ascii.go .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mandelbrot-server mandelbrot_ascii.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests (if needed)
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/mandelbrot-server .

# Create directory for gallery files and set ownership
RUN mkdir -p /app/data && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port 3000
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Default command - run server on port 3000
CMD ["./mandelbrot-server", "-server", "-port=:3000"]