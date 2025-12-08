#!/bin/bash
set -e

# Get the project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Read version from VERSION file
if [ -f "$PROJECT_ROOT/VERSION" ]; then
    VERSION=$(cat "$PROJECT_ROOT/VERSION" | tr -d ' \n')
else
    VERSION="2.0.0"
    echo "⚠️  VERSION file not found, using default: $VERSION"
fi

MACOS_DIST_DIR="$PROJECT_ROOT/macos/dist"
HOMEBREW_DIR="$PROJECT_ROOT/macos/homebrew"

echo "🍺 Generating Homebrew formula for FileManager v${VERSION}..."

# Check if both architecture builds exist
if [ ! -f "$MACOS_DIST_DIR/checksums-amd64.txt" ]; then
    echo "❌ AMD64 build not found. Please run: bash scripts/build-macos.sh amd64"
    exit 1
fi

if [ ! -f "$MACOS_DIST_DIR/checksums-arm64.txt" ]; then
    echo "❌ ARM64 build not found. Please run: bash scripts/build-macos.sh arm64"
    exit 1
fi

# Read checksums
source "$MACOS_DIST_DIR/checksums-amd64.txt"
source "$MACOS_DIST_DIR/checksums-arm64.txt"

# Create homebrew directory
mkdir -p "$HOMEBREW_DIR"

# Generate the formula for CLI installation (from tar.gz)
echo "📝 Creating Homebrew formula (CLI version)..."
cat > "$HOMEBREW_DIR/filemanager.rb" << EOFFORMULA
class Filemanager < Formula
  desc "Modern file manager with Rust+Go backend and web interface"
  homepage "https://github.com/devonionrouting4Moses/fileManager"
  version "${VERSION}"
  license "MIT"

  if Hardware::CPU.intel?
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v${VERSION}/filemanager-${VERSION}-macos-amd64.tar.gz"
    sha256 "${TAR_SHA256_AMD64}"
  elsif Hardware::CPU.arm?
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v${VERSION}/filemanager-${VERSION}-macos-arm64.tar.gz"
    sha256 "${TAR_SHA256_ARM64}"
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

echo "✅ Formula created: macos/homebrew/filemanager.rb"

# Generate the Cask formula for GUI installation (from DMG)
echo "📝 Creating Homebrew Cask formula (GUI version)..."
mkdir -p "$HOMEBREW_DIR/Casks"
cat > "$HOMEBREW_DIR/Casks/filemanager.rb" << EOFCASK
cask "filemanager" do
  version "${VERSION}"
  
  if Hardware::CPU.intel?
    sha256 "${DMG_SHA256_AMD64}"
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v#{version}/FileManager-#{version}-macos-amd64.dmg"
  else
    sha256 "${DMG_SHA256_ARM64}"
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v#{version}/FileManager-#{version}-macos-arm64.dmg"
  end

  name "FileManager"
  desc "Modern file manager with Rust+Go backend and web interface"
  homepage "https://github.com/devonionrouting4Moses/fileManager"

  livecheck do
    url :url
    strategy :github_latest
  end

  app "FileManager.app"

  uninstall quit: "com.devchigarlicmoses.filemanager"

  zap trash: [
    "~/Library/Application Support/FileManager",
    "~/Library/Preferences/com.devchigarlicmoses.filemanager.plist",
    "~/Library/Caches/FileManager",
    "~/Library/Logs/FileManager",
    "~/Library/Saved Application State/com.devchigarlicmoses.filemanager.savedState",
  ]

  caveats <<~EOS
    FileManager has been installed!
    
    To launch the app:
      - Open from Applications folder, or
      - Run: open -a FileManager
    
    The web interface will be available at:
      http://localhost:8080
    
    Note: On first launch, you may need to:
      1. Right-click the app and select "Open"
      2. Grant Full Disk Access in System Preferences → Privacy & Security
  EOS
end
EOFCASK

echo "✅ Cask formula created: macos/homebrew/Casks/filemanager.rb"

# Create a README for the homebrew directory
cat > "$HOMEBREW_DIR/README.md" << 'EOFREADME'
# FileManager Homebrew Formulas

This directory contains Homebrew formulas for installing FileManager on macOS.

