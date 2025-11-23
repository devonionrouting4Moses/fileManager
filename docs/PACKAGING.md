# 📦 FileManager Packaging & Distribution Guide

## Current Status: v0.1.0

This guide walks you through creating distributable packages for FileManager across all platforms.

---

## 🚀 Quick Start: Create Your First Release

### Step 1: Prepare Your Repository

```bash
# Initialize git repository (if not done)
cd filemanager
git init
git add .
git commit -m "Initial commit - FileManager v0.1.0"

# Create GitHub repository and push
git remote add origin https://github.com/YOUR_USERNAME/filemanager.git
git branch -M main
git push -u origin main
```

### Step 2: Set Up GitHub Actions

```bash
# Create workflow directory
mkdir -p .github/workflows

# Copy the release.yml (provided in artifacts)
# Place it at: .github/workflows/release.yml

# Commit and push
git add .github/
git commit -m "Add CI/CD workflow"
git push
```

### Step 3: Update Version URLs

Edit `version.go` and replace placeholders:

```go
const (
	Version     = "0.1.0"
	AppName     = "FileManager"
	ReleaseURL  = "https://api.github.com/repos/YOUR_GITHUB_USERNAME/filemanager/releases/latest"
	DownloadURL = "https://github.com/YOUR_GITHUB_USERNAME/filemanager/releases/latest"
)
```

### Step 4: Create Release

```bash
# Make sure everything is committed
git add .
git commit -m "Prepare v0.1.0 release"
git push

# Create and push tag
git tag -a v0.1.0 -m "Release v0.1.0 - Initial Release"
git push origin v0.1.0
```

This will automatically:
- Build for all platforms (Linux, macOS, Windows)
- Create distribution packages (.tar.gz, .zip)
- Generate checksums
- Create GitHub release with all files

---

## 📂 Project Structure for Packaging

Your project should have:

```
filemanager/
├── .github/
│   └── workflows/
│       └── release.yml          # CI/CD automation
├── src/
│   └── lib.rs                   # Rust library
├── main.go                      # Go main
├── operations.go                # Go FFI bindings
├── templates.go                 # Project templates
├── version.go                   # Version & update check
├── Cargo.toml                   # Rust config
├── go.mod                       # Go module
├── Makefile                     # Build automation
├── build.sh                     # Cross-platform builder
├── run.sh                       # Local runner
├── README.md                    # Main docs
├── INSTALL.md                   # Installation guide
├── USAGE_GUIDE.md               # Usage documentation
├── ADVANCED_FEATURES.md         # Advanced features
├── RELEASE.md                   # Release process
└── LICENSE                      # License file
```

---

## 🛠️ Manual Build Process

### Build for All Platforms

```bash
# Make build script executable
chmod +x build.sh

# Run build
./build.sh
```

Output structure:
```
dist/
├── linux-amd64/
│   ├── filemanager
│   ├── libfilemanager.so
│   ├── install.sh
│   ├── uninstall.sh
│   └── README.md
├── darwin-amd64/
│   ├── filemanager
│   ├── libfilemanager.dylib
│   ├── install.sh
│   ├── uninstall.sh
│   └── README.md
├── windows-amd64/
│   ├── filemanager.exe
│   ├── filemanager.dll
│   ├── install.bat
│   ├── uninstall.bat
│   └── README.md
├── filemanager-0.1.0-linux-amd64.tar.gz
├── filemanager-0.1.0-darwin-amd64.tar.gz
├── filemanager-0.1.0-windows-amd64.zip
└── checksums.txt
```

---

## 📝 Version Management

### Updating Version Number

When preparing a new release, update version in:

**1. version.go:**
```go
const Version = "0.2.0"  // Update here
```

**2. build.sh:**
```bash
VERSION="0.2.0"  // Update here
```

**3. Cargo.toml:**
```toml
[package]
version = "0.2.0"  # Update here
```

**4. README.md:**
Update all version references

### Version Numbering Scheme

Follow Semantic Versioning (SemVer):
- **MAJOR.MINOR.PATCH** (e.g., 1.2.3)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

Examples:
- `0.1.0` → Initial release
- `0.1.1` → Bug fixes
- `0.2.0` → New features
- `1.0.0` → Stable release

---

## 🔄 Release Workflow

### Standard Release Process

```bash
# 1. Update version numbers (see above)

# 2. Update CHANGELOG.md
cat >> CHANGELOG.md << EOF

## [0.2.0] - $(date +%Y-%m-%d)

### Added
- New feature X
- New feature Y

### Fixed
- Bug fix A
- Bug fix B

### Changed
- Improvement C
EOF

# 3. Commit changes
git add .
git commit -m "Bump version to 0.2.0"
git push

# 4. Create tag
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# 5. GitHub Actions will automatically:
#    - Build all platforms
#    - Create release
#    - Upload packages
```

### Pre-release (Beta/RC)

```bash
# Beta release
git tag -a v0.2.0-beta.1 -m "Release v0.2.0-beta.1"
git push origin v0.2.0-beta.1

# Release Candidate
git tag -a v0.2.0-rc.1 -m "Release v0.2.0-rc.1"
git push origin v0.2.0-rc.1
```

---

## 🧪 Testing Packages

### Before Release

Test each platform package:

**Linux:**
```bash
cd dist/linux-amd64
./filemanager --version
./filemanager --help
./filemanager  # Test interactive mode
```

