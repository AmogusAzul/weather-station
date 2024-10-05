FROM golang:alpine3.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the all the files for data-server
COPY ./server/data-server ./server/data-server
COPY ./common ./common



# Install required packages for Cgo
RUN apk add --no-cache gcc musl-dev

# Enable CGO explicitly
ENV CGO_ENABLED=1



# Download Go module dependencies
RUN cd ./server/data-server && go mod download

# Build the Go application
WORKDIR /app/server/data-server
RUN go build -o /app/data-server



# Final stage to create a lightweight container
FROM alpine:3.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/data-server /app/data-server



# Command to run the application
CMD ["./data-server"]
