# Start from a lightweight Go base image
FROM golang:1.21-alpine as builder

# Set up working directory for the backend build
WORKDIR /app

# Copy Go module files and build the backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/*.go ./

# Build the Go backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

# Download PocketBase and make it executable
RUN wget -O pocketbase.zip https://github.com/pocketbase/pocketbase/releases/download/v0.9.3/pocketbase_0.9.3_linux_amd64.zip && \
    unzip pocketbase.zip && \
    chmod +x pocketbase

# Final stage
FROM alpine:latest

# Set up working directory
WORKDIR /app

# Copy the backend binary and PocketBase from the builder stage
COPY --from=builder /app/backend .
COPY --from=builder /app/pocketbase .

# Copy frontend files
COPY frontend/index.html ./frontend/index.html

# Expose ports for Go backend and PocketBase
EXPOSE 8080 8090

# Create entrypoint script to run both PocketBase and the Go backend concurrently
COPY <<EOF /app/start.sh
#!/bin/sh
# Start PocketBase in the background
./pocketbase serve --http=0.0.0.0:8090 --dir=/pb_data &

# Start the Go backend
./backend
EOF

# Make start script executable
RUN chmod +x /app/start.sh

# Run entrypoint script
CMD ["/app/start.sh"]

