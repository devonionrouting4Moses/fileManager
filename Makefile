.PHONY: all clean rust go run install test help dev web frontend

# Variables
APP_NAME := filemanager
VERSION := 2.0.0
BUILD_DIR := dist
RUST_DIR := rust_ffi
GO_MAIN := ./cmd/app
INSTALL_PREFIX := /usr/local

# Platform detection
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    LIB_EXT := so
    LIB_NAME := libfilemanager.so
endif
ifeq ($(UNAME_S),Darwin)
    LIB_EXT := dylib
    LIB_NAME := libfilemanager.dylib
endif

# Default target
all: rust go

# Build Rust shared library
rust:
	@echo "🦀 Building Rust library..."
	cd $(RUST_DIR) && cargo build --release -p fs-operations-core
	@echo "✅ Rust library built: $(RUST_DIR)/target/release/$(LIB_NAME)"

# Build Go binary
go: rust
	@echo "🐹 Building Go binary..."
	cd file_manager && \
	CGO_ENABLED=1 \
	CGO_LDFLAGS="-L../$(RUST_DIR)/target/release -lfs_operations_core -ldl -lpthread -lm" \
	go build -ldflags="-s -w" -o ../$(APP_NAME).bin $(GO_MAIN)
	@chmod +x ../$(APP_NAME).sh
	@echo "✅ Go binary built: ./$(APP_NAME).bin (wrapper: ./$(APP_NAME).sh)"

# Development build (with debug symbols)
dev:
	@echo "🔧 Building in development mode..."
	cd $(RUST_DIR) && cargo build -p fs-operations-core
	cd file_manager && \
	CGO_ENABLED=1 \
	CGO_LDFLAGS="-L../$(RUST_DIR)/target/debug -lfs_operations_core" \
	go build -race -o ../$(APP_NAME).bin $(GO_MAIN)
	@chmod +x ../$(APP_NAME).sh
	@echo "✅ Development build complete"

# Run the application
run: all
	@echo "🚀 Running $(APP_NAME)..."
	@./$(APP_NAME).sh

# Run web server mode
web: all
	@echo "🌐 Starting web server..."
	@./$(APP_NAME).sh --web

# Install to system
install: all
	@echo "📦 Installing $(APP_NAME)..."
	sudo install -m 755 $(APP_NAME).bin $(INSTALL_PREFIX)/bin/$(APP_NAME).bin
	sudo install -m 755 $(APP_NAME).sh $(INSTALL_PREFIX)/bin/$(APP_NAME)
	sudo install -m 644 $(RUST_DIR)/target/release/$(LIB_NAME) $(INSTALL_PREFIX)/lib/
ifeq ($(UNAME_S),Linux)
	sudo ldconfig
endif
	@echo "✅ Installed to $(INSTALL_PREFIX)/bin/$(APP_NAME)"

# Uninstall from system
uninstall:
	@echo "🗑️  Uninstalling $(APP_NAME)..."
	sudo rm -f $(INSTALL_PREFIX)/bin/$(APP_NAME)
	sudo rm -f $(INSTALL_PREFIX)/lib/$(LIB_NAME)
ifeq ($(UNAME_S),Linux)
	sudo ldconfig
endif
	@echo "✅ Uninstalled"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	cd $(RUST_DIR) && cargo clean
	rm -f $(APP_NAME)
	rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running Rust tests..."
	cd $(RUST_DIR) && cargo test
	@echo "🧪 Running Go tests..."
	go test -v ./...

# Run Rust benchmarks
bench:
	@echo "⚡ Running benchmarks..."
	cd $(RUST_DIR) && cargo bench

# Format code
fmt:
	@echo "🎨 Formatting code..."
	cd $(RUST_DIR) && cargo fmt
	go fmt ./...
	@echo "✅ Formatting complete"

# Lint code
lint:
	@echo "🔍 Linting code..."
	cd $(RUST_DIR) && cargo clippy -- -D warnings
	golangci-lint run ./... || echo "⚠️  Install golangci-lint for Go linting"
	@echo "✅ Linting complete"

# Build for all platforms
build-all: clean
	@echo "🏗️  Building for all platforms..."
	@bash build.sh

# Create distribution packages
dist: build-all
	@echo "📦 Creating distribution packages..."
	@mkdir -p $(BUILD_DIR)
	@echo "✅ Distribution packages created in $(BUILD_DIR)/"

# Check for updates
update-check:
	@echo "🔄 Checking for updates..."
	@./$(APP_NAME) --update

# Show version
version:
	@echo "$(APP_NAME) v$(VERSION)"
	@echo "Rust: $$(rustc --version)"
	@echo "Go: $$(go version)"

# Development environment setup
setup-dev:
	@echo "🛠️  Setting up development environment..."
	@echo "Installing Rust..."
	@which rustc || curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
	@echo "Installing Go dependencies..."
	go mod download
	@echo "✅ Development environment ready"

# Generate documentation
docs:
	@echo "📚 Generating documentation..."
	cd $(RUST_DIR) && cargo doc --no-deps
	@echo "✅ Rust docs: $(RUST_DIR)/target/doc/fs_operations_core/index.html"

# Help target
help:
	@echo "FileManager v$(VERSION) - Available targets:"
	@echo ""
	@echo "Building:"
	@echo "  all         - Build both Rust library and Go binary (default)"
	@echo "  rust        - Build only Rust shared library"
	@echo "  go          - Build only Go binary (requires Rust library)"
	@echo "  dev         - Build in development mode with debug symbols"
	@echo "  build-all   - Build for all platforms (Linux, macOS, Windows)"
	@echo ""
	@echo "Running:"
	@echo "  run         - Build and run the application"
	@echo "  web         - Build and run web server mode"
	@echo ""
	@echo "Installation:"
	@echo "  install     - Install to system (requires sudo)"
	@echo "  uninstall   - Remove from system (requires sudo)"
	@echo ""
	@echo "Development:"
	@echo "  test        - Run all tests"
	@echo "  bench       - Run Rust benchmarks"
	@echo "  fmt         - Format all code"
	@echo "  lint        - Lint all code"
	@echo "  docs        - Generate documentation"
	@echo "  setup-dev   - Setup development environment"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean       - Remove all build artifacts"
	@echo "  version     - Show version information"
	@echo "  update-check - Check for updates"
	@echo "  help        - Show this help message"
	@echo ""