# 🔨 How to Build FileManager - Complete Guide

## Quick Overview

`build-all-platforms.sh` is an **orchestrator script** that:
- Accepts platform arguments
- Calls the appropriate platform-specific build scripts
- Tracks success/failure
- Shows final summary

Think of it as a **master controller** that delegates work to specialized builders.

---

## 🎯 How It Works

### Step-by-Step Flow

```
User runs: bash build-all-platforms.sh [platforms]
    ↓
Script reads VERSION file
    ↓
Script parses arguments (linux-amd64, linux-arm64, windows-amd64, harmonyos, all)
    ↓
For each requested platform:
    ├─ Set flag to 1 (BUILD_LINUX_AMD64=1, etc.)
    ↓
For each enabled platform:
    ├─ Call build-linux-amd64.sh
    ├─ Call build-linux-arm64.sh
    ├─ Call build-windows-amd64.sh
    └─ Call build-harmonyos.sh
    ↓
Show summary of created artifacts
```

---

## 📋 Usage Examples

### 1. Build for Single Platform

```bash
cd /path/to/FileManager/file_manager_v1
bash scripts/build-all-platforms.sh linux-amd64
```

**What happens**:
- Only `build-linux-amd64.sh` is called
- Creates:
  - `filemanager_2.0.0_amd64.deb`
  - `filemanager-2.0.0-linux-amd64.tar.gz`
  - `filemanager-apt_2.0.0_amd64.deb`
  - `arch/PKGBUILD`

---

### 2. Build for Multiple Platforms

```bash
bash scripts/build-all-platforms.sh linux-amd64 linux-arm64 windows-amd64
```

**What happens**:
- Calls `build-linux-amd64.sh`
- Calls `build-linux-arm64.sh`
- Calls `build-windows-amd64.sh`
- Skips `build-harmonyos.sh`

**Creates**:
- Linux AMD64 packages
- Linux ARM64 packages
- Windows packages

---

### 3. Build for All Platforms

```bash
bash scripts/build-all-platforms.sh all
```

**What happens**:
- Calls ALL platform-specific scripts
- Creates packages for:
  - Linux AMD64
  - Linux ARM64
  - Windows AMD64
  - HarmonyOS

**Output**:
```
🔨 FileManager v2.0.0 - Multi-Platform Build
==================================================

Building for Linux amd64...
✅ Linux amd64 build completed

Building for Linux arm64...
✅ Linux arm64 build completed

Building for Windows amd64...
✅ Windows amd64 build completed

Building for HarmonyOS...
✅ HarmonyOS build completed

🎉 All requested builds completed!

Build artifacts:
  - Linux amd64: filemanager_2.0.0_amd64.deb, filemanager-2.0.0-linux-amd64.tar.gz
  - Linux arm64: filemanager_2.0.0_arm64.deb, filemanager-2.0.0-linux-arm64.tar.gz
  - Windows: filemanager.exe, filemanager-2.0.0-windows-amd64.zip
  - HarmonyOS: harmonyos/entry/build/outputs/hap/
```

---

### 4. Show Help

```bash
bash scripts/build-all-platforms.sh
```

**Output**:
```
Usage: build-all-platforms.sh [linux-amd64] [linux-arm64] [windows-amd64] [harmonyos] [all]

Examples:
  build-all-platforms.sh linux-amd64          # Build for Linux amd64 only
  build-all-platforms.sh linux-amd64 windows-amd64  # Build for Linux and Windows
  build-all-platforms.sh all                  # Build for all platforms
```

---

## 🔍 Understanding the Code

### Argument Parsing (Lines 40-65)

```bash
for arg in "$@"; do
    case $arg in
        linux-amd64)
            BUILD_LINUX_AMD64=1
            ;;
        linux-arm64)
            BUILD_LINUX_ARM64=1
            ;;
        windows-amd64)
            BUILD_WINDOWS=1
            ;;
        harmonyos)
            BUILD_HARMONYOS=1
            ;;
        all)
            BUILD_LINUX_AMD64=1
            BUILD_LINUX_ARM64=1
            BUILD_WINDOWS=1
            BUILD_HARMONYOS=1
            ;;
    esac
done
```

**What it does**:
- Loops through all arguments you pass
- Sets flags (0 or 1) for each platform
- If you pass `all`, sets all flags to 1

**Example**:
```bash
bash build-all-platforms.sh linux-amd64 windows-amd64

# Results in:
# BUILD_LINUX_AMD64=1
# BUILD_LINUX_ARM64=0
# BUILD_WINDOWS=1
# BUILD_HARMONYOS=0
```

---

