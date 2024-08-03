# Start with a base image containing Go and the necessary dependencies
FROM golang:1.16-alpine

LABEL maintainer="shenghuo2"

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download all the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build .

# Expose the port the app runs on (If you need)
# EXPOSE 8000

# Command to run the executable
CMD ["./sleep-status", "--port=8000", "--host=0.0.0.0"]
