# macOS Build Guide

## Quick Start

### Build Everything (Recommended)
```bash
bash scripts/build-macos-all.sh
```

This will:
1. Build for Intel (amd64)
2. Build for Apple Silicon (arm64)
3. Generate Homebrew formulas
4. Organize everything in `macos/` directory

### Build Single Architecture
```bash
# Intel Macs
bash scripts/build-macos.sh amd64

# Apple Silicon (M1/M2/M3)
bash scripts/build-macos.sh arm64
```

### Generate Homebrew Formulas
After building both architectures:
```bash
bash scripts/generate-homebrew-formula.sh
```

---

## Output Structure

All macOS build artifacts are organized in the `macos/` directory:

```
macos/
├── filemanager-macos-amd64              # Intel binary
├── filemanager-macos-arm64              # ARM binary
├── dist/
│   ├── amd64/                           # Intel build artifacts
│   │   ├── FileManager.app/             # macOS app bundle
│   │   ├── FileManager-2.0.2-macos-amd64.dmg
│   │   └── filemanager-2.0.2-macos-amd64.tar.gz
│   ├── arm64/                           # ARM build artifacts
│   │   ├── FileManager.app/             # macOS app bundle
│   │   ├── FileManager-2.0.2-macos-arm64.dmg
│   │   └── filemanager-2.0.2-macos-arm64.tar.gz
│   ├── checksums-amd64.txt              # Intel checksums
│   └── checksums-arm64.txt              # ARM checksums
└── homebrew/
    ├── filemanager.rb                   # CLI Formula (tar.gz)
    ├── Casks/
    │   └── filemanager.rb               # GUI Cask (DMG)
    ├── CHECKSUMS.txt                    # All checksums
    └── README.md                        # Homebrew instructions
```

---

## What Gets Built

### For Each Architecture (amd64, arm64)

1. **Standalone Binary**
   - Location: `macos/filemanager-macos-{arch}`
   - Can be run directly from terminal

2. **Application Bundle (.app)**
   - Location: `macos/dist/{arch}/FileManager.app`
   - Contains binary, libraries, and frontend
   - Can be dragged to Applications folder

3. **DMG Installer**
   - Location: `macos/dist/{arch}/FileManager-*.dmg`
   - Professional installer with Applications symlink
   - Used for Homebrew Cask distribution

4. **TAR.GZ Archive**
   - Location: `macos/dist/{arch}/filemanager-*.tar.gz`
   - Contains binary, library, frontend, and install script
   - Used for Homebrew Formula (CLI) distribution

5. **Checksums File**
   - Location: `macos/dist/checksums-{arch}.txt`
   - SHA256 hashes for DMG and TAR.GZ files

---

## Homebrew Distribution

Two installation methods are generated:

### 1. Cask (GUI App) - Recommended for Most Users
```bash
brew install --cask filemanager
```
- Installs from DMG
- Puts FileManager.app in /Applications
- Users launch from Applications folder

### 2. Formula (CLI Tool) - For Command-Line Users
```bash
brew install filemanager
```
- Installs from TAR.GZ
- Puts binary in /usr/local/bin
- Users run `filemanager` from terminal

---

## Distribution Workflow

### 1. Build for Both Architectures
```bash
bash scripts/build-macos-all.sh
```

### 2. Test the Builds
```bash
# Test Intel build
open macos/dist/amd64/FileManager.app

# Test ARM build (if on Apple Silicon)
open macos/dist/arm64/FileManager.app

# Test binaries directly
./macos/filemanager-macos-amd64
./macos/filemanager-macos-arm64
```

### 3. Create GitHub Release

1. Go to your GitHub repository
2. Create new release: `v2.0.2`
3. Upload these files from `macos/dist/`:
   ```
   amd64/FileManager-2.0.2-macos-amd64.dmg
   amd64/filemanager-2.0.2-macos-amd64.tar.gz
   arm64/FileManager-2.0.2-macos-arm64.dmg
   arm64/filemanager-2.0.2-macos-arm64.tar.gz
   ```

### 4. Setup Homebrew Tap

#### Create Repository
```bash
# On GitHub, create a new repository named: homebrew-tap
```

#### Setup Tap Structure
```bash
cd /tmp
git clone https://github.com/yourusername/homebrew-tap
cd homebrew-tap

# Create structure
mkdir -p Formula Casks

# Copy formulas
cp /path/to/fileManager/macos/homebrew/filemanager.rb Formula/
cp /path/to/fileManager/macos/homebrew/Casks/filemanager.rb Casks/
cp /path/to/fileManager/macos/homebrew/README.md .

# Commit and push
git add .
git commit -m "Add FileManager formulas"
git push
```

