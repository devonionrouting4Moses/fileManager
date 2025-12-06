# 🛡️ Production Installers & Security Best Practices

## Overview
This document provides complete installer scripts and security configurations for deploying FileManager across all major package managers and platforms.

---

## 🐧 Linux - APT/DEB Package

### Directory Structure
```
debian/
├── control              # Package metadata
├── postinst             # Post-installation script
├── prerm                # Pre-removal script
├── postrm               # Post-removal script
├── conffiles            # Config file declarations
└── install              # File installation manifest
```

### `debian/control`
```
Package: filemanager
Version: 2.0.0
Section: utils
Priority: optional
Architecture: amd64
Maintainer: Your Name <you@example.com>
Depends: libc6 (>= 2.31)
Homepage: https://github.com/yourusername/filemanager
Description: Cross-platform file management tool
 FileManager provides both CLI and web interface for file operations.
 Features include batch operations, project templates, and more.
```

### `debian/postinst` (Post-Installation)
```bash
#!/bin/bash
set -e

# Create system frontend directory with restrictive permissions
FRONTEND_DIR="/usr/share/filemanager/frontend"
mkdir -p "$FRONTEND_DIR"

# Copy frontend files from package
cp -r /usr/share/filemanager-staging/frontend/* "$FRONTEND_DIR/" 2>/dev/null || true

# 🛡️ SECURITY: Set restrictive permissions
echo "🔒 Setting security permissions..."

# Directory permissions: rwxr-xr-x (755)
chmod 755 "$FRONTEND_DIR"
find "$FRONTEND_DIR" -type d -exec chmod 755 {} \;

# File permissions: rw-r--r-- (644)
find "$FRONTEND_DIR" -type f -exec chmod 644 {} \;

# Ownership: root:root
chown -R root:root "$FRONTEND_DIR"

# Create system config directory
mkdir -p /etc/filemanager
chown root:root /etc/filemanager
chmod 755 /etc/filemanager

# Create default system config (if not exists)
if [ ! -f /etc/filemanager/config.yaml ]; then
    cat > /etc/filemanager/config.yaml <<EOF
# FileManager System Configuration
# This file is read-only and managed by the package manager
frontend_dir: /usr/share/filemanager/frontend
EOF
    chmod 644 /etc/filemanager/config.yaml
    chown root:root /etc/filemanager/config.yaml
fi

# Create user config directory template (users will create their own)
mkdir -p /etc/skel/.config/filemanager
cat > /etc/skel/.config/filemanager/config.yaml <<EOF
# FileManager User Configuration
# Copy this file to ~/.config/filemanager/config.yaml to customize

# Frontend directory path (leave empty to use system default)
# frontend_dir: /usr/share/filemanager/frontend

# Or override with environment variable:
# export FILEMANAGER_FRONTEND_DIR=/custom/path
EOF
chmod 644 /etc/skel/.config/filemanager/config.yaml

echo "✅ FileManager frontend installed to: $FRONTEND_DIR"
echo "✅ System config: /etc/filemanager/config.yaml"
echo "✅ User config template: ~/.config/filemanager/config.yaml"

exit 0
```

### `debian/prerm` (Pre-Removal)
```bash
#!/bin/bash
set -e

# Backup user configs before removal
if [ -d ~/.config/filemanager ]; then
    echo "📦 Backing up user configuration..."
    tar czf ~/.config/filemanager-backup-$(date +%Y%m%d-%H%M%S).tar.gz ~/.config/filemanager
fi

exit 0
```

### `debian/postrm` (Post-Removal)
```bash
#!/bin/bash
set -e

# Remove system files
rm -rf /usr/share/filemanager/frontend 2>/dev/null || true
rm -f /etc/filemanager/config.yaml 2>/dev/null || true

# Keep user configs (users may want to preserve them)
# User configs in ~/.config/filemanager are left intact

echo "✅ FileManager uninstalled"
echo "ℹ️  User configurations preserved in ~/.config/filemanager"

exit 0
```

### `debian/conffiles` (Config File Declaration)
```
/etc/filemanager/config.yaml
```

