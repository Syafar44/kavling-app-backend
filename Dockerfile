# Multi-stage build: compile + runtime
# Stage 1 uses golang image (large) untuk build, Stage 2 uses alpine (minimal)

# ============ STAGE 1: BUILD ============
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs (agar package docs bisa di-import)
RUN go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# ============ STAGE 2: RUNTIME ============
FROM alpine:latest

# Install ca-certificates (HTTPS), curl (healthcheck), tzdata (Asia/Makassar)
RUN apk --no-cache add ca-certificates curl tzdata

ENV TZ=Asia/Makassar

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy migrations folder
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Health check (menggunakan curl yang sudah ada)
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run application
CMD ["./server"]
