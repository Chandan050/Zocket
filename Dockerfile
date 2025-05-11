# Use official Go image
FROM golang:1.23.0

# Set working directory
WORKDIR /app

# Copy the application files
COPY . .

# Download dependencies
RUN go mod tidy

# Build the application
RUN go build -o kvstore

# Expose the port
EXPOSE 8080

# Set environment variable for node name (default to node1, can be overridden)
ENV NODE_NAME=node1

# Command to run the application
CMD ["./kvstore"]
