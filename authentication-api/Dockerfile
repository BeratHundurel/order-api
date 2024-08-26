# Use a base image for Go
FROM golang:1.22.2 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY authentication-api/go.mod authentication-api/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY authentication-api ./

# Build the Go app with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -o auth-service .

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/auth-service .

# Expose port
EXPOSE 8081

# Command to run the executable
CMD ["./auth-service"]
