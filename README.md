# FileManager v2 - Dual-Mode File Management Tool

A powerful, modern file management tool that works in both **terminal** and **web browser** modes. Built with a clean backend-frontend architecture: Rust handles low-level file operations while Go provides the user interface.

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                   FileManager v2                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────────┐          ┌──────────────────┐   │
│  │  Terminal UI     │          │   Web Browser    │   │
│  │  (Interactive    │          │   (Modern UI)    │   │
│  │   CLI Menu)      │          │   on :8080       │   │
│  └────────┬─────────┘          └────────┬─────────┘   │
│           │                             │             │
│           └─────────────┬───────────────┘             │
│                         │                             │
│                    Go Frontend                        │
│              (file_manager/ directory)                │
│                                                       │
│                    CGO FFI Bridge                     │
│                                                       │
│                   Rust Backend                        │
│              (rust_ffi/ directory)                    │
│                                                       │
│  - File Operations    - Directory Management         │
│  - Permissions        - System Calls                 │
│  - Copy/Move/Delete   - Create/Rename                │
│                                                       │
└─────────────────────────────────────────────────────────┘
```

## 📦 Project Structure

```
file_manager_v2/
├── file_manager/                # Go Frontend (UI Layer)
│   ├── cmd/app/                 # Application entry point
│   ├── internal/
│   │   ├── ffi/                 # CGO bindings to Rust
│   │   ├── handler/             # HTTP handlers & web server
│   │   ├── service/             # Business logic
│   │   └── repository/          # Data access
│   ├── pkg/version/             # Version management
│   ├── scripts/                 # Build & install scripts
│   ├── go.mod                   # Go dependencies
│   └── README.md                # Frontend documentation
│
├── rust_ffi/                    # Rust Backend (Core Logic)
│   ├── crates/
│   │   ├── core/                # Core library with FFI
│   │   │   ├── src/
│   │   │   │   ├── lib.rs
│   │   │   │   ├── common/      # Shared types
│   │   │   │   ├── ffi/         # C FFI wrappers
│   │   │   │   └── operations/  # File operations
│   │   │   └── Cargo.toml
│   │   └── cli/                 # CLI tool (fsops)
│   │       ├── src/main.rs
│   │       └── Cargo.toml
│   ├── Cargo.toml               # Workspace config
│   └── README.md                # Backend documentation
│
├── filemanager_frontend/        # Web UI Assets
│   ├── index.html
│   ├── css/style.css
│   ├── js/main.js
│   └── images/
│
├── .github/workflows/
│   └── release.yml              # CI/CD Release Pipeline
│
├── Makefile                     # Build automation
├── build.sh                     # Cross-platform build script
└── README.md                    # This file
```

## 🚀 Quick Start

### Prerequisites

- **Rust**: 1.70+ (for backend)
- **Go**: 1.24+ (for frontend)
- **Make**: For build automation
- **Build tools**: `gcc`, `g++`, `build-essential` (Linux)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd file_manager_v2

# Build everything
make all

# Run in terminal mode
./filemanager

# Or run in web mode
./filemanager --web
```

## 🎯 Features

### Terminal Mode
- ✅ Interactive CLI with numbered menu (0-9)
- ✅ Single and batch file/folder operations
- ✅ 12 pre-built project templates
- ✅ Interactive structure builder
- ✅ Custom structure definitions
- ✅ Tree structure parsing
- ✅ Cross-platform support (Linux, macOS, Windows)

### Web Mode
- ✅ Modern, responsive web interface
- ✅ All terminal features available in browser
- ✅ Real-time operation results
- ✅ Beautiful card-based UI
- ✅ No installation required (just open browser)
- ✅ Automatic browser launch
- ✅ REST API endpoints

### Core Operations
- 📁 **Create** - Files and folders
- 🗑️ **Delete** - Remove files and directories
- ✏️ **Rename** - Rename files and directories
- 📦 **Copy** - Copy files and directories (recursive)
- 🚚 **Move** - Move files and directories
- 🔐 **Permissions** - Change file permissions (Unix)
- 🏗️ **Structures** - Create complex project structures