**macOS:**
```bash
cd dist/darwin-amd64
./filemanager --version
# Test all features
```

**Windows:**
```cmd
cd dist\windows-amd64
filemanager.exe --version
REM Test all features
```

### Automated Tests

```bash
# Run all tests
make test

# Run specific tests
cargo test
go test ./...
```

---

## 📦 Platform-Specific Notes

### Linux Packaging

**Debian/Ubuntu (.deb):**
```bash
# Create package structure
mkdir -p filemanager-deb/DEBIAN
mkdir -p filemanager-deb/usr/local/bin
mkdir -p filemanager-deb/usr/local/lib

# Copy files
cp dist/linux-amd64/filemanager filemanager-deb/usr/local/bin/
cp dist/linux-amd64/libfilemanager.so filemanager-deb/usr/local/lib/

# Create control file
cat > filemanager-deb/DEBIAN/control << EOF
Package: filemanager
Version: 0.1.0
Section: utils
Priority: optional
Architecture: amd64
Maintainer: Your Name <your@email.com>
Description: Hybrid Rust+Go File Manager
 A powerful file management tool with project scaffolding
EOF

# Build package
dpkg-deb --build filemanager-deb
```

**RPM (Fedora/RHEL):**
```bash
# Create spec file
cat > filemanager.spec << EOF
Name:           filemanager
Version:        0.1.0
Release:        1%{?dist}
Summary:        Hybrid File Manager

License:        MIT
URL:            https://github.com/YOUR_USERNAME/filemanager
Source0:        %{name}-%{version}.tar.gz

%description
A powerful file management tool with project scaffolding

%install
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/usr/local/lib
cp filemanager %{buildroot}/usr/local/bin/
cp libfilemanager.so %{buildroot}/usr/local/lib/

%files
/usr/local/bin/filemanager
/usr/local/lib/libfilemanager.so
EOF

# Build RPM
rpmbuild -ba filemanager.spec
```

### macOS Packaging

**Homebrew Formula:**
```ruby
# Create formula at: Formula/filemanager.rb
class Filemanager < Formula
  desc "Hybrid Rust+Go File Manager"
  homepage "https://github.com/YOUR_USERNAME/filemanager"
  url "https://github.com/YOUR_USERNAME/filemanager/archive/v0.1.0.tar.gz"
  sha256 "..." # Generate with: shasum -a 256 filemanager-0.1.0.tar.gz
  license "MIT"

  depends_on "rust" => :build
  depends_on "go" => :build

  def install
    system "make"
    bin.install "filemanager"
    lib.install "target/release/libfilemanager.dylib"
  end

  test do
    system "#{bin}/filemanager", "--version"
  end
end
```

### Windows Packaging

**Chocolatey Package:**
```powershell
# Create package structure
choco new filemanager

# Edit filemanager.nuspec
# Edit tools/chocolateyinstall.ps1
# Build and publish
choco pack
choco push filemanager.0.1.0.nupkg --source https://push.chocolatey.org/
```

---

## 🔐 Code Signing

### macOS

```bash
# Sign the binary
codesign --sign "Developer ID Application: Your Name" \
  --timestamp \
  --options runtime \
  dist/darwin-amd64/filemanager

# Verify
codesign --verify --verbose dist/darwin-amd64/filemanager
```

### Windows

```cmd
REM Sign with signtool
signtool sign /f certificate.pfx /p password /tr http://timestamp.digicert.com /td sha256 /fd sha256 filemanager.exe
```

---

## 📊 Release Checklist

- [ ] Version numbers updated in all files
- [ ] CHANGELOG.md updated
- [ ] All tests passing
- [ ] Documentation updated
- [ ] README.md has correct download links
- [ ] Build script tested locally
- [ ] All platform packages tested
- [ ] Git tag created
- [ ] GitHub release created
- [ ] Checksums generated
- [ ] Release notes written
- [ ] Social media announcement prepared

---

## 🎯 Distribution Channels

### Immediate (v0.1.0)
- ✅ GitHub Releases
- ✅ Direct download (tar.gz, zip)
- ✅ Manual installation scripts

### Short-term (v0.2.0)
- [ ] Homebrew (macOS)
- [ ] AUR (Arch Linux)
- [ ] Snap Store (Linux)

### Medium-term (v0.3.0)
- [ ] Chocolatey (Windows)
- [ ] apt repository (Debian/Ubuntu)
- [ ] RPM repository (Fedora/RHEL)
- [ ] Winget (Windows)

### Long-term (v1.0.0)
- [ ] Flatpak (Linux)
- [ ] Mac App Store
- [ ] Microsoft Store

---

## 📈 Monitoring Releases

### Download Statistics

GitHub provides automatic download statistics:
- Go to: Releases → Your Release
- View download counts for each asset

### User Feedback

Monitor:
- GitHub Issues
- GitHub Discussions
- Social media mentions
- Direct user emails

---

## 🤝 Contributing to Packaging

Want to help with packaging?

1. **Test on your platform**
2. **Report issues**
3. **Submit package definitions** (Homebrew, AUR, etc.)
4. **Improve documentation**

---

## 📞 Support

- **Packaging Issues:** https://github.com/YOUR_USERNAME/filemanager/issues
- **Documentation:** https://github.com/YOUR_USERNAME/filemanager/wiki

---

**Made with ❤️ for the open source community**