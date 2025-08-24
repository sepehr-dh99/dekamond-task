# Stage 1: Build
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Install build dependencies (optional: git if needed for go modules)
RUN apk add --no-cache git

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o dekamond-task main.go

# Stage 2: Run
FROM alpine:3.20

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/dekamond-task .

# Expose application port
EXPOSE 8080

# Run the application
CMD ["./dekamond-task"]
