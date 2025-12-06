# 🔨 Build System Integration Guide

## Overview
The FileManager build system is organized into platform-specific scripts that work together to create production packages with proper security permissions and frontend resolver integration.

---

## 📁 Build Script Architecture

```
scripts/
├── build.sh                    # Local development build
├── build-all-platforms.sh      # Multi-platform orchestrator
├── build-linux-amd64.sh        # Linux AMD64 (DEB + TAR.GZ + APT)
├── build-linux-arm64.sh        # Linux ARM64 (DEB + TAR.GZ + APT)
├── build-windows-amd64.sh      # Windows AMD64 (EXE + ZIP)
├── build-harmonyos.sh          # HarmonyOS build
├── build-deb.sh                # DEB wrapper (new)
├── verify-security.sh          # Security verification (new)
└── Set-FileManagerPermissions.ps1  # Windows ACL setup (new)

debian/
├── control                     # Package metadata
├── postinst                    # Post-installation script
├── prerm                       # Pre-removal script
├── postrm                      # Post-removal script
├── conffiles                   # Config file declarations
└── install                     # File manifest
```

---

## 🔄 Build Flow

### 1. Local Development Build
```bash
bash scripts/build.sh
```

**Output**:
- `scripts/dist/filemanager.bin` - Development binary
- `scripts/dist/libfs_operations_core.so` - Rust library
- Ready for testing with `filemanager.sh` wrapper

---

### 2. Multi-Platform Production Build
```bash
# Build all platforms
bash scripts/build-all-platforms.sh all

# Or specific platforms
bash scripts/build-all-platforms.sh linux-amd64 linux-arm64
bash scripts/build-all-platforms.sh windows-amd64
bash scripts/build-all-platforms.sh harmonyos
```

**Orchestrates**:
- `build-linux-amd64.sh` → Creates DEB + TAR.GZ + APT
- `build-linux-arm64.sh` → Creates DEB + TAR.GZ + APT
- `build-windows-amd64.sh` → Creates EXE + ZIP
- `build-harmonyos.sh` → Creates HAP bundle

---

### 3. Platform-Specific Builds

#### Linux AMD64
```bash
bash scripts/build-linux-amd64.sh
```

**Creates**:
- `filemanager_2.0.0_amd64.deb` - Debian package
- `filemanager-2.0.0-linux-amd64.tar.gz` - Portable archive
- `filemanager-apt_2.0.0_amd64.deb` - APT package
- `arch/PKGBUILD` - Arch Linux package

**Features**:
- ✓ Enhanced postinst script with security permissions
- ✓ Frontend resolver integration
- ✓ System config creation
- ✓ User config template
- ✓ Library cache update

#### Linux ARM64
```bash
bash scripts/build-linux-arm64.sh
```

**Creates**:
- `filemanager_2.0.0_arm64.deb` - Debian package
- `filemanager-2.0.0-linux-arm64.tar.gz` - Portable archive
- `filemanager-apt_2.0.0_arm64.deb` - APT package

**Features**:
- ✓ Cross-compilation support
- ✓ Same security enhancements as AMD64
- ✓ ARM64-specific toolchain configuration

#### Windows AMD64
```bash
bash scripts/build-windows-amd64.sh
```

**Creates**:
- `filemanager.exe` - Windows executable
- `filemanager-2.0.0-windows-amd64.zip` - Portable ZIP

**Post-Installation**:
```powershell
powershell -ExecutionPolicy Bypass -File scripts/Set-FileManagerPermissions.ps1
```

---

### 4. DEB Package Builder (Wrapper)
```bash
# Build DEB for specific architecture
bash scripts/build-deb.sh amd64
bash scripts/build-deb.sh arm64
```

**What it does**:
1. Calls appropriate platform-specific build script
2. Verifies DEB was created
3. Shows security features
4. Provides installation instructions

**Output**:
```
✅ DEB package created successfully!

📦 Package: filemanager_2.0.0_amd64.deb
📊 Size: 8.5M

🔒 Security Features:
   ✓ Frontend files in /usr/share/filemanager/frontend (read-only)
   ✓ System config in /etc/filemanager/config.yaml (root-owned)
   ✓ User config in ~/.config/filemanager/config.yaml (user-owned)
   ✓ Proper file permissions (755 dirs, 644 files)
   ✓ Frontend resolver for stable directory location
```

---

## 🛡️ Security Integration

### DEB Package Security Scripts

#### postinst (Post-Installation)
```bash
# Executed after package installation
# Creates frontend directory with restrictive permissions
# Sets up system and user config files
# Updates library cache
```

**Permissions Set**:
```bash
chmod 755 /usr/share/filemanager/frontend           # rwxr-xr-x
chmod 644 /usr/share/filemanager/frontend/*.html    # rw-r--r--
chmod 644 /usr/share/filemanager/frontend/css/*     # rw-r--r--
chmod 644 /usr/share/filemanager/frontend/js/*      # rw-r--r--
chown -R root:root /usr/share/filemanager/frontend
```

#### prerm (Pre-Removal)
```bash
# Executed before package removal
# Backs up user configurations
# Preserves user data
```

#### postrm (Post-Removal)
```bash
# Executed after package removal
# Removes system files
# Cleans up empty directories
# Preserves user configurations
```

### Windows Security Setup
```powershell
# Run as Administrator
powershell -ExecutionPolicy Bypass -File scripts/Set-FileManagerPermissions.ps1

# Sets ACLs:
# - Administrators: Full Control
# - Users: Read & Execute (no write)
```

---

## ✅ Verification & Testing

### Security Verification
```bash
bash scripts/verify-security.sh
```

