# 🗂️ FileManager - Hybrid Rust + Go File Management Tool

A high-performance, terminal-based file and folder manager that combines the safety and speed of Rust with the concurrency and simplicity of Go.

## 🏗️ Architecture

- **Rust Layer**: Handles all performance-critical file operations (create, delete, rename, copy, move, permissions)
- **Go Layer**: Provides the CLI interface, user input handling, and calls Rust functions via CGO/FFI
- **Communication**: Rust code compiled to `.so` shared library, linked with Go using CGO

## ✨ Features

1. **Create Folder** - Create directories with full path support
2. **Create File** - Create new empty files
3. **Rename** - Rename files or folders
4. **Delete** - Safely delete files or folders (with confirmation)
5. **Change Permissions** - Modify Unix file permissions (chmod)
6. **Move** - Move files or folders to new locations
7. **Copy** - Copy files or entire directory trees
8. **Exit** - Clean exit from the application

---

## 📥 Installation

### Ubuntu/Debian (Recommended)

#### Via Snap Store (Easiest - Auto Updates)
```bash
sudo snap install filemanager
```

#### Via DEB Package
```bash
# Download latest release
wget https://github.com/devonionrouting4Moses/fileManager/releases/download/v0.1.2/filemanager_0.1.2_amd64.deb

# Install
sudo dpkg -i filemanager_0.1.2_amd64.deb

# Fix dependencies if needed
sudo apt install -f
```

#### Via PPA (Coming Soon)
```bash
sudo add-apt-repository ppa:mosesmuranja/filemanager
sudo apt update
sudo apt install filemanager
```

### Windows

```bash
# Download the latest Windows release
# Visit: https://github.com/devonionrouting4Moses/fileManager/releases/latest
# Download: filemanager-0.1.2-windows-amd64.zip

# Extract and run
filemanager.exe
```

### Other Linux Distributions