### `debian/install` (File Manifest)
```
filemanager usr/bin/
frontend/* usr/share/filemanager-staging/frontend/
```

### Build DEB Package
```bash
#!/bin/bash
set -e

VERSION="2.0.0"
ARCH="amd64"

# Create debian package structure
mkdir -p debian/DEBIAN
mkdir -p debian/usr/bin
mkdir -p debian/usr/share/filemanager-staging/frontend

# Copy binary
cp filemanager.bin debian/usr/bin/filemanager
chmod 755 debian/usr/bin/filemanager

# Copy frontend
cp -r frontend/* debian/usr/share/filemanager-staging/frontend/

# Copy control files
cp debian/control debian/DEBIAN/
cp debian/postinst debian/DEBIAN/
cp debian/prerm debian/DEBIAN/
cp debian/postrm debian/DEBIAN/
cp debian/conffiles debian/DEBIAN/

# Make scripts executable
chmod 755 debian/DEBIAN/postinst
chmod 755 debian/DEBIAN/prerm
chmod 755 debian/DEBIAN/postrm

# Build package
dpkg-deb --build debian filemanager_${VERSION}_${ARCH}.deb

echo "✅ DEB package created: filemanager_${VERSION}_${ARCH}.deb"
```

---

## 📦 Snap Package

### `snapcraft.yaml`
```yaml
name: filemanager
version: '2.0.0'
summary: Cross-platform file management tool
description: |
  FileManager provides both CLI and web interface for file operations.
  Features include batch operations, project templates, and more.

base: core22
confinement: strict
grade: stable

apps:
  filemanager:
    command: bin/filemanager
    plugs:
      - home
      - network
      - network-bind
    environment:
      # Snap automatically sets $SNAP environment variable
      # Frontend resolver will use: $SNAP/frontend
      SNAP_FRONTEND: $SNAP/frontend

parts:
  filemanager:
    plugin: go
    source: .
    build-packages:
      - gcc
      - cargo
      - rustc
    stage-packages:
      - libc6
    
    override-build: |
      # Build Rust library
      cd rust_ffi
      cargo build --release
      cd ..
      
      # Build Go binary
      go build -o $SNAPCRAFT_PART_INSTALL/bin/filemanager ./cmd/app
      
      # Copy frontend to snap (read-only by default)
      mkdir -p $SNAPCRAFT_PART_INSTALL/frontend
      cp -r frontend/* $SNAPCRAFT_PART_INSTALL/frontend/
      
      # Set restrictive permissions
      chmod 755 $SNAPCRAFT_PART_INSTALL/frontend
      find $SNAPCRAFT_PART_INSTALL/frontend -type f -exec chmod 644 {} \;
      find $SNAPCRAFT_PART_INSTALL/frontend -type d -exec chmod 755 {} \;
    
    # Snap automatically handles permissions
    # Files are read-only by default in strict confinement
```

### Security Features
- **Strict Confinement**: App cannot access system files outside snap
- **Read-Only Frontend**: Frontend bundled in snap is immutable
- **Isolated Storage**: User data in `$SNAP_USER_DATA`
- **Automatic Updates**: Snap handles updates securely

---

## 📦 Flatpak Package

