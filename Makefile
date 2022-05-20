export GOBIN ?= $(shell pwd)/bin
PATH := ${PATH}:${GOBIN}

GOLINT = golangci-lint run
STATICCHECK = staticcheck

# Directories containing independent Go modules.
#
# We track coverage only for the main module.
MODULE_DIRS = .

# Many Go tools take file globs or directories as arguments instead of packages.
GO_FILES := $(shell \
	find . '(' -path '*/.*' -o -path './vendor' ')' -prune \
	-o -name '*.go' -print | cut -b3-)

.PHONY: all
all: lint test

.PHONY: lint
lint:
	@rm -rf lint.log
	@echo "Checking formatting..."
	@gofmt -d -s $(GO_FILES) 2>&1 | tee lint.log
	@echo "Checking vet..."
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && go vet ./... 2>&1) &&) true | tee -a lint.log
	@echo "Checking lint..."
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && $(GOLINT) ./... 2>&1) &&) true | tee -a lint.log
	@echo "Checking for unresolved FIXMEs..."
	@git grep -i fixme | grep -v -e Makefile | tee -a lint.log
	@[ ! -s lint.log ]
	@echo "Checking 'go mod tidy'..."
	@make tidy
	@if ! git diff --quiet; then \
		echo "'go mod tidy' resulted in changes or working tree is dirty:"; \
		git --no-pager diff; \
	fi

.PHONY: test
test:
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && go test -race ./...) &&) true

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: tidy
tidy:
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && go mod tidy) &&) true
