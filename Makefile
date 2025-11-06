.PHONY: all clean rust go run install test help

# Default target
all: rust go

# Build Rust shared library
rust:
	@echo "🦀 Building Rust library..."
	cargo build --release
	@echo "✅ Rust library built: target/release/libfilemanager.so"

# Build Go binary
go: rust
	@echo "🐹 Building Go binary..."
	go build -o filemanager .
	@echo "✅ Go binary built: ./filemanager"

# Run the application
run: all
	@echo "🚀 Running filemanager..."
	@./filemanager

# Install to /usr/local/bin (requires sudo)
install: all
	@echo "📦 Installing filemanager..."
	sudo cp filemanager /usr/local/bin/
	sudo cp target/release/libfilemanager.so /usr/local/lib/
	sudo ldconfig
	@echo "✅ Installed to /usr/local/bin/filemanager"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	cargo clean
	rm -f filemanager
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running Rust tests..."
	cargo test
	@echo "🧪 Running Go tests..."
	go test ./...

# Development build (with debug symbols)
dev:
	@echo "🔧 Building in development mode..."
	cargo build
	CGO_LDFLAGS="-L./target/debug -lfilemanager" go build -o filemanager .
	@echo "✅ Development build complete"

# Help target
help:
	@echo "Available targets:"
	@echo "  all      - Build both Rust library and Go binary (default)"
	@echo "  rust     - Build only Rust shared library"
	@echo "  go       - Build only Go binary (requires Rust library)"
	@echo "  run      - Build and run the application"
	@echo "  install  - Install to system (requires sudo)"
	@echo "  clean    - Remove all build artifacts"
	@echo "  test     - Run all tests"
	@echo "  dev      - Build in development mode with debug symbols"
	@echo "  help     - Show this help message"