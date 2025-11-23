# Dynamic Version Management System

## Overview

FileManager v1 now has a **fully dynamic version management system** that automatically updates all version references across the entire project with a single command. No more manual version updates!

## What's New

### ✅ Centralized Version Control
- Single `VERSION` file as the source of truth
- All scripts read from this file automatically
- No hardcoded versions in build scripts

### ✅ Automatic Updates
- Update version once, everywhere updates automatically
- Supports semantic versioning (MAJOR.MINOR.PATCH)
- Validates version format

### ✅ Smart Version Bumping
- Bump patch (2.0.0 → 2.0.1)
- Bump minor (2.0.0 → 2.1.0)
- Bump major (2.0.0 → 3.0.0)

## Files Created

### 1. `VERSION` (Root Directory)
Single source of truth for the project version:
```
2.0.0
```

### 2. `scripts/version-manager.sh`
Comprehensive version management tool with commands:
- `get` - Show current version
- `set X.Y.Z` - Set specific version
- `bump-patch` - Increment patch version
- `bump-minor` - Increment minor version
- `bump-major` - Increment major version
- `list` - Show all version references
- `validate` - Check version format

### 3. `VERSION_MANAGEMENT.md`
Complete guide with examples and best practices

## Quick Start

### Check Current Version
```bash
cd scripts
./version-manager.sh get
# Output: 2.0.0
```

### Update to New Version
```bash
# Set specific version
./version-manager.sh set 2.1.0

# Or bump automatically
./version-manager.sh bump-minor  # 2.0.0 → 2.1.0
./version-manager.sh bump-patch  # 2.0.0 → 2.0.1
./version-manager.sh bump-major  # 2.0.0 → 3.0.0
```

### Build with New Version
```bash
# Automatically uses version from VERSION file
./build-all-platforms.sh all

# Creates packages with correct version:
# - filemanager_2.1.0_amd64.deb
# - filemanager-2.1.0-linux-amd64.tar.gz
# - filemanager-2.1.0-windows-amd64.zip
```

## How It Works

### Version File Hierarchy

```
VERSION (root)
    ↓
    ├─→ build-linux-amd64.sh
    ├─→ build-linux-arm64.sh
    ├─→ build-windows-amd64.sh
    ├─→ build-harmonyos.sh
    ├─→ build-all-platforms.sh
    │
    └─→ version-manager.sh updates:
        ├─→ snap/snapcraft.yaml
        ├─→ file_manager/pkg/version/version.go
        ├─→ rust_ffi/Cargo.toml
        ├─→ rust_ffi/crates/*/Cargo.toml
        └─→ README.md
```

### Build Script Integration

All build scripts now read version dynamically:

```bash
# Read version from VERSION file
if [ -f "$PROJECT_ROOT/VERSION" ]; then
    VERSION=$(cat "$PROJECT_ROOT/VERSION" | tr -d ' \n')
else
    VERSION="2.0.0"
fi

echo "🔨 Building FileManager v${VERSION}..."
```

No need to edit build scripts for version changes!

## Real-World Examples

### Scenario 1: Bug Fix Release

```bash
# Current version: 2.0.0
cd scripts

# Bump patch version
./version-manager.sh bump-patch
# Version is now: 2.0.1

# Build all platforms
./build-all-platforms.sh all

# Packages created:
# - filemanager_2.0.1_amd64.deb
# - filemanager-2.0.1-linux-amd64.tar.gz
# - filemanager-2.0.1-windows-amd64.zip

# Create git tag
git tag v2.0.1
git push origin v2.0.1
```

### Scenario 2: Feature Release

```bash
# Current version: 2.0.0
cd scripts

# Bump minor version
./version-manager.sh bump-minor
# Version is now: 2.1.0

# Build all platforms
./build-all-platforms.sh all

# Packages created:
# - filemanager_2.1.0_amd64.deb
# - filemanager-2.1.0-linux-amd64.tar.gz
# - filemanager-2.1.0-windows-amd64.zip

# Create git tag
git tag v2.1.0
git push origin v2.1.0
```

### Scenario 3: Major Release

```bash
# Current version: 2.5.3
cd scripts

# Bump major version
./version-manager.sh bump-major
# Version is now: 3.0.0

# Build all platforms
./build-all-platforms.sh all

# Packages created:
# - filemanager_3.0.0_amd64.deb
# - filemanager-3.0.0-linux-amd64.tar.gz
# - filemanager-3.0.0-windows-amd64.zip

# Create git tag
git tag v3.0.0
git push origin v3.0.0
```

## Files Updated by Version Manager

When you run `./version-manager.sh set 2.1.0`, these files are automatically updated:

### 1. VERSION
```
2.1.0
```

### 2. snap/snapcraft.yaml
```yaml
version: '2.1.0'
```

### 3. file_manager/pkg/version/version.go
```go
const Version = "2.1.0"
```

