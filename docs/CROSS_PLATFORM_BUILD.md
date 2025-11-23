# 🔨 FileManager v1 - Cross-Platform Build Guide

Complete guide for building FileManager v1 for all supported platforms and package formats.

---

## 📋 Overview

### Supported Platforms & Architectures

| Platform | Architecture | Status | Binary | Packages |
|----------|--------------|--------|--------|----------|
| **Linux** | amd64 | ✅ Ready | filemanager | .deb, .tar.gz, .snap, apt, PKGBUILD |
| **Linux** | arm64 | ✅ Ready | filemanager | .deb, .tar.gz, .snap, apt, PKGBUILD |
| **Windows** | amd64 | ✅ Ready | filemanager.exe | .zip, .exe installer |
| **macOS** | amd64 | ⏳ Later | filemanager | .tar.gz, .dmg |
| **macOS** | arm64 | ⏳ Later | filemanager | .tar.gz, .dmg |
| **HarmonyOS** | arm64 | 📋 Planned | - | .hap |

---

## 🐧 Linux Builds (amd64 & arm64)

### Prerequisites

```bash
# Install build tools
sudo apt-get update
sudo apt-get install -y build-essential pkg-config libssl-dev

# For cross-compilation to ARM64 from amd64
sudo apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu

# Install Go and Rust
# Go: https://golang.org/dl/
# Rust: https://rustup.rs/
```

### Build Linux amd64

```bash
cd file_manager_v1

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

# Verify
./filemanager --version
```

### Build Linux arm64 (Cross-compilation)

```bash
cd file_manager_v1

# Build Rust library for ARM64
cd rust_ffi
rustup target add aarch64-unknown-linux-gnu
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

# Verify (will fail if not on ARM64, but file is created)
file filemanager-arm64
```

---

## 📦 Linux Package Formats

### 1. DEB Package (Debian/Ubuntu)

```bash
#!/bin/bash
# Create DEB package structure

VERSION="2.0.0"
ARCH="amd64"

# Create directory structure
mkdir -p deb/DEBIAN
mkdir -p deb/usr/local/bin
mkdir -p deb/usr/local/lib
mkdir -p deb/usr/share/doc/filemanager

# Copy binaries
cp filemanager deb/usr/local/bin/
cp rust_ffi/target/release/libfs_operations_core.so deb/usr/local/lib/
chmod +x deb/usr/local/bin/filemanager

# Create control file
cat > deb/DEBIAN/control << 'EOF'
Package: filemanager
Version: 2.0.0
Architecture: amd64
Maintainer: DevChigarlicMoses <moses.muranja@strathmore.edu>
Description: Modern file manager with Rust+Go backend
 FileManager v2 is a dual-mode file manager with terminal and web interfaces.
 Built with Rust for performance and Go for flexibility.
Homepage: https://github.com/devonionrouting4Moses/fileManager
Depends: libc6, libgcc-s1, libstdc++6
EOF

# Create postinst script (runs after installation)
cat > deb/DEBIAN/postinst << 'EOF'
#!/bin/bash
set -e
# Update library cache
ldconfig
# Create user config directory if needed
mkdir -p ~/.config/filemanager
echo "FileManager installed successfully!"
EOF
chmod +x deb/DEBIAN/postinst

# Create prerm script (runs before removal)
cat > deb/DEBIAN/prerm << 'EOF'
#!/bin/bash
set -e
echo "Removing FileManager..."
EOF
chmod +x deb/DEBIAN/prerm

# Build DEB package
dpkg-deb --build deb filemanager_${VERSION}_${ARCH}.deb

# Verify
dpkg -c filemanager_${VERSION}_${ARCH}.deb
```

### 2. TAR.GZ Archive (Universal)