**Checks**:
1. Binary installation
2. Frontend directory
3. File permissions
4. Configuration files
5. Frontend resolution
6. Duplicate detection
7. World-writable files

### Frontend Info
```bash
filemanager --frontend-info
```

**Shows**:
- Resolution priority
- Active frontend location
- Config file locations
- Path status (EXISTS/DOES NOT EXIST)

### Installation Verification
```bash
# Check package
dpkg -l | grep filemanager

# Check frontend location
ls -la /usr/share/filemanager/frontend

# Check permissions
stat /usr/share/filemanager/frontend
stat /usr/share/filemanager/frontend/index.html

# Check config
cat /etc/filemanager/config.yaml
```

---

## 📦 Package Contents

### DEB Package Structure
```
filemanager_2.0.0_amd64.deb
├── usr/
│   ├── bin/
│   │   └── filemanager          # Executable
│   ├── local/
│   │   └── lib/
│   │       └── libfs_operations_core.so  # Rust library
│   └── share/
│       └── doc/
│           └── filemanager/     # Documentation
└── DEBIAN/
    ├── control                  # Metadata
    ├── postinst                 # Post-install script
    ├── prerm                    # Pre-remove script
    ├── postrm                   # Post-remove script
    └── conffiles                # Config declarations
```

### Installation Result
```
After installation:
/usr/bin/filemanager                          # Binary
/usr/local/lib/libfs_operations_core.so       # Library
/usr/share/filemanager/frontend/              # Frontend (created by postinst)
/etc/filemanager/config.yaml                  # System config (created by postinst)
~/.config/filemanager/config.yaml             # User config (created on first run)
```

---

## 🚀 Deployment Workflow

### Step 1: Build Locally
```bash
cd /path/to/FileManager/file_manager_v1
bash scripts/build.sh
```

### Step 2: Test Locally
```bash
bash scripts/filemanager.sh --frontend-info
bash scripts/filemanager.sh --web
bash scripts/verify-security.sh
```

### Step 3: Build Production Packages
```bash
# Build for all platforms
bash scripts/build-all-platforms.sh all

# Or specific platforms
bash scripts/build-linux-amd64.sh
bash scripts/build-linux-arm64.sh
bash scripts/build-windows-amd64.sh
```

### Step 4: Test Packages
```bash
# Test DEB installation
sudo apt install ./filemanager_2.0.0_amd64.deb
filemanager --frontend-info
bash scripts/verify-security.sh

# Test uninstallation
sudo apt remove filemanager
```

### Step 5: Deploy to Repositories
```bash
# Upload to APT repository
# Upload to Snap store
# Upload to Homebrew
# Upload to Windows installers
# etc.
```

---

## 🔧 Customization

### Change Frontend Location
```bash
# Temporary override
export FILEMANAGER_FRONTEND_DIR=/custom/path
filemanager --web

# Permanent override (edit config)
echo "frontend_dir: /custom/path" >> ~/.config/filemanager/config.yaml
```

### Change Version
```bash
# Update VERSION file
echo "2.1.0" > VERSION

# Rebuild
bash scripts/build-all-platforms.sh all
```

### Add Custom Frontend Files
```bash
# Place files in frontend/ directory
# They will be included in all packages
cp -r /path/to/custom/frontend/* frontend/

# Rebuild
bash scripts/build-all-platforms.sh all
```

---

## 📋 Build Checklist

- [ ] Update VERSION file
- [ ] Update CHANGELOG
- [ ] Test locally with `build.sh`
- [ ] Verify with `filemanager --frontend-info`
- [ ] Run `verify-security.sh`
- [ ] Build all platforms with `build-all-platforms.sh all`
- [ ] Test DEB installation
- [ ] Test Windows installer
- [ ] Verify security permissions
- [ ] Create release notes
- [ ] Sign packages (if applicable)
- [ ] Upload to repositories
- [ ] Verify downloads work
- [ ] Test installation from repositories
- [ ] Document any issues

---

## 🐛 Troubleshooting

### Build Fails
```bash
# Check dependencies
cargo --version
go version
rustc --version

# Clean and rebuild
rm -rf rust_ffi/target file_manager/target
bash scripts/build.sh
```

### DEB Installation Fails
```bash
# Check dependencies
sudo apt-get install -y libc6 libgcc-s1 libstdc++6

# Check for conflicts
dpkg -l | grep filemanager

# Install with verbose output
sudo apt install -v ./filemanager_2.0.0_amd64.deb
```

### Frontend Not Found
```bash
# Check frontend location
filemanager --frontend-info

# Check file permissions
ls -la /usr/share/filemanager/frontend
stat /usr/share/filemanager/frontend

# Check config
cat /etc/filemanager/config.yaml
cat ~/.config/filemanager/config.yaml
```

### Security Verification Fails
```bash
# Run verification script
bash scripts/verify-security.sh

# Check permissions manually
stat /usr/share/filemanager/frontend
find /usr/share/filemanager/frontend -type f -exec stat {} \;

# Fix permissions if needed
sudo chmod 755 /usr/share/filemanager/frontend
sudo find /usr/share/filemanager/frontend -type f -exec chmod 644 {} \;
```

---

## 📚 Related Documentation

- **FRONTEND_RESOLVER_IMPLEMENTATION.md** - Frontend resolver details
- **PRODUCTION_INSTALLERS.md** - Installer scripts for all platforms
- **SECURITY_IMPLEMENTATION_SUMMARY.md** - Security overview

---

**Status**: ✅ Complete
**Version**: 2.0.0
**Last Updated**: December 6, 2025
