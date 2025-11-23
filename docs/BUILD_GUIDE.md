# 🔨 FileManager v2 Build Guide

Complete guide for building FileManager v2 for all supported platforms.

---

## 📋 Prerequisites

### All Platforms
- Go 1.24 or later
- Rust 1.70 or later
- Git

### Linux (amd64 & arm64)
```bash
# Install build tools
sudo apt-get update
sudo apt-get install -y build-essential pkg-config libssl-dev

# For cross-compilation to ARM64 from amd64
sudo apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu
```

### Windows (amd64)
- Visual Studio Build Tools or MinGW-w64
- MSYS2 (recommended)

### macOS
- Xcode Command Line Tools
- Homebrew (optional)

---

## 🏗️ Building for Linux

### Linux amd64

```bash
# Navigate to project root
cd file_manager_v2

# Build Rust library
cd rust_ffi
cargo build --release -p fs-operations-core
cd ..

# Build Go binary
cd file_manager
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
  CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
  go build -ldflags="-s -w" -o ../filemanager ./cmd/app
cd ..

# Verify binary
./filemanager --version
```

### Linux ARM64 (from amd64 machine)

```bash
# Navigate to project root
cd file_manager_v2

# Build Rust library for ARM64
cd rust_ffi
cargo build --release --target aarch64-unknown-linux-gnu -p fs-operations-core
cd ..

# Build Go binary for ARM64
cd file_manager
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
  CC=aarch64-linux-gnu-gcc \
  CXX=aarch64-linux-gnu-g++ \
  CGO_LDFLAGS="-L../rust_ffi/target/aarch64-unknown-linux-gnu/release -lfs_operations_core -ldl -lpthread -lm" \
  go build -ldflags="-s -w" -o ../filemanager-arm64 ./cmd/app
cd ..

# Verify binary (will fail if not on ARM64, but file is created)
file filemanager-arm64
```

### Linux Packages

#### 1. DEB Package (Debian/Ubuntu)

```bash
# Create package structure
mkdir -p deb/DEBIAN
mkdir -p deb/usr/local/bin
mkdir -p deb/usr/local/lib

# Copy binaries
cp filemanager deb/usr/local/bin/
cp rust_ffi/target/release/libfs_operations_core.so deb/usr/local/lib/
chmod +x deb/usr/local/bin/filemanager

# Create control file
cat > deb/DEBIAN/control << 'EOF'
Package: filemanager
Version: 0.1.2
Architecture: amd64
Maintainer: DevChigarlicMoses <moses.muranja@strathmore.edu>
Description: A powerful file manager with Rust backend and Go frontend
 FileManager v2 is a dual-mode file manager with terminal and web interfaces.
 Built with Rust for performance and Go for flexibility.
Homepage: https://github.com/DevChigarlicMoses/FileManager
EOF

# Create postinst script
cat > deb/DEBIAN/postinst << 'EOF'
#!/bin/bash
ldconfig
EOF
chmod +x deb/DEBIAN/postinst

# Build DEB package
dpkg-deb --build deb filemanager_0.1.2_amd64.deb
```

#### 2. TAR.GZ Package

```bash
# Create distribution directory
mkdir -p dist/filemanager-0.1.2-linux-amd64
cp filemanager dist/filemanager-0.1.2-linux-amd64/
cp rust_ffi/target/release/libfs_operations_core.so dist/filemanager-0.1.2-linux-amd64/
cp README.md dist/filemanager-0.1.2-linux-amd64/
cp LICENSE dist/filemanager-0.1.2-linux-amd64/

# Create archive
cd dist
tar czf filemanager-0.1.2-linux-amd64.tar.gz filemanager-0.1.2-linux-amd64/
cd ..
```

#### 3. Snap Package

```bash
# Create snapcraft.yaml (already exists in snap/snapcraft.yaml)
# Build snap
snapcraft

# The snap package will be created as filemanager_0.1.2_amd64.snap
```

#### 4. APT Repository Package

```bash
# This is typically done on a repository server
# For now, the DEB package can be installed directly:
sudo dpkg -i filemanager_0.1.2_amd64.deb
```

