operatingsystem := $(shell go env GOOS)
version := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
build_flags := -ldflags="-s -w -X main.version=$(version)"

.PHONY: all clean test lint build fmt release debug deps vet

all: build

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

build:
	@echo "Building release for '$(operatingsystem)'..."
	go build $(build_flags) -o tstat ./cmd/textstat

debug:
	@echo "Building debug version for '$(operatingsystem)'..."
	go build -gcflags="all=-N -l" -o tstat-debug ./cmd/textstat

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	@echo "Running linter..."
	golangci-lint run || echo "golangci-lint not installed, skipping..."

fmt:
	@echo "Formatting code..."
	go fmt ./...
	go vet ./...

vet:
	@echo "Vetting code..."
	go vet ./...

clean:
	@echo "Cleaning up..."
	@rm -rf ./*.a *.o core *.log tstat tstat-debug coverage.out coverage.html

release: fmt vet test build

install: build
	@echo "Installing textstat..."
	go install $(build_flags) ./cmd/textstat

uninstall:
	@echo "Uninstalling textstat..."
	go clean -i

dev-deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: all clean test lint build fmt release debug deps vet coverage test-bench install uninstall dev-deps
