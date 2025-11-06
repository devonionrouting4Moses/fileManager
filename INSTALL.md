# 📥 FileManager Installation Guide

## Quick Install

### Linux

```bash
# Download latest release
curl -L https://github.com/YOUR_USERNAME/filemanager/releases/download/v0.1.0/filemanager-0.1.0-linux-amd64.tar.gz | tar xz

# Install
cd linux-amd64
sudo ./install.sh

# Verify
filemanager --version
```

### macOS

```bash
# Download latest release
curl -L https://github.com/YOUR_USERNAME/filemanager/releases/download/v0.1.0/filemanager-0.1.0-darwin-amd64.tar.gz | tar xz

# Install
cd darwin-amd64
sudo ./install.sh

# Verify
filemanager --version
```

### Windows

1. Download `filemanager-0.1.0-windows-amd64.zip`
2. Extract to a folder
3. Right-click `install.bat` → "Run as Administrator"
4. Verify: `filemanager --version`

---

## Detailed Installation

### 🐧 Linux

#### Method 1: Automated Installer (Recommended)

**Ubuntu/Debian:**
```bash
wget https://github.com/YOUR_USERNAME/filemanager/releases/download/v0.1.0/filemanager-0.1.0-linux-amd64.tar.gz
tar -xzf filemanager-0.1.0-linux-amd64.tar.gz
cd linux-amd64
sudo ./install.sh
```

**Fedora/RHEL/CentOS:**
```bash
curl -L -O https://github.com/YOUR_USERNAME/filemanager/releases/download/v0.1.0/filemanager-0.1.0-linux-amd64.tar.gz
tar -xzf filemanager-0.1.0-linux-amd64.tar.gz
cd linux-amd64
sudo ./install.sh
```

**Arch Linux:**
```bash
wget https://github.com/YOUR_USERNAME/filemanager/releases/download/v0.1.0/filemanager-0.1.0-linux-amd64.tar.gz
tar -xzf filemanager-0.1.0-linux-amd64.tar.gz
cd linux-amd64
sudo ./install.sh
```

#### Method 2: Manual Installation

```bash
# Extract files
tar -xzf filemanager-0.1.0-linux-amd64.tar.gz
cd linux-amd64

# Copy binary
sudo cp filemanager /usr/local/bin/
sudo chmod +x /usr/local/bin/filemanager

# Copy library
sudo cp libfilemanager.so /usr/local/lib/

# Update library cache
sudo ldconfig

# Verify
filemanager --version
```

#### Method 3: User-only Installation (No sudo)

```bash
# Create local directories
mkdir -p ~/.local/bin ~/.local/lib

# Extract and copy
tar -xzf filemanager-0.1.0-linux-amd64.tar.gz
cd linux-amd64
cp filemanager ~/.local/bin/
cp libfilemanager.so ~/.local/lib/

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"
export LD_LIBRARY_PATH="$HOME/.local/lib:$LD_LIBRARY_PATH"

# Reload shell
source ~/.bashrc
```

#### ARM64 (Raspberry Pi, ARM servers)

```bash
wget https://github.com/YOUR_USERNAME/filemanager/releases/download/v0.1.0/filemanager-0.1.0-linux-arm64.tar.gz
tar -xzf filemanager-0.1.0-linux-arm64.tar.gz
cd linux-arm64
sudo ./install.sh
```

---

### 🍎 macOS

#### Method 1: Automated Installer (Recommended)

**Intel Macs:**
```bash
curl -L https://github.com/YOUR_USERNAME/filemanager/releases/download/v0.1.0/filemanager-0.1.0-darwin-amd64.tar.gz | tar xz
cd darwin-amd64
sudo ./install.sh
```

**Apple Silicon (M1/M2):**
```bash
curl -L https://github.com/YOUR_USERNAME/filemanager/releases/download/v0.1.0/filemanager-0.1.0-darwin-arm64.tar.gz | tar xz
cd darwin-arm64
sudo ./install.sh
```

#### Method 2: Homebrew (Coming Soon)

```bash
brew tap YOUR_USERNAME/filemanager
brew install filemanager
```

#### Method 3: Manual Installation

```bash
# Extract files
tar -xzf filemanager-0.1.0-darwin-amd64.tar.gz
cd darwin-amd64

# Copy binary
sudo cp filemanager /usr/local/bin/
sudo chmod +x /usr/local/bin/filemanager

# Copy library
sudo cp libfilemanager.dylib /usr/local/lib/

# Verify
filemanager --version
```

#### Troubleshooting: "Developer Cannot Be Verified"

If you see a security warning:

```bash
# Remove quarantine attribute
sudo xattr -r -d com.apple.quarantine /usr/local/bin/filemanager
sudo xattr -r -d com.apple.quarantine /usr/local/lib/libfilemanager.dylib
```

Or allow in System Preferences:
1. Go to System Preferences → Security & Privacy
2. Click "Allow Anyway" for filemanager

---

### 🪟 Windows

#### Method 1: Automated Installer (Recommended)

1. Download `filemanager-0.1.0-windows-amd64.zip`
2. Extract the ZIP file
3. Right-click `install.bat`
4. Select "Run as Administrator"
5. Follow the prompts

