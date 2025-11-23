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

ARCH="arm64"

echo "🔨 Building FileManager v${VERSION} for Linux arm64..."

# Setup ARM64 toolchain
echo "🛠️  Setting up ARM64 cross-compilation toolchain..."
rustup target add aarch64-unknown-linux-gnu 2>/dev/null || true

# Check if ARM64 linker is available
if ! command -v aarch64-linux-gnu-gcc &> /dev/null; then
    echo "❌ ERROR: aarch64-linux-gnu-gcc not found!"
    echo ""
    echo "Please install ARM64 cross-compilation toolchain:"
    echo "  sudo apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu binutils-aarch64-linux-gnu"
    echo ""
    echo "Or on other systems:"
    echo "  - Fedora: sudo dnf install gcc-aarch64-linux-gnu"
    echo "  - Arch: sudo pacman -S aarch64-linux-gnu-gcc"
    exit 1
fi

# Configure Rust to use the ARM64 linker
mkdir -p "$HOME/.cargo"

# Check if config already has ARM64 target
if ! grep -q "\[target.aarch64-unknown-linux-gnu\]" "$HOME/.cargo/config.toml" 2>/dev/null; then
    cat >> "$HOME/.cargo/config.toml" << 'EOFCONFIG'
[target.aarch64-unknown-linux-gnu]
linker = "aarch64-linux-gnu-gcc"
ar = "aarch64-linux-gnu-ar"
EOFCONFIG
fi

# Build Rust library for ARM64
echo "🦀 Building Rust library for ARM64..."
cd rust_ffi
cargo build --release --target aarch64-unknown-linux-gnu -p fs-operations-core --lib
cd ..

# Verify Rust library was built
if [ ! -f "rust_ffi/target/aarch64-unknown-linux-gnu/release/libfs_operations_core.so" ]; then
    echo "❌ ERROR: Rust library not found at rust_ffi/target/aarch64-unknown-linux-gnu/release/libfs_operations_core.so"
    echo ""
    echo "Available files in rust_ffi/target/aarch64-unknown-linux-gnu/release/:"
    ls -la rust_ffi/target/aarch64-unknown-linux-gnu/release/ 2>/dev/null || echo "Directory not found"
    exit 1
fi
echo "✅ Rust library built: rust_ffi/target/aarch64-unknown-linux-gnu/release/libfs_operations_core.so"

# Build Go binary for ARM64
echo "🐹 Building Go binary for ARM64..."
cd file_manager
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
  CC=aarch64-linux-gnu-gcc \
  CXX=aarch64-linux-gnu-g++ \
  CGO_LDFLAGS="-L../rust_ffi/target/aarch64-unknown-linux-gnu/release -lfs_operations_core -ldl -lpthread -lm" \
  go build -ldflags="-s -w" -o ../filemanager-arm64 ./cmd/app
cd ..

echo "✅ Binary built: ./filemanager-arm64"

# Create DEB package
echo "📦 Creating DEB package..."
mkdir -p deb-arm64/DEBIAN
mkdir -p deb-arm64/usr/local/bin
mkdir -p deb-arm64/usr/local/lib
mkdir -p deb-arm64/usr/share/doc/filemanager

cp filemanager-arm64 deb-arm64/usr/local/bin/filemanager
cp rust_ffi/target/aarch64-unknown-linux-gnu/release/libfs_operations_core.so deb-arm64/usr/local/lib/
chmod +x deb-arm64/usr/local/bin/filemanager

cat > deb-arm64/DEBIAN/control << 'EOF'
Package: filemanager
Version: 2.0.0
Architecture: arm64
Maintainer: DevChigarlicMoses <moses.muranja@strathmore.edu>
Description: Modern file manager with Rust+Go backend
 FileManager v2 is a dual-mode file manager with terminal and web interfaces.
 Built with Rust for performance and Go for flexibility.
Homepage: https://github.com/devonionrouting4Moses/fileManager
Depends: libc6, libgcc-s1, libstdc++6
EOF

cat > deb-arm64/DEBIAN/postinst << 'EOF'
#!/bin/bash
set -e
ldconfig
mkdir -p ~/.config/filemanager
echo "FileManager installed successfully!"
EOF
chmod +x deb-arm64/DEBIAN/postinst

