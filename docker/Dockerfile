# Start from the official golang image
FROM golang:1.21.6

# Set the working directory in the container
WORKDIR /app

# Copy go mod and sum files
COPY ./ ./

# Download all dependencies
RUN go mod tidy

# Copy the source code into the container
# COPY . ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
