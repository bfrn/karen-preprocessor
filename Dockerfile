# Build stage
FROM golang:1.20.5 AS builder

# Set the working directory in the builder container
WORKDIR /app

# Copy go.mod and go.sum files into the working directory
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the working directory
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go 

# Run stage
FROM alpine:latest

# Set the working directory in the run container
WORKDIR /app

# Copy the binary from the builder stage to the run stage
COPY --from=builder /app/main /app/main

# Run the app
Entrypoint ["./main"]
