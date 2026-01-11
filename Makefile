# Variables
BINARY := bin/dividr
CMD := ./cmd/dividr
DOCKER_TAG ?= dividr:latest
GOLANGCI_LINT ?= golangci-lint
MIGRATE ?= migrate
SQLC ?= sqlc
TEMPL ?= templ
TAILWINDCSS ?= tailwindcss
VERSION ?= dev

# Database configuration - reads password from secret file
DB_PASSWORD := $(shell cat docker/secrets/db_password.txt 2>/dev/null || echo "")
DATABASE_URL ?= postgres://dividr:$(DB_PASSWORD)@localhost:5432/dividr?sslmode=disable

# Phony targets (not real files)
.PHONY: all build dev test lint tidy fmt vet run clean frontend docker migrate-create migrate-up migrate-down migrate-force docker-build docker-clean generate sqlc css css-watch

# Default target
all: build

# GENERATE: Compiles Templ components and SQLC queries to Go code
generate:
	@echo "Generating Templ components..."
	$(TEMPL) generate
	@echo "Generating SQLC..."
	cd internal/database && $(SQLC) generate && cd ../..

# CSS: Build Tailwind CSS
css:
	@echo "Building Tailwind CSS..."
	@mkdir -p web/static/css
	$(TAILWINDCSS) -i web/css/input.css -o web/static/css/output.css --minify

# CSS-WATCH: Build and watch Tailwind CSS for changes
css-watch:
	@echo "Watching Tailwind CSS..."
	@mkdir -p web/static/css
	$(TAILWINDCSS) -i web/css/input.css -o web/static/css/output.css --watch

# BUILD: Depends on 'generate' so code exists before compilation
build: generate
	@echo "Building binary (version: $(VERSION))..."
	@mkdir -p $(dir $(BINARY))
	go build -ldflags "-X 'github.com/rhysmcneill/dividr/internal/http/handler.Version=$(VERSION)'" -o $(BINARY) $(CMD)


test: generate
	go test -v ./...

lint:
	$(GOLANGCI_LINT) run ./...

docker:
	docker build -t $(DOCKER_TAG) .

# --- Database Migrations ---

## migrate-create: Create a new migration file. Usage: make migrate-create name=init_schema
migrate-create:
	@echo "Creating migration files for: $(name)..."
	$(MIGRATE) create -ext sql -dir ./internal/database/migrations -seq $(name)

## migrate-up: Apply all up migrations
migrate-up:
	@echo "Applying migrations..."
	$(MIGRATE) -path ./internal/database/migrations -database "$(DATABASE_URL)" up

## migrate-down: Rollback the last migration step
migrate-down:
	@echo "Rolling back last migration..."
	$(MIGRATE) -path ./internal/database/migrations -database "$(DATABASE_URL)" down 1

## migrate-force: Fix dirty state (use if migration fails and locks the DB). Usage: make migrate-force version=1
migrate-force:
	@echo "Forcing migration version $(version)..."
	$(MIGRATE) -path ./internal/database/migrations -database "$(DATABASE_URL)" force $(version)

sqlc:
	cd internal/database && $(SQLC) generate

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

# TEMPL: Generate Go code
templ:
	@templ generate

# TEMPL-WATCH: Watch .templ files
templ-watch:
	@templ generate --watch --proxy="http://localhost:8080" --open-browser=false
#
 DEV: Run Air, Templ Watch, and Tailwind Watch in parallel
dev:
	@make -j3 css-watch templ-watch air

# Helper to run air (so we can reference it in 'dev')
air:
	@air
