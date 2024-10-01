# Use the official Go image.
FROM golang:1.17

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go module files.
COPY go.mod go.sum ./

# Download the Go dependencies.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the Go application.
RUN go build -o metrics_exporter ./cmd/metrics_exporter

# Expose the port the application runs on.
EXPOSE 8080

# Command to run the application.
CMD ["./metrics_exporter"]