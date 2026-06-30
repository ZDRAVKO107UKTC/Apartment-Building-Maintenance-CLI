# syntax=docker/dockerfile:1

# --- Build stage ---------------------------------------------------------
FROM golang:1.26-alpine AS build

WORKDIR /src

# Cache dependencies first for faster incremental builds.
COPY go.mod go.sum ./
RUN go mod download

# Build the statically-linked CLI binary. The SQLite driver (glebarez/sqlite)
# is pure Go, so the binary builds with CGO disabled.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/maintenance ./cmd

# --- Runtime stage -------------------------------------------------------
FROM alpine:3.20

# Non-root user and CA certificates for outbound HTTPS (SendGrid).
# /data holds the SQLite database file and is owned by the runtime user.
RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 appuser \
    && mkdir -p /data \
    && chown appuser:appuser /data

COPY --from=build /out/maintenance /usr/local/bin/maintenance

USER appuser

ENTRYPOINT ["maintenance"]
CMD ["--help"]
