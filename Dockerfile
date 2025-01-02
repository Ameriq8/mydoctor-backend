# Start with the Go image
FROM golang:1.21 AS builder

# Set up the working directory
WORKDIR /app

# Copy go modules and source
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the binary
RUN go build -o main ./cmd/app/main.go

# Debugging: Check if the binary is built
RUN ls -l /app

# Use a minimal image for the final binary
FROM alpine:latest
WORKDIR /root/

# Install necessary dependencies
RUN apk add --no-cache libc6-compat

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Ensure binary is executable
RUN chmod +x ./main

# Expose the port and set the command
EXPOSE 8080
CMD ["./main"]
