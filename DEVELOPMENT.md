# Development Guide

This guide is for developers who want to build FileManager from source, contribute to the project, or create custom builds.

## 📋 Prerequisites

### Required Tools

- **Rust** (1.70+)
  ```bash
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
  source $HOME/.cargo/env
  ```

- **Go** (1.21+)
  ```bash
  # Download from: https://golang.org/dl/
  # Or on Ubuntu/Debian:
  sudo apt install golang-go
  ```

- **GCC/Clang** (for CGO)
  ```bash
  sudo apt install build-essential
  ```

- **Make**
  ```bash
  sudo apt install make
  ```

### Optional Tools

- **Git** - Version control
- **Snapcraft** - For building Snap packages
- **Docker** - For isolated builds
- **VSCode** or **GoLand** - Recommended IDEs

## 🚀 Quick Start

### 1. Clone Repository

```bash
git clone https://github.com/devonionrouting4Moses/fileManager.git
cd fileManager
```

### 2. Build

```bash
# Build everything
make

# Or step-by-step
make rust    # Build Rust library
make go      # Build Go binary
```

### 3. Run

```bash
# Run directly
make run

# Or execute the binary
./filemanager
```

### 4. Install Locally

```bash
sudo make install
```

## 🛠️ Project Structure

```
filemanager/
├── main.go                    # CLI entry point and menu
├── operations.go              # Go wrappers for Rust FFI
├── helper.go                  # Helper functions
├── templates.go               # UI templates
├── webserver.go               # Web UI (if applicable)
├── version.go                 # Version management and updates
├── go.mod                     # Go dependencies
│
├── src/
│   └── lib.rs                 # Rust FFI functions
├── Cargo.toml                 # Rust configuration
├── Cargo.lock                 # Rust dependency lock
│
├── snap/
│   ├── snapcraft.yaml         # Snap package manifest
│   └── local/
│       ├── filemanager.desktop
│       └── filemanager.png
│
├── debian-package-build.sh    # DEB package builder
├── setup-packaging.sh         # Packaging setup script
├── fix-snap-build.sh          # Snap build fixer
│
├── Makefile                   # Build automation
├── README.md                  # User documentation
├── DEVELOPMENT.md             # This file
├── INSTALL.md                 # Installation guide
├── PACKAGING.md               # Publishing guide
└── LICENSE                    # MIT License
```

## 🔨 Build System

### Makefile Targets

```bash
make all        # Build everything (default)
make rust       # Build Rust library only
make go         # Build Go binary only
make run        # Build and run
make install    # Install system-wide
make clean      # Remove build artifacts
make test       # Run tests
make dev        # Development build (with debug symbols)
make release    # Release build (optimized)
make help       # Show all targets
```

### Build Flags

**Rust (Cargo.toml)**
```toml
[profile.release]
opt-level = 3           # Maximum optimization
lto = true             # Link-time optimization
codegen-units = 1      # Better optimization, slower build
strip = true           # Remove debug symbols
```

**Go (operations.go)**
```go
// #cgo LDFLAGS: -L${SRCDIR}/target/release -lfilemanager -ldl
// #cgo linux LDFLAGS: -Wl,-rpath,${ORIGIN}/../lib
```

### Environment Variables

```bash
# Rust
export CARGO_BUILD_JOBS=8      # Parallel build jobs
export RUSTFLAGS="-C target-cpu=native"  # CPU-specific optimizations

# Go
export CGO_ENABLED=1           # Enable CGO
export GOOS=linux              # Target OS
export GOARCH=amd64           # Target architecture
```

## 🧪 Testing

### Run Tests

```bash
# All tests
make test

# Rust tests only
cd src && cargo test

# Go tests only
go test -v ./...
```

### Manual Testing

```bash
# Build and run
make run

# Test each operation:
# 1. Create folder: /tmp/test-folder
# 2. Create file: /tmp/test-file.txt
# 3. Rename: /tmp/test-file.txt -> /tmp/renamed.txt
# 4. Copy: /tmp/renamed.txt -> /tmp/copy.txt
# 5. Move: /tmp/copy.txt -> /tmp/test-folder/
# 6. Change permissions: /tmp/renamed.txt (755)
# 7. Delete: /tmp/test-folder/ (with confirmation)
```

### Integration Tests

```bash
# Create test script
cat > test.sh << 'EOF'
#!/bin/bash
set -e

echo "Creating test folder..."
./filemanager <<< "1
/tmp/filemanager-test
8"

echo "Creating test file..."
./filemanager <<< "2
/tmp/filemanager-test/test.txt
8"

echo "Tests passed!"
rm -rf /tmp/filemanager-test
EOF

chmod +x test.sh
./test.sh
```

## 📦 Packaging

### Build DEB Package

```bash
# Run the build script
chmod +x debian-package-build.sh
./debian-package-build.sh

# Output: filemanager_0.1.2_amd64.deb

# Test installation
sudo dpkg -i filemanager_0.1.2_amd64.deb
filemanager --version
sudo apt remove filemanager
```

### Build Snap Package

```bash
# Fix any build issues first
./fix-snap-build.sh

# Build snap
snapcraft

# Or in clean environment
snapcraft --use-lxd

# Test installation
sudo snap install --dangerous filemanager_0.1.2_amd64.snap
filemanager --version
sudo snap remove filemanager
```

