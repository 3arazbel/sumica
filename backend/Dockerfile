# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod if it exists, or create a new one
COPY go.mod* ./
RUN if [ ! -f go.mod ]; then go mod init sumica; fi

# Download dependencies (if any)
RUN go mod tidy

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