## 📋 Available Operations

```
┌────────────────────────────────────────┐
│           Available Operations         │
├────────────────────────────────────────┤
│  1️⃣  Create Folder                     │
│  2️⃣  Create File                       │
│  3️⃣  Rename File/Folder                │
│  4️⃣  Delete File/Folder                │
│  5️⃣  Change Permissions                │
│  6️⃣  Move File/Folder                  │
│  7️⃣  Copy File/Folder                  │
│  8️⃣  Create Structure (Multi-entity)   │
│  9️⃣  Launch Web Interface              │
│  0️⃣  Exit                              │
└────────────────────────────────────────┘
```

## 🛠️ Building

### Using Make (Recommended)

```bash
# Build everything (Rust + Go)
make all

# Build only Rust backend
make rust

# Build only Go frontend
make go

# Development build (with debug symbols)
make dev

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

### Manual Build

```bash
# Build Rust library
cd rust_ffi
cargo build --release -p fs-operations-core

# Build Go binary
cd ../file_manager
CGO_ENABLED=1 \
CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
go build -ldflags="-s -w" -o ../filemanager ./cmd/app
```

## 🏃 Running

### Terminal Mode (Default)

```bash
# Interactive menu
./filemanager

# Show version
./filemanager --version

# Check for updates
./filemanager --update

# Show help
./filemanager --help
```

### Web Mode

```bash
# Start web server (opens browser automatically)
./filemanager --web

# Server runs on: http://localhost:8080
```

## 🌐 Web Interface

### Access

Open your browser to: `http://localhost:8080`

### API Endpoints

- `GET /api/health` - Server health check
- `GET /api/templates` - Get available templates
- `POST /api/operation` - Execute file operations

### Features

1. **Create Folder** - Single or multiple folders
2. **Create File** - Files with auto-directory creation
3. **Rename** - Rename files or folders
4. **Delete** - Delete with confirmation
5. **Permissions** - Change file/folder permissions
6. **Move** - Move files or folders
7. **Copy** - Copy files or folders
8. **Create Structure** - Three modes:
   - **Templates**: 12 pre-built project templates
   - **Custom**: Define with `d:` and `f:` prefixes
   - **Parse Tree**: Paste tree-format structure

## 📦 Project Templates

1. Java Traits Project Structure
2. Standard Go Project Structure
3. Rust Project with Workspace
4. Python Flask Web Application
5. Python FastAPI REST API
6. Java RMI Distributed Application
7. Java Swing Desktop Application
8. Java Spring Boot REST API
9. Flutter Mobile Application
10. React Frontend Application
11. Next.js Full-Stack Application
12. Simple HTML/CSS/JavaScript Website

## 💡 Examples

### Terminal Mode - Create Project Structure

```bash
./filemanager
Enter your choice: 8
Select option: 2  # Standard Go Project
📁 Enter root directory name: my-go-project
```

### Web Mode - Custom Structure

1. Click "Create Structure" card
2. Switch to "Custom" tab
3. Enter:
```
d:myproject
d:myproject/src
d:myproject/tests
f:myproject/README.md
f:myproject/src/main.go
```
4. Click "Create Structure"

### Web Mode - Parse Tree

1. Click "Create Structure" card
2. Switch to "Parse Tree" tab
3. Paste:
```
myproject/
├── src/
│   ├── main.go
│   └── utils.go
├── tests/
│   └── test_main.go
└── README.md
```
4. Click "Parse & Create"

## 🔗 Component Documentation

### Go Frontend (`file_manager/`)

The user-facing application layer that provides both terminal and web interfaces.

- **Terminal UI**: Interactive menu-driven interface
- **Web Server**: HTTP server with REST API
- **FFI Bridge**: CGO bindings to Rust backend

See [file_manager/README.md](./file_manager/README.md) for detailed documentation.

### Rust Backend (`rust_ffi/`)

Low-level file system operations with C FFI bindings.

