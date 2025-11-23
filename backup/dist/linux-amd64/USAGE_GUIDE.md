# 🗂️ FileManager - Complete Usage Guide

## 🎯 Features Overview

### Single Entity Operations (Options 1-7)
Traditional file/folder operations for individual items.

### Multi-Entity Operations (Option 8)
Create entire project structures with **12 predefined templates** or custom definitions.

---

## 📦 Available Project Templates

### 1. Java Traits Project
Complete Java Traits implementation with annotation processing
- Maven multi-module structure
- Runtime, Processor, and Plugin modules
- Trait interfaces and implementations

### 2. Standard Go Project
Modern Go application structure
- `cmd/`, `internal/`, `pkg/` organization
- API, configs, and documentation folders
- Makefile and go.mod included

### 3. Rust Workspace
Multi-crate Rust project
- Core, CLI, and Lib crates
- Examples and benchmarks
- Workspace-level Cargo.toml

### 4. Python Flask Web App
Full-featured Flask application
- MVC structure with blueprints
- Templates and static assets
- Configuration and migrations setup
- Requirements.txt with dependencies

### 5. Python FastAPI REST API
Modern async Python API
- API versioning (v1 structure)
- Pydantic schemas and models
- Configuration management
- Uvicorn server setup

### 6. Java RMI Distributed System
Remote Method Invocation project
- Server and client implementations
- Common interfaces
- Maven POM configuration

### 7. Java Swing Desktop App
GUI application structure
- MVC pattern organization
- UI, Model, Controller packages
- Main frame setup

### 8. Java Spring Boot API
Enterprise REST API
- Controller, Service, Repository layers
- JPA and H2 database
- Application properties
- Maven dependencies

### 9. Flutter Mobile App
Cross-platform mobile application
- Screens, Widgets, Models organization
- Services and utilities
- pubspec.yaml with dependencies
- Material Design setup

### 10. React Frontend
Single Page Application
- Component-based structure
- React Router ready
- Hooks and services folders
- CSS styling setup

### 11. Next.js Full-Stack App
Modern web framework
- App Router (Next.js 14+)
- API routes included
- Server and client components
- SEO-optimized

### 12. HTML/CSS/JavaScript Website
Simple static website
- Clean HTML5 structure
- Organized CSS and JS
- Assets and images folders
- Responsive design ready

---

## 🚀 Quick Start Guide

### Building the Project
```bash
# Clean and rebuild
make clean
make

# Run the application
./run.sh
```

### Creating a Project

#### Example 1: Java Spring Boot API
```
Enter your choice: 8
Select option: 8
📁 Enter root directory name: my-spring-api
```

Result:
```
my-spring-api/
├── src/
│   ├── main/java/com/example/demo/
│   │   ├── DemoApplication.java
│   │   ├── controller/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── model/
│   │   └── dto/
│   ├── resources/
│   │   └── application.properties
│   └── test/java/com/example/demo/
└── pom.xml
```

#### Example 2: React + FastAPI Full-Stack
Create both projects separately:

**Backend:**
```
Enter your choice: 8
Select option: 5
📁 Enter root directory name: backend-api
```

**Frontend:**
```
Enter your choice: 8
Select option: 10
📁 Enter root directory name: frontend-app
```

---

## 🛠️ Advanced Features

### Custom Structure Builder
Define your own project structure line by line:

```
Enter your choice: 8
Select option: 13

📝 Custom Structure Builder
> d:myapp
> d:myapp/core
> d:myapp/api
> d:myapp/tests
> f:myapp/core/__init__.py
> f:myapp/api/routes.py
> f:myapp/tests/test_main.py
> f:myapp/README.md
> done

🔨 Creating 8 entities...
  ✅ 📁 myapp
  ✅ 📁 myapp/core
  ✅ 📁 myapp/api
  ✅ 📁 myapp/tests
  ✅ 📄 myapp/core/__init__.py
  ✅ 📄 myapp/api/routes.py
  ✅ 📄 myapp/tests/test_main.py
  ✅ 📄 myapp/README.md

📊 Summary: 8 succeeded, 0 failed
```

### Parse Tree Structure (Paste Structure)
Paste an existing project structure directly:

```
Enter your choice: 8
Select option: 14

🌳 Parse Tree Structure
Paste your folder structure (tree format):

my-app/
├── src/
│   ├── components/
│   │   ├── Header.jsx
│   │   └── Footer.jsx
│   ├── pages/
│   │   ├── Home.jsx
│   │   └── About.jsx
│   └── App.jsx
├── public/
│   └── index.html
└── package.json
END

📊 Parsed structure: 5 directories, 7 files
Proceed with creation? (yes/no): yes

🔨 Creating structure...
  ✅ 📁 my-app
  ✅ 📁 my-app/src
  ✅ 📁 my-app/src/components
  ✅ 📁 my-app/src/pages
  ✅ 📁 my-app/public
  ✅ 📄 my-app/src/components/Header.jsx
  ✅ 📄 my-app/src/components/Footer.jsx
  ✅ 📄 my-app/src/pages/Home.jsx
  ✅ 📄 my-app/src/pages/About.jsx
  ✅ 📄 my-app/src/App.jsx
  ✅ 📄 my-app/public/index.html
  ✅ 📄 my-app/package.json

📊 Summary: 12 succeeded, 0 failed
```

### Interactive Builder
Build structures using shell-like commands:

```
Enter your choice: 8
Select option: 15

🛠️  Interactive Structure Builder
📁 Enter root directory: my-project
✅ Root directory created

Commands:
  mkdir <path>  - Create directory
  touch <path>  - Create file
  done          - Finish building
  cancel        - Abort and return

> mkdir src
  ✅ 📁 my-project/src
> mkdir tests
  ✅ 📁 my-project/tests
> touch src/main.py
  ✅ 📄 my-project/src/main.py
> touch tests/test_main.py
  ✅ 📄 my-project/tests/test_main.py
> touch README.md
  ✅ 📄 my-project/README.md
> done

📊 Summary: 6 succeeded, 0 failed
```

---

## 🔄 Navigation Features

### Back Navigation
Every submenu now supports **back navigation**:

- Type **`back`** or **`b`** to return to the previous menu
- Type **`cancel`** to abort current operation
- All inputs are validated before processing

### Navigation Flow
```
Main Menu
   └─> Option 8: Create Structure
         ├─> Template Selection Menu
         │     ├─> Template 1-12 (with back option)
         │     ├─> Custom Structure (with cancel)
         │     ├─> Parse Tree (with cancel)
         │     ├─> Interactive Builder (with cancel)
         │     └─> Back to Main Menu
         └─> Returns to Main Menu
```

### Example Navigation Session
```
Enter your choice: 8

Select option: 5  # FastAPI template

📁 Enter root directory name: back  # Changed mind
# Returns to template selection

Select option: 10  # React template instead

📁 Enter root directory name: my-react-app
# Proceeds with creation
```

---

## 💡 Pro Tips

### 1. Project Organization
```bash
# Create a projects directory
./run.sh
Choose: 1
Path: ~/projects

# Then create all your projects there
Choose: 8
Select: 8  # Spring Boot
Root: ~/projects/backend-api
```

### 2. Rapid Prototyping
Use templates to quickly scaffold:
- Microservices (Spring Boot + FastAPI)
- Full-stack apps (Next.js or React + any backend)
- Mobile apps (Flutter)
- Desktop apps (Java Swing)

### 3. Learning Projects
Perfect for learning new frameworks:
1. Create project structure
2. Get properly organized codebase
3. Focus on implementation, not setup

### 4. Team Standardization
Everyone uses the same structure:
```bash
# Share your filemanager binary
# Team members create consistent project layouts
# No more "where does this file go?" questions
```

### 5. Combining Operations
```bash
# Create structure
Choose: 8 → Template

# Then use single operations for modifications
Choose: 2  # Add new file
Choose: 1  # Add new folder
Choose: 5  # Set permissions
```

---

## 📋 Template Contents Summary

| Template | Languages | Build Tool | Key Features |
|----------|-----------|------------|--------------|
| Java Traits | Java | Maven | Annotation processing, Multi-module |
| Go Project | Go | Make | Internal packages, Clean architecture |
| Rust Workspace | Rust | Cargo | Multi-crate, Examples, Benchmarks |
| Flask | Python | pip | Blueprints, Templates, Static files |
| FastAPI | Python | pip | Async API, Pydantic, API versioning |
| Java RMI | Java | Maven | Client-server, Remote interfaces |
| Java Desktop | Java | Maven | Swing UI, MVC pattern |
| Spring Boot | Java | Maven | REST API, JPA, H2 database |
| Flutter | Dart | pub | Material Design, Cross-platform |
| React | JavaScript | npm | Components, Hooks, Router ready |
| Next.js | JavaScript | npm | SSR, API routes, App Router |
| HTML/CSS/JS | Web | - | Static site, Responsive design |

---

## 🔧 Customization

### Modifying Templates
Edit `templates.go` to customize any template:

```go
func getMyCustomTemplate() StructureTemplate {
    return StructureTemplate{
        Name:        "my-template",
        Description: "My Custom Project",
        Directories: []string{
            "src",
            "tests",
        },
        Files: map[string]string{
            "README.md": "# My Project\n",
            "src/main.go": "package main\n\nfunc main() {}\n",
        },
    }
}
```

Then add it to `GetAvailableTemplates()`.

---

## 🎯 Use Cases

### Startups & MVPs
Quickly scaffold full-stack applications:
- Backend: Spring Boot or FastAPI
- Frontend: React or Next.js
- Mobile: Flutter

### Education
Students learn proper project structure from day one.

### Hackathons
Spend time coding, not setting up.

### Open Source
Consistent structure across projects.

### Enterprise
Standardized project layouts for teams.

---

## 🚨 Troubleshooting

### "Path cannot be empty"
- Always provide a root directory name
- Use relative or absolute paths

### "Failed to create"
- Check permissions
- Verify parent directories exist (for single operations)
- Use absolute paths if relative fails

### Wrong template selected
- Use **`back`** to return and select again
- No changes made until creation starts

### Cancel during creation
- Use **`cancel`** or **`back`** before entering root path
- Once creation starts, it cannot be cancelled

---

## 📈 Performance

- **Rust core**: Blazing fast file operations
- **Batch creation**: Creates 100+ files in seconds
- **Memory efficient**: No memory leaks, proper cleanup
- **Concurrent safe**: Thread-safe operations

---

## 🤝 Contributing

Want to add more templates? Edit `templates.go`:
1. Create new template function
2. Add to `GetAvailableTemplates()`
3. Rebuild: `make clean && make`

---

**Happy Coding! 🚀**

*Never spend time setting up project structure again!*