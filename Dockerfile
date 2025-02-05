# Stage 1: Build the Go application
FROM golang:1.23-alpine as builder

LABEL maintainer="Juan Vargas <vargasm.jp@gmail.com>"

WORKDIR /app

# Copy dependency files and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN go build -o guilliman cmd/server/main.go

# Stage 2: Create a lightweight runtime image
FROM golang:1.23-alpine

WORKDIR /app

# Copy the compiled Go binary and required files
COPY --from=builder /app/guilliman .
COPY init_db.sql .
COPY seed_db.sql .
COPY entrypoint.sh .
COPY wsgi.py .

# Ensure entrypoint.sh is executable
RUN chmod +x ./entrypoint.sh

# Expose ports (Go app and sqlite_web)
EXPOSE 8080 8081

# Use entrypoint.sh to manage database initialization and app startup
ENTRYPOINT [ "./entrypoint.sh" ]
