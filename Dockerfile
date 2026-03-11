# syntax=docker/dockerfile:1

# ─── Stage 1: Builder ──────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Download dependencies first (layer cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s -X main.version=${VERSION:-dev}" \
    -o /bin/social_media ./cmd/...

# ─── Stage 2: Final (distroless) ───────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12 AS runner

# Copy timezone data and CA certs from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /bin/social_media /bin/social_media

EXPOSE 8080

ENTRYPOINT ["/bin/social_media"]