### 4. rust_ffi/Cargo.toml
```toml
version = "2.1.0"
```

### 5. rust_ffi/crates/core/Cargo.toml
```toml
version = "2.1.0"
```

### 6. rust_ffi/crates/cli/Cargo.toml
```toml
version = "2.1.0"
```

### 7. README.md
All version references updated:
```markdown
FileManager v2.1.0
```

## Version Manager Commands

### Get Current Version
```bash
./version-manager.sh get
```
Returns: `2.0.0`

### Set Specific Version
```bash
./version-manager.sh set 2.1.0
```
Updates all files to version 2.1.0

### Bump Patch Version
```bash
./version-manager.sh bump-patch
```
2.0.0 → 2.0.1 (bug fixes)

### Bump Minor Version
```bash
./version-manager.sh bump-minor
```
2.0.0 → 2.1.0 (new features)

### Bump Major Version
```bash
./version-manager.sh bump-major
```
2.0.0 → 3.0.0 (breaking changes)

### List All Versions
```bash
./version-manager.sh list
```
Shows all version references in the project

### Validate Version Format
```bash
./version-manager.sh validate 2.1.0
```
Checks if version follows MAJOR.MINOR.PATCH format

## Build System Integration

### Before (Hardcoded Versions)
```bash
# Old way - had to edit every script
VERSION="2.0.0"  # Had to change this manually
```

### After (Dynamic Versions)
```bash
# New way - reads from VERSION file automatically
VERSION=$(cat "$PROJECT_ROOT/VERSION" | tr -d ' \n')
```

### Build Command
```bash
# Automatically uses correct version
./build-all-platforms.sh all

# Output shows correct version:
# 🔨 Building FileManager v2.0.0 - Multi-Platform Build
```

## Semantic Versioning Guide

### When to Bump Patch (Z)
```
2.0.0 → 2.0.1
```
- Bug fixes
- Security patches
- Minor improvements
- No breaking changes

### When to Bump Minor (Y)
```
2.0.0 → 2.1.0
```
- New features
- Improvements
- Backwards compatible
- No breaking changes

### When to Bump Major (X)
```
2.0.0 → 3.0.0
```
- Breaking changes
- Major rewrites
- Significant architectural changes
- Requires user migration

## Workflow Integration

### Git Workflow
```bash
# 1. Bump version
./version-manager.sh bump-minor

# 2. Commit changes
git add -A
git commit -m "Bump version to 2.1.0"

# 3. Build
./build-all-platforms.sh all

# 4. Tag release
git tag v2.1.0
git push origin v2.1.0
```

### CI/CD Integration
```yaml
# GitHub Actions example
- name: Bump version
  run: |
    cd scripts
    ./version-manager.sh bump-patch

- name: Build
  run: |
    cd scripts
    ./build-all-platforms.sh all
```

## Troubleshooting

### Version not updating in a file
Check if the file exists and has the expected format:
```bash
# Check snapcraft.yaml
grep "version:" ../snap/snapcraft.yaml

# Check Go version
grep "Version = " ../file_manager/pkg/version/version.go

# Check Cargo.toml
grep "version = " ../rust_ffi/Cargo.toml
```

### Invalid version format
Valid formats: `2.0.0`, `1.0.0`, `0.1.0`
Invalid formats: `2.0`, `v2.0.0`, `2.0.0-beta`

### Build fails after version update
Ensure all dependencies support the new version:
```bash
# Check Cargo.toml dependencies
grep "^version" ../rust_ffi/Cargo.toml
```

## Benefits

✅ **Single Source of Truth** - One VERSION file for entire project
✅ **Automatic Updates** - All files updated with one command
✅ **No Manual Edits** - Eliminates human error
✅ **Semantic Versioning** - Follows industry standards
✅ **Easy Bumping** - Automatic patch/minor/major increments
✅ **Build Integration** - Builds use correct version automatically
✅ **Validation** - Checks version format
✅ **Listing** - See all version references

## Summary

| Task | Command | Result |
|------|---------|--------|
| Check version | `./version-manager.sh get` | Shows current version |
| Set version | `./version-manager.sh set 2.1.0` | Updates all files |
| Bug fix | `./version-manager.sh bump-patch` | 2.0.0 → 2.0.1 |
| New feature | `./version-manager.sh bump-minor` | 2.0.0 → 2.1.0 |
| Major release | `./version-manager.sh bump-major` | 2.0.0 → 3.0.0 |
| List versions | `./version-manager.sh list` | Shows all references |
| Validate | `./version-manager.sh validate 2.1.0` | Checks format |

## Next Steps

1. Use `./version-manager.sh get` to check current version
2. Use `./version-manager.sh bump-*` to increment version
3. Use `./build-all-platforms.sh all` to build with new version
4. Tag release with `git tag v2.1.0`

For detailed information, see [VERSION_MANAGEMENT.md](VERSION_MANAGEMENT.md)
