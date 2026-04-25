# Root-me Stats Makefile

.PHONY: all build test clean install server cli badges help

# Default target
all: build

# Build all binaries
build: server cli badges
	@echo "✅ All binaries built successfully"

# Build server
server:
	@echo "🔨 Building server..."
	@go build -ldflags="-s -w" -o bin/rootme-server cmd/server/main.go

# Build CLI
cli:
	@echo "🔨 Building CLI..."
	@go build -ldflags="-s -w" -o bin/rootme-cli cmd/cli/main.go

# Build simple badge generator
badges:
	@echo "🔨 Building badge generator..."
	@go build -ldflags="-s -w" -o bin/rootme-badges cmd/badges/main.go

# Cross-compile for different platforms
build-all:
	@echo "🌍 Cross-compiling for all platforms..."
	@mkdir -p dist

	# Linux
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/rootme-server-linux-amd64 cmd/server/main.go
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/rootme-cli-linux-amd64 cmd/cli/main.go
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/rootme-badges-linux-amd64 cmd/badges/main.go

	# Windows
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/rootme-server-windows-amd64.exe cmd/server/main.go
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/rootme-cli-windows-amd64.exe cmd/cli/main.go
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/rootme-badges-windows-amd64.exe cmd/badges/main.go

	# macOS
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/rootme-server-darwin-amd64 cmd/server/main.go
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/rootme-cli-darwin-amd64 cmd/cli/main.go
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/rootme-badges-darwin-amd64 cmd/badges/main.go

	# macOS ARM64 (M1/M2)
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/rootme-server-darwin-arm64 cmd/server/main.go
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/rootme-cli-darwin-arm64 cmd/cli/main.go
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/rootme-badges-darwin-arm64 cmd/badges/main.go

	@echo "✅ Cross-compilation complete! Check dist/ directory"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Test with a real profile
test-real:
	@echo "🧪 Testing with real profile..."
	@go run cmd/badges/main.go --nickname=g0uZ --theme=dark --output=test-badge.svg
	@echo "✅ Test badge generated: test-badge.svg"

# Run benchmarks
bench:
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	@go vet ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/ dist/ .cache/ test-badge.svg
	@go clean

# Install binaries to $GOPATH/bin
install: build
	@echo "📦 Installing binaries..."
	@cp bin/rootme-server $(GOPATH)/bin/
	@cp bin/rootme-cli $(GOPATH)/bin/
	@cp bin/rootme-badges $(GOPATH)/bin/
	@echo "✅ Binaries installed to $(GOPATH)/bin/"

# Run server locally
run-server: server
	@echo "🚀 Starting server on http://localhost:3000"
	@./bin/rootme-server

# Quick development test
dev-test: badges
	@echo "🏃‍♂️ Quick development test..."
	@./bin/rootme-badges --nickname=g0uZ --theme=dark --stats | head -3

# Create release packages
release: build-all
	@echo "📦 Creating release packages..."
	@mkdir -p release

	# Linux
	@tar -czf release/rootme-stats-linux-amd64.tar.gz -C dist rootme-server-linux-amd64 rootme-cli-linux-amd64 rootme-badges-linux-amd64

	# Windows
	@zip -j release/rootme-stats-windows-amd64.zip dist/rootme-*-windows-amd64.exe

	# macOS Intel
	@tar -czf release/rootme-stats-darwin-amd64.tar.gz -C dist rootme-server-darwin-amd64 rootme-cli-darwin-amd64 rootme-badges-darwin-amd64

	# macOS ARM
	@tar -czf release/rootme-stats-darwin-arm64.tar.gz -C dist rootme-server-darwin-arm64 rootme-cli-darwin-arm64 rootme-badges-darwin-arm64

	@echo "✅ Release packages created in release/ directory"

# Docker build
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t rootme-stats .

# Docker run
docker-run: docker-build
	@echo "🐳 Running Docker container..."
	@docker run -p 3000:3000 rootme-stats

# Show help
help:
	@echo "🎯 Root-me Stats Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build        Build all binaries (server, cli, badges)"
	@echo "  server       Build server binary only"
	@echo "  cli          Build CLI binary only"
	@echo "  badges       Build simple badge generator only"
	@echo "  build-all    Cross-compile for all platforms"
	@echo "  test         Run tests"
	@echo "  test-real    Test with real Root-me profile"
	@echo "  bench        Run benchmarks"
	@echo "  fmt          Format Go code"
	@echo "  lint         Lint Go code"
	@echo "  clean        Remove build artifacts"
	@echo "  install      Install binaries to GOPATH/bin"
	@echo "  run-server   Build and run server locally"
	@echo "  dev-test     Quick development test"
	@echo "  release      Create release packages"
	@echo "  docker-build Build Docker image"
	@echo "  docker-run   Build and run Docker container"
	@echo "  help         Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build           # Build all binaries"
	@echo "  make run-server      # Start server on localhost:3000"
	@echo "  make test-real       # Test badge generation"
	@echo "  make build-all       # Cross-compile for Linux/Windows/macOS"

# Ensure bin directory exists
bin:
	@mkdir -p bin

# Add bin dependency to build targets
server: | bin
cli: | bin
badges: | bin