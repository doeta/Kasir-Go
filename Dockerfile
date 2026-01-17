# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# -o main matches the command we want to run
RUN go build -o main main.go

# Stage 2: Create a minimal image for running the app
FROM alpine:latest

WORKDIR /app

# Install certificates (for making HTTPS requests if needed)
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Command to run (matches the built binary name)
CMD ["./main"]