```bash
#!/bin/bash
# Create TAR.GZ distribution

VERSION="2.0.0"
ARCH="amd64"

# Create distribution directory
mkdir -p dist/filemanager-${VERSION}-linux-${ARCH}
cp filemanager dist/filemanager-${VERSION}-linux-${ARCH}/
cp rust_ffi/target/release/libfs_operations_core.so dist/filemanager-${VERSION}-linux-${ARCH}/
cp README.md dist/filemanager-${VERSION}-linux-${ARCH}/
cp LICENSE dist/filemanager-${VERSION}-linux-${ARCH}/

# Create install script
cat > dist/filemanager-${VERSION}-linux-${ARCH}/install.sh << 'EOF'
#!/bin/bash
set -e
echo "Installing FileManager..."

# Install binary
sudo cp filemanager /usr/local/bin/
sudo chmod +x /usr/local/bin/filemanager

# Install library
sudo cp libfs_operations_core.so /usr/local/lib/
sudo ldconfig

echo "✅ FileManager installed successfully!"
echo "Run 'filemanager' to start"
EOF
chmod +x dist/filemanager-${VERSION}-linux-${ARCH}/install.sh

# Create archive
cd dist
tar czf filemanager-${VERSION}-linux-${ARCH}.tar.gz filemanager-${VERSION}-linux-${ARCH}/
cd ..

echo "✅ Created: filemanager-${VERSION}-linux-${ARCH}.tar.gz"
```

### 3. Snap Package

```bash
# Build snap (see SNAP_BUILD_GUIDE.md for details)
snapcraft

# Result: filemanager_2.0.0_amd64.snap
```

### 4. APT Repository Package

```bash
# The DEB package can be added to an APT repository
# For now, the DEB package can be installed directly:

sudo dpkg -i filemanager_2.0.0_amd64.deb

# Or with apt:
sudo apt install ./filemanager_2.0.0_amd64.deb
```

### 5. Arch Linux Package (PKGBUILD)

```bash
#!/bin/bash
# Create PKGBUILD for Arch Linux

mkdir -p arch
cat > arch/PKGBUILD << 'EOF'
pkgname=filemanager
pkgver=2.0.0
pkgrel=1
pkgdesc="Modern file manager with Rust+Go backend"
arch=('x86_64' 'aarch64')
url="https://github.com/devonionrouting4Moses/fileManager"
license=('MIT' 'Apache-2.0')
depends=('glibc' 'gcc-libs')
makedepends=('go' 'rust')
source=("https://github.com/devonionrouting4Moses/fileManager/archive/v${pkgver}.tar.gz")

build() {
    cd "$srcdir/FileManager-${pkgver}"
    
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
    cd "$srcdir/FileManager-${pkgver}"
    
    # Install binary
    install -Dm755 filemanager "$pkgdir/usr/bin/filemanager"
    
    # Install library
    install -Dm755 rust_ffi/target/release/libfs_operations_core.so "$pkgdir/usr/lib/libfs_operations_core.so"
    
    # Install documentation
    install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
EOF

# Build package
cd arch
makepkg -si
```

---

## 🪟 Windows Build (amd64)

### Prerequisites

#### Option A: Using MSYS2 (Recommended)

```bash
# Download MSYS2 from https://www.msys2.org/
# Install and run MSYS2 MinGW 64-bit terminal

# Install dependencies
pacman -S mingw-w64-x86_64-toolchain mingw-w64-x86_64-rust mingw-w64-x86_64-go

# Verify installations
rustc --version
go version
gcc --version
```

#### Option B: Using Visual Studio Build Tools

```bash
# Install Visual Studio Build Tools
# Download from: https://visualstudio.microsoft.com/downloads/
# Select "Desktop development with C++"

# Install Rust and Go separately
# Rust: https://rustup.rs/
# Go: https://golang.org/dl/
```

### Build Windows amd64

```bash
# Open MSYS2 MinGW 64-bit terminal

cd file_manager_v1

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

# Verify
filemanager.exe --version
```

---

## 📦 Windows Package Formats

### 1. ZIP Archive

```bash
#!/bin/bash
# Create ZIP distribution

VERSION="2.0.0"

# Create distribution directory
mkdir -p dist/filemanager-${VERSION}-windows-amd64
cp filemanager.exe dist/filemanager-${VERSION}-windows-amd64/
cp rust_ffi/target/release/fs_operations_core.dll dist/filemanager-${VERSION}-windows-amd64/
cp README.md dist/filemanager-${VERSION}-windows-amd64/
cp LICENSE dist/filemanager-${VERSION}-windows-amd64/

# Create install script
cat > dist/filemanager-${VERSION}-windows-amd64/install.bat << 'EOF'
@echo off
REM FileManager Installation Script

echo Installing FileManager...
set INSTALL_DIR=%ProgramFiles%\FileManager

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copy executable
copy /Y filemanager.exe "%INSTALL_DIR%\"

REM Copy Rust library
copy /Y fs_operations_core.dll "%INSTALL_DIR%\"

REM Add to PATH
setx /M PATH "%PATH%;%INSTALL_DIR%"

echo.
echo ✅ Installation complete!
echo FileManager installed to: %INSTALL_DIR%
echo.
echo Run 'filemanager' from command prompt to start
pause
EOF

# Create ZIP archive
cd dist
7z a filemanager-${VERSION}-windows-amd64.zip filemanager-${VERSION}-windows-amd64/
cd ..

echo "✅ Created: filemanager-${VERSION}-windows-amd64.zip"
```

