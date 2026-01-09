ARG GO_VERSION=1.25
ARG VERSION=dev

# ==============================================================================
# STAGE 1: The Builder
# Use a full Go image to compile the application and generate templates.
# ==============================================================================
FROM golang:${GO_VERSION}-alpine AS builder

# Pass VERSION arg to builder stage
ARG VERSION

# 1. Install system dependencies
# 'git' is often needed for go mod download.
# 'make' is useful if you want to use your Makefile commands.
RUN apk add --no-cache git make

# 2. Install Templ CLI
# We install it inside the container so it can run 'templ generate'
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# 3. Cache Dependencies
# We copy go.mod/sum first to leverage Docker layer caching.
# If these files haven't changed, this step is skipped.
COPY go.mod go.sum ./
RUN go mod download

# 4. Copy Source Code
COPY . .

# 5. Generate Templ Code
# This converts your .templ files into .go files inside the container
RUN templ generate

# 6. Build the Binary
# -ldflags="-w -s" strips debug information to reduce binary size
# -X injects the version at build time (from build arg or git tags as fallback)
RUN if [ "$VERSION" = "dev" ]; then \
      VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "dev"); \
    fi && \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X 'github.com/rhysmcneill/dividr/internal/http/handler.Version=${VERSION}'" \
    -o ./bin/dividr ./cmd/dividr

# ==============================================================================
# STAGE 2: The Runner (Production)
# Use a tiny Alpine image. No Go, No Templ, just the binary.
# ==============================================================================
FROM alpine:latest

# 1. Install CA Certificates (must be done as root before switching users)
# Crucial if your app makes HTTPS requests (e.g. to AWS S3, Stripe, Auth providers)
RUN apk --no-cache add ca-certificates

# 2. Create non-root user (Alpine syntax)
RUN addgroup -g 1001 -S dividr && \
    adduser -u 1001 -S dividr -G dividr

# 3. Set working directory and ownership
WORKDIR /home/dividr
RUN chown -R dividr:dividr /home/dividr

# 4. Switch to non-root user
USER dividr

# 5. Copy the binary from the builder stage
COPY --from=builder --chown=dividr:dividr /app/bin/dividr .

# 6. Copy other necessary files (if any)
# e.g., if you have a .env.prod or assets folder
# COPY --from=builder --chown=dividr:dividr /app/assets ./assets

# 7. Expose the port
EXPOSE 8080

# 8. Run the binary
CMD ["./dividr"]
