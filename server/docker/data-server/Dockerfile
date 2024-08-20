FROM golang:alpine3.20 AS builder

WORKDIR /app

# Copy the Go module files
COPY ./data-server/go.mod ./
#COPY ./go.sum ./    

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY ./data-server .

# Build the Go application
RUN go build -o data-server .

FROM alpine:3.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/data-server /app/data-server

# Command to run the application
CMD ["./data-server"]