### Cross-Compilation

**Linux ARM64 from AMD64:**
```bash
# Install cross-compilation tools
rustup target add aarch64-unknown-linux-gnu
sudo apt install gcc-aarch64-linux-gnu

# Build
export CARGO_TARGET_AARCH64_UNKNOWN_LINUX_GNU_LINKER=aarch64-linux-gnu-gcc
cargo build --release --target aarch64-unknown-linux-gnu

# Build Go binary
export GOARCH=arm64
export CC=aarch64-linux-gnu-gcc
go build -o filemanager-arm64
```

**Windows from Linux:**
```bash
# Install mingw
sudo apt install mingw-w64

# Add Rust Windows target
rustup target add x86_64-pc-windows-gnu

# Build Rust
cargo build --release --target x86_64-pc-windows-gnu

# Build Go
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
export CC=x86_64-w64-mingw32-gcc
go build -o filemanager.exe
```

## 🐛 Debugging

### Enable Debug Logging

```go
// In main.go
import "log"

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    log.Println("Starting FileManager...")
    // ... rest of code
}
```

### Rust Debug Build

```bash
# Build with debug symbols
cargo build  # No --release flag

# Run with backtrace
RUST_BACKTRACE=1 ./filemanager
```

### Go Debug with Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug
```

### GDB Debugging

```bash
# Build with debug symbols
make dev

# Debug with GDB
gdb ./filemanager
(gdb) break main.main
(gdb) run
(gdb) backtrace
```

## 🔧 Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/my-feature
```

### 2. Make Changes

Edit code in your preferred IDE:
- **VSCode**: Install Rust Analyzer + Go extensions
- **GoLand**: Native Go support + Rust plugin
- **Vim/Neovim**: rust.vim + vim-go

### 3. Test Changes

```bash
make clean
make test
make run
```

### 4. Format Code

```bash
# Rust
cargo fmt

# Go
go fmt ./...
gofmt -s -w .
```

### 5. Lint Code

```bash
# Rust
cargo clippy

# Go
go vet ./...

# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
golangci-lint run
```

### 6. Commit Changes

```bash
git add .
git commit -m "feat: add new feature"
git push origin feature/my-feature
```

### 7. Create Pull Request

Go to GitHub and create a PR from your branch.

## 📝 Code Style Guidelines

### Rust Style

```rust
// Use descriptive names
fn create_folder_safe(path: &str) -> Result<(), String>

// Document public functions
/// Creates a folder at the specified path
/// 
/// # Arguments
/// * `path` - The path where the folder should be created
/// 
/// # Returns
/// * `Ok(())` on success
/// * `Err(String)` with error message on failure
```

### Go Style

```go
// Use camelCase for unexported
func createFolder(path string) error

// Use PascalCase for exported
func CreateFolder(path string) error

// Document exported functions
// CreateFolder creates a new directory at the specified path.
// It returns an error if the operation fails.
func CreateFolder(path string) error
```

## 🚢 Release Process

### 1. Update Version

```bash
# Update version in:
# - version.go: const Version = "0.1.3"
# - snap/snapcraft.yaml: version: '0.1.3'
# - debian-package-build.sh: VERSION="0.1.3"
# - Cargo.toml: version = "0.1.3"
```

### 2. Update Changelog

```bash
# Add to CHANGELOG.md
## [0.1.3] - 2025-11-15
### Added
- New feature X
### Fixed
- Bug Y
### Changed
- Improvement Z
```

### 3. Create Git Tag

```bash
git add .
git commit -m "Release v0.1.3"
git tag -a v0.1.3 -m "Release version 0.1.3"
git push origin main
git push origin v0.1.3
```

### 4. GitHub Actions Build

GitHub Actions will automatically:
- Build binaries for all platforms
- Create DEB and Snap packages
- Upload to GitHub Releases
- Publish Snap to store

### 5. Manual Publish (if needed)

```bash
# Build DEB
./debian-package-build.sh

# Build Snap
snapcraft

# Upload to Snap Store
snapcraft upload filemanager_0.1.3_amd64.snap --release=stable

# Upload to GitHub Release
gh release upload v0.1.3 filemanager_0.1.3_amd64.deb
gh release upload v0.1.3 filemanager_0.1.3_amd64.snap
```

## 🤝 Contributing

We welcome contributions! Please:

1. Check [open issues](https://github.com/devonionrouting4Moses/fileManager/issues)
2. Follow code style guidelines
3. Add tests for new features
4. Update documentation
5. Create detailed PR descriptions

## 📚 Additional Resources

- [Rust Book](https://doc.rust-lang.org/book/)
- [Go Documentation](https://go.dev/doc/)
- [CGO Documentation](https://pkg.go.dev/cmd/cgo)
- [Snapcraft Documentation](https://snapcraft.io/docs)
- [Debian Packaging Guide](https://www.debian.org/doc/manuals/maint-guide/)

## 🆘 Getting Help

- **Discussions**: [GitHub Discussions](https://github.com/devonionrouting4Moses/fileManager/discussions)
- **Issues**: [GitHub Issues](https://github.com/devonionrouting4Moses/fileManager/issues)
- **Email**: moses.muranja@strathmore.edu

---

Happy coding!🚀 From Moses Muranja ❤️ 