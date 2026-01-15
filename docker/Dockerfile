ARG GO_VERSION=1.25
ARG VERSION=dev

# ==============================================================================
# STAGE 1: The Builder
# ==============================================================================
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS builder

ARG VERSION
ARG TARGETOS
ARG TARGETARCH

# 1. Install system dependencies
RUN apt-get update && apt-get install -y git make curl

# 2. Install Templ CLI and SQLC
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app

# 3. Cache Dependencies
COPY go.mod go.sum ./
RUN go mod download

# 4. Install Tailwind CSS (Globally in the Container)
# This allows us to generate CSS without needing a Node.js stage
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
    mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss && \
    chmod +x /usr/local/bin/tailwindcss

# 5. Copy Source Code
# This copies everything, including your 'internal/database/migrations' for embedding
COPY . .

# 6. Run the Build Pipeline
# We generate Templ files -> Generate SQLC -> Generate CSS -> Build Binary
# This ensures 'go:embed' finds the CSS and Migration files exactly where it expects them.
RUN templ generate && \
    sqlc generate -f internal/database/sqlc.yaml && \
    make css && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH make build

# ==============================================================================
# STAGE 2: The Runner (Production)
# ==============================================================================
FROM alpine:latest

# 1. Setup certificates
RUN apk --no-cache add ca-certificates

# 2. Create non-root user
RUN addgroup -g 1001 -S dividr && \
    adduser -u 1001 -S dividr -G dividr

WORKDIR /home/dividr
RUN chown -R dividr:dividr /home/dividr
USER dividr

# 3. Copy Binary
COPY --from=builder --chown=dividr:dividr /app/bin/dividr .

# 4. Copy Static Assets (Optional but recommended)
COPY --from=builder --chown=dividr:dividr /app/web/static ./web/static

EXPOSE 8080
CMD ["./dividr"]
