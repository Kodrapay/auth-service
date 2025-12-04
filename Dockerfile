# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk --no-cache add ca-certificates git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o auth-service ./cmd/auth-service

# Runtime stage - use distroless for smaller image
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
COPY --from=builder /app/auth-service .
EXPOSE 7001
CMD ["./auth-service"]
