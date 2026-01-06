# Variables
BINARY := bin/dividr
CMD := ./cmd/dividr
DOCKER_TAG ?= dividr:latest
GOLANGCI_LINT ?= golangci-lint
MIGRATE ?= migrate
SQLC ?= sqlc
TEMPL ?= templ

# Phony targets (not real files)
.PHONY: all build dev test lint tidy fmt vet run clean frontend docker migrate-up sqlc generate

# Default target
all: build

# 1. GENERATE: Compiles Templ components and SQLC queries to Go code
generate:
	@echo "Generating Templ components..."
	$(TEMPL) generate
	@echo "Generating SQLC..."
	# $(SQLC) generate  <-- Uncomment this when you start Story 0.2.1

# 2. BUILD: Depends on 'generate' so code exists before compilation
build: generate
	@echo "Building binary..."
	@mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) $(CMD)

# 3. DEV: Using 'air' is recommended for live reloading Templ changes.
# If you don't use air, this runs generate then go run.
dev:
	@echo "Starting Watched Dev Server..."
	# This runs the air command, which reads .air.toml
	air

test: generate
	go test -v ./...

lint:
	$(GOLANGCI_LINT) run ./...

docker:
	docker build -t $(DOCKER_TAG) .

migrate-up:
	$(MIGRATE) -path ./migrations -database $(DATABASE_URL) up

sqlc:
	$(SQLC) generate

tidy:
	go mod tidy

fmt:
	gofmt -w .
	$(TEMPL) fmt .  # Also format your templ files!

vet:
	go vet ./...

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	# Optional: Clean generated templ files if you want a fresh start
	# find . -name "*_templ.go" -delete

frontend:
	cd frontend && npm ci && npm run build

docker-build:
	docker build -f ./docker/Dockerfile -t $(DOCKER_TAG) .

docker-clean:
	docker image rm $(DOCKER_TAG) || true
