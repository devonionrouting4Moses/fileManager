# 📋 Packaging Implementation Checklist

## Overview
This checklist tracks the implementation of frontend bundling for each platform/package type.

## 1. Snap (Linux - All Distros)

**Priority**: 🔴 HIGH (Already partially done)

**Files to Modify**:
- `snap/snapcraft.yaml`

**Checklist**:
- [x] Code supports Snap detection (`isSystemLevelInstall()`)
- [x] Code supports bundled frontend at `$SNAP/usr/share/filemanager/frontend`
- [x] Code supports fallback to `$SNAP_USER_COMMON/frontend`
- [ ] Verify `snapcraft.yaml` bundles frontend correctly
- [ ] Test Snap build locally
- [ ] Test Snap installation and startup
- [ ] Verify frontend loads from bundled location
- [ ] Test fallback to writable directory if needed

**Current Status**: 🟡 Partially Complete - Code ready, needs testing

---

## 2. Flatpak (Linux - All Distros)

**Priority**: 🔴 HIGH

**Files to Create/Modify**:
- Create: `flatpak/com.filemanager.FileManager.yaml`
- Modify: `frontend_resolver.go` (already done)
- Modify: `webserver.go` (already done)

**Checklist**:
- [x] Code supports Flatpak detection (`isSystemLevelInstall()`)
- [x] Code supports bundled frontend at `/app/share/filemanager/frontend`
- [x] Code supports fallback to `~/.var/app/FLATPAK_ID/data/frontend`
- [ ] Create Flatpak manifest file
- [ ] Configure frontend bundling in manifest
- [ ] Test Flatpak build locally
- [ ] Test Flatpak installation and startup
- [ ] Verify frontend loads from bundled location

**Current Status**: 🔴 Not Started - Needs manifest file

**Implementation**:
```bash
# Create directory
mkdir -p flatpak

# Create manifest file (see CROSS_PLATFORM_FRONTEND_STRATEGY.md for content)
touch flatpak/com.filemanager.FileManager.yaml
```

---

## 3. APT/DEB (Linux - Debian/Ubuntu)

**Priority**: 🔴 HIGH

**Files to Modify**:
- `debian/install`
- `debian/postinst`
- `debian/control` (if needed)

**Checklist**:
- [x] Code supports system install detection (`isSystemLevelInstall()`)
- [x] Code supports bundled frontend at `/usr/share/filemanager/frontend`
- [ ] Update `debian/install` to include frontend files
- [ ] Update `debian/postinst` to set permissions
- [ ] Test DEB build locally
- [ ] Test DEB installation and startup
- [ ] Verify frontend loads from system location
- [ ] Verify permissions are correct (read-only for users)

**Current Status**: 🔴 Not Started - Needs debian/ updates

**Implementation Steps**:

1. **Update debian/install**:
```
filemanager usr/bin/
frontend/* usr/share/filemanager/frontend/
```

2. **Update debian/postinst**:
```bash
#!/bin/bash
set -e

FRONTEND_DIR="/usr/share/filemanager/frontend"
mkdir -p "$FRONTEND_DIR"
chmod -R 755 "$FRONTEND_DIR"
find "$FRONTEND_DIR" -type f -exec chmod 644 {} \;
chown -R root:root "$FRONTEND_DIR"
```

---

## 4. RPM/Fedora (Linux - Red Hat/Fedora)

**Priority**: 🟡 MEDIUM

**Files to Create**:
- Create: `filemanager.spec`

**Checklist**:
- [x] Code supports system install detection (`isSystemLevelInstall()`)
- [x] Code supports bundled frontend at `/usr/share/filemanager/frontend`
- [ ] Create RPM spec file
- [ ] Configure frontend bundling in spec
- [ ] Test RPM build locally
- [ ] Test RPM installation and startup
- [ ] Verify frontend loads from system location
- [ ] Verify permissions are correct

**Current Status**: 🔴 Not Started - Needs spec file

**Implementation**:
```bash
# Create spec file
touch filemanager.spec
```

---

## 5. Homebrew (macOS)

**Priority**: 🟡 MEDIUM

**Files to Create**:
- Create: `homebrew/filemanager.rb`

**Checklist**:
- [x] Code supports macOS system paths (`isSystemLevelInstall()`)
- [x] Code supports Homebrew paths (`/usr/local/opt/`, `/opt/homebrew/`)
- [ ] Create Homebrew formula file
- [ ] Configure frontend bundling in formula
- [ ] Test Homebrew build locally
- [ ] Test Homebrew installation and startup
- [ ] Verify frontend loads from Homebrew location
- [ ] Test on both Intel and Apple Silicon Macs

