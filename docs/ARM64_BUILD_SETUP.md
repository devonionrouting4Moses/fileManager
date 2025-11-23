# ARM64 Cross-Compilation Setup Guide

## Overview

Building FileManager for ARM64 (aarch64) on an x86_64 system requires cross-compilation toolchains. This guide covers the setup for all major Linux distributions.

## Prerequisites

- Rust toolchain with ARM64 target support
- ARM64 cross-compilation linker and tools
- Go 1.21+ (for Go binary compilation)

## Installation by Distribution

### Ubuntu / Debian

```bash
sudo apt-get update
sudo apt-get install -y \
  gcc-aarch64-linux-gnu \
  g++-aarch64-linux-gnu \
  binutils-aarch64-linux-gnu \
  pkg-config
```

### Fedora / RHEL / CentOS

```bash
sudo dnf install -y \
  gcc-aarch64-linux-gnu \
  gcc-c++-aarch64-linux-gnu \
  binutils-aarch64-linux-gnu \
  pkgconfig
```

### Arch Linux

```bash
sudo pacman -S aarch64-linux-gnu-gcc aarch64-linux-gnu-binutils
```

### Alpine Linux

```bash
apk add gcc-aarch64 g++-aarch64 binutils-aarch64
```

## Verify Installation

After installation, verify the toolchain is available:

```bash
aarch64-linux-gnu-gcc --version
aarch64-linux-gnu-g++ --version
aarch64-linux-gnu-ar --version
```

You should see version information for each tool.

## Rust Setup

Add the ARM64 target to Rust:

```bash
rustup target add aarch64-unknown-linux-gnu
```

Configure Rust to use the ARM64 linker by adding to `~/.cargo/config.toml`:

```toml
[target.aarch64-unknown-linux-gnu]
linker = "aarch64-linux-gnu-gcc"
ar = "aarch64-linux-gnu-ar"
```

## Building for ARM64

### Using the Build Script

Simply run the ARM64 build script:

```bash
cd file_manager_v1/scripts
./build-all-platforms.sh linux-arm64
```

The script will:
1. Check for the ARM64 linker
2. Configure Rust for ARM64
3. Build the Rust library for ARM64
4. Build the Go binary for ARM64
5. Create DEB, TAR.GZ, and PKGBUILD packages

### Manual Build

If you prefer to build manually:

```bash
cd file_manager_v1

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
  go build -ldflags="-s -w" -o ../filemanager-arm64 ./cmd/app
cd ..
```

## Troubleshooting

### Error: "aarch64-linux-gnu-gcc not found"

**Solution:** Install the ARM64 cross-compilation toolchain for your distribution (see Installation section above).

### Error: "Relocations in generic ELF (EM: 183)"

**Cause:** The linker can't process ARM64 object files because the cross-compilation toolchain isn't installed.

**Solution:** Install the ARM64 toolchain and ensure `aarch64-linux-gnu-gcc` is in your PATH.

### Error: "cannot find -lc"

**Cause:** Missing ARM64 C library files.

**Solution:** Install `libc6-arm64-cross` or equivalent for your distribution:

```bash
# Ubuntu/Debian
sudo apt-get install -y libc6-arm64-cross

# Fedora
sudo dnf install -y glibc-aarch64-linux-gnu
```

### Error: "RUSTFLAGS not recognized"

**Solution:** Make sure you're using bash, not sh:

```bash
bash ./build-linux-arm64.sh
```

## Output Files

After successful build, you'll find:

- `filemanager-arm64` - ARM64 binary
- `filemanager_2.0.0_arm64.deb` - Debian package
- `filemanager-2.0.0-linux-arm64.tar.gz` - Universal archive
- `arch-arm64/PKGBUILD` - Arch Linux package definition

## Testing ARM64 Builds

### On ARM64 Hardware

Simply run the binary:

```bash
./filemanager-arm64
```

### Using QEMU (on x86_64)

Install QEMU and run the binary under ARM64 emulation:

```bash
# Ubuntu/Debian
sudo apt-get install -y qemu-user-static

# Run ARM64 binary
qemu-aarch64-static ./filemanager-arm64
```

### In Docker

Use an ARM64 Docker image to test:

```bash
docker run --rm -v $(pwd):/app arm64v8/ubuntu:24.04 /app/filemanager-arm64
```

## Performance Notes

- Cross-compilation is slower than native compilation
- First build may take 5-10 minutes
- Subsequent builds are faster due to caching
- The resulting binary runs at full native speed on ARM64 hardware

## References

- [Rust Platform Support](https://doc.rust-lang.org/nightly/rustc/platform-support.html)
- [GNU Arm Embedded Toolchain](https://developer.arm.com/tools-and-software/open-source-software/gnu-toolchain)
- [Cross-Compilation Guide](https://wiki.osdev.org/Cross-Compiler)
