# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOBIN=./bin
BINARY_NAME=ecom

# Docker parameters
DOCKER_IMAGE=ecom-app
DOCKER_BUILD_CMD=docker build -t $(DOCKER_IMAGE) .
DOCKER_RUN_CMD=docker run --rm -it -p 8080:8080 --env-file $(ENV_FILE) $(DOCKER_IMAGE)

# Environment-specific .env file
ENV_FILE=.env

all: test build

build:
	$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) -v ./cmd

docker-build:
	$(DOCKER_BUILD_CMD)

test:
	$(GOTEST) -v ./...

docker-test:
	docker run --rm -it --env-file $(ENV_FILE) $(DOCKER_IMAGE) $(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(GOBIN)

docker-clean:
	docker rmi $(DOCKER_IMAGE)

run:
	$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) -v ./cmd
	$(GOBIN)/$(BINARY_NAME)

docker-run:
	docker run --rm -it -p 8080:8080 --env-file $(ENV_FILE) -v $(PWD)/$(ENV_FILE):/app/$(ENV_FILE) $(DOCKER_IMAGE)


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
