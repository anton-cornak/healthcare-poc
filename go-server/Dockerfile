# Dockerfile

# Start from the latest Node base image
FROM golang:1.21.1-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Install all dependencies. Dependencies will be cached if the package files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY .envrc ./.envrc
COPY cmd/ ./cmd
COPY docs/ ./docs
COPY handlers/ ./handlers
COPY models/ ./models

# Build the Go app
RUN go build -o main ./cmd

# Run tests
RUN go test ./...

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to start the server
CMD ["./main"]
