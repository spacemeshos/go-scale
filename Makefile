all: build test
.PHONY: all

build:
	go build -o bin/go-scale ./scalegen
.PHONY: build

test:
	gotestsum -- -timeout 5m -p 1 -race ./...
.PHONY: test

install:
	go mod download
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.0
	go install gotest.tools/gotestsum@v1.9.0
	go install honnef.co/go/tools/cmd/staticcheck@v0.4.3
.PHONY: install

tidy:
	go mod tidy
.PHONY: tidy

test-tidy:
	# Working directory must be clean, or this test would be destructive
	git diff --quiet || (echo "\033[0;31mWorking directory not clean!\033[0m" && git --no-pager diff && exit 1)
	# We expect `go mod tidy` not to change anything, the test should fail otherwise
	make tidy
	git diff --exit-code || (git --no-pager diff && git checkout . && exit 1)
.PHONY: test-tidy

test-fmt:
	git diff --quiet || (echo "\033[0;31mWorking directory not clean!\033[0m" && git --no-pager diff && exit 1)
	# We expect `go fmt` not to change anything, the test should fail otherwise
	go fmt ./...
	git diff --exit-code || (git --no-pager diff && git checkout . && exit 1)
.PHONY: test-fmt

clear-test-cache:
	go clean -testcache
.PHONY: clear-test-cache

lint:
	go vet ./...
	./bin/golangci-lint run --config .golangci.yml
.PHONY: lint

# Auto-fixes golangci-lint issues where possible.
lint-fix:
	./bin/golangci-lint run --config .golangci.yml --fix
.PHONY: lint-fix

lint-github-action:
	go vet ./...
	./bin/golangci-lint run --config .golangci.yml --out-format=github-actions
.PHONY: lint-github-action

cover:
	go test -coverprofile=cover.out -timeout 0 -p 1 $(UNIT_TESTS)
.PHONY: cover

staticcheck:
	staticcheck ./...
.PHONY: staticcheck
