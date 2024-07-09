# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21.0-alpine AS build-stage
WORKDIR /app

# Install any build tools needed
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and configuration files
COPY . .
COPY ./config ./config

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/server/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage
WORKDIR /app

# Copy the built binary and configuration files from the previous stage
COPY --from=build-stage /api /api
COPY --from=build-stage /app/config ./config
COPY --from=build-stage /app/.env .env

# Expose the necessary port
EXPOSE 8080

# Health check to ensure the container is healthy
HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD curl --fail http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["/api"]
