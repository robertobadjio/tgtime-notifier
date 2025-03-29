#!/usr/bin/make
.DEFAULT_GOAL := help
.PHONY: help

DOCKER_COMPOSE ?= docker compose -f docker-compose.yml
PPROF_DIR = ./pprof
PPROF_COMMAND = go tool pprof
PANDORA_COMMAND = ./bin/pandora
BIN_DIR = ./bin
MINIMOCK_COMMAND = minimock
GO_TEST_COMMAND = go test
TEST_COVER_FILENAME = c.out

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-deps-mac: ## Install dependencies for MAC
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v1.64.5
	wget https://github.com/yandex/pandora/releases/download/v0.5.32/pandora_0.5.32_darwin_amd64 -O ./bin/pandora
	chmod a+x ./bin/pandora
	wget -qO- https://github.com/gojuno/minimock/releases/download/v3.4.5/minimock_3.4.5_darwin_amd64.tar.gz | gunzip | tar xvf - -C ./bin minimock

fmt: ## Automatically format source code
	go fmt ./...
.PHONY:fmt

lint: fmt lint-config-verify  ## Check code (lint)
	./bin/golangci-lint run ./... --config .golangci.pipeline.yaml
.PHONY:lint

lint-config-verify: fmt ## Verify config (lint)
	./bin/golangci-lint config verify --config .golangci.pipeline.yaml

vet: fmt ## Check code (vet)
	go vet -vettool=$(which shadow) ./...
.PHONY:vet

vet-shadow: fmt ## Check code with detect shadow (vet)
	go vet -vettool=$(which shadow) ./...
.PHONY:vet

build: ## Build service containers
	$(DOCKER_COMPOSE) build

up: vet ## Start services
	$(DOCKER_COMPOSE) up -d $(SERVICES)

down: ## Down services
	$(DOCKER_COMPOSE) down

clean: ## Delete all containers
	$(DOCKER_COMPOSE) down --remove-orphans

load-testing: ## Run load testing with pandora
	$(PANDORA_COMMAND) ./pandora/pandora.yml

prof-open: ## Open profile CLI version from ./pprof dir: filename=pprof.tgtime-notifier.samples.cpu.001.pb.gz
	# top
	# peek
	# disasm
	# tree
	# web
	$(PPROF_COMMAND) $(PPROF_DIR)/$(filename)

prof-cpu: ## CPU profiling UI version
	PPROF_TMPDIR=$(PPROF_DIR) $(PPROF_COMMAND) -http :8082 -seconds 20 http://127.0.0.1:8081/debug/pprof/profile

prof-mem: ## Memory profiling UI version
	PPROF_TMPDIR=$(PPROF_DIR) $(PPROF_COMMAND) -http :8082 http://127.0.0.1:8081/debug/pprof/heap

prof-trace: ## Trace profiling CLI version
	curl http://127.0.0.1:8081/debug/pprof/trace\?seconds\=20 -o $(PPROF_DIR)/trace.out

minimock-gen: ## Generate mocks
	$(BIN_DIR)/$(MINIMOCK_COMMAND) -i ./internal/config.OS -o ./internal/config/os_mock.go

test-unit: ## Run unit tests
	$(GO_TEST_COMMAND) \
		./internal/... \
		-count=1 \
		-cover -coverprofile=$(TEST_COVER_FILENAME)

test-unit-race: ## Run unit tests with -race flag
	$(GO_TEST_COMMAND) ./internal/... -count=1 -race