- **Core Library**: `fs-operations-core` - Handles all file operations
- **CLI Tool**: `fsops` - Command-line interface to operations
- **FFI Layer**: C-compatible function signatures
- **Operations**: Create, delete, rename, move, copy, permissions

See [rust_ffi/README.md](./rust_ffi/README.md) for detailed documentation.

### Web Frontend (`filemanager_frontend/`)

Modern, responsive web interface.

- **HTML**: Semantic structure
- **CSS**: TailwindCSS styling
- **JavaScript**: Interactive functionality
- **API Integration**: REST API communication

See [filemanager_frontend/README.md](./filemanager_frontend/README.md) for detailed documentation.

## 🔧 Development

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Lint Code

```bash
make lint
```

### Generate Documentation

```bash
make docs
```

### Setup Development Environment

```bash
make setup-dev
```

## 📦 Installation

### From Source (Recommended for Development)

```bash
# Clone the repository
git clone https://github.com/DevChigarlicMoses/FileManager.git
cd file_manager_v2

# Build and install
make install

# Uninstall
make uninstall
```

### System-wide Installation (Manual)

```bash
# Build first
make all

# Copy binary
sudo cp filemanager /usr/local/bin/

# Copy Rust library
sudo cp rust_ffi/target/release/libfs_operations_core.so /usr/local/lib/

# Update library cache (Linux)
sudo ldconfig

# Make executable
sudo chmod +x /usr/local/bin/filemanager

# Verify installation
filemanager --version
```

### Platform-Specific Installation

#### Ubuntu/Debian
```bash
# Via DEB package (when available)
wget https://github.com/DevChigarlicMoses/FileManager/releases/download/v2.0.1/filemanager_0.1.2_amd64.deb
sudo dpkg -i filemanager_0.1.2_amd64.deb
sudo apt install -f  # Fix dependencies if needed
```

#### Snap Store (Coming Soon)
```bash
sudo snap install filemanager
# Auto-updates every 6 hours
```

#### Windows
```bash
# Download latest release
# Visit: https://github.com/DevChigarlicMoses/FileManager/releases/latest
# Download: filemanager-0.1.2-windows-amd64.zip
# Extract and run: filemanager.exe
```

#### macOS
```bash
# Via Homebrew (coming soon)
brew install filemanager
```

### Verify Installation

```bash
filemanager --version
which filemanager
ldd $(which filemanager)  # Check library dependencies
```

## 🚀 CI/CD Pipeline

### GitHub Actions Release Workflow

The project includes a comprehensive release workflow (`.github/workflows/release.yml`) that:

1. **Builds** for multiple platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)

2. **Packages** distributions with:
   - Go binary
   - Rust library
   - CLI tool
   - Documentation
   - Installer scripts

3. **Creates** GitHub releases with:
   - Platform-specific archives
   - SHA256 checksums
   - Installation instructions
   - Release notes

### Triggering a Release

```bash
# Create a version tag
git tag v2.0.1
git push origin v2.0.1

# GitHub Actions automatically:
# - Builds for all platforms
# - Creates release packages
# - Publishes to GitHub Releases
```

## ⚙️ Configuration

### Web Server Port

Edit `file_manager/internal/handler/webserver.go`:
```go
port := "8080"  // Change to your preferred port
```

### Frontend Location

Edit `file_manager/internal/handler/webserver.go`:
```go
staticDir := "./filemanager_frontend"  // Change path
```

### Build Version

Edit `Makefile`:
```makefile
VERSION := 0.1.2  # Update version
```

## 🐛 Troubleshooting

### Build Issues

**CGO Errors**
```bash
# Ensure Rust library is built
cd rust_ffi && cargo build --release -p fs-operations-core

# Verify CGO_ENABLED=1
export CGO_ENABLED=1

# Check library path
ls -la rust_ffi/target/release/libfs_operations_core.so
```

**Missing Dependencies**
```bash
# Linux
sudo apt-get install build-essential

# macOS
xcode-select --install

# Windows
# Install Visual Studio Build Tools
```

### Runtime Issues

**Web Server Won't Start**
```bash
# Check if port is in use
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill process or change port
```