#### 5. Arch Linux Package

```bash
# Create PKGBUILD
mkdir -p arch
cat > arch/PKGBUILD << 'EOF'
pkgname=filemanager
pkgver=0.1.2
pkgrel=1
pkgdesc="A powerful file manager with Rust backend and Go frontend"
arch=('x86_64' 'aarch64')
url="https://github.com/DevChigarlicMoses/FileManager"
license=('MIT' 'Apache-2.0')
depends=('glibc')
makedepends=('go' 'rust')

build() {
    cd "$srcdir"
    
    # Build Rust library
    cd rust_ffi
    cargo build --release -p fs-operations-core
    cd ..
    
    # Build Go binary
    cd file_manager
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
      CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
      go build -ldflags="-s -w" -o ../filemanager ./cmd/app
    cd ..
}

package() {
    install -Dm755 filemanager "$pkgdir/usr/bin/filemanager"
    install -Dm755 rust_ffi/target/release/libfs_operations_core.so "$pkgdir/usr/lib/libfs_operations_core.so"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
EOF

# Build package
cd arch
makepkg -si
```

#### 6. HarmonyOS Package

```bash
# HarmonyOS uses OpenHarmony SDK
# This requires the OpenHarmony SDK to be installed

# Create HAP (HarmonyOS Ability Package) structure
mkdir -p harmonyos/entry/src/main

# Note: HarmonyOS requires specific toolchain and SDK
# This is a placeholder for future implementation
# Full implementation requires OpenHarmony SDK setup
```

---

## 🪟 Building for Windows

### Windows amd64 (.exe)

#### Using MSYS2 (Recommended)

```bash
# Install MSYS2 from https://www.msys2.org/

# Open MSYS2 MinGW 64-bit terminal
# Install dependencies
pacman -S mingw-w64-x86_64-toolchain mingw-w64-x86_64-rust

# Navigate to project
cd file_manager_v2

# Build Rust library
cd rust_ffi
cargo build --release -p fs-operations-core
cd ..

# Build Go binary
cd file_manager
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
set CGO_LDFLAGS=-L../rust_ffi/target/release -lfs_operations_core -lws2_32 -luserenv
go build -ldflags="-s -w" -o ../filemanager.exe ./cmd/app
cd ..

# Verify binary
filemanager.exe --version
```

#### Using Visual Studio Build Tools

```bash
# Open Visual Studio Developer Command Prompt

# Navigate to project
cd file_manager_v2

# Build Rust library
cd rust_ffi
cargo build --release -p fs-operations-core
cd ..

# Build Go binary
cd file_manager
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
set CGO_LDFLAGS=-L../rust_ffi/target/release -lfs_operations_core -lws2_32 -luserenv
go build -ldflags="-s -w" -o ../filemanager.exe ./cmd/app
cd ..
```

### Windows Installer

```batch
@echo off
REM Create installer directory structure
mkdir installer\bin
mkdir installer\lib

REM Copy binaries
copy filemanager.exe installer\bin\
copy rust_ffi\target\release\fs_operations_core.dll installer\lib\

REM Create install script (install.bat)
REM Already created in build process

REM Create ZIP archive
REM Use 7-Zip or similar tool
7z a filemanager-0.1.2-windows-amd64.zip installer\
```

---

## 🍎 Building for macOS

### macOS amd64 (Intel)

```bash
# Install Xcode Command Line Tools (if not already installed)
xcode-select --install

# Navigate to project
cd file_manager_v2

# Build Rust library
cd rust_ffi
cargo build --release -p fs-operations-core
cd ..

# Build Go binary
cd file_manager
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
  CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core" \
  go build -ldflags="-s -w" -o ../filemanager ./cmd/app
cd ..

# Verify binary
./filemanager --version
```

### macOS arm64 (Apple Silicon)

```bash
# Navigate to project
cd file_manager_v2

# Build Rust library for ARM64
cd rust_ffi
cargo build --release --target aarch64-apple-darwin -p fs-operations-core
cd ..

# Build Go binary for ARM64
cd file_manager
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
  CGO_LDFLAGS="-L../rust_ffi/target/aarch64-apple-darwin/release -lfs_operations_core" \
  go build -ldflags="-s -w" -o ../filemanager-arm64 ./cmd/app
cd ..
```

