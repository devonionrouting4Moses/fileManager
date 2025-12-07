#!/bin/bash
set -e

# Get the project root (parent of scripts directory)
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Read version from VERSION file
if [ -f "$PROJECT_ROOT/VERSION" ]; then
    VERSION=$(cat "$PROJECT_ROOT/VERSION" | tr -d ' \n')
else
    VERSION="2.0.0"
    echo "⚠️  VERSION file not found, using default: $VERSION"
fi

# Detect architecture or use parameter
if [ -z "$1" ]; then
    ARCH=$(uname -m)
    if [ "$ARCH" = "x86_64" ]; then
        ARCH="amd64"
        TARGET="x86_64-apple-darwin"
    elif [ "$ARCH" = "arm64" ]; then
        TARGET="aarch64-apple-darwin"
    else
        echo "❌ Unknown architecture: $ARCH"
        exit 1
    fi
else
    ARCH="$1"
    if [ "$ARCH" = "amd64" ]; then
        TARGET="x86_64-apple-darwin"
    elif [ "$ARCH" = "arm64" ]; then
        TARGET="aarch64-apple-darwin"
    else
        echo "❌ Unknown architecture: $ARCH"
        echo "Usage: $0 [amd64|arm64]"
        exit 1
    fi
fi

echo "🔨 Building FileManager v${VERSION} for macOS ${ARCH}..."

# Check if we're on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo "⚠️  This script is designed to run on macOS"
    echo "Cross-compilation from other platforms to macOS requires additional setup"
    exit 1
fi

# Add Rust target if needed
echo "📦 Adding macOS ${ARCH} target to Rust..."
rustup target add ${TARGET} 2>/dev/null || true

# Build Rust library for macOS
echo "🦀 Building Rust library for macOS ${ARCH}..."
cd rust_ffi
cargo build --release --target ${TARGET} -p fs-operations-core
cd ..

# Verify Rust library was built
if [ ! -f "rust_ffi/target/${TARGET}/release/libfs_operations_core.dylib" ]; then
    echo "❌ ERROR: Rust library not found at rust_ffi/target/${TARGET}/release/libfs_operations_core.dylib"
    echo ""
    echo "Available files in rust_ffi/target/${TARGET}/release/:"
    ls -la "rust_ffi/target/${TARGET}/release/" 2>/dev/null || echo "Directory not found"
    exit 1
fi
echo "✅ Rust library built: rust_ffi/target/${TARGET}/release/libfs_operations_core.dylib"

# Build Go binary for macOS
echo "🐹 Building Go binary for macOS ${ARCH}..."
cd file_manager

if [ "$ARCH" = "amd64" ]; then
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
      CGO_LDFLAGS="-L../rust_ffi/target/${TARGET}/release -lfs_operations_core -framework CoreFoundation -framework Security" \
      go build -ldflags="-s -w" -o ../filemanager-macos-${ARCH} ./cmd/app
else
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
      CGO_LDFLAGS="-L../rust_ffi/target/${TARGET}/release -lfs_operations_core -framework CoreFoundation -framework Security" \
      go build -ldflags="-s -w" -o ../filemanager-macos-${ARCH} ./cmd/app
fi
cd ..

echo "✅ Binary built: ./filemanager-macos-${ARCH}"

# Create macOS .app bundle
echo "📦 Creating macOS application bundle..."
APP_NAME="FileManager.app"
APP_DIR="dist/${APP_NAME}"
mkdir -p "${APP_DIR}/Contents/MacOS"
mkdir -p "${APP_DIR}/Contents/Resources"
mkdir -p "${APP_DIR}/Contents/Frameworks"

# Copy binary
cp "filemanager-macos-${ARCH}" "${APP_DIR}/Contents/MacOS/filemanager"
chmod +x "${APP_DIR}/Contents/MacOS/filemanager"

# Copy Rust library
cp "rust_ffi/target/${TARGET}/release/libfs_operations_core.dylib" "${APP_DIR}/Contents/Frameworks/"

# Update library paths in binary
install_name_tool -change \
    "@rpath/libfs_operations_core.dylib" \
    "@executable_path/../Frameworks/libfs_operations_core.dylib" \
    "${APP_DIR}/Contents/MacOS/filemanager"

