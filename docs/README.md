# 📚 FileManager v2 Documentation

Welcome to the FileManager v2 documentation hub. This directory contains comprehensive guides for users, developers, and maintainers.

---

## 📖 Documentation Index

### 🚀 Getting Started

- **[Installation Guide](./INSTALL.md)** - How to install FileManager on all platforms
  - Quick install for Linux, macOS, Windows
  - Detailed installation methods
  - Platform-specific notes
  - Troubleshooting

- **[Quick Start](../README.md#-quick-start)** - Get up and running in 5 minutes
  - Prerequisites
  - Basic installation
  - First run

### 📖 User Guides

- **[Usage Guide](./USAGE_GUIDE.md)** - Complete feature documentation
  - Terminal mode operations
  - Web mode features
  - API endpoints
  - Examples and workflows

- **[Advanced Features](./ADVANCED_FEATURES.md)** - Power user features
  - Multiple file/folder creation
  - Automatic parent directory creation
  - Path validation
  - Interactive builder with navigation
  - Advanced usage patterns

### 🔄 Updates & Upgrades

- **[Update Strategy](./UPDATE_STRATEGY.md)** - How updates work
  - Semantic Versioning (SemVer) explained
  - PATCH, MINOR, MAJOR update types
  - Notification strategies
  - Release notes best practices
  - Security considerations

- **[Automatic Update System](./AUTO_UPDATE_SYSTEM.md)** - How automatic updates work
  - PATCH updates (automatic installation)
  - MINOR updates (user prompt)
  - MAJOR updates (explicit consent)
  - Update triggers and configuration
  - Security considerations
  - Troubleshooting

- **[Migration Guide](./MIGRATION.md)** - Upgrading between major versions
  - v1.x → v2.0 migration
  - Pre-migration checklist
  - Step-by-step migration
  - Configuration migration
  - Troubleshooting
  - Rollback instructions

- **[Release Notes Examples](./RELEASE_NOTES_EXAMPLES.md)** - Writing effective release notes
  - Templates for each update type
  - Real-world examples
  - Writing guidelines
  - Best practices

### 👨‍💻 Development

- **[Development Guide](./DEVELOPMENT.md)** - Building from source
  - Prerequisites
  - Development setup
  - Building the project
  - Running tests
  - Code structure
  - Contributing guidelines

- **[Architecture](./ARCHITECTURE.md)** - System design and components
  - Rust backend (fs-operations-core)
  - Go frontend
  - FFI bridge
  - Web interface
  - API design

### 📦 Packaging & Distribution

- **[Packaging Guide](./PACKAGING.md)** - Creating distributions
  - Building release packages
  - Platform-specific packaging
  - Creating installers
  - Publishing releases
  - Distribution channels

- **[CI/CD Pipeline](../.github/workflows/release.yml)** - Automated builds and releases
  - GitHub Actions workflow
  - Multi-platform builds
  - Automated testing
  - Release automation

### 🔧 Troubleshooting

- **[Troubleshooting Guide](./TROUBLESHOOTING.md)** - Common issues and solutions
  - Installation issues
  - Runtime errors
  - Performance problems
  - Platform-specific issues
  - Getting help

### 📚 Reference

- **[API Documentation](./API.md)** - REST API reference
  - Endpoints
  - Request/response formats
  - Authentication
  - Error handling
  - Examples

- **[Configuration](./CONFIGURATION.md)** - Configuration options
  - Config file format
  - Available settings
  - Environment variables
  - Default values

- **[Changelog](../CHANGELOG.md)** - Version history
  - All releases
  - Breaking changes
  - New features
  - Bug fixes

---

## 🎯 Quick Navigation by Role

### 👤 For Users

1. Start with [Installation Guide](./INSTALL.md)
2. Read [Quick Start](../README.md#-quick-start)
3. Explore [Usage Guide](./USAGE_GUIDE.md)
4. Try [Advanced Features](./ADVANCED_FEATURES.md)
5. Check [Update Strategy](./UPDATE_STRATEGY.md) for updates

### 👨‍💻 For Developers

1. Read [Development Guide](./DEVELOPMENT.md)
2. Understand [Architecture](./ARCHITECTURE.md)
3. Check [API Documentation](./API.md)
4. Review [Code Structure](./DEVELOPMENT.md#project-structure)
5. See [Contributing Guidelines](./DEVELOPMENT.md#contributing)

### 📦 For Maintainers

1. Review [Packaging Guide](./PACKAGING.md)
2. Understand [CI/CD Pipeline](../.github/workflows/release.yml)
3. Learn [Release Process](./PACKAGING.md#release-process)
4. Check [Update Strategy](./UPDATE_STRATEGY.md)
5. Read [Release Notes Examples](./RELEASE_NOTES_EXAMPLES.md)

---

## 📋 Document Overview

### Installation & Setup (3 docs)
- **INSTALL.md** (465 lines) - Complete installation guide
- **DEVELOPMENT.md** - Development environment setup
- **CONFIGURATION.md** - Configuration options

### User Guides (2 docs)
- **USAGE_GUIDE.md** - Feature documentation
- **ADVANCED_FEATURES.md** (459 lines) - Power user features

### Updates & Releases (3 docs)
- **UPDATE_STRATEGY.md** (450+ lines) - Semantic versioning and notifications
- **MIGRATION.md** (400+ lines) - Major version migration
- **RELEASE_NOTES_EXAMPLES.md** (400+ lines) - Release notes templates

### Development (2 docs)
- **DEVELOPMENT.md** - Building and contributing
- **ARCHITECTURE.md** - System design

### Operations (2 docs)
- **PACKAGING.md** - Distribution packaging
- **TROUBLESHOOTING.md** - Problem solving

### Reference (2 docs)
- **API.md** - REST API reference
- **CHANGELOG.md** - Version history

---

## 🔑 Key Concepts

### Semantic Versioning (SemVer)

FileManager uses **MAJOR.MINOR.PATCH** versioning:

```
v2.1.3
│ │ └─ PATCH: Bug fixes, security patches (silent update)
│ └─── MINOR: New features, improvements (subtle notification)
└───── MAJOR: Breaking changes, redesigns (modal notification)
```

See [UPDATE_STRATEGY.md](./UPDATE_STRATEGY.md) for details.

### Hybrid Notification Strategy

Updates are communicated based on their impact:

- **PATCH**: Silent/background installation
- **MINOR**: Subtle in-app banner
- **MAJOR**: Full-screen modal with details

See [UPDATE_STRATEGY.md](./UPDATE_STRATEGY.md) for examples.

### Architecture

FileManager has a clean backend-frontend architecture:

```
┌─────────────────────────────────────┐
│  Terminal UI    │    Web Browser    │
├─────────────────────────────────────┤
│       Go Frontend (file_manager)    │
├─────────────────────────────────────┤
│       CGO FFI Bridge                │
├─────────────────────────────────────┤
│  Rust Backend (rust_ffi)            │
│  - fs-operations-core (library)     │
│  - fs-operations-cli (tool)         │
└─────────────────────────────────────┘
```

See [ARCHITECTURE.md](./ARCHITECTURE.md) for details.

---

## 🚀 Common Tasks

### Installing FileManager
→ See [INSTALL.md](./INSTALL.md)

### Using FileManager
→ See [USAGE_GUIDE.md](./USAGE_GUIDE.md)

### Advanced Operations
→ See [ADVANCED_FEATURES.md](./ADVANCED_FEATURES.md)

### Updating FileManager
→ See [UPDATE_STRATEGY.md](./UPDATE_STRATEGY.md)

### Migrating from v1.x
→ See [MIGRATION.md](./MIGRATION.md)

### Building from Source
→ See [DEVELOPMENT.md](./DEVELOPMENT.md)

### Creating a Release
→ See [PACKAGING.md](./PACKAGING.md)

### Troubleshooting Issues
→ See [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)

### Using the API
→ See [API.md](./API.md)

---

## 📞 Getting Help

### Documentation
- Check the relevant guide above
- Search [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)
- Review [FAQ](./TROUBLESHOOTING.md#faq)

### Community
- [GitHub Issues](https://github.com/DevChigarlicMoses/FileManager/issues) - Report bugs
- [GitHub Discussions](https://github.com/DevChigarlicMoses/FileManager/discussions) - Ask questions
- [GitHub Wiki](https://github.com/DevChigarlicMoses/FileManager/wiki) - Community knowledge

### Direct Support
- **Email**: moses.muranja@strathmore.edu
- **GitHub**: [@DevChigarlicMoses](https://github.com/DevChigarlicMoses)

---

## 📊 Documentation Statistics

| Category | Documents | Lines | Topics |
|----------|-----------|-------|--------|
| Installation | 1 | 465 | 10+ |
| User Guides | 2 | 750+ | 20+ |
| Updates | 3 | 1,250+ | 15+ |
| Development | 2 | 500+ | 10+ |
| Operations | 2 | 400+ | 10+ |
| Reference | 2 | 300+ | 5+ |
| **Total** | **12** | **3,665+** | **70+** |

---

## 🔄 Documentation Maintenance

### Last Updated
- **README.md**: November 23, 2025
- **UPDATE_STRATEGY.md**: November 23, 2025
- **MIGRATION.md**: November 23, 2025
- **RELEASE_NOTES_EXAMPLES.md**: November 23, 2025

### Contributing to Documentation

Found an error or want to improve the docs?

1. Fork the repository
2. Edit the relevant document
3. Submit a pull request
4. We'll review and merge

See [DEVELOPMENT.md](./DEVELOPMENT.md#contributing) for details.

---

## 📚 Related Resources

- **Main README**: [../README.md](../README.md)
- **GitHub Repository**: [github.com/DevChigarlicMoses/FileManager](https://github.com/DevChigarlicMoses/FileManager)
- **Releases**: [GitHub Releases](https://github.com/DevChigarlicMoses/FileManager/releases)
- **Issues**: [GitHub Issues](https://github.com/DevChigarlicMoses/FileManager/issues)

---

**FileManager Documentation - Your complete guide to file management** 📚

Last updated: November 23, 2025
