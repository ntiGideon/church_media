# Use the official Golang image
FROM golang:1.23.4

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project directory
COPY . .

# Build the Go app from ./cmd/web/
RUN go build -o app ./cmd/web/

# Railway sets the PORT env variable automatically — expose it
EXPOSE 8080

# Start the app (Railway will pass in PORT)
CMD ["./app"]
