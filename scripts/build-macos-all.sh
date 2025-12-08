#!/bin/bash
set -e

# Get the project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPTS_DIR="$PROJECT_ROOT/scripts"

echo "🔨 Building FileManager for All macOS Architectures"
echo "===================================================="
echo ""

# Build for AMD64 (Intel)
echo "🔨 Building for macOS AMD64 (Intel)..."
bash "$SCRIPTS_DIR/build-macos.sh" amd64
echo ""

# Build for ARM64 (Apple Silicon)
echo "🔨 Building for macOS ARM64 (Apple Silicon)..."
bash "$SCRIPTS_DIR/build-macos.sh" arm64
echo ""

# Generate Homebrew formulas
echo "🍺 Generating Homebrew formulas..."
bash "$SCRIPTS_DIR/generate-homebrew-formula.sh"
echo ""

# Display summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 All macOS builds completed successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📦 Build artifacts organized in macos/ directory:"
echo ""
echo "  macos/"
echo "  ├── filemanager-macos-amd64          (Intel binary)"
echo "  ├── filemanager-macos-arm64          (ARM binary)"
echo "  ├── dist/"
echo "  │   ├── amd64/"
echo "  │   │   ├── FileManager.app"
echo "  │   │   ├── FileManager-*-macos-amd64.dmg"
echo "  │   │   └── filemanager-*-macos-amd64.tar.gz"
echo "  │   ├── arm64/"
echo "  │   │   ├── FileManager.app"
echo "  │   │   ├── FileManager-*-macos-arm64.dmg"
echo "  │   │   └── filemanager-*-macos-arm64.tar.gz"
echo "  │   ├── checksums-amd64.txt"
echo "  │   └── checksums-arm64.txt"
echo "  └── homebrew/"
echo "      ├── filemanager.rb               (CLI Formula)"
echo "      ├── Casks/filemanager.rb         (GUI Cask)"
echo "      ├── CHECKSUMS.txt"
echo "      └── README.md"
echo ""
echo "🚀 Ready to distribute!"
echo ""
echo "Next steps:"
echo "  1. Test the builds:"
echo "     • open macos/dist/amd64/FileManager.app"
echo "     • open macos/dist/arm64/FileManager.app"
echo ""
echo "  2. Create GitHub Release (v$(cat $PROJECT_ROOT/VERSION)):"
echo "     • Upload all DMG files from macos/dist/"
echo "     • Upload all TAR.GZ files from macos/dist/"
echo ""
echo "  3. Setup Homebrew distribution:"
echo "     • Create 'homebrew-tap' repository on GitHub"
echo "     • Copy macos/homebrew/ contents to the tap"
echo "     • Users can then: brew tap yourname/tap && brew install --cask filemanager"
echo ""
echo "  4. Optional - Submit to official Homebrew:"
echo "     • Fork homebrew/homebrew-cask"
echo "     • Add macos/homebrew/Casks/filemanager.rb"
echo "     • Submit pull request"
echo ""

# Check if running on the native architecture
CURRENT_ARCH=$(uname -m)
if [ "$CURRENT_ARCH" = "arm64" ]; then
    echo "💡 Tip: You built on Apple Silicon. Test the ARM64 build first:"
    echo "   open macos/dist/arm64/FileManager.app"
elif [ "$CURRENT_ARCH" = "x86_64" ]; then
    echo "💡 Tip: You built on Intel Mac. Test the AMD64 build first:"
    echo "   open macos/dist/amd64/FileManager.app"
fi
echo ""
