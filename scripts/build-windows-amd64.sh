#!/bin/bash
set -e

# Get the project root (parent of scripts directory)
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Read version from VERSION file
if [ -f "$PROJECT_ROOT/VERSION" ]; then
    VERSION=$(cat "$PROJECT_ROOT/VERSION" | tr -d ' \n')
else
    VERSION="2.0.0"
    echo "⚠️  VERSION file not found, using default: $VERSION"
fi

ARCH="amd64"

echo "🔨 Building FileManager v${VERSION} for Windows amd64..."

# Check if we're on Windows or have MinGW
if [[ "$OSTYPE" != "msys" && "$OSTYPE" != "cygwin" && "$OSTYPE" != "win32" ]]; then
    # We're on Linux/macOS - check for MinGW
    if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
        echo "⚠️  Windows cross-compilation requires MinGW toolchain"
        echo ""
        echo "To build for Windows on Linux, install MinGW:"
        echo "  Ubuntu/Debian: sudo apt-get install -y mingw-w64"
        echo "  Fedora: sudo dnf install -y mingw64-gcc mingw64-gcc-c++"
        echo "  Arch: sudo pacman -S mingw-w64"
        echo ""
        echo "Or build on Windows with:"
        echo "  - MSYS2: https://www.msys2.org/"
        echo "  - Visual Studio Build Tools"
        echo ""
        echo "Skipping Windows build..."
        exit 0
    fi
    
    # Add Windows target to Rust
    echo "📦 Adding Windows target to Rust..."
    rustup target add x86_64-pc-windows-gnu 2>/dev/null || true
fi

# Build Rust library for Windows
echo "🦀 Building Rust library for Windows..."
cd rust_ffi
cargo build --release --target x86_64-pc-windows-gnu -p fs-operations-core
cd ..

# Build Go binary for Windows
echo "🐹 Building Go binary for Windows..."
cd file_manager
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  CC=x86_64-w64-mingw32-gcc \
  CXX=x86_64-w64-mingw32-g++ \
  CGO_LDFLAGS="-L../rust_ffi/target/x86_64-pc-windows-gnu/release -lfs_operations_core -lws2_32 -luserenv" \
  go build -ldflags="-s -w" -o ../filemanager.exe ./cmd/app
cd ..

echo "✅ Binary built: ./filemanager.exe"

# Create ZIP archive
echo "📦 Creating ZIP archive..."
mkdir -p dist/filemanager-${VERSION}-windows-${ARCH}
cp filemanager.exe dist/filemanager-${VERSION}-windows-${ARCH}/
cp rust_ffi/target/x86_64-pc-windows-gnu/release/fs_operations_core.dll dist/filemanager-${VERSION}-windows-${ARCH}/ 2>/dev/null || true
cp README.md dist/filemanager-${VERSION}-windows-${ARCH}/ 2>/dev/null || true
cp LICENSE dist/filemanager-${VERSION}-windows-${ARCH}/ 2>/dev/null || true

# Create Windows installer script
cat > dist/filemanager-${VERSION}-windows-${ARCH}/install.bat << 'EOFBAT'
@echo off
REM FileManager Installation Script

echo Installing FileManager...
set INSTALL_DIR=%ProgramFiles%\FileManager

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copy executable
copy /Y filemanager.exe "%INSTALL_DIR%\"

REM Copy Rust library
copy /Y fs_operations_core.dll "%INSTALL_DIR%\" 2>nul

REM Add to PATH
setx /M PATH "%PATH%;%INSTALL_DIR%"

echo.
echo ✅ Installation complete!
echo FileManager installed to: %INSTALL_DIR%
echo.
echo Run 'filemanager' from command prompt to start
pause
EOFBAT

# Create ZIP
if command -v 7z &> /dev/null; then
    cd dist
    7z a filemanager-${VERSION}-windows-${ARCH}.zip filemanager-${VERSION}-windows-${ARCH}/
    cd ..
    echo "✅ ZIP archive created: filemanager-${VERSION}-windows-${ARCH}.zip"
elif command -v zip &> /dev/null; then
    cd dist
    zip -r filemanager-${VERSION}-windows-${ARCH}.zip filemanager-${VERSION}-windows-${ARCH}/
    cd ..
    echo "✅ ZIP archive created: filemanager-${VERSION}-windows-${ARCH}.zip"
else
    echo "⚠️  7z or zip not found. Skipping ZIP creation."
    echo "Install 7-Zip or zip utility to create archives."
fi

echo ""
echo "🎉 Windows amd64 build completed successfully!"
echo ""
echo "Files created:"
echo "  - filemanager.exe (binary)"
echo "  - dist/filemanager-${VERSION}-windows-${ARCH}/ (distribution folder)"
echo "  - filemanager-${VERSION}-windows-${ARCH}.zip (if zip available)"
echo ""
echo "To install on Windows:"
echo "  1. Extract the ZIP file"
echo "  2. Run install.bat as Administrator"
echo "  3. Or manually copy files to Program Files"
