# 🌍 Cross-Platform Frontend Strategy

## Overview

The FileManager application needs to handle frontend files across multiple platforms and package types. This document outlines the complete strategy for ensuring the frontend is always available and properly served.

## 🎯 The Golden Rule

### System-Level Installs (requires root/admin)
- ✅ **Bundle frontend files during packaging** (not at runtime)
- ❌ **Don't create at runtime** (read-only for users)
- ✅ **Serve pre-bundled files** (reading is allowed)

### User-Level Installs (no root required)
- ✅ **Create frontend at runtime if missing**
- ✅ **User can modify files** (their own directory)

## 📊 Platform/Package Matrix

| Platform | Frontend Location | Writable? | Strategy | Status |
|----------|------------------|-----------|----------|--------|
| **Snap** | `$SNAP/usr/share/filemanager/frontend` | ❌ No | Bundle during build | ✅ Complete |
| **Flatpak** | `/app/share/filemanager/frontend` | ❌ No | Bundle during build | 🔄 In Progress |
| **APT/DEB** | `/usr/share/filemanager/frontend` | ❌ No (system) | Bundle in package | 🔄 In Progress |
| **RPM/Fedora** | `/usr/share/filemanager/frontend` | ❌ No (system) | Bundle in package | 🔄 In Progress |
| **Homebrew (macOS)** | `/usr/local/opt/filemanager/share/frontend` | ❌ No (system) | Bundle in formula | 🔄 In Progress |
| **Windows MSI** | `C:\Program Files\FileManager\frontend` | ⚠️ Admin only | Bundle in installer | 🔄 In Progress |
| **Manual/Portable** | `~/.local/share/filemanager/frontend` (Linux) | ✅ Yes | Create at runtime | ✅ Complete |
| **Development** | `./frontend` or `~/.local/share/` | ✅ Yes | Create at runtime | ✅ Complete |

## 🔧 Code Implementation

### 1. Frontend Resolver (`frontend_resolver.go`)

The resolver detects the correct frontend location based on:

1. **Environment Variable** - `FILEMANAGER_FRONTEND_DIR`
2. **User Config** - `~/.config/filemanager/config.yaml`
3. **System Config** - `/etc/filemanager/config.yaml`
4. **OS Default** - Platform-specific system locations
5. **Portable Fallback** - Next to executable

**Current Implementation**: ✅ Handles Snap, Flatpak, Linux, macOS, Windows, HarmonyOS

### 2. Web Server (`webserver.go`)

The server determines if it's a system-level install and acts accordingly:

```go
isSystemInstall := isSystemLevelInstall(frontendDir)

if isSystemInstall {
    // System installs: frontend MUST be pre-bundled
    if !isValidDir(frontendDir) {
        return fmt.Errorf("system installation error: frontend not found")
    }
    // Verify critical files exist
    // Serve files (read-only is OK)
} else {
    // User installs: create if missing
    ensureFrontendExists(frontendDir)
}
```

**Current Implementation**: ✅ Detects system vs user installs

### 3. System Level Detection

The `isSystemLevelInstall()` function checks:

- **Snap**: `$SNAP` env var + `/snap/` in path
- **Flatpak**: `$FLATPAK_ID` env var + `/app/` prefix
- **Linux System**: Paths starting with `/usr/share/`, `/usr/local/share/`, `/opt/`
- **macOS System**: Paths starting with `/usr/local/opt/`, `/opt/homebrew/`, `/Library/Application Support/`
- **Windows System**: Paths starting with `%ProgramFiles%`

## 📦 Platform-Specific Implementation Guide

### 1. Snap (Linux - All Distros)

**Status**: ✅ **COMPLETE**

**Location**: `$SNAP/usr/share/filemanager/frontend`

