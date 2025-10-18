FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

# Build main application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/app ./cmd/main.go

# Build migration binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/migrate ./cmd/migrate

FROM alpine:latest

# Install ca-certificates for HTTPS requests, PostgreSQL client, and timezone data
RUN apk --no-cache add ca-certificates postgresql-client tzdata

# Copy binaries
COPY --from=builder /app/app /app/migrate /app/

# Copy migrations
COPY internal/migrations /app/migrations

# No entrypoint script needed

WORKDIR /app

# Default entrypoint for main application
ENTRYPOINT ["/app/app"]