### `com.filemanager.FileManager.yaml`
```yaml
app-id: com.filemanager.FileManager
runtime: org.freedesktop.Platform
runtime-version: '23.08'
sdk: org.freedesktop.Sdk
command: filemanager

finish-args:
  - --share=network
  - --socket=wayland
  - --socket=fallback-x11
  - --filesystem=home
  - --filesystem=~/.var/app/com.filemanager.FileManager/data

modules:
  - name: filemanager
    buildsystem: simple
    build-commands:
      # Build Rust library
      - cd rust_ffi && cargo build --release && cd ..
      
      # Build Go binary
      - go build -o /app/bin/filemanager ./cmd/app
      
      # Install frontend (read-only in app bundle)
      - mkdir -p /app/share/filemanager/frontend
      - cp -r frontend/* /app/share/filemanager/frontend/
      
      # 🛡️ SECURITY: Set restrictive permissions
      - chmod 755 /app/share/filemanager/frontend
      - find /app/share/filemanager/frontend -type f -exec chmod 644 {} \;
      - find /app/share/filemanager/frontend -type d -exec chmod 755 {} \;
      
      # Create default config
      - mkdir -p /app/etc/filemanager
      - |
        cat > /app/etc/filemanager/config.yaml <<EOF
        # FileManager Flatpak Configuration
        frontend_dir: /app/share/filemanager/frontend
        EOF
      - chmod 644 /app/etc/filemanager/config.yaml
    
    sources:
      - type: file
        path: filemanager
      - type: dir
        path: frontend
      - type: dir
        path: rust_ffi
      - type: dir
        path: cmd

cleanup:
  - /include
  - /lib/pkgconfig
  - /share/doc
  - /share/man
```

### Security Features
- **Sandboxed**: App runs in isolated environment
- **Read-Only Frontend**: Frontend in `/app/share` is immutable
- **User Data Isolation**: Data in `~/.var/app/com.filemanager.FileManager/`
- **Permission Control**: Only requested permissions granted

---

## 🍺 Homebrew (macOS/Linux)

### `filemanager.rb`
```ruby
class Filemanager < Formula
  desc "Cross-platform file management tool with CLI and web interface"
  homepage "https://github.com/yourusername/filemanager"
  url "https://github.com/yourusername/filemanager/archive/v2.0.0.tar.gz"
  sha256 "YOUR_SHA256_HERE"
  license "MIT"

  depends_on "go" => :build
  depends_on "rust" => :build

  def install
    # Build Rust library
    system "cd", "rust_ffi", "cargo", "build", "--release"
    
    # Build Go binary
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/app"

    # Install frontend files
    (share/"filemanager/frontend").install Dir["frontend/*"]

    # 🛡️ SECURITY: Set permissions for system installation
    # On Homebrew, files are owned by the user who installed
    # But we set restrictive permissions for safety
    system "chmod", "-R", "755", share/"filemanager/frontend"
    system "find", share/"filemanager/frontend", "-type", "f", "-exec", "chmod", "644", "{}", "+"
    system "find", share/"filemanager/frontend", "-type", "d", "-exec", "chmod", "755", "{}", "+"

    # Create config file
    (etc/"filemanager").mkpath
    (etc/"filemanager/config.yaml").write <<~EOS
      # FileManager Configuration (Homebrew Installation)
      frontend_dir: #{share}/filemanager/frontend
    EOS
    
    # Set config permissions
    system "chmod", "644", etc/"filemanager/config.yaml"
  end

  def post_install
    # Ensure frontend directory permissions
    system "chmod", "755", share/"filemanager/frontend"
  end

  test do
    system "#{bin}/filemanager", "--version"
  end
end
```

### Installation Paths
- **Binary**: `/usr/local/bin/filemanager` (Intel) or `/opt/homebrew/bin/filemanager` (Apple Silicon)
- **Frontend**: `/usr/local/opt/filemanager/share/frontend` or `/opt/homebrew/opt/filemanager/share/frontend`
- **Config**: `/usr/local/etc/filemanager/config.yaml` or `/opt/homebrew/etc/filemanager/config.yaml`

---

## 🪟 Windows Installer

