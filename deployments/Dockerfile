# Stage 1: Build the Go application
FROM --platform=linux/amd64 golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Set environment variables for a statically linked binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the Go application
RUN go build -o main ./cmd/main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the application ports
EXPOSE 7002

# Command to run the executable
CMD ["./main"]