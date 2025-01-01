# Start with the Go image
FROM golang:1.20 AS builder

# Set up the working directory
WORKDIR /app

# Copy go modules and source
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the binary
RUN go build -o app cmd/app/main.go

# Use a minimal image for the final binary
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .

# Expose the port and set the command
EXPOSE 8080
CMD ["./app"]
