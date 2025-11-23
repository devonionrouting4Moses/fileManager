# 📦 Building FileManager Snap Package

Complete guide for building and testing the FileManager snap package.

---

## Prerequisites

### Install Snapcraft

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y snapcraft

# Or using snap
sudo snap install snapcraft --classic
```

### Install LXD (for isolated builds)

```bash
# Ubuntu
sudo apt-get install -y lxd

# Initialize LXD
lxd init

# Add your user to lxd group
sudo usermod -aG lxd $USER
newgrp lxd
```

### Install Build Dependencies

```bash
sudo apt-get install -y \
  build-essential \
  golang-go \
  rustc \
  cargo \
  pkg-config \
  libssl-dev
```

---

## Building the Snap Package

### Option 1: Using LXD (Recommended - Isolated Build)

```bash
cd file_manager_v1

# Build in isolated LXD container
snapcraft

# This will:
# 1. Create an LXD container
# 2. Install all dependencies
# 3. Build Rust library
# 4. Build Go binary
# 5. Package everything into .snap file
```

### Option 2: Using --destructive-mode (Local Build)

```bash
cd file_manager_v1

# Build directly on your system (modifies your environment)
snapcraft --destructive-mode

# WARNING: This modifies your system packages
# Only use if LXD is not available
```

### Option 3: Using Docker (Alternative)

```bash
# Build using Docker instead of LXD
snapcraft --use-docker
```

---

## Build Output

After successful build, you'll see:

```
✅ Snap build completed successfully
Snapped filemanager_2.0.0_amd64.snap
```

The snap package will be created as:
```
filemanager_2.0.0_amd64.snap
```

---

## Installing the Snap Package

### Local Installation (Testing)

```bash
# Install from local file
sudo snap install --dangerous --classic filemanager_2.0.0_amd64.snap

# The --dangerous flag allows installation of unsigned snaps
# The --classic flag allows access to system resources
```

### Verify Installation

```bash
# Check snap is installed
snap list | grep filemanager

# Run the application
filemanager

# Check version
filemanager --version
```

### Uninstall

```bash
sudo snap remove filemanager
```

---

## Troubleshooting Build Issues

### Issue: "could not find `Cargo.toml`"

**Cause:** Build script is looking in wrong directory

**Solution:** The fixed `snapcraft.yaml` now:
- Changes to `rust_ffi` directory before building
- Uses correct paths for Rust library
- Verifies library exists before continuing

### Issue: "Go binary not found"

**Cause:** Go build failed or binary not in expected location

**Solution:**
```bash
# Check build logs
snapcraft --verbose

# Verify Go binary was created
ls -la file_manager/

# Check CGO_LDFLAGS
echo $CGO_LDFLAGS
```

### Issue: "libfs_operations_core.so not found"

**Cause:** Rust library wasn't built or is in wrong location

**Solution:**
```bash
# Build Rust library manually
cd rust_ffi
cargo build --release -p fs-operations-core
cd ..

# Verify library exists
ls -la rust_ffi/target/release/libfs_operations_core.so
```

### Issue: "Snap confinement error"

**Cause:** App doesn't have required permissions

**Solution:** The `snapcraft.yaml` includes these plugs:
- `home` - Access to home directory
- `removable-media` - Access to USB drives
- `network` - Network access
- `desktop` - Desktop integration
- `x11`, `wayland` - Display server access

If you need more permissions, add to the `plugs` section.

### Issue: "Build takes too long"

**Cause:** First build compiles everything from scratch

**Solution:**
```bash
# Subsequent builds are faster (cached)
# Or use --use-lxd to use persistent container
snapcraft --use-lxd
```

---

## Build Configuration Details

### snapcraft.yaml Structure

```yaml
name: filemanager              # Package name
base: core24                   # Ubuntu 24.04 base
version: '2.0.0'               # Version
confinement: strict            # Security confinement
platforms:                     # Supported architectures
  amd64:                       # 64-bit Intel/AMD
  arm64:                       # 64-bit ARM (Raspberry Pi, etc.)