### Platform Building (Lines 70-115)

```bash
if [ $BUILD_LINUX_AMD64 -eq 1 ]; then
    echo -e "${BLUE}Building for Linux amd64...${NC}"
    if bash "${SCRIPT_DIR}/build-linux-amd64.sh"; then
        echo -e "${GREEN}✅ Linux amd64 build completed${NC}"
    else
        echo -e "${RED}❌ Linux amd64 build failed${NC}"
        exit 1
    fi
    echo ""
fi
```

**What it does**:
- Checks if flag is set to 1
- If yes, calls the platform-specific script
- If build succeeds, shows ✅
- If build fails, shows ❌ and exits

---

## 🏗️ What Each Platform Script Does

### build-linux-amd64.sh
```bash
bash scripts/build-linux-amd64.sh
```

**Creates**:
1. Builds Rust library for AMD64
2. Builds Go binary for Linux AMD64
3. Creates DEB package with security scripts
4. Creates TAR.GZ portable archive
5. Creates APT package
6. Creates Arch PKGBUILD

**Output Files**:
- `filemanager_2.0.0_amd64.deb`
- `filemanager-2.0.0-linux-amd64.tar.gz`
- `filemanager-apt_2.0.0_amd64.deb`
- `arch/PKGBUILD`

---

### build-linux-arm64.sh
```bash
bash scripts/build-linux-arm64.sh
```

**Creates**:
1. Sets up ARM64 cross-compilation toolchain
2. Builds Rust library for ARM64
3. Builds Go binary for Linux ARM64
4. Creates DEB package with security scripts
5. Creates TAR.GZ portable archive
6. Creates APT package

**Output Files**:
- `filemanager_2.0.0_arm64.deb`
- `filemanager-2.0.0-linux-arm64.tar.gz`
- `filemanager-apt_2.0.0_arm64.deb`

---

### build-windows-amd64.sh
```bash
bash scripts/build-windows-amd64.sh
```

**Creates**:
1. Builds Go binary for Windows AMD64
2. Creates Windows executable
3. Creates ZIP portable archive
4. (Optional) Creates NSIS installer

**Output Files**:
- `filemanager.exe`
- `filemanager-2.0.0-windows-amd64.zip`

---

### build-harmonyos.sh
```bash
bash scripts/build-harmonyos.sh
```

**Creates**:
1. Builds for HarmonyOS platform
2. Creates HAP (HarmonyOS App Package)

**Output Files**:
- `harmonyos/entry/build/outputs/hap/`

---

## 🚀 Complete Build Workflow

### Scenario: Build for Production Release

```bash
# Step 1: Navigate to project
cd /home/devchigarlicmosesrouting/Documents/LIVE_PROJECTS/FileManager/file_manager_v1

# Step 2: Update version (optional)
echo "2.1.0" > VERSION

# Step 3: Build for all platforms
bash scripts/build-all-platforms.sh all

# Step 4: Verify builds
ls -lh filemanager_*.deb
ls -lh filemanager-*.tar.gz
ls -lh filemanager.exe
ls -lh filemanager-*.zip

# Step 5: Test DEB installation
sudo apt install ./filemanager_2.0.0_amd64.deb
filemanager --frontend-info
bash scripts/verify-security.sh

# Step 6: Test uninstallation
sudo apt remove filemanager

# Step 7: Upload to repositories
# - Upload .deb files to APT repository
# - Upload .tar.gz to releases
# - Upload .exe and .zip to Windows releases
# - Upload HAP to HarmonyOS app store
```

---

## 📊 Build Artifacts Summary

### After Running: `bash build-all-platforms.sh all`

```
Project Root/
├── filemanager_2.0.0_amd64.deb              ← DEB for AMD64
├── filemanager_2.0.0_arm64.deb              ← DEB for ARM64
├── filemanager-apt_2.0.0_amd64.deb          ← APT for AMD64
├── filemanager-apt_2.0.0_arm64.deb          ← APT for ARM64
├── filemanager-2.0.0-linux-amd64.tar.gz     ← Portable AMD64
├── filemanager-2.0.0-linux-arm64.tar.gz     ← Portable ARM64
├── filemanager.exe                          ← Windows executable
├── filemanager-2.0.0-windows-amd64.zip      ← Windows portable
├── arch/
│   └── PKGBUILD                             ← Arch Linux package
└── harmonyos/
    └── entry/build/outputs/hap/             ← HarmonyOS package
```

---

## 🔧 Customization

### Change Version
```bash
# Edit VERSION file
echo "2.1.0" > VERSION

# Rebuild
bash scripts/build-all-platforms.sh all
```

