BINARY := bin/dividr
CMD := ./cmd/dividr
DOCKER_TAG ?= dividr:latest
GOLANGCI_LINT ?= golangci-lint
MIGRATE ?= migrate
SQLC ?= sqlc

.PHONY: all build dev test lint tidy fmt vet run clean frontend docker migrate docker-build migrate-up sqlc
all: build

build:
	@mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) $(CMD)

dev:
	@echo "Starting dev (go run)..."
	go run $(CMD)

test:
	go test ./...

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

vet:
	go vet ./...

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

frontend:
	cd frontend && npm ci && npm run build

docker-build:
	docker build -t $(DOCKER_TAG) .