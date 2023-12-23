VERSION :=$(shell cat .version)
YAMLS   :=$(shell find . -type f -regex ".*y*ml" -print)
COMMIT  :=$(shell git rev-parse --short HEAD)
DATE    :=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

## Variable assertions
ifndef VERSION
	$(error RELEASE_VERSION is not set)
endif

all: help

.PHONY: version
version: ## Prints the current version
	@echo $(VERSION)

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependancies 
	go mod tidy
	go mod vendor

.PHONY: upgrade
upgrade: ## Upgrades all dependancies 
	go get -d -u ./...
	go mod tidy
	go mod vendor

.PHONY: test
test: tidy ## Runs unit tests
	go test -short -count=1 -race -covermode=atomic -coverprofile=cover.out ./...

.PHONY: cover
cover: test ## Runs unit tests and putputs coverage
	go tool cover -func=cover.out

.PHONY: lint
lint: lint-go lint-yaml ## Lints the entire project 
	@echo "Completed Go and YAML lints"

.PHONY: lint
lint-go: ## Lints the entire project using go 
	golangci-lint run -c .golangci.yaml

.PHONY: lint-yaml
lint-yaml: ## Runs yamllint on all yaml files (brew install yamllint)
	yamllint -c .yamllint $(YAMLS)

.PHONY: run
run: tidy ## Runs uncompiled version of the CLI
	go run main.go

.PHONY: build
build: tidy ## Builds CLI binary
	mkdir -p ./bin
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w \
    -X github.com/mchmarny/aictl/pkg/cli.version=$(VERSION) \
	-X github.com/mchmarny/aictl/pkg/cli.commit=$(COMMIT) \
	-X github.com/mchmarny/aictl/pkg/cli.date=$(DATE) \
	-extldflags '-static'" \
    -mod vendor -o bin/aictl main.go

.PHONY: tag
tag: ## Creates release tag 
	git tag -s -m "release $(VERSION)" $(VERSION)
	git push origin $(VERSION)

.PHONY: tagless
tagless: ## Delete the current release tag 
	git tag -d $(VERSION)
	git push --delete origin $(VERSION)

.PHONY: clean
clean: ## Cleans bin and temp directories
	go clean
	rm -fr ./vendor
	rm -fr ./bin

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