The installer will:
- Copy files to `C:\Program Files\FileManager\`
- Add to System PATH
- Make `filemanager` available globally

#### Method 2: Chocolatey (Coming Soon)

```cmd
choco install filemanager
```

#### Method 3: Winget (Coming Soon)

```cmd
winget install filemanager
```

#### Method 4: Manual Installation

1. Extract ZIP to `C:\Program Files\FileManager\`
2. Add to PATH:
   - Open System Properties (Win+Pause)
   - Click "Advanced system settings"
   - Click "Environment Variables"
   - Under "System variables", find "Path"
   - Click "Edit" → "New"
   - Add: `C:\Program Files\FileManager`
   - Click OK
3. Restart terminal
4. Verify: `filemanager --version`

#### Method 5: Portable Installation

No installation required:
1. Extract ZIP to any folder
2. Run `filemanager.exe` directly from that folder

Note: Make sure `filemanager.dll` is in the same directory.

---

## Verification

After installation, verify everything works:

```bash
# Check version
filemanager --version

# Check for updates
filemanager --update

# Show help
filemanager --help

# Start interactive mode
filemanager
```

Expected output:
```
FileManager v0.1.0
Platform: linux/amd64
Go version: go1.24.9
```

---

## Updating

### Linux/macOS

```bash
# Check for updates
filemanager --update

# If update available, download and reinstall
curl -L [NEW_VERSION_URL] | tar xz
cd [extracted-folder]
sudo ./install.sh
```

### Windows

```bash
# Check for updates
filemanager --update

# If update available, download new ZIP
# Run install.bat as Administrator
```

### Auto-update (Coming Soon)

```bash
filemanager --update --auto
```

---

## Uninstallation

### Linux/macOS

**Using uninstall script:**
```bash
# Navigate to installation folder
cd [installation-folder]
sudo ./uninstall.sh
```

**Manual removal:**
```bash
sudo rm /usr/local/bin/filemanager
sudo rm /usr/local/lib/libfilemanager.so    # Linux
sudo rm /usr/local/lib/libfilemanager.dylib  # macOS
sudo ldconfig  # Linux only
```

### Windows

**Using uninstall script:**
1. Navigate to installation folder
2. Right-click `uninstall.bat`
3. Select "Run as Administrator"

**Manual removal:**
1. Delete `C:\Program Files\FileManager\`
2. Remove from PATH (Environment Variables)

---

## Platform-Specific Notes

### Linux

**Dependencies:**
- glibc 2.27+ (Ubuntu 18.04+, Debian 10+, Fedora 28+)
- No other runtime dependencies

**Supported Distributions:**
- Ubuntu 18.04, 20.04, 22.04, 24.04
- Debian 10, 11, 12
- Fedora 35+
- RHEL/CentOS 8+
- Arch Linux (current)
- openSUSE Leap 15+

### macOS

**Requirements:**
- macOS 10.15 (Catalina) or later
- Intel or Apple Silicon (M1/M2)

**Known Issues:**
- First run may show "Developer cannot be verified" warning
- Solution: See troubleshooting above

### Windows

**Requirements:**
- Windows 10 version 1909 or later
- Windows 11 (all editions)
- x64 processor

**Known Issues:**
- Windows Defender may flag the application (false positive)
- Solution: Add to exclusions or use Windows Security

---

## Building from Source

If you prefer to build from source:

### Prerequisites

```bash
# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Install Go
# Download from https://golang.org/dl/

# Linux: Install build tools
sudo apt install build-essential  # Debian/Ubuntu
sudo dnf install gcc              # Fedora

# macOS: Install Xcode Command Line Tools
xcode-select --install
```

### Build Steps

```bash
# Clone repository
git clone https://github.com/YOUR_USERNAME/filemanager.git
cd filemanager

# Build
make

# Install
sudo make install

# Or build for all platforms
make build-all
```

---

## Troubleshooting

### "Command not found" after installation

**Linux/macOS:**
```bash
# Check if in PATH
which filemanager

# If not found, add to PATH
export PATH="/usr/local/bin:$PATH"

# Make permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**Windows:**
- Restart terminal
- Verify PATH includes installation directory

### "Library not found" error

**Linux:**
```bash
# Update library cache
sudo ldconfig

# Or set LD_LIBRARY_PATH
export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
```

**macOS:**
```bash
# Set DYLD_LIBRARY_PATH
export DYLD_LIBRARY_PATH=/usr/local/lib:$DYLD_LIBRARY_PATH
```

**Windows:**
- Ensure `.dll` is in same directory as `.exe`

### Permission Denied

```bash
# Make executable
chmod +x filemanager

# Or reinstall with proper permissions
sudo ./install.sh
```

---

## Support

- **Documentation:** [GitHub Wiki](https://github.com/YOUR_USERNAME/filemanager/wiki)
- **Issues:** [GitHub Issues](https://github.com/YOUR_USERNAME/filemanager/issues)
- **Discussions:** [GitHub Discussions](https://github.com/YOUR_USERNAME/filemanager/discussions)

---

**Happy File Managing! 🚀**