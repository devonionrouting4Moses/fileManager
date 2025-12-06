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

echo "🔨 Building FileManager v${VERSION} for Linux amd64..."

# Build Rust library
echo "🦀 Building Rust library..."
cd rust_ffi
cargo build --release -p fs-operations-core
cd ..

# Build Go binary
echo "🐹 Building Go binary..."
cd file_manager
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
  CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
  go build -ldflags="-s -w" -o ../filemanager ./cmd/app
cd ..

echo "✅ Binary built: ./filemanager"

# Create DEB package
echo "📦 Creating DEB package..."
mkdir -p deb/DEBIAN
mkdir -p deb/usr/local/bin
mkdir -p deb/usr/local/lib
mkdir -p deb/usr/share/doc/filemanager

cp filemanager deb/usr/local/bin/
cp rust_ffi/target/release/libfs_operations_core.so deb/usr/local/lib/
chmod +x deb/usr/local/bin/filemanager

# Copy debian/ files (single source of truth)
cp debian/control deb/DEBIAN/control
cp debian/postinst deb/DEBIAN/postinst
cp debian/prerm deb/DEBIAN/prerm
cp debian/postrm deb/DEBIAN/postrm

# Make scripts executable
chmod +x deb/DEBIAN/postinst
chmod +x deb/DEBIAN/prerm
chmod +x deb/DEBIAN/postrm

dpkg-deb --build deb filemanager_${VERSION}_${ARCH}.deb
echo "✅ DEB package created: filemanager_${VERSION}_${ARCH}.deb"

# Create TAR.GZ archive
echo "📦 Creating TAR.GZ archive..."
mkdir -p dist/filemanager-${VERSION}-linux-${ARCH}
cp filemanager dist/filemanager-${VERSION}-linux-${ARCH}/
cp rust_ffi/target/release/libfs_operations_core.so dist/filemanager-${VERSION}-linux-${ARCH}/
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
mkdir -p arch
cat > arch/PKGBUILD << 'EOFPKGBUILD'
pkgname=filemanager
pkgver=2.0.0
pkgrel=1
pkgdesc="Modern file manager with Rust+Go backend"
arch=('x86_64' 'aarch64')
url="https://github.com/devonionrouting4Moses/fileManager"
license=('MIT' 'Apache-2.0')
depends=('glibc' 'gcc-libs')
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
    install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md" 2>/dev/null || true
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE" 2>/dev/null || true
}
EOFPKGBUILD
echo "✅ PKGBUILD created: arch/PKGBUILD"

echo ""
echo "🎉 All Linux amd64 packages created successfully!"
echo ""
echo "Files created:"
echo "  - filemanager (binary)"
echo "  - filemanager_${VERSION}_${ARCH}.deb"
echo "  - filemanager-${VERSION}-linux-${ARCH}.tar.gz"
echo "  - filemanager-apt_${VERSION}_${ARCH}.deb"
echo "  - arch/PKGBUILD"
