# 🚀 Build FileManager Now

Quick start guide to build for different platforms.

## Prerequisites

```bash
# Install build tools
sudo apt-get update
sudo apt-get install -y build-essential pkg-config libssl-dev

# For ARM64 cross-compilation
sudo apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu

# For Windows builds on Linux
# Install MinGW-w64 (optional, for cross-compilation)
```

## Build Scripts

### Linux amd64 (DEB, TAR.GZ, APT, PKGBUILD)

```bash
chmod +x build-linux-amd64.sh
./build-linux-amd64.sh
```

Creates:
- `filemanager_2.0.0_amd64.deb` - Debian/Ubuntu package
- `filemanager-2.0.0-linux-amd64.tar.gz` - Portable archive
- `filemanager-apt_2.0.0_amd64.deb` - APT package
- `arch/PKGBUILD` - Arch Linux package

### Linux arm64 (DEB, TAR.GZ, APT, PKGBUILD)

```bash
chmod +x build-linux-arm64.sh
./build-linux-arm64.sh
```

Creates:
- `filemanager_2.0.0_arm64.deb` - Debian/Ubuntu package
- `filemanager-2.0.0-linux-arm64.tar.gz` - Portable archive
- `filemanager-apt_2.0.0_arm64.deb` - APT package
- `arch-arm64/PKGBUILD` - Arch Linux package

### Windows amd64 (EXE, ZIP)

```bash
chmod +x build-windows-amd64.sh
./build-windows-amd64.sh
```

Creates:
- `filemanager.exe` - Windows executable
- `filemanager-2.0.0-windows-amd64.zip` - Portable archive with installer

### HarmonyOS (HAP)

```bash
chmod +x build-harmonyos.sh
./build-harmonyos.sh
```

Creates:
- `harmonyos/` - HarmonyOS project structure
- `harmonyos/entry/build/outputs/hap/entry-default-signed.hap` - HarmonyOS package

## Build All Platforms

```bash
chmod +x build-all-platforms.sh

# Build all
./build-all-platforms.sh all

# Build specific platforms
./build-all-platforms.sh linux-amd64 windows-amd64
./build-all-platforms.sh linux-arm64
./build-all-platforms.sh harmonyos
```

## Install Packages

### DEB Package

```bash
sudo dpkg -i filemanager_2.0.0_amd64.deb
```

### TAR.GZ Archive

```bash
tar xzf filemanager-2.0.0-linux-amd64.tar.gz
cd filemanager-2.0.0-linux-amd64
./install.sh
```

### Snap Package

```bash
snapcraft pack
sudo snap install --dangerous --classic filemanager_2.0.0_amd64.snap
```

### Arch Linux

```bash
cd arch
makepkg -si
```

### Windows

1. Extract `filemanager-2.0.0-windows-amd64.zip`
2. Run `install.bat` as Administrator
3. Or manually copy files to `C:\Program Files\FileManager`

## Verify Installation

```bash
filemanager --version
filemanager --help
```

## Build Status

- ✅ Linux amd64 - Ready
- ✅ Linux arm64 - Ready (requires cross-compilation tools)
- ✅ Windows amd64 - Ready
- ✅ Snap - Ready (use `snapcraft pack`)
- 📋 HarmonyOS - Requires OpenHarmony SDK
- ⏳ macOS - Build at school

## Troubleshooting

### "Command not found: aarch64-linux-gnu-gcc"

Install ARM64 cross-compilation tools:
```bash
sudo apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu
```

### "Cargo not found"

Install Rust:
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
```

### "Go not found"

Install Go:
```bash
# Download from https://golang.org/dl/
# Or use package manager
sudo apt-get install -y golang-go
```

### Windows build on Linux

The scripts use cross-compilation. Make sure you have the Windows toolchain:
```bash
# MinGW-w64 (optional)
sudo apt-get install -y mingw-w64
```

## Next Steps

1. Run build scripts for your target platforms
2. Test the binaries
3. Create releases on GitHub
4. Upload to package managers (Snap, AUR, Homebrew, etc.)

---

**Happy Building! 🚀**
