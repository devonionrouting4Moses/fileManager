I've created a comprehensive macOS build script that:
Key Features:

Dual Architecture Support: Builds for both Intel (amd64) and Apple Silicon (arm64)
Multiple Distribution Formats:

.app bundle for drag-and-drop installation
.dmg installer with Applications symlink
.tar.gz archive for manual installation
Homebrew formula for package managers


macOS-Specific Features:

Proper .app bundle structure with Info.plist
Framework bundling with correct @rpath configuration
Optional code signing (if Developer ID certificate available)
Frontend resource bundling


Integration:

Follows the same pattern as your Linux/Windows scripts
Reads version from VERSION file
Compatible with the main build.sh orchestrator



Usage:
bash# Build for current architecture
./scripts/build-macos.sh

# Build for specific architecture
./scripts/build-macos.sh amd64  # Intel
./scripts/build-macos.sh arm64  # Apple Silicon

# Build all platforms including macOS
./scripts/build.sh all

# Build only macOS versions
./scripts/build.sh macos-amd64 macos-arm64
Output Files:

filemanager-macos-{arch} - Standalone binary
dist/FileManager.app - macOS application bundle
dist/FileManager-{version}-macos-{arch}.dmg - DMG installer
dist/filemanager-{version}-macos-{arch}.tar.gz - Archive
homebrew/filemanager.rb - Homebrew formula

The script handles macOS-specific requirements like framework linking, properly updates the library paths, and creates user-friendly installation options!

Perfect! I've merged the best features from both versions. Here's what the enhanced Homebrew formula now includes:
Key Improvements:
From Your Original Formula:
✅ Config Management: Creates config.yaml in etc/filemanager/
✅ Post-Install Permissions: Properly sets frontend directory and file permissions
✅ Detailed Configuration: Includes host, port, and theme settings
✅ Better Test Suite: Checks for binary existence and executability
From My Version:
✅ Multi-Architecture Support: Automatic Intel/ARM64 detection
✅ Pre-Built Binaries: Downloads from GitHub releases (faster than building from source)
✅ Rust Library Integration: Properly installs and manages libfs_operations_core.dylib
✅ Service Support: Can run as a background service with brew services
New Combined Features:

Complete Installation:

Binary in /usr/local/bin
Library in /usr/local/lib
Frontend in /usr/local/share/filemanager/frontend
Config in /usr/local/etc/filemanager/config.yaml


Service Management:

bash   brew services start filemanager   # Run in background
   brew services stop filemanager    # Stop service
   brew services restart filemanager # Restart

Logging:

Stdout: /usr/local/var/log/filemanager/stdout.log
Stderr: /usr/local/var/log/filemanager/stderr.log


Comprehensive Caveats: Shows users all important paths and how to use the app

To Use This Formula:
bash# After building both architectures
shasum -a 256 dist/filemanager-2.0.2-macos-amd64.tar.gz
shasum -a 256 dist/filemanager-2.0.2-macos-arm64.tar.gz

# Update SHA256 values in the formula, then:
brew install --build-from-source homebrew/filemanager.rb

# Or for testing locally:
brew install --formula homebrew/filemanager.rb
This unified version gives you the best of both approaches! 🎉