**Current Status**: 🔴 Not Started - Needs formula file

**Implementation**:
```bash
# Create directory
mkdir -p homebrew

# Create formula file (see CROSS_PLATFORM_FRONTEND_STRATEGY.md for content)
touch homebrew/filemanager.rb
```

---

## 6. Windows Installer (NSIS)

**Priority**: 🟡 MEDIUM

**Files to Create**:
- Create: `windows/installer.nsi`
- Create: `windows/Set-FileManagerPermissions.ps1`

**Checklist**:
- [x] Code supports Windows system paths (`isSystemLevelInstall()`)
- [x] Code supports `%ProgramFiles%` detection
- [ ] Create NSIS installer script
- [ ] Configure frontend bundling in installer
- [ ] Create PowerShell permission script
- [ ] Test Windows installer build locally
- [ ] Test Windows installation and startup
- [ ] Verify frontend loads from Program Files
- [ ] Verify permissions are correct (admin-only write)

**Current Status**: 🔴 Not Started - Needs installer files

**Implementation**:
```bash
# Create directory
mkdir -p windows

# Create installer files (see CROSS_PLATFORM_FRONTEND_STRATEGY.md for content)
touch windows/installer.nsi
touch windows/Set-FileManagerPermissions.ps1
```

---

## 7. Manual/Portable Install

**Priority**: ✅ COMPLETE

**Status**: ✅ Already implemented

**Verification**:
- [x] Code creates frontend at runtime
- [x] Code uses user-writable directories
- [x] Code handles all platforms (Linux, macOS, Windows)
- [x] Code creates basic frontend structure if missing

---

## Summary Table

| Platform | Status | Priority | Files | Effort |
|----------|--------|----------|-------|--------|
| Snap | 🟡 Partial | 🔴 HIGH | 1 | Low |
| Flatpak | 🔴 Not Started | 🔴 HIGH | 1 | Medium |
| APT/DEB | 🔴 Not Started | 🔴 HIGH | 2 | Low |
| RPM | 🔴 Not Started | 🟡 MEDIUM | 1 | Medium |
| Homebrew | 🔴 Not Started | 🟡 MEDIUM | 1 | Medium |
| Windows | 🔴 Not Started | 🟡 MEDIUM | 2 | High |
| Manual | ✅ Complete | ✅ DONE | 0 | 0 |

---

## Testing Strategy

### Unit Tests
```bash
# Test system-level detection
go test -v ./internal/handler -run TestIsSystemLevelInstall

# Test frontend resolver
go test -v ./internal/handler -run TestGetFrontendDir
```

### Integration Tests
```bash
# Test each platform installation
./test-snap.sh
./test-flatpak.sh
./test-deb.sh
./test-rpm.sh
./test-homebrew.sh
./test-windows.sh
```

### Manual Testing Checklist
For each platform:
1. [ ] Build package
2. [ ] Install package
3. [ ] Start application
4. [ ] Verify frontend loads
5. [ ] Check file permissions
6. [ ] Verify read-only status (system installs)
7. [ ] Test config file creation (user installs)

---

## Dependencies

### Required Tools
- [ ] Snapcraft (for Snap builds)
- [ ] Flatpak builder (for Flatpak builds)
- [ ] dpkg-deb (for DEB builds)
- [ ] rpmbuild (for RPM builds)
- [ ] Homebrew (for macOS testing)
- [ ] NSIS (for Windows builds)

### Build Scripts Needed
- [ ] `scripts/build-snap.sh`
- [ ] `scripts/build-flatpak.sh`
- [ ] `scripts/build-deb.sh` (update existing)
- [ ] `scripts/build-rpm.sh`
- [ ] `scripts/build-homebrew.sh`
- [ ] `scripts/build-windows.sh`

---

## Notes

- All code changes are complete and tested
- Packaging implementations follow the same pattern
- Each platform bundles frontend during build, not at runtime
- System installs are read-only; user installs are writable
- Environment variable `FILEMANAGER_FRONTEND_DIR` can override any detection

---

## Next Steps

1. **Immediate** (This session):
   - [ ] Verify Snap implementation works
   - [ ] Create Flatpak manifest
   - [ ] Update DEB packaging

2. **Short-term** (Next session):
   - [ ] Create RPM spec file
   - [ ] Create Homebrew formula
   - [ ] Create Windows installer

3. **Long-term**:
   - [ ] Test all platforms
   - [ ] Create CI/CD pipeline for builds
   - [ ] Document installation instructions