**Frontend Not Loading**
- Check browser console (F12)
- Verify `filemanager_frontend` directory exists
- Check file permissions

**Permission Errors**
```bash
# Use sudo for protected directories
sudo ./filemanager

# Or use user-writable location
./filemanager /home/user/my_files
```

## � Security Features

- **Safe FFI**: All string conversions checked for UTF-8 validity
- **Memory Safety**: Rust prevents buffer overflows and memory leaks
- **Permission Validation**: Unix permission checks before operations
- **Confirmation Prompts**: Destructive operations require confirmation
- **Sandboxed Execution**: Snap confinement with controlled permissions (when available)
- **No Runtime Dependencies**: Compiled native binaries with minimal attack surface

## ⚡ Performance Benefits

- **Rust Core**: Zero-cost abstractions, no garbage collection overhead
- **Concurrent Go**: Easy to extend with goroutines for batch operations
- **Native Code**: Compiled binaries with no runtime dependencies
- **Small Binary**: Optimized builds under 15MB
- **Efficient FFI**: Minimal overhead in Rust-Go communication

### Benchmarks

Run Rust benchmarks:
```bash
make bench
```

### Optimization

- Rust backend uses release optimizations (`opt-level = 3`)
- Go binary stripped of debug symbols (`-ldflags="-s -w"`)
- Link-time optimization enabled in Rust
- Efficient FFI communication layer
- Minimal memory footprint

## 📋 System Requirements

- **OS**: Linux (Ubuntu 20.04+, Debian 11+, any modern distro), macOS 10.15+, Windows 10+
- **Architecture**: x86_64 (amd64) or ARM64
- **RAM**: 512 MB minimum, 1 GB recommended
- **Disk Space**: 50 MB for installation
- **Dependencies**: 
  - Linux: libc6 >= 2.31 (usually pre-installed)
  - macOS: None (standalone executable)
  - Windows: None (standalone executable)

## 🔄 Updates & Upgrade Strategy

FileManager uses **Semantic Versioning (SemVer)** with a **hybrid notification approach** to keep you informed about updates while respecting your workflow.

### Understanding Version Numbers

```
Version Format: MAJOR.MINOR.PATCH
Example: v2.0.1

🔧 PATCH (v1.5.Z) - Bug fixes, security patches
   → Silent/automatic installation
   → Minimal disruption

✨ MINOR (v1.Y.0) - New features, improvements
   → Subtle in-app notification
   → Update at your convenience

🚀 MAJOR (vX.0.0) - Breaking changes, redesigns
   → Full-screen notification
   → Review required before updating
```

### Checking for Updates

```bash
# Check for available updates
filemanager --update

# The app will:
# 1. Detect your current version
# 2. Compare with latest release
# 3. Display appropriate notification based on change type
# 4. Provide download and installation instructions
```

### Update Methods

#### From Source
```bash
# Pull latest changes
git pull origin main

# Rebuild
make clean && make all
```

#### Snap Store (Coming Soon)
```bash
# Auto-updates every 6 hours
# Force update
sudo snap refresh filemanager
```

#### DEB Package
```bash
# Download new version
wget https://github.com/DevChigarlicMoses/FileManager/releases/download/v2.0.1/filemanager_0.1.3_amd64.deb
sudo dpkg -i filemanager_0.1.3_amd64.deb
```

### Update Notifications

**PATCH Updates** (e.g., v2.0.1 → v2.0.1)
```
🔧 PATCH UPDATE AVAILABLE: v2.0.1 → v2.0.1
─ Security & Bug Fixes ─

✅ This is a safe, backwards-compatible update.
💡 It will be installed automatically on next restart.
```

**MINOR Updates** (e.g., v2.0.1 → v2.0.1)
```
════════════════════════════════════════════════════════════
║ ✨ NEW FEATURES AVAILABLE: v2.0.1 → v2.0.1
║ ────────────────────────────────────────────────────────
║ 📊 Update Type: MINOR (New Features & Improvements)
║ 📈 User Impact: Low to Moderate
║ 🔄 Update Strategy: Subtle In-App Notification
║ ────────────────────────────────────────────────────────
║
║ 💡 Tip: Check the release notes to see what's new!
║ 🔗 You can update at your convenience.
════════════════════════════════════════════════════════════
```