### NSIS Script: `installer.nsi`
```nsis
!define APPNAME "FileManager"
!define APPVERSION "2.0.0"
!define INSTALLDIR "$PROGRAMFILES\FileManager"

Name "${APPNAME}"
OutFile "FileManager-Setup-${APPVERSION}.exe"
InstallDir "${INSTALLDIR}"
RequestExecutionLevel admin

Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "Install"
    SetOutPath "$INSTDIR"
    
    # Install binary
    File "filemanager.exe"
    
    # Install frontend (read-only for security)
    SetOutPath "$INSTDIR\frontend"
    File /r "frontend\*.*"
    
    # Create config directory
    CreateDirectory "$INSTDIR"
    
    # Create config file
    FileOpen $0 "$INSTDIR\config.yaml" w
    FileWrite $0 "# FileManager Configuration$\r$\n"
    FileWrite $0 "frontend_dir: $INSTDIR\frontend$\r$\n"
    FileClose $0
    
    # Create Start Menu shortcut
    CreateDirectory "$SMPROGRAMS\${APPNAME}"
    CreateShortcut "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk" "$INSTDIR\filemanager.exe"
    CreateShortcut "$SMPROGRAMS\${APPNAME}\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
    
    # Add to PATH (optional)
    EnVar::SetHKCU
    EnVar::AddValue "PATH" "$INSTDIR"
    
    # 🛡️ SECURITY: Set file permissions (Windows ACLs)
    # Note: NSIS has limited ACL support, use PowerShell for advanced permissions
    
    WriteUninstaller "$INSTDIR\Uninstall.exe"
SectionEnd

Section "Uninstall"
    Delete "$INSTDIR\filemanager.exe"
    RMDir /r "$INSTDIR\frontend"
    Delete "$INSTDIR\config.yaml"
    Delete "$INSTDIR\Uninstall.exe"
    RMDir "$INSTDIR"
    
    Delete "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk"
    Delete "$SMPROGRAMS\${APPNAME}\Uninstall.lnk"
    RMDir "$SMPROGRAMS\${APPNAME}"
    
    # Remove from PATH
    EnVar::SetHKCU
    EnVar::RemoveValue "PATH" "$INSTDIR"
SectionEnd
```

### PowerShell Security Script: `set-permissions.ps1`
```powershell
# Run as Administrator
param(
    [string]$InstallPath = "$env:ProgramFiles\FileManager"
)

# Verify admin rights
if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Error "This script must be run as Administrator"
    exit 1
}

Write-Host "🔒 Setting security permissions for FileManager..." -ForegroundColor Green

# Get the frontend directory
$FrontendDir = Join-Path $InstallPath "frontend"

if (-not (Test-Path $FrontendDir)) {
    Write-Error "Frontend directory not found: $FrontendDir"
    exit 1
}

# Remove inheritance and set explicit permissions
$Acl = Get-Acl $FrontendDir
$Acl.SetAccessRuleProtection($true, $false)
Set-Acl $FrontendDir $Acl

# Grant Administrators full control
$AdminRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
    "BUILTIN\Administrators",
    "FullControl",
    "ContainerInherit,ObjectInherit",
    "None",
    "Allow"
)
$Acl.AddAccessRule($AdminRule)

# Grant Users read and execute only
$UserRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
    "BUILTIN\Users",
    "ReadAndExecute",
    "ContainerInherit,ObjectInherit",
    "None",
    "Allow"
)
$Acl.AddAccessRule($UserRule)

# Apply permissions
Set-Acl $FrontendDir $Acl

Write-Host "✅ Permissions set successfully" -ForegroundColor Green
Write-Host "   Administrators: Full Control" -ForegroundColor Gray
Write-Host "   Users: Read & Execute (no write)" -ForegroundColor Gray
```

