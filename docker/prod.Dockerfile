# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set the environment variable for the Go application
ENV PORT=8080

# Install tzdata package for timezone configuration
RUN apk add --no-cache tzdata

# Set the timezone to Asia/Bangkok
ENV TZ=Asia/Bangkok

# Set the current working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files for dependency resolution
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o task-nexus-task-management .

# Stage 2: Create the final image
FROM alpine:latest

# Install necessary packages for running Go apps (if needed)
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/task-nexus-task-management .

# Copy the timezone information from the builder stage
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip

# Set the PORT environment variable
EXPOSE $PORT

# Command to run the executable
CMD ["./task-nexus-task-management"]

