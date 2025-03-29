# Use the official Golang 1.24.1 image as base
FROM golang:1.24.1 

# Set the working directory inside the container
WORKDIR /usr/src/app

# Installing Air live-reloading command line utility 
RUN go install github.com/air-verse/air@latest

# Copy all files from current directory to working directory
COPY . .

# Download and install any required dependencies
RUN go mod tidy 