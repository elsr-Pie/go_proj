# Stage 1: Build the Go application using the official Golang image
FROM golang:1.19 as builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build the app
COPY . .
RUN go build -o main .

# Stage 2: Use a minimal base image to run the compiled binary
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]