```

### Build Parts

1. **rust-deps**
   - Installs Rust toolchain
   - Runs once, cached for subsequent builds

2. **filemanager**
   - Builds Rust library
   - Builds Go binary
   - Installs files to snap directory
   - Includes documentation

3. **desktop-file**
   - Adds desktop integration
   - Provides application icon

### Environment Variables

During build:
```bash
CGO_ENABLED=1                  # Enable CGO
CGO_LDFLAGS="-L... -l..."      # Link Rust library
LD_LIBRARY_PATH=...            # Runtime library path
```

---

## Publishing to Snap Store

### Prerequisites

1. Create Ubuntu One account: https://login.ubuntu.com
2. Register snap name: `snapcraft register filemanager`
3. Login to snapcraft: `snapcraft login`

### Build for Release

```bash
# Build without --dangerous flag
snapcraft

# This creates a signed snap (if you have signing keys)
```

### Upload to Store

```bash
# Upload snap to store
snapcraft upload filemanager_2.0.0_amd64.snap --release=stable

# Or release to different channels
snapcraft upload filemanager_2.0.0_amd64.snap --release=candidate
snapcraft upload filemanager_2.0.0_amd64.snap --release=beta
snapcraft upload filemanager_2.0.0_amd64.snap --release=edge
```

### Check Upload Status

```bash
snapcraft status filemanager
```

---

## Testing the Snap

### Basic Functionality Test

```bash
# Test CLI mode
filemanager

# Test web mode
filemanager --web

# Test version
filemanager --version

# Test help
filemanager --help
```

### File Operations Test

```bash
# Create test directory
mkdir -p ~/snap-test

# Test file creation
filemanager create-file ~/snap-test/test.txt

# Test folder creation
filemanager create-folder ~/snap-test/subfolder

# Test listing
filemanager list ~/snap-test
```

### Confinement Test

```bash
# Test home directory access
filemanager list ~/

# Test removable media access
filemanager list /media/

# Test network access (if applicable)
# Test desktop integration
```

---

## Build Optimization

### Reduce Snap Size

Current size: ~50-100 MB (depends on Rust dependencies)

To reduce:
```yaml
# In snapcraft.yaml, add:
compression: lz4  # Faster decompression
```

### Faster Builds

```bash
# Use persistent container
snapcraft --use-lxd

# This keeps the container between builds
# Subsequent builds are much faster
```

### Parallel Builds

```bash
# Build for multiple architectures
snapcraft --build-for=amd64,arm64
```

---

## Snap Lifecycle

### Build Stages

1. **Pull** - Download sources
2. **Build** - Compile code
3. **Stage** - Prepare files
4. **Prime** - Final packaging
5. **Snap** - Create .snap file

### Clean Build

```bash
# Remove all build artifacts
snapcraft clean

# Remove specific part
snapcraft clean filemanager

# Remove and rebuild
snapcraft clean filemanager
snapcraft
```

---

## Debugging

### Verbose Output

```bash
snapcraft --verbose
```

### Check Build Logs

```bash
# View full build log
cat ~/.local/state/snapcraft/log/snapcraft-*.log

# Or use tail for real-time
tail -f ~/.local/state/snapcraft/log/snapcraft-*.log
```

### Enter Build Container

```bash
# For LXD builds, you can enter the container
lxc list

# Find the snapcraft container
lxc exec snapcraft-filemanager -- bash
```

### Manual Build Steps

```bash
# Build Rust library manually
cd rust_ffi
cargo build --release -p fs-operations-core

# Build Go binary manually
cd file_manager
CGO_ENABLED=1 \
  CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
  go build -ldflags="-s -w" -o ../filemanager ./cmd/app
```

---

## Common Commands

```bash
# Build snap
snapcraft

# Build with verbose output
snapcraft --verbose

# Clean and rebuild
snapcraft clean && snapcraft

# Install locally
sudo snap install --dangerous --classic filemanager_2.0.0_amd64.snap

# Test installed snap
filemanager --version

# Remove snap
sudo snap remove filemanager

# Check snap info
snap info filemanager

# View snap logs
snap logs filemanager -f

# List all snaps
snap list

# Search snap store
snap search filemanager
```

---

## Next Steps

1. ✅ Build snap package locally
2. ✅ Test all functionality
3. ✅ Test on different systems (if possible)
4. 📋 Create Ubuntu One account
5. 📋 Register snap name
6. 📋 Upload to snap store
7. 📋 Publish to stable channel

---

## Resources

- [Snapcraft Documentation](https://snapcraft.io/docs)
- [Snapcraft.yaml Reference](https://snapcraft.io/docs/snapcraft-yaml-reference)
- [Snap Store](https://snapcraft.io/store)
- [Ubuntu One Account](https://login.ubuntu.com)

---

**Happy Snapping! 📦**
