FROM golang:alpine3.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the all the files for data-server
COPY ./server/data-server/ /app/

# Download Go module dependencies
RUN go mod download

# Run the Go build for the binary
RUN go build -o /app/data-server /app/

# Final stage to create a lightweight container
FROM alpine:3.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/data-server /app/data-server

# Use Argument for creating TOKEN_PATH and then coping value to env variable
ARG TOKEN_PATH
COPY ./server/docker/data-server/token-file-init.sh  .
RUN sh token-file-init.sh $TOKEN_PATH
ENV TOKEN_PATH=${TOKEN_PATH}

# Command to run the application
CMD ["./data-server"]
# CMD ["sleep", "600"]