**Implementation** (in `snapcraft.yaml`):
```yaml
parts:
  filemanager:
    override-build: |
      # ... build binary ...
      
      # ⭐ CREATE FRONTEND STRUCTURE IN SNAP
      echo "Creating frontend structure for snap..."
      FRONTEND_DIR="${SNAPCRAFT_PART_INSTALL}/usr/share/filemanager/frontend"
      mkdir -p "${FRONTEND_DIR}"/{css,js,images}
      
      # Copy existing frontend if it exists
      if [ -d scripts/dist/filemanager_frontend ]; then
        echo "Copying existing frontend..."
        cp -r scripts/dist/filemanager_frontend/* "${FRONTEND_DIR}/"
      else
        echo "❌ ERROR: No frontend source found!"
        exit 1
      fi
      
      # Verify frontend was created
      if [ ! -f "${FRONTEND_DIR}/index.html" ]; then
        echo "❌ ERROR: index.html not found after copy!"
        exit 1
      fi
      
      echo "✅ Frontend bundled successfully"

environment:
  FILEMANAGER_FRONTEND_DIR: $SNAP/usr/share/filemanager/frontend
```

### 2. Flatpak (Linux - All Distros)

**Status**: 🔄 **TODO**

**Location**: `/app/share/filemanager/frontend`

**Implementation** (in `com.filemanager.FileManager.yaml`):
```yaml
modules:
  - name: filemanager
    buildsystem: simple
    build-commands:
      # Install binary
      - install -D filemanager /app/bin/filemanager
      
      # Install frontend (bundled in Flatpak)
      - mkdir -p /app/share/filemanager/frontend
      - cp -r frontend/* /app/share/filemanager/frontend/
      
      # Set environment variable
      - mkdir -p /app/etc/filemanager
      - echo 'frontend_dir: /app/share/filemanager/frontend' > /app/etc/filemanager/config.yaml
    
    sources:
      - type: file
        path: filemanager
      - type: dir
        path: frontend
```

### 3. APT/DEB (Linux - Debian/Ubuntu)

**Status**: 🔄 **TODO**

**Location**: `/usr/share/filemanager/frontend`

**Implementation**:

**debian/install file:**
```
filemanager usr/bin/
frontend/* usr/share/filemanager/frontend/
```

**debian/postinst script:**
```bash
#!/bin/bash
set -e

# System frontend directory
FRONTEND_DIR="/usr/share/filemanager/frontend"

# Create directory
mkdir -p "$FRONTEND_DIR"

# Set permissions (read-only for users)
chmod -R 755 "$FRONTEND_DIR"
find "$FRONTEND_DIR" -type f -exec chmod 644 {} \;
chown -R root:root "$FRONTEND_DIR"

echo "✅ Frontend installed to: $FRONTEND_DIR"
```

### 4. RPM/Fedora (Linux - Red Hat/Fedora)

**Status**: 🔄 **TODO**

**Location**: `/usr/share/filemanager/frontend`

**Implementation** (in `.spec` file):
```spec
%files
%{_bindir}/filemanager
%{_datadir}/filemanager/frontend/

%post
# Set permissions
chmod -R 755 %{_datadir}/filemanager/frontend
find %{_datadir}/filemanager/frontend -type f -exec chmod 644 {} \;

%postun
# Cleanup if needed
```

### 5. Homebrew (macOS)

**Status**: 🔄 **TODO**

**Location**: `/usr/local/opt/filemanager/share/frontend` or `/opt/homebrew/opt/filemanager/share/frontend`

**Implementation** (in `filemanager.rb`):
```ruby
class Filemanager < Formula
  desc "File manager with web interface"
  homepage "https://github.com/devonionrouting4Moses/fileManager"
  url "https://github.com/devonionrouting4Moses/fileManager/archive/v2.0.2.tar.gz"
  
  def install
    # Install binary
    bin.install "filemanager"
    
    # Install frontend (bundled)
    (share/"filemanager/frontend").install Dir["frontend/*"]
    
    # Create config
    (etc/"filemanager").mkpath
    (etc/"filemanager/config.yaml").write <<~EOS
      frontend_dir: #{share}/filemanager/frontend
    EOS
  end
  
  test do
    system "#{bin}/filemanager", "--version"
  end
end
```

### 6. Windows Installer (NSIS)

