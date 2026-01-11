ARG GO_VERSION=1.25
ARG VERSION=dev

# ==============================================================================
# STAGE 1: The Builder
# ==============================================================================
FROM golang:${GO_VERSION}-bookworm AS builder

ARG VERSION

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
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
    mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss && \
    chmod +x /usr/local/bin/tailwindcss

# 5. Copy Source Code
COPY . .

# 6. Run the Build Pipeline
RUN make css && CGO_ENABLED=0 GOOS=linux make build

# ==============================================================================
# STAGE 2: The Runner (Production)
# ==============================================================================
FROM alpine:latest

# 1. Setup certificates and user
RUN apk --no-cache add ca-certificates
# 2. Create non-root user (Alpine syntax)
RUN addgroup -g 1001 -S dividr && \
    adduser -u 1001 -S dividr -G dividr

WORKDIR /home/dividr
RUN chown -R dividr:dividr /home/dividr
USER dividr

# 3. Copy Binary
COPY --from=builder --chown=dividr:dividr /app/bin/dividr .

EXPOSE 8080
CMD ["./dividr"]