### WiX Toolset (MSI): `Product.wxs`
```xml
<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
  <Product Id="*" Name="FileManager" Language="1033" Version="2.0.0" 
           Manufacturer="YourCompany" UpgradeCode="YOUR-GUID-HERE">
    
    <Package InstallerVersion="200" Compressed="yes" InstallScope="perMachine" />
    
    <Media Id="1" Cabinet="filemanager.cab" EmbedCab="yes" />
    
    <Directory Id="TARGETDIR" Name="SourceDir">
      <Directory Id="ProgramFilesFolder">
        <Directory Id="INSTALLFOLDER" Name="FileManager">
          
          <!-- Binary -->
          <Component Id="MainExecutable" Guid="YOUR-GUID-1">
            <File Id="FileManagerEXE" Source="filemanager.exe" KeyPath="yes" />
          </Component>
          
          <!-- Frontend Directory -->
          <Directory Id="FrontendDir" Name="frontend">
            <Component Id="FrontendFiles" Guid="YOUR-GUID-2">
              <File Source="frontend\index.html" />
              <File Source="frontend\css\style.css" />
              <File Source="frontend\js\main.js" />
            </Component>
          </Directory>
          
          <!-- Config File -->
          <Component Id="ConfigFile" Guid="YOUR-GUID-3">
            <File Source="config.yaml" />
          </Component>
          
        </Directory>
      </Directory>
      
      <!-- Start Menu -->
      <Directory Id="ProgramMenuFolder">
        <Directory Id="ApplicationProgramsFolder" Name="FileManager">
          <Component Id="ApplicationShortcut" Guid="YOUR-GUID-4">
            <Shortcut Id="ApplicationStartMenuShortcut"
                     Name="FileManager"
                     Description="File Management Tool"
                     Target="[INSTALLFOLDER]filemanager.exe"
                     WorkingDirectory="INSTALLFOLDER"/>
            <RemoveFolder Id="CleanUpShortCut" On="uninstall"/>
            <RegistryValue Root="HKCU" Key="Software\FileManager" 
                          Name="installed" Type="integer" Value="1" KeyPath="yes"/>
          </Component>
        </Directory>
      </Directory>
    </Directory>
    
    <Feature Id="ProductFeature" Title="FileManager" Level="1">
      <ComponentRef Id="MainExecutable" />
      <ComponentRef Id="FrontendFiles" />
      <ComponentRef Id="ConfigFile" />
      <ComponentRef Id="ApplicationShortcut" />
    </Feature>
  </Product>
</Wix>
```

---

## 🔐 Security Checklist

### File Permissions
- [ ] Directories: `755` (rwxr-xr-x)
- [ ] Files: `644` (rw-r--r--)
- [ ] Ownership: `root:root` (system) or user (user install)
- [ ] No world-writable files

### Configuration Files
- [ ] System config: `/etc/filemanager/config.yaml` (644, root:root)
- [ ] User config: `~/.config/filemanager/config.yaml` (644, user:user)
- [ ] Config not world-readable if contains sensitive data

### Frontend Security
- [ ] Frontend files read-only for regular users
- [ ] No executable permissions on HTML/CSS/JS files
- [ ] No symlinks to sensitive locations
- [ ] Validate all file paths to prevent directory traversal

### Installation Security
- [ ] Verify checksums before installation
- [ ] Run installers with minimal required privileges
- [ ] Backup user configs before removal
- [ ] Log all installation actions
- [ ] Test uninstallation thoroughly

### Runtime Security
- [ ] Drop privileges after startup if running as root
- [ ] Validate user input in web interface
- [ ] Use HTTPS for web interface (if exposed)
- [ ] Implement CSRF protection
- [ ] Rate limit API endpoints

---

## 📋 Installation Verification

### Verify APT Installation
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

### Verify Snap Installation
```bash
# Check snap
snap list | grep filemanager

# Check frontend location
ls -la /snap/filemanager/current/frontend

# Check permissions
stat /snap/filemanager/current/frontend
```

### Verify Homebrew Installation
```bash
# Check formula
brew list filemanager

# Check frontend location
ls -la $(brew --prefix)/opt/filemanager/share/frontend

# Check permissions
stat $(brew --prefix)/opt/filemanager/share/frontend
```

### Verify Windows Installation
```powershell
# Check registry
Get-ItemProperty "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\FileManager"

# Check frontend location
Get-ChildItem "C:\Program Files\FileManager\frontend"

# Check permissions
Get-Acl "C:\Program Files\FileManager\frontend"
```

---

## 🚀 Deployment Checklist

- [ ] Build and test on each target platform
- [ ] Verify frontend location with `--frontend-info`
- [ ] Test installation from scratch
- [ ] Test upgrade from previous version
- [ ] Test uninstallation and cleanup
- [ ] Verify permissions are correct
- [ ] Test web interface functionality
- [ ] Test CLI operations
- [ ] Verify no duplicate frontend folders created
- [ ] Check logs for errors
- [ ] Document any platform-specific issues
- [ ] Create release notes
- [ ] Sign packages (if applicable)
- [ ] Upload to package repositories
