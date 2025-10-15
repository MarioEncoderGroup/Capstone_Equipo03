# Build stage
FROM public.ecr.aws/docker/library/golang:alpine AS builder
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN cd cmd/api && go build -o /misviaticos-api

# Final stage
FROM public.ecr.aws/docker/library/golang:alpine AS release
WORKDIR /app

# Install necessary packages including migrate CLI
RUN apk add --no-cache ca-certificates tzdata postgresql-client curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

# Copy necessary files from builder stage
COPY --from=builder /app/email_templates /app/email_templates
COPY --from=builder /app/db /app/db
COPY --from=builder /misviaticos-api .

# Copy entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Set timezone for Chile (MisViaticos market)
ENV TZ="America/Santiago"

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/ || exit 1

# Use entrypoint script to run migrations before starting API
ENTRYPOINT ["/app/entrypoint.sh"]
