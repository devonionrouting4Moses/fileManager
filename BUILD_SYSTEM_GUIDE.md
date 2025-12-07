# 🚀 FileManager Build System Guide

## Version Management

### Using version-manager.sh

All version updates across the entire project are handled automatically by `scripts/version-manager.sh`.

```bash
# Get current version
bash scripts/version-manager.sh get

# Set new version (updates ALL files automatically)
bash scripts/version-manager.sh set 2.1.0

# Or bump automatically
bash scripts/version-manager.sh bump-patch   # 2.0.2 → 2.0.3
bash scripts/version-manager.sh bump-minor   # 2.0.2 → 2.1.0
bash scripts/version-manager.sh bump-major   # 2.0.2 → 3.0.0

# List all version files
bash scripts/version-manager.sh list
```

### Files Updated Automatically

When you run `bash scripts/version-manager.sh set X.Y.Z`, these files are updated:

- VERSION
- snap/snapcraft.yaml
- file_manager/pkg/version/version.go
- rust_ffi/Cargo.toml (all files)
- README.md
- debian/control
- flatpak/com.filemanager.FileManager.yaml
- homebrew/filemanager.rb
- windows/installer.nsi
- arch/PKGBUILD
- arch-arm64/PKGBUILD
- Documentation files (if they exist)

---

## Building for Different Platforms

### Linux (AMD64 & ARM64)

```bash
# Build both architectures
bash scripts/build-all-platforms.sh linux-amd64 linux-arm64

# Or individually
bash scripts/build-linux-amd64.sh
bash scripts/build-linux-arm64.sh
```

**Output**:
- `filemanager` / `filemanager-arm64` (binary)
- `filemanager_X.Y.Z_amd64.deb` / `filemanager_X.Y.Z_arm64.deb` (DEB package)
- `filemanager-X.Y.Z-linux-amd64.tar.gz` / `filemanager-X.Y.Z-linux-arm64.tar.gz` (TAR archive)
- `arch/PKGBUILD` / `arch-arm64/PKGBUILD` (Arch Linux package)

### Windows (AMD64)

```bash
bash scripts/build-all-platforms.sh windows-amd64
# Or
bash scripts/build-windows-amd64.sh
```

**Output**:
- `filemanager.exe` (binary)
- `FileManager-X.Y.Z-Setup.exe` (NSIS installer)
- `filemanager-X.Y.Z-windows-amd64.zip` (ZIP archive)

### macOS (on your Mac)

```bash
# Build both architectures
bash scripts/build-macos.sh amd64
bash scripts/build-macos.sh arm64
```

**Output**:
- `filemanager-macos-amd64` / `filemanager-macos-arm64` (binary)
- `dist/FileManager.app` (application bundle)
- `dist/FileManager-X.Y.Z-macos-amd64.dmg` / `dist/FileManager-X.Y.Z-macos-arm64.dmg` (DMG installer)
- `dist/filemanager-X.Y.Z-macos-amd64.tar.gz` / `dist/filemanager-X.Y.Z-macos-arm64.tar.gz` (TAR archive)
- `homebrew/filemanager.rb` (Homebrew formula - needs SHA256 update)

### All Platforms (except Snap & macOS)

```bash
bash scripts/build-all-platforms.sh linux-amd64 linux-arm64 windows-amd64
```

---

## macOS Homebrew Formula

After building on macOS, update the Homebrew formula:

```bash
# Get checksums
shasum -a 256 dist/filemanager-X.Y.Z-macos-amd64.tar.gz
shasum -a 256 dist/filemanager-X.Y.Z-macos-arm64.tar.gz

# Update homebrew/filemanager.rb with the checksums
# Then test:
brew install --formula homebrew/filemanager.rb
```

---

## Frontend Strategy

All builds follow the **Golden Rule**:

### System-Level Installs (DEB, RPM, Homebrew, NSIS, Flatpak)
- ✅ Bundle frontend during packaging
- ✅ Don't create at runtime
- ✅ Serve pre-bundled files with restrictive permissions

### User-Level Installs (Manual, Portable)
- ✅ Create frontend at runtime if missing
- ✅ User can modify files in `~/.local/share/filemanager/frontend`

---

## Build Verification

### Check all versions are consistent

```bash
bash scripts/version-manager.sh list
```

### Verify build scripts

- `scripts/build-all-platforms.sh` - Multi-platform orchestrator
- `scripts/build-linux-amd64.sh` - Linux AMD64 build
- `scripts/build-linux-arm64.sh` - Linux ARM64 build (cross-compilation)
- `scripts/build-windows-amd64.sh` - Windows AMD64 build
- `scripts/build-macos.sh` - macOS build (Intel & Apple Silicon)

---

## Quick Commands

```bash
# Update version
bash scripts/version-manager.sh set 2.1.0

# Build everything (Linux + Windows)
bash scripts/build-all-platforms.sh linux-amd64 linux-arm64 windows-amd64

# Build Linux only
bash scripts/build-all-platforms.sh linux-amd64 linux-arm64

# Build Windows only
bash scripts/build-all-platforms.sh windows-amd64

# On macOS: Build both architectures
bash scripts/build-macos.sh amd64 && bash scripts/build-macos.sh arm64
```

---

## Platform Support Matrix

| Platform | Build Ready | Package Format | Notes |
|----------|-------------|----------------|-------|
| Linux AMD64 | ✅ Yes | DEB, TAR.GZ, Arch | Ready to build |
| Linux ARM64 | ✅ Yes | DEB, TAR.GZ, Arch | Cross-compilation ready |
| Windows AMD64 | ✅ Yes | EXE, ZIP, NSIS | Ready to build |
| macOS Intel | ✅ Yes | DMG, TAR.GZ, Homebrew | Build on Mac, update SHA256 |
| macOS ARM64 | ✅ Yes | DMG, TAR.GZ, Homebrew | Build on Mac, update SHA256 |
| Flatpak | ✅ Yes | YAML manifest | Ready to build |
| Arch x86_64 | ✅ Yes | PKGBUILD | Ready to build |
| Arch ARM64 | ✅ Yes | PKGBUILD | Ready to build |

---

**All systems ready for production builds.**
