ARG GO_VERSION=1.25
ARG VERSION=dev

# ==============================================================================
# STAGE 1: The Builder
# ==============================================================================
FROM golang:${GO_VERSION}-alpine AS builder

ARG VERSION

# 1. Install system dependencies
RUN apk add --no-cache git make curl

# 2. Install Templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# 3. Cache Dependencies
COPY go.mod go.sum ./
RUN go mod download

# 4. Install Tailwind CSS (Globally in the Container)
# We download it directly to /usr/local/bin/tailwindcss so it's in the PATH.
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
    mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss && \
    chmod +x /usr/local/bin/tailwindcss

# 5. Copy Source Code
COPY . .

# 6. Run the Build Pipeline
# Now 'make build' just works, because 'tailwindcss' is in the PATH
RUN make build

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