### 2. Windows Installer (.exe)

```bash
# Using NSIS (Nullsoft Scriptable Install System)
# Download from: https://nsis.sourceforge.io/

# Create installer script (filemanager.nsi)
cat > installer/filemanager.nsi << 'EOF'
; FileManager Installer Script

!include "MUI2.nsh"

Name "FileManager v2.0.0"
OutFile "filemanager-2.0.0-installer.exe"
InstallDir "$PROGRAMFILES\FileManager"

!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

Section "Install"
    SetOutPath "$INSTDIR"
    File "filemanager.exe"
    File "fs_operations_core.dll"
    File "README.md"
    File "LICENSE"
    
    ; Add to PATH
    EnvVarUpdate::AddValue "HKLM" "System\CurrentControlSet\Control\Session Manager\Environment" "PATH" "$INSTDIR"
    
    ; Create Start Menu shortcut
    CreateDirectory "$SMPROGRAMS\FileManager"
    CreateShortcut "$SMPROGRAMS\FileManager\FileManager.lnk" "$INSTDIR\filemanager.exe"
SectionEnd

Section "Uninstall"
    Delete "$INSTDIR\filemanager.exe"
    Delete "$INSTDIR\fs_operations_core.dll"
    Delete "$INSTDIR\README.md"
    Delete "$INSTDIR\LICENSE"
    RMDir "$INSTDIR"
    
    ; Remove from PATH
    EnvVarUpdate::RemoveValue "HKLM" "System\CurrentControlSet\Control\Session Manager\Environment" "PATH" "$INSTDIR"
SectionEnd
EOF

# Build installer
makensis installer/filemanager.nsi

# Result: filemanager-2.0.0-installer.exe
```

---

## 🍎 macOS Builds (amd64 & arm64)

### Prerequisites

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Rust and Go
# Rust: https://rustup.rs/
# Go: https://golang.org/dl/
```

### Build macOS amd64 (Intel)

```bash
cd file_manager_v1

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

# Verify
./filemanager --version
```

### Build macOS arm64 (Apple Silicon)

```bash
cd file_manager_v1

# Build Rust library for ARM64
cd rust_ffi
rustup target add aarch64-apple-darwin
cargo build --release --target aarch64-apple-darwin -p fs-operations-core
cd ..

# Build Go binary for ARM64
cd file_manager
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
  CGO_LDFLAGS="-L../rust_ffi/target/aarch64-apple-darwin/release -lfs_operations_core" \
  go build -ldflags="-s -w" -o ../filemanager-arm64 ./cmd/app
cd ..
```

### macOS Package Formats

#### TAR.GZ Archive

```bash
# Same as Linux TAR.GZ, but with macOS binary
mkdir -p dist/filemanager-2.0.0-macos-amd64
cp filemanager dist/filemanager-2.0.0-macos-amd64/
cp rust_ffi/target/release/libfs_operations_core.dylib dist/filemanager-2.0.0-macos-amd64/
tar czf filemanager-2.0.0-macos-amd64.tar.gz dist/filemanager-2.0.0-macos-amd64/
```

#### DMG (Disk Image)

```bash
# Create DMG using macOS tools
mkdir -p dmg-staging
cp filemanager dmg-staging/
cp rust_ffi/target/release/libfs_operations_core.dylib dmg-staging/

# Create DMG
hdiutil create -volname "FileManager" \
  -srcfolder dmg-staging \
  -ov -format UDZO \
  filemanager-2.0.0-macos-amd64.dmg
