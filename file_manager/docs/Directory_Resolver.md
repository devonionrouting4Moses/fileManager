# 📦 FileManager Installation Configurations

## 🐧 Linux - APT/DEB Package

### `debian/postinst`
```bash
#!/bin/bash
set -e

# Create system frontend directory
FRONTEND_DIR="/usr/share/filemanager/frontend"
mkdir -p "$FRONTEND_DIR"

# Copy frontend files from package
cp -r /usr/share/filemanager-staging/frontend/* "$FRONTEND_DIR/"
chmod -R 755 "$FRONTEND_DIR"

# Create system config directory
mkdir -p /etc/filemanager
if [ ! -f /etc/filemanager/config.yaml ]; then
    cat > /etc/filemanager/config.yaml <<EOF
# FileManager System Configuration
frontend_dir: /usr/share/filemanager/frontend
EOF
    chmod 644 /etc/filemanager/config.yaml
fi

echo "✅ FileManager frontend installed to: $FRONTEND_DIR"
exit 0
```

### `debian/control`
```
Package: filemanager
Version: 2.0.0
Section: utils
Priority: optional
Architecture: amd64
Maintainer: Your Name <you@example.com>
Description: Cross-platform file management tool
 FileManager provides both CLI and web interface for file operations
```

### `debian/install`
```
filemanager usr/bin/
frontend/* usr/share/filemanager-staging/frontend/
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

parts:
  filemanager:
    plugin: go
    source: .
    build-packages:
      - gcc
    stage-packages:
      - libc6
    override-build: |
      go build -o $SNAPCRAFT_PART_INSTALL/bin/filemanager
      
      # Copy frontend to snap directory
      mkdir -p $SNAPCRAFT_PART_INSTALL/frontend
      cp -r frontend/* $SNAPCRAFT_PART_INSTALL/frontend/
      
      # The app will automatically use $SNAP/frontend via GetFrontendDir()
```