**Status**: 🔄 **TODO**

**Location**: `C:\Program Files\FileManager\frontend`

**Implementation** (in `installer.nsi`):
```nsis
Section "Install"
    SetOutPath "$INSTDIR"
    
    # Install binary
    File "filemanager.exe"
    
    # Install frontend (bundled)
    SetOutPath "$INSTDIR\frontend"
    File /r "frontend\*.*"
    
    # Set read-only for non-admin users
    ExecWait 'powershell -ExecutionPolicy Bypass -File "$INSTDIR\Set-Permissions.ps1"'
    
    # Create config
    FileOpen $0 "$INSTDIR\config.yaml" w
    FileWrite $0 "frontend_dir: $INSTDIR\frontend$\r$\n"
    FileClose $0
SectionEnd
```

**Set-Permissions.ps1:**
```powershell
# Set frontend directory permissions
$frontendPath = "$PSScriptRoot\frontend"
$acl = Get-Acl $frontendPath

# Admins: Full Control
$adminRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
    "Administrators", "FullControl", "ContainerInherit,ObjectInherit", "None", "Allow"
)
$acl.AddAccessRule($adminRule)

# Users: Read & Execute only
$userRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
    "Users", "ReadAndExecute", "ContainerInherit,ObjectInherit", "None", "Allow"
)
$acl.AddAccessRule($userRule)

Set-Acl $frontendPath $acl
```

### 7. Manual/Portable Install

**Status**: ✅ **COMPLETE**

**Location**: `~/.local/share/filemanager/frontend` (Linux), `~/Library/Application Support/FileManager/frontend` (macOS), `%LOCALAPPDATA%\FileManager\frontend` (Windows)

**Implementation**: Already handled by `frontend_resolver.go` and `webserver.go`

## ✅ Checklist for Implementation

### Code Changes (✅ DONE)
- [x] Add `isSystemLevelInstall()` to `webserver.go`
- [x] Update `StartWebServer()` to use system-level detection
- [x] Update `frontend_resolver.go` for all platforms
- [x] Handle Snap bundled frontend
- [x] Handle Flatpak bundled frontend

### Packaging Changes (🔄 TODO)
- [ ] **Snap**: Verify `snapcraft.yaml` is correct
- [ ] **Flatpak**: Create `com.filemanager.FileManager.yaml`
- [ ] **APT/DEB**: Update `debian/install` and `debian/postinst`
- [ ] **RPM**: Create `.spec` file with frontend bundling
- [ ] **Homebrew**: Create `filemanager.rb` formula
- [ ] **Windows**: Create NSIS installer with frontend bundling

### Testing (🔄 TODO)
- [ ] Test Snap installation and frontend loading
- [ ] Test Flatpak installation and frontend loading
- [ ] Test DEB installation and frontend loading
- [ ] Test RPM installation and frontend loading
- [ ] Test Homebrew installation and frontend loading
- [ ] Test Windows MSI installation and frontend loading
- [ ] Test manual/portable installation

## 🚀 Next Steps

1. **Verify Snap Implementation**: Test the current Snap build
2. **Create Flatpak Manifest**: Add Flatpak support
3. **Update DEB Packaging**: Bundle frontend in DEB packages
4. **Create RPM Spec**: Add RPM support
5. **Create Homebrew Formula**: Add macOS support
6. **Create Windows Installer**: Add Windows MSI/NSIS support
7. **Run Integration Tests**: Test all platforms

## 📝 Notes

- The code now automatically detects system vs user installs
- System installs expect pre-bundled frontend (no runtime creation)
- User installs create frontend at runtime if missing
- All platforms follow the same pattern for consistency
- Environment variables can override any auto-detection

## 🔗 Related Files

- `file_manager/internal/handler/webserver.go` - Web server with system detection
- `file_manager/internal/handler/frontend_resolver.go` - Frontend location resolver
- `snap/snapcraft.yaml` - Snap build configuration
- `debian/install` - DEB package file list
- `debian/postinst` - DEB post-installation script