---

## 🔧 Troubleshooting Build Issues

### Issue: "cannot find -lfs_operations_core"

**Solution:**
1. Verify Rust library was built: `ls rust_ffi/target/release/libfs_operations_core.*`
2. Check CGO_LDFLAGS path is correct
3. On Linux, run: `ldconfig -p | grep fs_operations_core`

### Issue: "CGO_ENABLED not recognized"

**Solution:**
- On Windows, use `set` instead of `export`
- On macOS/Linux, use `export`

### Issue: Cross-compilation ARM64 fails

**Solution:**
1. Install ARM64 toolchain: `sudo apt-get install gcc-aarch64-linux-gnu`
2. Set CC and CXX environment variables
3. Use `--target aarch64-unknown-linux-gnu` for Rust

### Issue: Rust library not found at runtime

**Solution:**
1. Set `LD_LIBRARY_PATH` before running:
   ```bash
   export LD_LIBRARY_PATH=./rust_ffi/target/release:$LD_LIBRARY_PATH
   ./filemanager
   ```
2. Or copy library to system path:
   ```bash
   sudo cp rust_ffi/target/release/libfs_operations_core.so /usr/local/lib/
   sudo ldconfig
   ```

---

## 📦 Automated Build Script

```bash
#!/bin/bash
# build-all.sh - Build for all platforms

set -e

VERSION="0.1.2"
PROJECT_ROOT=$(pwd)

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Building FileManager v${VERSION}${NC}"

# Build Rust library once
echo -e "${BLUE}Building Rust library...${NC}"
cd rust_ffi
cargo build --release -p fs-operations-core
cd ..

# Linux amd64
echo -e "${BLUE}Building Linux amd64...${NC}"
cd file_manager
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
  CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
  go build -ldflags="-s -w" -o ../filemanager-linux-amd64 ./cmd/app
cd ..
echo -e "${GREEN}✓ Linux amd64 built${NC}"

# Linux arm64
echo -e "${BLUE}Building Linux arm64...${NC}"
cd rust_ffi
cargo build --release --target aarch64-unknown-linux-gnu -p fs-operations-core
cd ..
cd file_manager
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
  CC=aarch64-linux-gnu-gcc \
  CXX=aarch64-linux-gnu-g++ \
  CGO_LDFLAGS="-L../rust_ffi/target/aarch64-unknown-linux-gnu/release -lfs_operations_core -ldl -lpthread -lm" \
  go build -ldflags="-s -w" -o ../filemanager-linux-arm64 ./cmd/app
cd ..
echo -e "${GREEN}✓ Linux arm64 built${NC}"

# Windows amd64
echo -e "${BLUE}Building Windows amd64...${NC}"
cd file_manager
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -lws2_32 -luserenv" \
  go build -ldflags="-s -w" -o ../filemanager.exe ./cmd/app
cd ..
echo -e "${GREEN}✓ Windows amd64 built${NC}"

echo -e "${GREEN}All builds completed!${NC}"
```

---

## 📝 Summary

### Supported Platforms

| Platform | Status | Binary | Package |
|----------|--------|--------|---------|
| Linux amd64 | ✅ Ready | filemanager | .deb, .tar.gz, .snap |
| Linux arm64 | ✅ Ready | filemanager | .deb, .tar.gz |
| Windows amd64 | ✅ Ready | filemanager.exe | .zip, .exe installer |
| macOS amd64 | ⏳ Later | filemanager | .tar.gz, .dmg |
| macOS arm64 | ⏳ Later | filemanager | .tar.gz, .dmg |
| HarmonyOS | 📋 Planned | - | .hap |

### Next Steps

1. **Test builds locally** before pushing to GitHub
2. **Update GitHub Actions workflow** with proper cross-compilation setup
3. **Create release packages** for each platform
4. **Test installation** on target systems
5. **Document platform-specific requirements**

---

**Happy Building! 🚀**