### Build Only Specific Platforms
```bash
# Linux only
bash scripts/build-all-platforms.sh linux-amd64 linux-arm64

# Windows only
bash scripts/build-all-platforms.sh windows-amd64

# Linux AMD64 and Windows
bash scripts/build-all-platforms.sh linux-amd64 windows-amd64
```

### Call Individual Scripts Directly
```bash
# If you only want to build for one platform
bash scripts/build-linux-amd64.sh

# Or for ARM64
bash scripts/build-linux-arm64.sh

# Or for Windows
bash scripts/build-windows-amd64.sh

# Or for HarmonyOS
bash scripts/build-harmonyos.sh
```

---

## ✅ Verification Checklist

After building, verify everything:

```bash
# 1. Check DEB packages exist
dpkg -I filemanager_2.0.0_amd64.deb

# 2. Check TAR.GZ archives exist
tar -tzf filemanager-2.0.0-linux-amd64.tar.gz | head -10

# 3. Check Windows executable exists
file filemanager.exe

# 4. Check Windows ZIP exists
unzip -l filemanager-2.0.0-windows-amd64.zip | head -10

# 5. Test DEB installation
sudo apt install ./filemanager_2.0.0_amd64.deb
filemanager --frontend-info
bash scripts/verify-security.sh

# 6. Uninstall
sudo apt remove filemanager
```

---

## 🐛 Troubleshooting

### Build Fails for One Platform
```bash
# The script will stop at first failure
# Check the error message
# Fix the issue
# Run again (it will rebuild all)

bash scripts/build-all-platforms.sh all
```

### Missing Dependencies
```bash
# For Linux builds
sudo apt-get install -y build-essential cargo rustc golang-go

# For Windows builds (on Linux)
# Windows cross-compilation may require additional setup

# For HarmonyOS builds
# Requires HarmonyOS SDK installed
```

### Partial Build
```bash
# If you only want to rebuild one platform
bash scripts/build-linux-amd64.sh

# This is faster than rebuilding everything
```

---

## 📚 Related Scripts

| Script | Purpose | When to Use |
|--------|---------|------------|
| `build.sh` | Local development build | Testing locally |
| `build-all-platforms.sh` | Multi-platform orchestrator | Production releases |
| `build-linux-amd64.sh` | AMD64 builder | Single platform build |
| `build-linux-arm64.sh` | ARM64 builder | Single platform build |
| `build-windows-amd64.sh` | Windows builder | Windows releases |
| `build-harmonyos.sh` | HarmonyOS builder | HarmonyOS releases |
| `build-deb.sh` | DEB wrapper | Quick DEB builds |
| `verify-security.sh` | Security verification | Post-install checks |

---

## 🎯 Quick Reference

```bash
# Show help
bash scripts/build-all-platforms.sh

# Build for Linux AMD64 only
bash scripts/build-all-platforms.sh linux-amd64

# Build for Linux ARM64 only
bash scripts/build-all-platforms.sh linux-arm64

# Build for Windows only
bash scripts/build-all-platforms.sh windows-amd64

# Build for HarmonyOS only
bash scripts/build-all-platforms.sh harmonyos

# Build for Linux (both AMD64 and ARM64)
bash scripts/build-all-platforms.sh linux-amd64 linux-arm64

# Build for all platforms
bash scripts/build-all-platforms.sh all

# Build directly (skip orchestrator)
bash scripts/build-linux-amd64.sh
bash scripts/build-linux-arm64.sh
bash scripts/build-windows-amd64.sh
bash scripts/build-harmonyos.sh
```

---

## 🔄 Typical Release Process

```bash
# 1. Prepare
cd /path/to/FileManager/file_manager_v1
git pull origin main

# 2. Update version
echo "2.1.0" > VERSION
git add VERSION
git commit -m "Release v2.1.0"
git tag v2.1.0

# 3. Build all platforms
bash scripts/build-all-platforms.sh all

# 4. Test builds
sudo apt install ./filemanager_2.0.0_amd64.deb
filemanager --frontend-info
bash scripts/verify-security.sh
sudo apt remove filemanager

# 5. Create release notes
# Document changes, new features, bug fixes

# 6. Upload to repositories
# - GitHub Releases: .deb, .tar.gz, .exe, .zip
# - APT Repository: .deb files
# - Snap Store: snapcraft.yaml
# - Homebrew: filemanager.rb
# - Windows: .exe, .msi

# 7. Announce release
# - GitHub release notes
# - Website
# - Social media
```

---

**Status**: ✅ Complete
**Version**: 2.0.0
**Last Updated**: December 6, 2025