#### Users Install
```bash
# Add your tap
brew tap yourusername/tap

# Install GUI app
brew install --cask filemanager

# Or install CLI tool
brew install filemanager
```

### 5. Submit to Official Homebrew (Optional)

#### For Cask (GUI):
```bash
# Fork homebrew/homebrew-cask
git clone https://github.com/yourusername/homebrew-cask
cd homebrew-cask
git checkout -b filemanager

# Add your cask
cp /path/to/fileManager/macos/homebrew/Casks/filemanager.rb Casks/f/

# Test it
brew audit --cask --strict Casks/f/filemanager.rb
brew install --cask --debug Casks/f/filemanager.rb

# Commit and create PR
git add Casks/f/filemanager.rb
git commit -m "Add FileManager cask"
git push origin filemanager
```

#### For Formula (CLI):
```bash
# Fork homebrew/homebrew-core
git clone https://github.com/yourusername/homebrew-core
cd homebrew-core
git checkout -b filemanager

# Add your formula
cp /path/to/fileManager/macos/homebrew/filemanager.rb Formula/f/

# Test it
brew audit --strict Formula/f/filemanager.rb
brew install --debug Formula/f/filemanager.rb

# Commit and create PR
git add Formula/f/filemanager.rb
git commit -m "Add FileManager formula"
git push origin filemanager
```

---

## Code Signing & Notarization

For official distribution, you need to sign and notarize your builds.

### Prerequisites
- Apple Developer Account ($99/year)
- Developer ID Application certificate

### Sign the App
```bash
# Sign the app bundle
codesign --force --deep --sign "Developer ID Application: Your Name (TEAM_ID)" \
  macos/dist/amd64/FileManager.app

codesign --force --deep --sign "Developer ID Application: Your Name (TEAM_ID)" \
  macos/dist/arm64/FileManager.app
```

### Notarize the DMG
```bash
# For Intel
xcrun notarytool submit macos/dist/amd64/FileManager-2.0.2-macos-amd64.dmg \
  --apple-id "your@email.com" \
  --password "app-specific-password" \
  --team-id "YOUR_TEAM_ID" \
  --wait

xcrun stapler staple macos/dist/amd64/FileManager-2.0.2-macos-amd64.dmg

# For ARM
xcrun notarytool submit macos/dist/arm64/FileManager-2.0.2-macos-arm64.dmg \
  --apple-id "your@email.com" \
  --password "app-specific-password" \
  --team-id "YOUR_TEAM_ID" \
  --wait

xcrun stapler staple macos/dist/arm64/FileManager-2.0.2-macos-arm64.dmg
```

### Update Checksums
After signing and notarizing, recalculate checksums:
```bash
shasum -a 256 macos/dist/amd64/FileManager-2.0.2-macos-amd64.dmg
shasum -a 256 macos/dist/arm64/FileManager-2.0.2-macos-arm64.dmg

# Update the checksums in macos/homebrew/Casks/filemanager.rb
```

---

## Troubleshooting

### Build Fails: "cargo: command not found"
```bash
# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source "$HOME/.cargo/env"
```

### Build Fails: "go: command not found"
```bash
# Install Go
brew install go
```

### "Library not found" Error
```bash
# Check library path
otool -L macos/filemanager-macos-amd64

# Fix with install_name_tool (usually done by build script)
install_name_tool -change @rpath/libfs_operations_core.dylib \
  @executable_path/../Frameworks/libfs_operations_core.dylib \
  macos/dist/amd64/FileManager.app/Contents/MacOS/filemanager
```

### Homebrew Formula Fails Audit
```bash
# Common issues and fixes:
# 1. SHA256 mismatch - recalculate checksums
# 2. URL not accessible - ensure GitHub release is public
# 3. Cask naming - should be lowercase with hyphens
# 4. Version mismatch - ensure VERSION file matches release tag
```

---

## Clean Build

To start fresh:
```bash
# Remove all macOS build artifacts
rm -rf macos/

# Remove Rust build cache
rm -rf rust_ffi/target/

# Remove Go build cache
cd file_manager && go clean && cd ..

# Rebuild everything
bash scripts/build-macos-all.sh
```

---

## Testing Checklist

Before distributing:

- [ ] Both architectures build successfully
- [ ] DMGs open correctly
- [ ] Apps launch without errors
- [ ] Web interface loads (http://localhost:8080)
- [ ] File operations work (copy, move, delete)
- [ ] Frontend resources load correctly
- [ ] Homebrew formulas install successfully
- [ ] Uninstall works cleanly
- [ ] Code signing valid (if applicable)
- [ ] Notarization successful (if applicable)

---

## Support

For issues with the build system:
- Check `scripts/build-macos.sh` for build logic
- Check `scripts/generate-homebrew-formula.sh` for formula generation
- Review logs in terminal output
- Test individual components separately
