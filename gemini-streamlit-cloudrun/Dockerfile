# Base image: official Go 1.22
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go.mod and go.sum from the root of the workshop
COPY go.mod go.sum ./
COPY adk-go ./adk-go
COPY gemini-streamlit-cloudrun ./gemini-streamlit-cloudrun

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/app ./gemini-streamlit-cloudrun

# Final runtime image: minimal Alpine
FROM alpine:latest

EXPOSE 8080

WORKDIR /app

# Copy the compiled binary from builder
COPY --from=builder /build/app .

# Copy .env.example for reference (not actually used at runtime)
COPY .env.example .

ENTRYPOINT ["./app"]
