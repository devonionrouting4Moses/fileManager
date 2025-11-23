# Complete Setup and Integration Guide

## Quick Start

```bash
# Clone the repository
git clone <your-repo-url> file_manager
cd file_manager

# Build everything
make all

# Run the application
make run
```

## Detailed Setup

### 1. Install Prerequisites

#### Rust
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
rustc --version  # Verify: rustc 1.70+ recommended
```

#### Go
```bash
# Download from https://golang.org/dl/
# Or use your package manager

# Verify
go version  # Verify: go 1.21+ recommended
```

#### Build Tools

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get update
sudo apt-get install build-essential pkg-config
```

**Linux (Fedora/RHEL):**
```bash
sudo dnf install gcc make
```

**macOS:**
```bash
xcode-select --install
```

**Windows:**
```bash
# Install MinGW-w64 from https://www.mingw-w64.org/
# Or use chocolatey:
choco install mingw
```

### 2. Project Structure Setup

```bash
# Create the integrated structure
mkdir -p file_manager/{cmd/app,internal/{ffi,handler,service,repository},pkg/{version,utils},api,configs,docs,scripts,filemanager_frontend}

# Initialize Go module
cd file_manager
go mod init file_manager
go mod tidy
```

### 3. Set Up Rust Workspace

```bash
# Create Rust workspace
mkdir -p rust_ffi/crates/{core/src,cli/src}

# Copy Rust files from artifacts:
# - rust_ffi/Cargo.toml (workspace config)
# - rust_ffi/crates/core/Cargo.toml
# - rust_ffi/crates/core/src/ (all Rust source files)
```

### 4. Build Process

#### Option A: Using Makefile (Recommended)

```bash
# Build everything
make all

# Or step by step:
make rust  # Build Rust library
make go    # Build Go binary

# Development build with debug symbols
make dev
```

#### Option B: Manual Build

```bash
# 1. Build Rust library
cd rust_ffi
cargo build --release -p fs-operations-core
cd ..

# 2. Build Go application
CGO_ENABLED=1 \
CGO_LDFLAGS="-L./rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
go build -o filemanager ./cmd/app
```

### 5. Running the Application

#### Interactive CLI Mode
```bash
# With library path set
export LD_LIBRARY_PATH=$PWD/rust_ffi/target/release:$LD_LIBRARY_PATH
./filemanager

# Or use the helper script
./filemanager.sh

# Or use make
make run
```

#### Web Server Mode
```bash
./filemanager --web
# Or
make web
```

### 6. Installation

#### System-wide Installation
```bash
sudo make install
```

This installs:
- Binary to `/usr/local/bin/filemanager`
- Library to `/usr/local/lib/`
- Updates library cache (Linux)

#### Uninstall
```bash
sudo make uninstall
```

## Platform-Specific Setup

### Linux

#### Ubuntu/Debian Package
```bash
# Build DEB package
./scripts/debian-package-build.sh

# Install
sudo dpkg -i dist/filemanager_*.deb
```

#### Snap Package
```bash
# Build snap
snapcraft

# Install locally
sudo snap install filemanager_*.snap --dangerous
```

### macOS

#### Homebrew (Future)
```bash
# Coming soon
brew install filemanager
```

#### Manual Installation
```bash
# Build
make all

# Install
sudo make install

# May need to allow in System Preferences > Security & Privacy
```

### Windows

#### Build with MinGW
```bash
# In Git Bash or WSL
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
export CC=x86_64-w64-mingw32-gcc

# Build
make all

# Run
filemanager.exe
```

## Development Workflow

### 1. Modify Rust Code

```bash
# Navigate to Rust workspace
cd rust_ffi/crates/core

# Edit files in src/operations/

# Run Rust tests
cargo test

# Build
cargo build --release
```

### 2. Modify Go Code

```bash
# Edit files in internal/, cmd/, pkg/

# Run Go tests
go test ./...

# Format code
go fmt ./...

# Rebuild
make go
```

### 3. Testing

```bash
# Run all tests
make test

# Run only Rust tests
cd rust_ffi && cargo test

# Run only Go tests
go test -v ./...

# Run benchmarks
make bench
```

### 4. Code Quality

