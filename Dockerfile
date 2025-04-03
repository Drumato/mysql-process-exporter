# Build stage
FROM golang:1.24.0 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY ./main.go ./
COPY ./metrics ./metrics

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# Final stage
FROM gcr.io/distroless/static:nonroot

# Set the working directory
WORKDIR /

# Copy the built binary
COPY --from=builder /app/server /server

# Expose the application port
EXPOSE 8080

# Set the default command
CMD ["/server"]