Download the appropriate package from [Releases](https://github.com/devonionrouting4Moses/fileManager/releases/latest):
- **Arch Linux**: `filemanager-0.1.2-linux-amd64.tar.gz`
- **Fedora/RHEL**: `filemanager-0.1.2-linux-amd64.tar.gz`
- **Any Linux**: Snap Store (works everywhere)

### Verify Installation

```bash
filemanager --version
```

---

## 🚀 Usage

Simply run the command:

```bash
filemanager
```

Then follow the interactive menu to perform file operations.

## 📖 Usage Examples

### Creating a Folder
```
Enter your choice: 1
📁 Enter folder path: /tmp/myproject
✅ Folder created: /tmp/myproject
```

### Changing Permissions
```
Enter your choice: 5
🔒 Enter path: /tmp/script.sh
Enter permissions (octal, e.g., 755): 755
✅ Permissions changed: /tmp/script.sh (0755)
```

### Copying Files
```
Enter your choice: 7
📋 Enter source path: /tmp/original.txt
Enter destination path: /tmp/backup.txt
✅ Copied: /tmp/original.txt -> /tmp/backup.txt
```

---

## 🔄 Updates

- **Snap**: Updates automatically every 6 hours
  ```bash
  # Force update
  sudo snap refresh filemanager
  ```

- **DEB Package**: Download new version from [releases](https://github.com/devonionrouting4Moses/fileManager/releases/latest)
  ```bash
  wget https://github.com/devonionrouting4Moses/fileManager/releases/download/v0.1.3/filemanager_0.1.3_amd64.deb
  sudo dpkg -i filemanager_0.1.3_amd64.deb
  ```

- **PPA**: 
  ```bash
  sudo apt update && sudo apt upgrade
  ```

---

## 🔧 For Developers

Want to build from source or contribute? See [DEVELOPMENT.md](DEVELOPMENT.md)

### Quick Development Setup

```bash
# Prerequisites
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
# Install Go from: https://golang.org/dl/

# Clone and build
git clone https://github.com/devonionrouting4Moses/fileManager.git
cd fileManager
make
```

---

## 🛠️ Project Structure

```
filemanager/
├── main.go              # CLI menu and user interaction
├── operations.go        # Go wrappers for Rust FFI calls
├── go.mod              # Go module definition
├── src/
│   └── lib.rs          # Rust FFI functions
├── Cargo.toml          # Rust dependencies and build config
├── Makefile            # Build automation
├── snap/               # Snap package configuration
│   ├── snapcraft.yaml  # Snap build manifest
│   └── local/          # Desktop files and icons
├── debian-package-build.sh  # DEB package builder
└── README.md           # This file
```

---

## 🔒 Security Features

- **Safe FFI**: All string conversions checked for UTF-8 validity
- **Memory Safety**: Rust prevents buffer overflows and memory leaks
- **Permission Validation**: Unix permission checks before operations
- **Confirmation Prompts**: Destructive operations require confirmation
- **Snap Confinement**: Sandboxed execution with controlled permissions

---

## ⚡ Performance Benefits

- **Rust Core**: Zero-cost abstractions, no garbage collection overhead
- **Concurrent Go**: Easy to extend with goroutines for batch operations
- **Native Code**: Compiled binaries with no runtime dependencies
- **Small Binary**: Optimized builds under 10MB

---

## 🐛 Troubleshooting

### Snap Permission Issues

If you can't access certain directories:
```bash
# Connect home directory access
sudo snap connect filemanager:home

# Connect removable media access
sudo snap connect filemanager:removable-media

# View all connections
snap connections filemanager
```

### DEB Installation Issues

```bash
# If dependencies are missing
sudo apt install -f

# If library not found
sudo ldconfig

# Check installation
which filemanager
ldd $(which filemanager)
```

### Windows Issues

- Ensure you have the latest Visual C++ Redistributable installed
- Extract the entire ZIP contents before running
- Run as Administrator if permission errors occur

### General Issues

- Check [GitHub Issues](https://github.com/devonionrouting4Moses/fileManager/issues)
- Read [INSTALL.md](INSTALL.md) for detailed installation help
- View logs: `snap logs filemanager` (for Snap installations)

---

## 📊 System Requirements

- **OS**: Linux (Ubuntu 20.04+, Debian 11+, any modern distro), Windows 10+
- **Architecture**: x86_64 (amd64) or ARM64
- **RAM**: 512 MB minimum, 1 GB recommended
- **Disk Space**: 50 MB for installation
- **Dependencies**: 
  - Linux: libc6 >= 2.31 (usually pre-installed)
  - Windows: None (standalone executable)

---

## 🤝 Contributing

We welcome contributions! Here's how to help:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests: `make test`
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

---

## 📄 License

MIT License - feel free to use in your projects!

See [LICENSE](LICENSE) for full details.

---

## 🙏 Acknowledgments

- Built with **Rust** for safety and performance
- Powered by **Go** for simplicity and concurrency
- Uses **CGO/FFI** for seamless interop
- Packaged with **Snapcraft** for universal Linux distribution
- Available via **APT** for Debian/Ubuntu users

---

## 📚 Documentation

- [Installation Guide](INSTALL.md) - Detailed installation instructions
- [Usage Guide](USAGE_GUIDE.md) - Complete feature documentation
- [Development Guide](DEVELOPMENT.md) - Build from source
- [Advanced Features](ADVANCED_FEATURES.md) - Power user tips
- [Publishing Guide](PACKAGING.md) - For package maintainers

---

## 🌟 Star History

If you find FileManager useful, please consider giving it a star! ⭐

---

## 📞 Support

- **Bug Reports**: [GitHub Issues](https://github.com/devonionrouting4Moses/fileManager/issues)
- **Feature Requests**: [GitHub Discussions](https://github.com/devonionrouting4Moses/fileManager/discussions)
- **Email**: moses.muranja@strathmore.edu
- **Documentation**: [Wiki](https://github.com/devonionrouting4Moses/fileManager/wiki)

---

## 🔗 Links

- **Snap Store**: [snapcraft.io/filemanager](https://snapcraft.io/filemanager)
- **GitHub**: [github.com/devonionrouting4Moses/fileManager](https://github.com/devonionrouting4Moses/fileManager)
- **Releases**: [Latest Release](https://github.com/devonionrouting4Moses/fileManager/releases/latest)
- **Changelog**: [CHANGELOG.md](CHANGELOG.md)

---

## 🎯 Roadmap

- [x] Basic file operations (create, delete, rename, copy, move)
- [x] Permission management
- [x] Snap Store distribution
- [x] DEB package distribution
- [ ] PPA for easier APT installation
- [ ] Batch operations support
- [ ] File search functionality
- [ ] Archive support (zip, tar.gz)
- [ ] GUI version
- [ ] macOS support via Homebrew
- [ ] File preview capabilities

See [ROADMAP.md](ROADMAP.md) for detailed future plans.

---

**Happy File Managing! 🚀**

Made with ❤️ by Moses Muranja and contributors