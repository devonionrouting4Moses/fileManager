# FileManager - Dual Mode (Terminal + Web)

A powerful file management tool that works both in the terminal and through a web browser interface.

## 🚀 Features

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

## 📋 Menu Structure

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

## 📁 Project Structure

```
filemanager/
├── main.go                 # Main entry point with menu
├── webserver.go           # HTTP server and API handlers
├── helper.go              # Unix/Linux/Mac utilities
├── helper_windows.go      # Windows-specific utilities
├── operations.go          # File operation functions
├── templates.go           # Project templates
├── parser.go              # Tree structure parser
├── version.go             # Version management
└── filemanager_frontend/  # Web interface
    ├── index.html
    ├── css/
    │   └── style.css
    └── js/
        └── main.js
```

## 🛠️ Installation & Setup

### 1. Build the Application

```bash
# Build for your current platform
go build -o filemanager .

# Or for specific platforms
# Linux
GOOS=linux GOARCH=amd64 go build -o filemanager .

# macOS
GOOS=darwin GOARCH=amd64 go build -o filemanager .

# Windows
GOOS=windows GOARCH=amd64 go build -o filemanager.exe .
```

### 2. Set Up Frontend

Create the frontend directory structure:

```bash
mkdir -p filemanager_frontend/css
mkdir -p filemanager_frontend/js
mkdir -p filemanager_frontend/images
```

Copy the provided files:
- `index.html` → `filemanager_frontend/`
- `style.css` → `filemanager_frontend/css/`
- `main.js` → `filemanager_frontend/js/`

### 3. Make Executable (Linux/macOS)

```bash
chmod +x filemanager
```

## 🎯 Usage

### Terminal Mode (Default)

```bash
# Start interactive terminal mode
./filemanager

# Show version
./filemanager --version

# Check for updates
./filemanager --update

# Show help
./filemanager --help
```

### Web Mode

**Option 1: From Terminal Menu**
```bash
./filemanager
# Then select option 9 to launch web interface
```

**Option 2: Direct Launch**
```bash
./filemanager --web
```

The web interface will automatically:
- Start HTTP server on `http://localhost:8080`
- Try to open your default browser
- Display all available operations

## 🌐 Web Interface Features

### Operations Available in Browser

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

### API Endpoints

The web interface communicates with these endpoints:

- `GET /api/health` - Server health check
- `GET /api/templates` - Get available templates
- `POST /api/operation` - Execute file operations

## 💡 Examples

### Terminal Mode Examples

#### Create Multiple Folders
```
Enter your choice: 1
📁 Enter folder path(s): project/src project/docs project/tests
```

#### Create Files with Auto-Directory
```
Enter your choice: 2
📄 Enter file path(s): project/src/main.go project/README.md
```

#### Create from Template
```
Enter your choice: 8
Select option: 4  # Flask Web Application
📁 Enter root directory name: my-flask-app
```

### Web Mode Examples

#### Using Custom Structure
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

#### Using Tree Parser
1. Click "Create Structure" card
2. Switch to "Parse Tree" tab
3. Paste:
```
myproject/
├── src/
│   ├── main.go
│   └── utils.go
├── docs/
│   └── README.md
└── tests/
    └── test_main.go
```
4. Click "Parse & Create"

## 🔄 Switching Between Modes

### Terminal → Web
```bash
# While in terminal mode, select option 9
Enter your choice: 9
🌐 Launching Web Interface...
```

### Web → Terminal
```bash
# Press Ctrl+C in terminal where server is running
# Then restart in terminal mode
./filemanager
```

## ⚙️ Configuration

### Change Web Server Port

Edit `webserver.go`:
```go
port := "8080"  // Change to your preferred port
```

### Custom Frontend Location

Edit `webserver.go`:
```go
staticDir := "./filemanager_frontend"  // Change path
```

## 🐛 Troubleshooting

### Web Server Won't Start
- Check if port 8080 is already in use
- Try running: `lsof -i :8080` (Mac/Linux) or `netstat -ano | findstr :8080` (Windows)
- Kill the process or change the port

### Frontend Not Loading
- Ensure `filemanager_frontend` directory exists
- Check that HTML, CSS, and JS files are in correct locations
- Check browser console for errors (F12)

### API Calls Failing
- Verify server is running on `http://localhost:8080`
- Check browser console for CORS errors
- Ensure firewall isn't blocking local connections

### Permission Errors (Unix/Linux/Mac)
- Use `sudo` if creating files in protected directories
- Or change target directory to user-writable location

## 📊 Project Templates Available

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

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## 📄 License

MIT License - feel free to use this project for any purpose.

## 🔗 Links

- GitHub: [Your Repository URL]
- Documentation: [Your Docs URL]
- Issues: [Your Issues URL]

---

**Enjoy using FileManager in both terminal and web modes! 🎉**