**MAJOR Updates** (e.g., v2.0.1 → v2.0.1)
```
══════════════════════════════════════════════════════════════════
║                  🚀 MAJOR UPGRADE AVAILABLE                     ║
║                                                                  ║
║  Current Version: v2.0.1
║  Available Version: v2.0.1
║                                                                  ║
║  ⚠️  IMPORTANT: This is a major upgrade with breaking changes.  ║
║  ✅ Action Required: Please review release notes before updating.
║  🔗 You may need to reconfigure settings or migrate data.
══════════════════════════════════════════════════════════════════
```

### Migration for Major Updates

When upgrading between major versions (e.g., v1.x → v2.0), see [MIGRATION.md](./docs/MIGRATION.md) for:
- Pre-migration checklist
- Step-by-step migration guide
- Configuration migration
- Troubleshooting
- Rollback instructions

### Update Strategy Details

For comprehensive information about how FileManager handles updates, see [UPDATE_STRATEGY.md](./docs/UPDATE_STRATEGY.md):
- Semantic versioning explained
- Notification strategies by update type
- Release notes best practices
- Security considerations
- Implementation details

## 📄 License

MIT OR Apache-2.0

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## 📞 Support

- **Bug Reports**: [GitHub Issues](https://github.com/DevChigarlicMoses/FileManager/issues)
- **Feature Requests**: [GitHub Discussions](https://github.com/DevChigarlicMoses/FileManager/discussions)
- **Email**: moses.muranja@strathmore.edu
- **Documentation**: See component README files
- **Wiki**: [GitHub Wiki](https://github.com/DevChigarlicMoses/FileManager/wiki)

## 🔗 Links

- **GitHub**: [github.com/DevChigarlicMoses/FileManager](https://github.com/DevChigarlicMoses/FileManager)
- **Releases**: [Latest Release](https://github.com/DevChigarlicMoses/FileManager/releases/latest)
- **Snap Store**: [snapcraft.io/filemanager](https://snapcraft.io/filemanager) (coming soon)
- **Changelog**: [CHANGELOG.md](./CHANGELOG.md)

## 🎯 Roadmap

### Completed ✅
- [x] Basic file operations (create, delete, rename, copy, move)
- [x] Permission management
- [x] Terminal mode with interactive menu
- [x] Web mode with modern UI
- [x] Dual-mode (terminal + web) support
- [x] Multi-platform builds (Linux, macOS, Windows)
- [x] Project templates (12 templates)
- [x] Custom structure definitions
- [x] Tree structure parsing

### In Progress 🚀
- [ ] Snap Store distribution
- [ ] DEB package distribution
- [ ] PPA for easier APT installation
- [ ] Homebrew support for macOS

### Planned 📋
- [ ] Batch operations support
- [ ] File search functionality
- [ ] Archive support (zip, tar.gz)
- [ ] File preview capabilities
- [ ] Undo/Redo functionality
- [ ] Favorites/Bookmarks
- [ ] File comparison tool
- [ ] Sync operations
- [ ] Cloud storage integration

See [ROADMAP.md](./ROADMAP.md) for detailed future plans.

## 🌟 Related Projects

- **v1**: [FileManager v1](../file_manager_v1/) - Original terminal-only version
- **Rust FFI**: [fs-operations-core](./rust_ffi/) - Core library documentation
- **Web Frontend**: [filemanager_frontend](./filemanager_frontend/) - Web UI assets

## 📈 Version History

- **v2.0.1** - Current version (dual-mode with web interface)
- **v2.0.1** - Previous release
- **v2.0.1** - Initial release

---

**FileManager v2 - Making file management simple, powerful, and accessible** 🎉

Built with ❤️ by Moses Muranja and contributors

Powered by **Rust** for safety and performance, **Go** for simplicity and concurrency, and modern web technologies for the UI.
