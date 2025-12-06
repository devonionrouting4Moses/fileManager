# 🚀 FileManager Build Guide - Quick Start

## TL;DR - Just Run This

```bash
cd /path/to/FileManager/file_manager_v1

# Build for all platforms
bash scripts/build-all-platforms.sh all

# Or specific platforms
bash scripts/build-all-platforms.sh linux-amd64
bash scripts/build-all-platforms.sh linux-amd64 windows-amd64
```

---

## 📊 How build-all-platforms.sh Works

### The Script is an **ORCHESTRATOR**

It's like a **master conductor** that:
1. Takes your platform requests
2. Calls the right builder scripts
3. Tracks success/failure
4. Shows you what was created

### Visual Flow

```
You run:
  bash build-all-platforms.sh linux-amd64 windows-amd64
                                    ↓
Script parses arguments and sets flags:
  BUILD_LINUX_AMD64=1
  BUILD_WINDOWS=1
                                    ↓
Script calls appropriate builders:
  bash build-linux-amd64.sh  ──→ Creates DEB, TAR.GZ, APT, PKGBUILD
  bash build-windows-amd64.sh ──→ Creates EXE, ZIP
                                    ↓
Script shows summary:
  ✅ All requested builds completed!
  Build artifacts:
    - Linux amd64: filemanager_2.0.0_amd64.deb, ...
    - Windows: filemanager.exe, ...
```

---

## 🎯 Common Commands

### Build for Single Platform
```bash
# Linux AMD64
bash scripts/build-all-platforms.sh linux-amd64

# Linux ARM64
bash scripts/build-all-platforms.sh linux-arm64

# Windows
bash scripts/build-all-platforms.sh windows-amd64

# HarmonyOS
bash scripts/build-all-platforms.sh harmonyos
```

### Build for Multiple Platforms
```bash
# Linux (both architectures)
bash scripts/build-all-platforms.sh linux-amd64 linux-arm64

# Linux and Windows
bash scripts/build-all-platforms.sh linux-amd64 windows-amd64

# All platforms
bash scripts/build-all-platforms.sh all
```

### Direct Builder Calls (Skip Orchestrator)
```bash
# Same as: build-all-platforms.sh linux-amd64
bash scripts/build-linux-amd64.sh

# Same as: build-all-platforms.sh linux-arm64
bash scripts/build-linux-arm64.sh

# Same as: build-all-platforms.sh windows-amd64
bash scripts/build-windows-amd64.sh

# Same as: build-all-platforms.sh harmonyos
bash scripts/build-harmonyos.sh
```

---

## 📦 What Gets Created

### After: `bash build-all-platforms.sh all`

```
Project Root/
├── filemanager_2.0.0_amd64.deb              ← Install on Linux AMD64
├── filemanager_2.0.0_arm64.deb              ← Install on Linux ARM64
├── filemanager-2.0.0-linux-amd64.tar.gz     ← Portable Linux AMD64
├── filemanager-2.0.0-linux-arm64.tar.gz     ← Portable Linux ARM64
├── filemanager.exe                          ← Run on Windows
├── filemanager-2.0.0-windows-amd64.zip      ← Portable Windows
├── arch/PKGBUILD                            ← Arch Linux package
└── harmonyos/entry/build/outputs/hap/       ← HarmonyOS package
```

---

## 🔒 Security Features Built-In

✅ **Frontend Directory**
- Single stable location: `/usr/share/filemanager/frontend`
- Read-only for regular users
- Configurable via environment variable

✅ **File Permissions**
- Directories: `755` (rwxr-xr-x)
- Files: `644` (rw-r--r--)
- Ownership: `root:root`

✅ **Configuration**
- System config: `/etc/filemanager/config.yaml`
- User config: `~/.config/filemanager/config.yaml`
- Backed up before removal
- Preserved after uninstall

✅ **Verification**
- Run: `filemanager --frontend-info`
- Run: `bash scripts/verify-security.sh`

---

## 🚀 Complete Workflow

```bash
# Step 1: Navigate to project
cd /path/to/FileManager/file_manager_v1

# Step 2: Update version (optional)
echo "2.1.0" > VERSION

# Step 3: Build all platforms
bash scripts/build-all-platforms.sh all

# Step 4: Test DEB installation
sudo apt install ./filemanager_2.0.0_amd64.deb

# Step 5: Verify installation
filemanager --frontend-info
bash scripts/verify-security.sh

# Step 6: Uninstall
sudo apt remove filemanager

# Step 7: Upload to repositories
# - GitHub Releases
# - APT Repository
# - Snap Store
# - Homebrew
# - Windows releases
```

---

## 📚 Detailed Documentation

For more information, see:

- **HOW_TO_BUILD.md** - Complete build guide with examples
- **BUILD_ARCHITECTURE.txt** - Visual diagrams and flow charts
- **BUILD_SYSTEM_INTEGRATION.md** - Build system details
- **FRONTEND_RESOLVER_IMPLEMENTATION.md** - Frontend resolver details
- **PRODUCTION_INSTALLERS.md** - Installer scripts for all platforms
- **SECURITY_IMPLEMENTATION_SUMMARY.md** - Security overview

---

## ✅ Quick Checklist

- [ ] Navigate to project root
- [ ] Update VERSION file (if needed)
- [ ] Run: `bash scripts/build-all-platforms.sh all`
- [ ] Verify builds: `ls -lh filemanager_*.deb filemanager.exe`
- [ ] Test DEB: `sudo apt install ./filemanager_2.0.0_amd64.deb`
- [ ] Verify: `filemanager --frontend-info`
- [ ] Security check: `bash scripts/verify-security.sh`
- [ ] Uninstall: `sudo apt remove filemanager`
- [ ] Upload to repositories

---

## 🐛 Troubleshooting

**Build fails?**
```bash
# Check error message
# Fix the issue
# Run again
bash scripts/build-all-platforms.sh all
```

**Missing dependencies?**
```bash
sudo apt-get install -y build-essential cargo rustc golang-go
```

**Want to rebuild one platform?**
```bash
bash scripts/build-linux-amd64.sh
```

**Need help?**
```bash
bash scripts/build-all-platforms.sh
```

---

## 🎯 Key Takeaway

**build-all-platforms.sh** is a **smart orchestrator** that:
- Accepts platform names as arguments
- Calls the right builder for each platform
- Tracks success/failure
- Shows you what was created

**You can also call builders directly** if you only need one platform (faster).

---

**Version**: 2.0.0  
**Status**: ✅ Production Ready  
**Last Updated**: December 6, 2025
