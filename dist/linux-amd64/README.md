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

## 📋 Prerequisites

- **Rust** (1.70+): `curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh`
- **Go** (1.21+): Download from [golang.org](https://golang.org/dl/)
- **GCC/Clang**: Required for CGO compilation
- **Make**: Build automation tool

## 🚀 Quick Start

### 1. Clone and Build

```bash
# Build everything
make

# Or build step-by-step
make rust    # Build Rust library
make go      # Build Go binary
```

### 2. Run

```bash
# Run directly
make run

# Or execute the binary
./filemanager
```

### 3. Install System-Wide (Optional)

```bash
sudo make install
```

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
└── README.md           # This file
```

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

## 🧪 Testing

```bash
make test
```

## 🧹 Cleaning

```bash
make clean
```

## 🔧 Development Build

For debugging with symbols:

```bash
make dev
```

## 📝 Makefile Targets

- `make all` - Build everything (default)
- `make rust` - Build Rust library only
- `make go` - Build Go binary only
- `make run` - Build and run
- `make install` - Install system-wide
- `make clean` - Remove build artifacts
- `make test` - Run tests
- `make dev` - Development build
- `make help` - Show all targets

## 🔒 Security Features

- **Safe FFI**: All string conversions checked for UTF-8 validity
- **Memory Safety**: Rust prevents buffer overflows and memory leaks
- **Permission Validation**: Unix permission checks before operations
- **Confirmation Prompts**: Destructive operations require confirmation

## ⚡ Performance Benefits

- **Rust Core**: Zero-cost abstractions, no garbage collection overhead
- **Concurrent Go**: Easy to extend with goroutines for batch operations
- **Native Code**: Compiled binaries with no runtime dependencies

## 🐛 Troubleshooting

### "cannot find -lfilemanager"
Make sure the Rust library is built first:
```bash
make rust
```

### "undefined reference to dlsym"
Add `-ldl` to CGO_LDFLAGS in operations.go (already included)

### Permissions errors
For system-wide installation:
```bash
sudo make install
sudo ldconfig
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## 📄 License

MIT License - feel free to use in your projects!

## 🙏 Acknowledgments

- Built with **Rust** for safety and performance
- Powered by **Go** for simplicity and concurrency
- Uses **CGO/FFI** for seamless interop

---

**Happy File Managing! 🚀**