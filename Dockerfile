# Stage 1: Build the Go application
FROM golang:1.23-alpine as builder

LABEL maintainer="Juan Vargas <vargasm.jp@gmail.com>"

WORKDIR /app

# Enable multi-arch builds
ARG TARGETARCH
ENV GOARCH=$TARGETARCH

# Copy dependency files and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary for the correct architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o guilliman cmd/server/main.go

# Stage 2: Create a lightweight runtime image
FROM alpine:latest

RUN apk add --no-cache jq

WORKDIR /app

# Copy only necessary files from builder stage
COPY --from=builder /app/guilliman .
COPY init_db.sql .
COPY seed_db.sql .
COPY entrypoint.sh .
COPY wsgi.py .

# Ensure entrypoint.sh is executable
RUN chmod +x ./entrypoint.sh

# Expose ports for Go app
EXPOSE 8080 8081

# Use entrypoint.sh to manage database initialization and app startup
ENTRYPOINT [ "./entrypoint.sh" ]