## Installation Options

### Option 1: GUI App (Cask - Recommended for most users)

Installs the full macOS application bundle:

```bash
# From your tap
brew tap devonionrouting4Moses/tap
brew install --cask filemanager

# Or directly (once published to official Homebrew)
brew install --cask filemanager
```

### Option 2: CLI Tool (Formula)

Installs the command-line binary:

```bash
# From your tap
brew tap devonionrouting4Moses/tap
brew install filemanager

# Or directly (once published to official Homebrew)
brew install filemanager
```

## Usage

### GUI App
- Launch from Applications folder
- Or run: `open -a FileManager`

### CLI Tool
```bash
# Start the file manager
filemanager

# Access web interface at http://localhost:8080
```

## Uninstall

### GUI App
```bash
brew uninstall --cask filemanager
```

### CLI Tool
```bash
brew uninstall filemanager
```

## Development

### Testing Formulas Locally

```bash
# Test the Cask
brew install --cask --debug ./Casks/filemanager.rb

# Test the Formula
brew install --formula --debug ./filemanager.rb
```

### Audit Before Submission

```bash
# Audit the Cask
brew audit --cask --strict ./Casks/filemanager.rb

# Audit the Formula
brew audit --strict ./filemanager.rb
```

## Publishing

### To Your Own Tap

1. Create a repository: `homebrew-tap`
2. Copy formulas:
   - `filemanager.rb` → `Formula/filemanager.rb` (for CLI)
   - `Casks/filemanager.rb` → `Casks/filemanager.rb` (for GUI)
3. Users install via: `brew tap yourusername/tap`

### To Official Homebrew

#### For Cask (GUI):
1. Fork `homebrew/homebrew-cask`
2. Add your cask to `Casks/f/filemanager.rb`
3. Submit pull request

#### For Formula (CLI):
1. Fork `homebrew/homebrew-core`
2. Add your formula to `Formula/f/filemanager.rb`
3. Submit pull request

## Checksums

Checksums are automatically generated during build and embedded in these formulas.

To verify:
```bash
shasum -a 256 FileManager-*.dmg
shasum -a 256 filemanager-*-macos-*.tar.gz
```
EOFREADME

echo "✅ README created: macos/homebrew/README.md"

# Create a summary file with all checksums
cat > "$HOMEBREW_DIR/CHECKSUMS.txt" << EOFCHECKSUMS
FileManager v${VERSION} - Build Checksums
Generated: $(date)

=== DMG Files (for Cask) ===
AMD64 (Intel):    ${DMG_SHA256_AMD64}
ARM64 (M1/M2/M3): ${DMG_SHA256_ARM64}

=== TAR.GZ Files (for Formula) ===
AMD64 (Intel):    ${TAR_SHA256_AMD64}
ARM64 (M1/M2/M3): ${TAR_SHA256_ARM64}

=== File Locations ===
DMG AMD64:    macos/dist/amd64/FileManager-${VERSION}-macos-amd64.dmg
DMG ARM64:    macos/dist/arm64/FileManager-${VERSION}-macos-arm64.dmg
TAR.GZ AMD64: macos/dist/amd64/filemanager-${VERSION}-macos-amd64.tar.gz
TAR.GZ ARM64: macos/dist/arm64/filemanager-${VERSION}-macos-arm64.tar.gz
EOFCHECKSUMS

echo "✅ Checksums saved: macos/homebrew/CHECKSUMS.txt"

echo ""
echo "🎉 Homebrew formulas generated successfully!"
echo ""
echo "Generated files:"
echo "  📝 macos/homebrew/filemanager.rb (CLI Formula)"
echo "  📝 macos/homebrew/Casks/filemanager.rb (GUI Cask)"
echo "  📝 macos/homebrew/README.md"
echo "  📝 macos/homebrew/CHECKSUMS.txt"
echo ""
echo "Next steps:"
echo "  1. Test locally:"
echo "     brew install --cask macos/homebrew/Casks/filemanager.rb"
echo "  2. Create GitHub release with DMG and TAR.GZ files"
echo "  3. Create homebrew-tap repository"
echo "  4. Copy formulas to your tap"
echo ""
