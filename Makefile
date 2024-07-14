# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOBIN=./bin
BINARY_NAME=ecom

# Docker parameters
DOCKER_COMPOSE=docker-compose
DOCKER_IMAGE=ecom-app

# Environment-specific .env file
ENV_FILE=.env

all: test build

# Define variables
GOFMT = gofmt
FIND = find

# Default target
all: gofmtcheck

# Gofmt check and fix target
gofmtcheck:
	@echo "Checking and fixing gofmt issues..."
	@need_fmt=$$($(GOFMT) -l $$($(FIND) . -type f -name '*.go' -not -path './vendor/*')); \
	if [ "$$need_fmt" = "" ]; then \
		echo "All files are properly formatted!"; \
	else \
		echo "Files that need formatting and will be fixed:"; \
		echo "$$need_fmt"; \
		$(GOFMT) -w $$need_fmt; \
		echo "All files have been formatted."; \
	fi

.PHONY: all gofmtcheck

build:
	$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) -v ./cmd/server

docker-build:
	$(DOCKER_COMPOSE) build

test:
	$(GOTEST) -v ./...

docker-test:
	$(DOCKER_COMPOSE) run --rm app $(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(GOBIN)

docker-clean:
	docker rmi $(DOCKER_IMAGE)

run:
	$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) -v ./cmd/server
	$(GOBIN)/$(BINARY_NAME)

docker-run:
	$(DOCKER_COMPOSE) up --build

migrate-up:
	@$(GOCMD) run cmd/migrate/main.go up

migrate-down:
	@$(GOCMD) run cmd/migrate/main.go down

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

# Load environment-specific .env file
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

# Environment-specific targets
.PHONY: local dev prod docker-build docker-run docker-clean docker-test

local: ENV_FILE=.env.local
local: run

dev: ENV_FILE=.env.dev
dev: run

prod: ENV_FILE=.env.prod
prod: run