**Note**: Snap apps automatically set `$SNAP` environment variable, which our resolver detects.

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
      - install -D filemanager /app/bin/filemanager
      - mkdir -p /app/share/filemanager/frontend
      - cp -r frontend/* /app/share/filemanager/frontend/
      
      # Create default config
      - mkdir -p /app/etc/filemanager
      - echo 'frontend_dir: /app/share/filemanager/frontend' > /app/etc/filemanager/config.yaml
    
    sources:
      - type: file
        path: filemanager
      - type: dir
        path: frontend
```

**Note**: Flatpak apps set `$FLATPAK_ID`, which our resolver uses to place user data in `~/.var/app/<app-id>/data/frontend`.

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

  def install
    # Build the binary
    system "go", "build", *std_go_args(ldflags: "-s -w")

    # Install frontend files
    (share/"filemanager/frontend").install Dir["frontend/*"]

    # Create config file
    (etc/"filemanager").mkpath
    (etc/"filemanager/config.yaml").write <<~EOS
      # FileManager Configuration
      frontend_dir: #{share}/filemanager/frontend
    EOS
  end

  def post_install
    # Ensure frontend directory permissions
    chmod 0755, share/"filemanager/frontend"
  end

  test do
    system "#{bin}/filemanager", "--version"
  end
end
```

**Installation paths**:
- Binary: `/usr/local/bin/filemanager` (Intel) or `/opt/homebrew/bin/filemanager` (Apple Silicon)
- Frontend: `/usr/local/opt/filemanager/share/frontend` or `/opt/homebrew/opt/filemanager/share/frontend`
- Config: `/usr/local/etc/filemanager/config.yaml` or `/opt/homebrew/etc/filemanager/config.yaml`

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

Section "Install"
    SetOutPath "$INSTDIR"
    
    # Install binary
    File "filemanager.exe"
    
    # Install frontend
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
    
    # Add to PATH (optional)
    EnVar::SetHKCU
    EnVar::AddValue "PATH" "$INSTDIR"
    
    WriteUninstaller "$INSTDIR\Uninstall.exe"
SectionEnd

Section "Uninstall"
    Delete "$INSTDIR\filemanager.exe"
    RMDir /r "$INSTDIR\frontend"
    Delete "$INSTDIR\config.yaml"
    Delete "$INSTDIR\Uninstall.exe"
    RMDir "$INSTDIR"
    
    Delete "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk"
    RMDir "$SMPROGRAMS\${APPNAME}"
SectionEnd
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

## 🎯 Portable Installation (Manual)

### `install.sh` (Linux/macOS)
```bash
#!/bin/bash
set -e

echo "📦 Installing FileManager (User Installation)..."

# Determine OS
OS=$(uname -s)

if [ "$OS" = "Linux" ]; then
    INSTALL_DIR="$HOME/.local/share/filemanager"
    BIN_DIR="$HOME/.local/bin"
    CONFIG_DIR="$HOME/.config/filemanager"
elif [ "$OS" = "Darwin" ]; then
    INSTALL_DIR="$HOME/Library/Application Support/FileManager"
    BIN_DIR="/usr/local/bin"
    CONFIG_DIR="$HOME/Library/Preferences/filemanager"
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

# Create directories
mkdir -p "$INSTALL_DIR/frontend"
mkdir -p "$BIN_DIR"
mkdir -p "$CONFIG_DIR"

# Copy binary
cp filemanager "$BIN_DIR/"
chmod +x "$BIN_DIR/filemanager"

# Copy frontend
cp -r frontend/* "$INSTALL_DIR/frontend/"

# Create config
cat > "$CONFIG_DIR/config.yaml" <<EOF
# FileManager Configuration
frontend_dir: $INSTALL_DIR/frontend
EOF

echo "✅ FileManager installed successfully!"
echo "📍 Binary: $BIN_DIR/filemanager"
echo "📍 Frontend: $INSTALL_DIR/frontend"
echo "📍 Config: $CONFIG_DIR/config.yaml"
echo ""
echo "Run 'filemanager' to start!"
```

### `install.bat` (Windows)
```batch
@echo off
setlocal

echo Installing FileManager (User Installation)...

set INSTALL_DIR=%LOCALAPPDATA%\FileManager
set CONFIG_DIR=%APPDATA%\filemanager

:: Create directories
mkdir "%INSTALL_DIR%\frontend" 2>nul
mkdir "%CONFIG_DIR%" 2>nul

:: Copy binary
copy filemanager.exe "%INSTALL_DIR%\" >nul
echo Binary installed to: %INSTALL_DIR%\filemanager.exe

:: Copy frontend
xcopy /E /I /Y frontend "%INSTALL_DIR%\frontend" >nul
echo Frontend installed to: %INSTALL_DIR%\frontend

:: Create config
(
echo # FileManager Configuration
echo frontend_dir: %INSTALL_DIR%\frontend
) > "%CONFIG_DIR%\config.yaml"
echo Config created: %CONFIG_DIR%\config.yaml

:: Add to PATH (optional)
setx PATH "%PATH%;%INSTALL_DIR%" >nul

echo.
echo Installation complete!
echo Run 'filemanager' from Command Prompt to start.
pause
```

---

## 🔧 Testing Installation

After installing via any method, verify with:

```bash
# Run diagnostic command
filemanager --frontend-info

# Expected output:
🔍 Frontend Directory Resolution Info:
──────────────────────────────────────────────────
1️⃣  ENV FILEMANAGER_FRONTEND_DIR: (not set)
2️⃣  User Config (~/.config/filemanager/config.yaml): /home/user/.local/share/filemanager/frontend [EXISTS]
3️⃣  System Config (/etc/filemanager/config.yaml): (not configured)
4️⃣  OS Default: /usr/share/filemanager/frontend [DOES NOT EXIST]
5️⃣  Portable Fallback: /opt/filemanager/frontend [DOES NOT EXIST]
──────────────────────────────────────────────────
✅ Active Frontend: /home/user/.local/share/filemanager/frontend
```

---

## 🎯 Summary of Installation Paths

| Package Type | Frontend Location | Config Location |
|-------------|------------------|----------------|
| **APT/DEB** | `/usr/share/filemanager/frontend` | `/etc/filemanager/config.yaml` |
| **Snap** | `$SNAP/frontend` (e.g., `/snap/filemanager/current/frontend`) | N/A (uses snap path) |
| **Flatpak** | `/app/share/filemanager/frontend` | `~/.var/app/com.filemanager/config/config.yaml` |
| **Homebrew (Intel)** | `/usr/local/opt/filemanager/share/frontend` | `/usr/local/etc/filemanager/config.yaml` |
| **Homebrew (M1/M2)** | `/opt/homebrew/opt/filemanager/share/frontend` | `/opt/homebrew/etc/filemanager/config.yaml` |
| **Windows (MSI/NSIS)** | `C:\Program Files\FileManager\frontend` | `C:\Program Files\FileManager\config.yaml` |
| **User Install (Linux)** | `~/.local/share/filemanager/frontend` | `~/.config/filemanager/config.yaml` |
| **User Install (macOS)** | `~/Library/Application Support/FileManager/frontend` | `~/Library/Preferences/filemanager/config.yaml` |
| **User Install (Windows)** | `%LOCALAPPDATA%\FileManager\frontend` | `%APPDATA%\filemanager\config.yaml` |
| **Portable** | `<exe-dir>/frontend` | N/A (uses portable path) |

---

## 🚀 Key Benefits

✅ **Single Source of Truth** - Frontend exists in ONE location per installation  
✅ **No Duplication** - Never creates multiple frontend folders  
✅ **Cross-Platform** - Works identically on Linux, macOS, Windows, HarmonyOS  
✅ **Package-Manager Friendly** - Proper integration with all package managers  
✅ **User Overridable** - Environment variable for custom paths  
✅ **Config-Based** - YAML configuration for flexibility  
✅ **Portable Mode** - Falls back to executable directory  
✅ **Production-Ready** - Enterprise-grade patterns used by VS Code, Docker, etc.






the code implementation is as follows:

frontend_resolver.go - Production Frontend Directory Resolver

package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the YAML configuration structure
type Config struct {
	FrontendDir string `yaml:"frontend_dir"`
}

// GetFrontendDir returns the single, stable frontend directory location
// Priority: ENV → User Config → System Config → OS Default → Portable Fallback
func GetFrontendDir() (string, error) {
	// 1️⃣ Environment Variable Override (Highest Priority)
	if envDir := os.Getenv("FILEMANAGER_FRONTEND_DIR"); envDir != "" {
		absPath, err := filepath.Abs(envDir)
		if err == nil && isValidDir(absPath) {
			log.Printf("📍 Using frontend from ENV: %s\n", absPath)
			return absPath, nil
		}
		if err == nil {
			log.Printf("⚠️  ENV path exists but is not valid directory: %s\n", absPath)
		}
	}

	// 2️⃣ User Config File
	if userConfigDir := getUserConfigFrontendDir(); userConfigDir != "" {
		if isValidDir(userConfigDir) {
			log.Printf("📍 Using frontend from user config: %s\n", userConfigDir)
			return userConfigDir, nil
		}
	}

	// 3️⃣ System Config File
	if systemConfigDir := getSystemConfigFrontendDir(); systemConfigDir != "" {
		if isValidDir(systemConfigDir) {
			log.Printf("📍 Using frontend from system config: %s\n", systemConfigDir)
			return systemConfigDir, nil
		}
	}

	// 4️⃣ OS-Specific Default Directory
	defaultDir := getOSDefaultFrontendDir()
	if isValidDir(defaultDir) {
		log.Printf("📍 Using OS default frontend: %s\n", defaultDir)
		return defaultDir, nil
	}

	// 5️⃣ Portable Fallback (directory next to executable)
	portableDir := getPortableFrontendDir()
	if isValidDir(portableDir) {
		log.Printf("📍 Using portable frontend: %s\n", portableDir)
		return portableDir, nil
	}

	// If nothing exists, return the most appropriate default to create
	preferredDir := getPreferredCreationDir()
	log.Printf("📍 Frontend will be created at: %s\n", preferredDir)
	return preferredDir, nil
}

// getPreferredCreationDir returns the best location to create frontend if it doesn't exist
func getPreferredCreationDir() string {
	// For user installs, prefer user directories
	// For system installs, prefer system directories
	
	// Check if running with elevated privileges (system install)
	if isElevated() {
		return getOSDefaultFrontendDir()
	}

	// For regular users, prefer user-writable locations
	switch runtime.GOOS {
	case "linux":
		// Check for Snap/Flatpak first
		if snapDir := os.Getenv("SNAP"); snapDir != "" {
			return filepath.Join(snapDir, "frontend")
		}
		if flatpakID := os.Getenv("FLATPAK_ID"); flatpakID != "" {
			home, _ := os.UserHomeDir()
			return filepath.Join(home, ".var/app", flatpakID, "data/frontend")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".local/share/filemanager/frontend")
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library/Application Support/FileManager/frontend")
	case "windows":
		appData := os.Getenv("LOCALAPPDATA")
		if appData != "" {
			return filepath.Join(appData, "FileManager", "frontend")
		}
		return getPortableFrontendDir()
	default:
		// HarmonyOS or unknown - use portable
		return getPortableFrontendDir()
	}
}

// isElevated checks if the process is running with elevated privileges
func isElevated() bool {
	switch runtime.GOOS {
	case "linux", "darwin":
		return os.Geteuid() == 0
	case "windows":
		// On Windows, check if we can write to Program Files
		testPath := filepath.Join(os.Getenv("ProgramFiles"), ".filemanager_write_test")
		if err := os.WriteFile(testPath, []byte("test"), 0644); err == nil {
			os.Remove(testPath)
			return true
		}
		return false
	default:
		return false
	}
}

// getUserConfigFrontendDir reads frontend_dir from user config file
func getUserConfigFrontendDir() string {
	configPath := getUserConfigPath()
	return parseFrontendDirFromConfig(configPath)
}

// getSystemConfigFrontendDir reads frontend_dir from system config file
func getSystemConfigFrontendDir() string {
	configPath := getSystemConfigPath()
	return parseFrontendDirFromConfig(configPath)
}

// parseFrontendDirFromConfig parses YAML config and returns frontend_dir
func parseFrontendDirFromConfig(configPath string) string {
	if configPath == "" {
		return ""
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Printf("⚠️  Failed to parse config %s: %v\n", configPath, err)
		return ""
	}

	if config.FrontendDir != "" {
		absPath, err := filepath.Abs(config.FrontendDir)
		if err == nil {
			return absPath
		}
	}

	return ""
}

// getUserConfigPath returns the user config file path based on OS
func getUserConfigPath() string {
	switch runtime.GOOS {
	case "linux":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config/filemanager/config.yaml")
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library/Preferences/filemanager/config.yaml")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "filemanager", "config.yaml")
		}
		return ""
	default:
		// HarmonyOS - requires UID
		uid := os.Getuid()
		return fmt.Sprintf("/data/accounts/%d/appdata/filemanager/config.yaml", uid)
	}
}

// getSystemConfigPath returns the system config file path based on OS
func getSystemConfigPath() string {
	switch runtime.GOOS {
	case "linux":
		return "/etc/filemanager/config.yaml"
	case "darwin":
		return "/Library/Application Support/FileManager/config.yaml"
	case "windows":
		programFiles := os.Getenv("ProgramFiles")
		if programFiles != "" {
			return filepath.Join(programFiles, "FileManager", "config.yaml")
		}
		return ""
	default:
		// HarmonyOS
		return "/system/app/FileManager/config.yaml"
	}
}

// getOSDefaultFrontendDir returns the OS-specific default frontend directory
func getOSDefaultFrontendDir() string {
	switch runtime.GOOS {
	case "linux":
		return getLinuxDefaultFrontendDir()
	case "darwin":
		return getMacOSDefaultFrontendDir()
	case "windows":
		return getWindowsDefaultFrontendDir()
	default:
		return getHarmonyOSDefaultFrontendDir()
	}
}

// getLinuxDefaultFrontendDir returns Linux-specific frontend directory
func getLinuxDefaultFrontendDir() string {
	// Check for Snap
	if snapDir := os.Getenv("SNAP"); snapDir != "" {
		return filepath.Join(snapDir, "frontend")
	}

	// Check for Flatpak
	if flatpakID := os.Getenv("FLATPAK_ID"); flatpakID != "" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".var/app", flatpakID, "data/frontend")
	}

	// Check system installation
	systemDir := "/usr/share/filemanager/frontend"
	if isValidDir(systemDir) {
		return systemDir
	}

	// Fallback to user directory
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local/share/filemanager/frontend")
}

// getMacOSDefaultFrontendDir returns macOS-specific frontend directory
func getMacOSDefaultFrontendDir() string {
	// Check Homebrew installation
	brewPaths := []string{
		"/usr/local/opt/filemanager/share/frontend",
		"/opt/homebrew/opt/filemanager/share/frontend", // Apple Silicon
	}
	for _, brewDir := range brewPaths {
		if isValidDir(brewDir) {
			return brewDir
		}
	}

	// Check system installation
	systemDir := "/Library/Application Support/FileManager/frontend"
	if isValidDir(systemDir) {
		return systemDir
	}

	// Fallback to user directory
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library/Application Support/FileManager/frontend")
}

// getWindowsDefaultFrontendDir returns Windows-specific frontend directory
func getWindowsDefaultFrontendDir() string {
	// Check Program Files installation
	programFiles := os.Getenv("ProgramFiles")
	if programFiles != "" {
		systemDir := filepath.Join(programFiles, "FileManager", "frontend")
		if isValidDir(systemDir) {
			return systemDir
		}
	}

	// Fallback to AppData
	appData := os.Getenv("LOCALAPPDATA")
	if appData != "" {
		return filepath.Join(appData, "FileManager", "frontend")
	}

	// Last resort: portable mode
	return getPortableFrontendDir()
}

// getHarmonyOSDefaultFrontendDir returns HarmonyOS-specific frontend directory
func getHarmonyOSDefaultFrontendDir() string {
	// Check system installation
	systemDir := "/system/app/FileManager/frontend"
	if isValidDir(systemDir) {
		return systemDir
	}

	// Fallback to user directory
	uid := os.Getuid()
	return fmt.Sprintf("/data/accounts/%d/applications/filemanager/frontend", uid)
}

// getPortableFrontendDir returns the portable frontend directory (next to executable)
func getPortableFrontendDir() string {
	exePath, err := os.Executable()
	if err == nil {
		// Resolve symlinks
		realPath, err := filepath.EvalSymlinks(exePath)
		if err == nil {
			exePath = realPath
		}
		return filepath.Join(filepath.Dir(exePath), "frontend")
	}
	// Absolute fallback
	return "./frontend"
}

// isValidDir checks if a path exists and is a directory
func isValidDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ensureFrontendExists ensures the frontend directory exists
// If it doesn't exist, creates it with the basic frontend structure
func ensureFrontendExists(frontendDir string) error {
	if isValidDir(frontendDir) {
		// Verify critical files exist
		indexPath := filepath.Join(frontendDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			log.Printf("✅ Frontend directory verified: %s\n", frontendDir)
			return nil
		}
		log.Printf("⚠️  Frontend directory exists but missing files: %s\n", frontendDir)
	}

	log.Printf("📦 Creating frontend structure at: %s\n", frontendDir)
	
	// Create parent directories if needed
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return fmt.Errorf("failed to create frontend directory: %w", err)
	}

	// Create the basic frontend structure
	createBasicFrontend(frontendDir)
	
	log.Printf("✅ Frontend created successfully at: %s\n", frontendDir)
	return nil
}

// CreateDefaultUserConfig creates a default user config file with the provided frontend directory
func CreateDefaultUserConfig(frontendDir string) error {
	configPath := getUserConfigPath()
	if configPath == "" {
		return fmt.Errorf("unable to determine user config path")
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Don't overwrite existing config
	if _, err := os.Stat(configPath); err == nil {
		log.Printf("ℹ️  Config file already exists: %s\n", configPath)
		return nil
	}

	// Create config with frontend_dir
	config := Config{
		FrontendDir: frontendDir,
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add helpful comments
	configContent := fmt.Sprintf(`# FileManager Configuration
# This file controls where the web interface frontend files are located
# You can override this by setting the FILEMANAGER_FRONTEND_DIR environment variable

# Frontend directory path
%s

# Example configurations:
# frontend_dir: /usr/share/filemanager/frontend
# frontend_dir: ~/.local/share/filemanager/frontend
# frontend_dir: /custom/path/to/frontend
`, string(data))

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	log.Printf("✅ Created user config: %s\n", configPath)
	return nil
}

// PrintFrontendInfo prints diagnostic information about frontend location
func PrintFrontendInfo() {
	fmt.Println("🔍 Frontend Directory Resolution Info:")
	fmt.Println(strings.Repeat("─", 50))
	
	// 1. Environment variable
	if envDir := os.Getenv("FILEMANAGER_FRONTEND_DIR"); envDir != "" {
		fmt.Printf("1️⃣  ENV FILEMANAGER_FRONTEND_DIR: %s [%s]\n", envDir, getPathStatus(envDir))
	} else {
		fmt.Println("1️⃣  ENV FILEMANAGER_FRONTEND_DIR: (not set)")
	}

	// 2. User config
	userConfigPath := getUserConfigPath()
	userConfigDir := getUserConfigFrontendDir()
	if userConfigDir != "" {
		fmt.Printf("2️⃣  User Config (%s): %s [%s]\n", userConfigPath, userConfigDir, getPathStatus(userConfigDir))
	} else {
		fmt.Printf("2️⃣  User Config (%s): (not configured)\n", userConfigPath)
	}

	// 3. System config
	systemConfigPath := getSystemConfigPath()
	systemConfigDir := getSystemConfigFrontendDir()
	if systemConfigDir != "" {
		fmt.Printf("3️⃣  System Config (%s): %s [%s]\n", systemConfigPath, systemConfigDir, getPathStatus(systemConfigDir))
	} else {
		fmt.Printf("3️⃣  System Config (%s): (not configured)\n", systemConfigPath)
	}

	// 4. OS default
	defaultDir := getOSDefaultFrontendDir()
	fmt.Printf("4️⃣  OS Default: %s [%s]\n", defaultDir, getPathStatus(defaultDir))

	// 5. Portable fallback
	portableDir := getPortableFrontendDir()
	fmt.Printf("5️⃣  Portable Fallback: %s [%s]\n", portableDir, getPathStatus(portableDir))

	fmt.Println(strings.Repeat("─", 50))

	// Show which one will be used
	frontendDir, _ := GetFrontendDir()
	fmt.Printf("✅ Active Frontend: %s\n", frontendDir)
}

func getPathStatus(path string) string {
	if isValidDir(path) {
		return "EXISTS"
	}
	if _, err := os.Stat(path); err == nil {
		return "EXISTS (not a directory)"
	}
	return "DOES NOT EXIST"
}



Now update your webserver.go to use this resolver:

package handler

import (
	"filemanager/pkg/version"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// StartWebServer starts the HTTP server
func StartWebServer() error {
	// Get the stable, single frontend directory
	frontendDir, err := GetFrontendDir()
	if err != nil {
		return fmt.Errorf("failed to determine frontend directory: %w", err)
	}

	// Ensure frontend exists (creates if missing)
	if err := ensureFrontendExists(frontendDir); err != nil {
		return fmt.Errorf("failed to setup frontend: %w", err)
	}

	// Create user config for future reference (non-fatal if it fails)
	CreateDefaultUserConfig(frontendDir)

	// Serve static files from the stable directory
	fs := http.FileServer(http.Dir(frontendDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Add cache control headers
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Set content type for common files
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		} else if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		} else if strings.HasSuffix(r.URL.Path, ".html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}

		// Serve the file
		fs.ServeHTTP(w, r)
	})

	// API endpoints
	http.HandleFunc("/api/operation", HandleOperation)
	http.HandleFunc("/api/templates", HandleTemplates)
	http.HandleFunc("/api/health", HandleHealth)

	port := "8080"
	url := fmt.Sprintf("http://localhost:%s", port)

	// Try to open browser automatically
	openBrowser(url)

	log.Printf("🌐 Web server starting on %s\n", url)
	log.Printf("📂 Serving frontend from: %s\n", frontendDir)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return
	}

	// Try to open browser, but don't fail if it doesn't work
	go func() {
		if err := exec.Command(cmd, args...).Start(); err != nil {
			// Silently fail - user can open manually
		}
	}()
}

// Keep your existing createBasicFrontend function exactly as is
// (I've omitted it here to save space, but it stays unchanged)





Finally, add a CLI command to check frontend status: this is another 


package main

import (
	"filemanager/internal/handler"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Add CLI flags
	frontendInfo := flag.Bool("frontend-info", false, "Show frontend directory resolution information")
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Handle --frontend-info flag
	if *frontendInfo {
		handler.PrintFrontendInfo()
		os.Exit(0)
	}

	// Handle --version flag
	if *versionFlag {
		fmt.Println("FileManager v2.0.0")
		os.Exit(0)
	}

	// Your existing main menu logic here...
	showMainMenu()
}

func showMainMenu() {
	// Your existing menu code
	fmt.Println("┌────────────────────────────────────────────────┐")
	fmt.Println("│ FILEMANAGER - HYBRID FILE OPERATIONS SYSTEM    │")
	fmt.Println("│ Optimized for Linux, macOS, Windows & HarmonyOS │")
	fmt.Println("└────────────────────────────────────────────────┘")
	// ... rest of your menu
}

// Add this to your go.mod dependencies:
// gopkg.in/yaml.v3 v3.0.1