# Bundle frontend files
if [ -d "filemanager_frontend" ]; then
    echo "📂 Bundling frontend files..."
    mkdir -p "${APP_DIR}/Contents/Resources/frontend"
    cp -r filemanager_frontend/* "${APP_DIR}/Contents/Resources/frontend/"
else
    echo "⚠️  Frontend directory not found, skipping frontend bundling"
fi

# Create Info.plist
cat > "${APP_DIR}/Contents/Info.plist" << EOFPLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>en</string>
    <key>CFBundleExecutable</key>
    <string>filemanager</string>
    <key>CFBundleIdentifier</key>
    <string>com.devchigarlicmoses.filemanager</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>FileManager</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>${VERSION}</string>
    <key>CFBundleVersion</key>
    <string>${VERSION}</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSHumanReadableCopyright</key>
    <string>Copyright © 2024 DevChigarlicMoses. All rights reserved.</string>
    <key>LSApplicationCategoryType</key>
    <string>public.app-category.utilities</string>
</dict>
</plist>
EOFPLIST

echo "✅ Application bundle created: ${APP_DIR}"

# Create DMG installer
echo "📦 Creating DMG installer..."
DMG_NAME="FileManager-${VERSION}-macos-${ARCH}.dmg"
DMG_TEMP_DIR="dist/dmg-temp"
mkdir -p "${DMG_TEMP_DIR}"

# Copy app to temp directory
cp -R "${APP_DIR}" "${DMG_TEMP_DIR}/"

# Create Applications symlink
ln -s /Applications "${DMG_TEMP_DIR}/Applications"

# Create README
cat > "${DMG_TEMP_DIR}/README.txt" << EOFREADME
FileManager v${VERSION}

Installation:
1. Drag FileManager.app to the Applications folder
2. Open FileManager from Applications or Launchpad

Uninstallation:
1. Move FileManager.app to Trash
2. Clean up: rm -rf ~/Library/Application\ Support/FileManager

For more information, visit:
https://github.com/devonionrouting4Moses/fileManager
EOFREADME

# Create DMG
if command -v hdiutil &> /dev/null; then
    # Create temporary DMG
    hdiutil create -volname "FileManager" -srcfolder "${DMG_TEMP_DIR}" -ov -format UDZO "dist/${DMG_NAME}"
    echo "✅ DMG installer created: dist/${DMG_NAME}"
else
    echo "⚠️  hdiutil not found, skipping DMG creation"
fi

# Clean up temp directory
rm -rf "${DMG_TEMP_DIR}"

# Create TAR.GZ archive
echo "📦 Creating TAR.GZ archive..."
mkdir -p "dist/filemanager-${VERSION}-macos-${ARCH}"
cp "filemanager-macos-${ARCH}" "dist/filemanager-${VERSION}-macos-${ARCH}/filemanager"
cp "rust_ffi/target/${TARGET}/release/libfs_operations_core.dylib" "dist/filemanager-${VERSION}-macos-${ARCH}/"
cp README.md "dist/filemanager-${VERSION}-macos-${ARCH}/" 2>/dev/null || true
cp LICENSE "dist/filemanager-${VERSION}-macos-${ARCH}/" 2>/dev/null || true

# Bundle frontend in tar.gz
if [ -d "filemanager_frontend" ]; then
    mkdir -p "dist/filemanager-${VERSION}-macos-${ARCH}/frontend"
    cp -r filemanager_frontend/* "dist/filemanager-${VERSION}-macos-${ARCH}/frontend/"
fi

# Create install script
cat > "dist/filemanager-${VERSION}-macos-${ARCH}/install.sh" << 'EOFINSTALL'
#!/bin/bash
set -e
echo "Installing FileManager..."

# Check if running with sudo
if [ "$EUID" -ne 0 ]; then 
    echo "Please run with sudo: sudo ./install.sh"
    exit 1
fi

# Copy binary to /usr/local/bin
cp filemanager /usr/local/bin/
chmod +x /usr/local/bin/filemanager

# Copy library to /usr/local/lib
cp libfs_operations_core.dylib /usr/local/lib/

# Update library search path
if [ ! -f /etc/paths.d/filemanager ]; then
    echo "/usr/local/lib" > /etc/paths.d/filemanager
fi

echo "✅ FileManager installed successfully!"
echo "Run 'filemanager' from terminal to start"
EOFINSTALL
chmod +x "dist/filemanager-${VERSION}-macos-${ARCH}/install.sh"

cd dist
tar czf "filemanager-${VERSION}-macos-${ARCH}.tar.gz" "filemanager-${VERSION}-macos-${ARCH}/"
cd ..
echo "✅ TAR.GZ archive created: dist/filemanager-${VERSION}-macos-${ARCH}.tar.gz"

# Create Homebrew formula
echo "📦 Creating Homebrew formula..."
mkdir -p homebrew
cat > "homebrew/filemanager.rb" << EOFFORMULA
class Filemanager < Formula
  desc "Modern file manager with Rust+Go backend and web interface"
  homepage "https://github.com/devonionrouting4Moses/fileManager"
  version "${VERSION}"
  license "MIT"

  if Hardware::CPU.intel?
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v${VERSION}/filemanager-${VERSION}-macos-amd64.tar.gz"
    sha256 "REPLACE_WITH_AMD64_SHA256"
  elsif Hardware::CPU.arm?
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v${VERSION}/filemanager-${VERSION}-macos-arm64.tar.gz"
    sha256 "REPLACE_WITH_ARM64_SHA256"
  end

  def install
    # Install binary
    bin.install "filemanager"
    
    # Install Rust dynamic library
    lib.install "libfs_operations_core.dylib"
    
    # Install frontend if present
    if File.directory?("frontend")
      (share/"filemanager/frontend").install Dir["frontend/*"]
    end
    
    # Create config directory
    (etc/"filemanager").mkpath
    
    # Create default system config
    (etc/"filemanager/config.yaml").write <<~EOS
      # FileManager System Configuration
      # Frontend directory path
      frontend_dir: #{share}/filemanager/frontend
      
      # Server configuration
      host: 127.0.0.1
      port: 8080
      
      # Default settings
      theme: auto
      show_hidden: false
    EOS
  end

  def post_install
    # Set frontend directory permissions
    frontend_dir = "#{share}/filemanager/frontend"
    
    if File.directory?(frontend_dir)
      # Make directories readable
      system "chmod", "-R", "755", frontend_dir
      
      # Make files readable
      system "find", frontend_dir, "-type", "f", "-exec", "chmod", "644", "{}", "+"
    end
    
    # Verify library is accessible
    dylib_path = "#{lib}/libfs_operations_core.dylib"
    if File.exist?(dylib_path)
      system "chmod", "755", dylib_path
    end
  end

  def caveats
    <<~EOS
      FileManager has been installed!
      
      To start the file manager:
        filemanager
      
      The web interface will be available at:
        http://localhost:8080
      
      Configuration file:
        #{etc}/filemanager/config.yaml
      
      Frontend files:
        #{share}/filemanager/frontend
      
      To run as a background service:
        brew services start filemanager
    EOS
  end

  service do
    run [opt_bin/"filemanager"]
    keep_alive true
    log_path var/"log/filemanager/stdout.log"
    error_log_path var/"log/filemanager/stderr.log"
  end

  test do
    # Test that binary exists and is executable
    assert_predicate bin/"filemanager", :exist?
    assert_predicate bin/"filemanager", :executable?
    
    # Test that library exists
    assert_predicate lib/"libfs_operations_core.dylib", :exist?
    
    # Test version output (if implemented)
    system "#{bin}/filemanager", "--version" rescue nil
  end
end
EOFFORMULA
echo "✅ Homebrew formula created: homebrew/filemanager.rb"
echo ""
echo "To calculate SHA256 checksums after building:"
echo "  shasum -a 256 dist/filemanager-${VERSION}-macos-amd64.tar.gz"
echo "  shasum -a 256 dist/filemanager-${VERSION}-macos-arm64.tar.gz"
echo ""
echo "Then update the formula with the checksums before publishing."

# Code signing (optional, if developer certificate is available)
if security find-identity -v -p codesigning | grep -q "Developer ID Application"; then
    echo "🔐 Code signing application..."
    codesign --force --deep --sign "Developer ID Application" "${APP_DIR}"
    echo "✅ Application signed"
else
    echo "⚠️  No Developer ID certificate found, skipping code signing"
    echo "   Users may need to allow the app in System Preferences > Security & Privacy"
fi

echo ""
echo "🎉 macOS ${ARCH} build completed successfully!"
echo ""
echo "Files created:"
echo "  - filemanager-macos-${ARCH} (binary)"
echo "  - dist/${APP_NAME} (application bundle)"
echo "  - dist/${DMG_NAME} (DMG installer)"
echo "  - dist/filemanager-${VERSION}-macos-${ARCH}.tar.gz"
echo "  - homebrew/filemanager.rb (Homebrew formula)"
echo ""
echo "Installation options:"
echo "  1. DMG: Open ${DMG_NAME} and drag to Applications"
echo "  2. Binary: Run dist/filemanager-${VERSION}-macos-${ARCH}/install.sh"
echo "  3. Homebrew: brew install --cask /path/to/filemanager.rb"
echo ""
echo "To build for the other architecture:"
echo "  ./build-macos.sh [amd64|arm64]"