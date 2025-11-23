# FileManager - Go Frontend

The Go frontend application that provides a user-friendly interface for file system operations. This is the **frontend/UI layer** that consumes the Rust FFI backend for all file operations.

## 🏗️ Architecture

This is the **frontend** component in a backend-frontend architecture:

```
┌─────────────────────────────────────┐
│   FileManager (Go Frontend)         │
│  - Terminal UI                      │
│  - Web Server & API                 │
│  - User Interface                   │
└────────────┬────────────────────────┘
             │ (CGO FFI)
┌────────────▼────────────────────────┐
│   Rust FFI Backend                  │
│  - File Operations                  │
│  - Directory Management             │
│  - System Calls                     │
└─────────────────────────────────────┘
```

## 📁 Project Structure

```
file_manager/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── internal/
│   ├── ffi/                     # CGO bindings to Rust library
│   │   ├── operations.go        # Unix/Linux/macOS FFI
│   │   └── operations_windows.go # Windows FFI
│   ├── handler/                 # HTTP handlers & web server
│   │   ├── handler.go
│   │   └── webserver.go
│   ├── service/                 # Business logic
│   │   └── templates.go
│   └── repository/              # Data access layer
├── pkg/
│   └── version/                 # Version management
├── scripts/                      # Build & installation scripts
├── go.mod                        # Go module definition
└── README.md
```

## 🎯 Features

### Terminal Mode
- Interactive CLI with numbered menu (0-9)
- Single and batch file/folder operations
- 12 project templates
- Interactive structure builder
- Custom structure definitions
- Tree structure parsing

### Web Mode
- Modern, responsive web interface
- All terminal features available in browser
- Real-time operation results
- Beautiful card-based UI
- No installation required (just open browser)

## 🚀 Quick Start

### Build

```bash
# Build with Rust backend (from project root)
make all

# Or manually:
cd ../rust_ffi && cargo build --release -p fs-operations-core
cd ../file_manager && CGO_ENABLED=1 CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" go build -o ../filemanager ./cmd/app
```

### Run

```bash
# Terminal mode
./filemanager

# Web mode
./filemanager --web

# Show version
./filemanager --version

# Show help
./filemanager --help
```

## 📋 Available Operations

1. **Create Folder** - Create single or multiple folders
2. **Create File** - Create files with auto-directory creation
3. **Rename** - Rename files or folders
4. **Delete** - Delete files or folders (with confirmation)
5. **Permissions** - Change file/folder permissions (Unix/Linux)
6. **Move** - Move files or folders
7. **Copy** - Copy files or folders
8. **Create Structure** - Three modes:
   - **Templates**: 12 pre-built project templates
   - **Custom**: Define structure with `d:` and `f:` prefixes
   - **Parse Tree**: Paste tree-format structure
9. **Launch Web Interface** - Start the web server
10. **Exit** - Close the application

## 🌐 Web Interface

### Start Web Server

```bash
./filemanager --web
```

The server will:
- Start on `http://localhost:8080`
- Automatically open your default browser
- Display all available operations

### API Endpoints

- `GET /api/health` - Server health check
- `GET /api/templates` - Get available templates
- `POST /api/operation` - Execute file operations

## 💡 Examples

### Terminal Mode

```bash
# Create multiple folders
Enter your choice: 1
📁 Enter folder path(s): project/src project/docs project/tests

# Create files with auto-directory
Enter your choice: 2
📄 Enter file path(s): project/src/main.go project/README.md

# Create from template
Enter your choice: 8
Select option: 4  # Flask Web Application
📁 Enter root directory name: my-flask-app
```

### Web Mode

1. Click "Create Structure" card
2. Switch to "Custom" tab
3. Enter:
```
d:myproject
d:myproject/src
d:myproject/docs
f:myproject/src/main.go
f:myproject/README.md
```
4. Click "Create Structure"

## 🔗 Dependencies

### Go Dependencies
- `github.com/gorilla/mux` - HTTP router for web server
- `golang.org/x/sys` - System calls

### Rust Dependencies (via FFI)
- `fs-operations-core` - Rust backend library for file operations

## 🛠️ Development

### Build for Development

```bash
make dev
```

This builds with debug symbols and race condition detection.

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

## ⚙️ Configuration

### Change Web Server Port

Edit `internal/handler/webserver.go`:
```go
port := "8080"  // Change to your preferred port
```

### Custom Frontend Location

Edit `internal/handler/webserver.go`:
```go
staticDir := "./filemanager_frontend"  // Change path
```

## 🐛 Troubleshooting

### Build Fails with CGO Errors
- Ensure Rust library is built: `cd ../rust_ffi && cargo build --release -p fs-operations-core`
- Check that `CGO_ENABLED=1` is set
- Verify `CGO_LDFLAGS` points to correct library path

### Web Server Won't Start
- Check if port 8080 is already in use
- Try: `lsof -i :8080` (Mac/Linux) or `netstat -ano | findstr :8080` (Windows)
- Kill the process or change the port

### Frontend Not Loading
- Ensure `filemanager_frontend` directory exists
- Check that HTML, CSS, and JS files are in correct locations
- Check browser console for errors (F12)

### Permission Errors
- Use `sudo` if creating files in protected directories
- Or change target directory to user-writable location

## 📄 License

MIT License - feel free to use this project for any purpose.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

---

**Part of FileManager v2 - A powerful dual-mode file management tool** 🎉
