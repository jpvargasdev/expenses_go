# Stage 1: Build the Go application
FROM golang:1.23-alpine as builder

LABEL maintainer="Juan Vargas <vargasm.jp@gmail.com>"

WORKDIR /app

ENV CGO_ENABLED=1
ENV GOOS=linux

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy dependency files and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN go build -o guilliman cmd/server/main.go

# Stage 2: Create a lightweight runtime image
FROM python:3.7-alpine3.17

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates musl sqlite jq

# Install Python dependencies
RUN pip install gevent sqlite_web

# Copy the compiled Go binary and required files
COPY --from=builder /app/guilliman .
COPY init_db.sql .
COPY seed_db.sql .
COPY entrypoint.sh .
COPY wsgi.py .

# Ensure entrypoint.sh is executable
RUN chmod +x ./entrypoint.sh
RUN echo "$GOOGLE_APPLICATION_CREDENTIALS_JSON" | jq '.' > /app/firebase-config.json

# Expose ports (Go app and sqlite_web)
EXPOSE 8080 8081

# Use entrypoint.sh to manage database initialization and app startup
ENTRYPOINT [ "./entrypoint.sh" ]