```

#### Homebrew Formula

```bash
# Create Homebrew formula (filemanager.rb)
cat > filemanager.rb << 'EOF'
class Filemanager < Formula
  desc "Modern file manager with Rust+Go backend"
  homepage "https://github.com/devonionrouting4Moses/fileManager"
  url "https://github.com/devonionrouting4Moses/fileManager/archive/v2.0.0.tar.gz"
  sha256 "YOUR_SHA256_HERE"
  license "MIT"

  depends_on "go" => :build
  depends_on "rust" => :build

  def install
    # Build Rust library
    system "cd rust_ffi && cargo build --release -p fs-operations-core"
    
    # Build Go binary
    system "cd file_manager && GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 " \
           "CGO_LDFLAGS='-L../rust_ffi/target/release -lfs_operations_core' " \
           "go build -ldflags='-s -w' -o ../filemanager ./cmd/app"
    
    # Install binary
    bin.install "filemanager"
    
    # Install library
    lib.install "rust_ffi/target/release/libfs_operations_core.dylib"
  end

  test do
    system "#{bin}/filemanager", "--version"
  end
end
EOF

# Install locally for testing
brew install --build-from-source ./filemanager.rb
```

---

## 🌐 HarmonyOS Package (Planned)

### Prerequisites

```bash
# Install OpenHarmony SDK
# Download from: https://developer.harmonyos.com/

# Set up environment
export OHOS_SDK_HOME=/path/to/ohos-sdk
```

### Build HarmonyOS Package

```bash
# Create HAP (HarmonyOS Ability Package) structure
mkdir -p harmonyos/entry/src/main
mkdir -p harmonyos/entry/src/main/ets
mkdir -p harmonyos/entry/src/main/resources

# Create module.json5
cat > harmonyos/entry/src/main/module.json5 << 'EOF'
{
  "module": {
    "name": "entry",
    "type": "entry",
    "srcEntry": "./ets/Application/AbilityStage.ts",
    "description": "$string:module_desc",
    "mainElement": "MainAbility",
    "deviceTypes": ["phone", "tablet"],
    "abilities": [
      {
        "name": "MainAbility",
        "srcEntry": "./ets/pages/Index.ets",
        "description": "$string:MainAbility_desc",
        "icon": "$media:icon",
        "label": "$string:MainAbility_label",
        "startWindowIcon": "$media:icon",
        "startWindowBackground": "$color:start_window_background",
        "exported": true,
        "skills": [
          {
            "entities": ["entity.system.home"],
            "actions": ["action.system.home"]
          }
        ]
      }
    ]
  }
}
EOF

# Build HAP
hvigorw build --mode module

# Result: entry-default-signed.hap
```

---

## 🚀 Automated Build Script

Create `build-all.sh` (or `build-all.bat` for Windows):

```bash
#!/bin/bash

VERSION="2.0.0"
PROJECT_ROOT=$(pwd)

echo "🔨 Building FileManager v${VERSION} for all platforms..."

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Build Rust library (once for all platforms)
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
echo ""
echo "Binaries created:"
echo "  - filemanager-linux-amd64"
echo "  - filemanager-linux-arm64"
echo "  - filemanager.exe"
```

---

## 📊 Build Summary

### Quick Reference

```bash
# Linux amd64
make all

# Linux arm64 (cross-compile)
GOARCH=arm64 CC=aarch64-linux-gnu-gcc make all

# Windows amd64
GOOS=windows GOARCH=amd64 make all

# macOS amd64
GOOS=darwin GOARCH=amd64 make all

# macOS arm64
GOOS=darwin GOARCH=arm64 make all
```

### Using Makefile

The included `Makefile` supports:

```bash
make all          # Build for current platform
make rust         # Build Rust library only
make go           # Build Go binary only
make clean        # Clean build artifacts
make install      # Install to system
make test         # Run tests
make build-all    # Build for all platforms (requires cross-compilation tools)
```

---

## 🎯 Next Steps

1. ✅ Install prerequisites for your platform
2. ✅ Build binaries
3. ✅ Create packages
4. ✅ Test installations
5. ✅ Create releases on GitHub
6. ✅ Upload to package managers (Snap, Homebrew, AUR, etc.)

---

## 📚 Resources

- [Go Cross-Compilation](https://golang.org/doc/install/source#environment)
- [Rust Cross-Compilation](https://rust-lang.github.io/rustup/cross-compilation.html)
- [DEB Package Format](https://wiki.debian.org/DebianPackageFormat)
- [PKGBUILD Reference](https://wiki.archlinux.org/title/PKGBUILD)
- [Snap Documentation](https://snapcraft.io/docs)
- [Homebrew Formula](https://docs.brew.sh/Formula-Cookbook)

---

**Happy Building! 🚀**