```bash
# Format all code
make fmt

# Lint all code
make lint

# Generate documentation
make docs
```

## Environment Variables

### Build Time

```bash
# Go build
export CGO_ENABLED=1
export CGO_LDFLAGS="-L./rust_ffi/target/release -lfs_operations_core"

# Cross-compilation
export GOOS=linux      # or darwin, windows
export GOARCH=amd64    # or arm64, 386
```

### Runtime

```bash
# Linux
export LD_LIBRARY_PATH=$PWD/rust_ffi/target/release:$LD_LIBRARY_PATH

# macOS
export DYLD_LIBRARY_PATH=$PWD/rust_ffi/target/release:$DYLD_LIBRARY_PATH

# Windows (PowerShell)
$env:PATH="$PWD\rust_ffi\target\release;$env:PATH"
```

## File Structure Reference

```
file_manager/
├── rust_ffi/                    # Rust FFI library
│   ├── Cargo.toml               # Workspace config
│   └── crates/
│       └── core/                # Core library
│           ├── Cargo.toml
│           └── src/
│               ├── lib.rs       # Entry point
│               ├── common/      # Shared types
│               ├── ffi/         # C exports
│               └── operations/  # File operations
│
├── cmd/app/main.go              # Go application entry
├── internal/
│   ├── ffi/operations.go        # Go FFI bindings
│   ├── handler/                 # HTTP handlers
│   ├── service/                 # Business logic
│   └── repository/              # Data access
├── pkg/
│   ├── version/                 # Version info
│   └── utils/                   # Utilities
│
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
├── README.md                    # Main docs
└── INTEGRATION.md               # This file
```

## Common Issues

### Issue: "Library not found"

**Solution:**
```bash
# Add to library path
export LD_LIBRARY_PATH=$PWD/rust_ffi/target/release:$LD_LIBRARY_PATH

# Or install system-wide
sudo make install
```

### Issue: "CGO not enabled"

**Solution:**
```bash
export CGO_ENABLED=1
make clean
make all
```

### Issue: "Undefined symbol"

**Solution:**
```bash
# Rebuild Rust library
cd rust_ffi
cargo clean
cargo build --release

# Rebuild Go app
cd ..
make clean
make all
```

### Issue: Windows DLL not found

**Solution:**
```bash
# Copy DLL to same directory as exe
copy rust_ffi\target\release\fs_operations_core.dll .

# Or add to PATH
set PATH=%CD%\rust_ffi\target\release;%PATH%
```

## CI/CD Setup

### GitHub Actions

See `.github/workflows/release.yml` for automated builds on:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Local Multi-Platform Build

```bash
# Build for all platforms
make build-all

# Creates packages in dist/
# - filemanager-VERSION-linux-amd64.tar.gz
# - filemanager-VERSION-darwin-amd64.tar.gz
# - filemanager-VERSION-windows-amd64.zip
```

## Performance Tuning

### Release Builds
```bash
# Always use release builds for production
cargo build --release
go build -ldflags="-s -w"  # Strip debug info
```

### Link-Time Optimization
```toml
# rust_ffi/Cargo.toml
[profile.release]
lto = true
codegen-units = 1
opt-level = 3
```

### Static Linking (Linux)
```bash
# For portable binaries
export CGO_ENABLED=1
export CGO_LDFLAGS="-L./rust_ffi/target/release -static -lfs_operations_core"
go build -ldflags '-extldflags "-static"' ./cmd/app
```

## Next Steps

1. **Read INTEGRATION.md** - Understand the architecture
2. **Run Tests** - Verify everything works
3. **Try Examples** - Use the CLI and web interface
4. **Contribute** - See CONTRIBUTING.md
5. **Deploy** - Package for your platform

## Resources

- [Rust FFI Book](https://doc.rust-lang.org/nomicon/ffi.html)
- [CGO Documentation](https://pkg.go.dev/cmd/cgo)
- [Project Repository](https://github.com/yourusername/file_manager)
- [Issue Tracker](https://github.com/yourusername/file_manager/issues)

## Getting Help

- **Issues:** Open a GitHub issue
- **Discussions:** Use GitHub Discussions
- **Email:** your.email@example.com
- **Docs:** Read the wiki

---

**Happy Coding! 🚀**