cat > deb-arm64/DEBIAN/prerm << 'EOF'
#!/bin/bash
set -e
echo "Removing FileManager..."
EOF
chmod +x deb-arm64/DEBIAN/prerm

dpkg-deb --build deb-arm64 filemanager_${VERSION}_${ARCH}.deb
echo "✅ DEB package created: filemanager_${VERSION}_${ARCH}.deb"

# Create TAR.GZ archive
echo "📦 Creating TAR.GZ archive..."
mkdir -p dist/filemanager-${VERSION}-linux-${ARCH}
cp filemanager-arm64 dist/filemanager-${VERSION}-linux-${ARCH}/filemanager
cp rust_ffi/target/aarch64-unknown-linux-gnu/release/libfs_operations_core.so dist/filemanager-${VERSION}-linux-${ARCH}/
cp README.md dist/filemanager-${VERSION}-linux-${ARCH}/ 2>/dev/null || true
cp LICENSE dist/filemanager-${VERSION}-linux-${ARCH}/ 2>/dev/null || true

cat > dist/filemanager-${VERSION}-linux-${ARCH}/install.sh << 'EOFINSTALL'
#!/bin/bash
set -e
echo "Installing FileManager..."
sudo cp filemanager /usr/local/bin/
sudo chmod +x /usr/local/bin/filemanager
sudo cp libfs_operations_core.so /usr/local/lib/
sudo ldconfig
echo "✅ FileManager installed successfully!"
echo "Run 'filemanager' to start"
EOFINSTALL
chmod +x dist/filemanager-${VERSION}-linux-${ARCH}/install.sh

cd dist
tar czf filemanager-${VERSION}-linux-${ARCH}.tar.gz filemanager-${VERSION}-linux-${ARCH}/
cd ..
echo "✅ TAR.GZ archive created: filemanager-${VERSION}-linux-${ARCH}.tar.gz"

# Create APT package (same as DEB for now)
cp filemanager_${VERSION}_${ARCH}.deb filemanager-apt_${VERSION}_${ARCH}.deb
echo "✅ APT package created: filemanager-apt_${VERSION}_${ARCH}.deb"

# Create Arch PKGBUILD
echo "📦 Creating Arch PKGBUILD..."
mkdir -p arch-arm64
cat > arch-arm64/PKGBUILD << 'EOFPKGBUILD'
pkgname=filemanager
pkgver=2.0.0
pkgrel=1
pkgdesc="Modern file manager with Rust+Go backend"
arch=('aarch64')
url="https://github.com/devonionrouting4Moses/fileManager"
license=('MIT' 'Apache-2.0')
depends=('glibc' 'gcc-libs')
makedepends=('go' 'rust')

build() {
    cd "$srcdir"
    
    # Build Rust library
    cd rust_ffi
    cargo build --release --target aarch64-unknown-linux-gnu -p fs-operations-core
    cd ..
    
    # Build Go binary
    cd file_manager
    GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
      CC=aarch64-linux-gnu-gcc \
      CXX=aarch64-linux-gnu-g++ \
      CGO_LDFLAGS="-L../rust_ffi/target/aarch64-unknown-linux-gnu/release -lfs_operations_core -ldl -lpthread -lm" \
      go build -ldflags="-s -w" -o ../filemanager ./cmd/app
    cd ..
}

package() {
    install -Dm755 filemanager "$pkgdir/usr/bin/filemanager"
    install -Dm755 rust_ffi/target/aarch64-unknown-linux-gnu/release/libfs_operations_core.so "$pkgdir/usr/lib/libfs_operations_core.so"
    install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md" 2>/dev/null || true
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE" 2>/dev/null || true
}
EOFPKGBUILD
echo "✅ PKGBUILD created: arch-arm64/PKGBUILD"

echo ""
echo "🎉 All Linux arm64 packages created successfully!"
echo ""
echo "Files created:"
echo "  - filemanager-arm64 (binary)"
echo "  - filemanager_${VERSION}_${ARCH}.deb"
echo "  - filemanager-${VERSION}-linux-${ARCH}.tar.gz"
echo "  - filemanager-apt_${VERSION}_${ARCH}.deb"
echo "  - arch-arm64/PKGBUILD"
