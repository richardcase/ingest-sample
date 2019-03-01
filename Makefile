BUILDCOMMIT := $(shell git describe --dirty --always)
BUILDDATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VER_FLAGS=-X main.commit=$(BUILDCOMMIT) -X main.date=$(BUILDDATE)

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL:=help

##@ Build

.PHONY: build
build: ## Build the services
	@go build -ldflags "$(VER_FLAGS)" ./cmd/ingest-svc
	@go build -ldflags "$(VER_FLAGS)" ./cmd/person-svc

.PHONY: release
release: ## Build a release version of services
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s $(VER_FLAGS)" -o $(GOPATH)/bin/ingest-svc ./cmd/ingest-svc
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s $(VER_FLAGS)" -o $(GOPATH)/bin/person-svc ./cmd/person-svc

.PHONY: docker-build
docker-build: ## Build a release version of service via multi-stage docker file
	docker build -f Dockerfile-ingest  .
	docker build -f Dockerfile-personsvc .

##@ Testing & CI

.PHONY: test
test:   ## Run unit tests
	#@git diff --exit-code ./pkg/api/v1/person-service.proto > /dev/null || (echo "Proto changed, update generated code"; exit 1)
	#@git diff --exit-code ./pkg/repository/mocks > /dev/null || (echo "Mocked interface changed, update mocks"; exit 1)
	@go test -v -covermode=count -coverprofile=coverage.out ./pkg/... ./cmd/...

.PHONY: lint
lint: ## Run linting over the codebase
	golangci-lint run

.PHONY: ci
ci: test lint ## Target for CI system to invoke to run tests and linting

##@ Code Generation

.PHONY: codegen
codegen:
	protoc --proto_path=pkg/api/proto --go_out=plugins=grpc:pkg/api person-service.proto
	go generate ./pkg/repository/mocks
	go generate ./pkg/api/mocks

##@ Utility

.PHONY: fmt
fmt: ## Format all the source code using gofmt
	@gofmt -l -w $(SRC)

.PHONY: help
help:  ## Display this help. Thanks to https://suva.sh/posts/well-documented-makefiles/
@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
