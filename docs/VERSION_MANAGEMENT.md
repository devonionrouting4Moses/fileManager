# FileManager Version Management Guide

## Overview

FileManager uses a centralized version management system that automatically updates all version references across the entire project. This ensures consistency and eliminates manual version updates.

## Version Format

FileManager follows **Semantic Versioning (SemVer)**:

```
MAJOR.MINOR.PATCH
  ↓      ↓      ↓
  3      2      1
```

- **MAJOR**: Breaking changes (e.g., 2.0.0 → 3.0.0)
- **MINOR**: New features, backwards compatible (e.g., 2.0.0 → 2.1.0)
- **PATCH**: Bug fixes, security patches (e.g., 2.0.0 → 2.0.1)

## Quick Start

### Check Current Version

```bash
cd file_manager_v1/scripts
./version-manager.sh get
```

Output: `2.0.0`

### Update to Specific Version

```bash
./version-manager.sh set 2.1.0
```

This automatically updates:
- `VERSION` file
- `snap/snapcraft.yaml`
- `file_manager/pkg/version/version.go`
- All `Cargo.toml` files
- `README.md` version references

### Bump Versions Automatically

```bash
# Bump patch version (2.0.0 → 2.0.1)
./version-manager.sh bump-patch

# Bump minor version (2.0.0 → 2.1.0)
./version-manager.sh bump-minor

# Bump major version (2.0.0 → 3.0.0)
./version-manager.sh bump-major
```

## Detailed Usage

### Get Current Version

```bash
./version-manager.sh get
```

Returns the current version from the `VERSION` file.

### Set New Version

```bash
./version-manager.sh set 2.1.0
```

**What it updates:**
1. `VERSION` - Core version file
2. `snap/snapcraft.yaml` - Snap package metadata
3. `file_manager/pkg/version/version.go` - Go version constant
4. `rust_ffi/Cargo.toml` - Rust library version
5. `rust_ffi/crates/*/Cargo.toml` - Rust crate versions
6. `README.md` - All version references

### Bump Patch Version

```bash
./version-manager.sh bump-patch
```

Increments the patch version:
- 2.0.0 → 2.0.1
- 2.0.5 → 2.0.6
- 1.9.9 → 1.9.10

Use for: Bug fixes, security patches, minor improvements

### Bump Minor Version

```bash
./version-manager.sh bump-minor
```

Increments the minor version and resets patch to 0:
- 2.0.0 → 2.1.0
- 2.5.3 → 2.6.0
- 1.9.9 → 1.10.0

Use for: New features, improvements, backwards-compatible changes

### Bump Major Version

```bash
./version-manager.sh bump-major
```

Increments the major version and resets minor and patch to 0:
- 2.0.0 → 3.0.0
- 1.9.9 → 2.0.0
- 0.1.0 → 1.0.0

Use for: Breaking changes, major rewrites, significant architectural changes

### List All Version References

```bash
./version-manager.sh list
```

Shows all version files in the project and their current values:

```
Version files in the project:

Core Version:
  VERSION: 2.0.0

snapcraft.yaml:
  version: '2.0.0'

Go version.go:
  Version = "2.0.0"

Cargo.toml files:
  version = "2.0.0"
```

### Validate Version Format

```bash
./version-manager.sh validate 2.1.0
```

Checks if the version follows the MAJOR.MINOR.PATCH format:
- ✅ Valid: `2.1.0`, `1.0.0`, `0.1.0`
- ❌ Invalid: `2.1`, `v2.1.0`, `2.1.0-beta`

## Workflow Examples

### Releasing a Bug Fix

```bash
# Current version: 2.0.0
./version-manager.sh bump-patch
# New version: 2.0.1

# Build and test
./build-all-platforms.sh all

# Create release
git tag v2.0.1
git push origin v2.0.1
```

### Releasing New Features

```bash
# Current version: 2.0.0
./version-manager.sh bump-minor
# New version: 2.1.0

# Build and test
./build-all-platforms.sh all

# Create release
git tag v2.1.0
git push origin v2.1.0
```

### Major Release

```bash
# Current version: 2.5.3
./version-manager.sh bump-major
# New version: 3.0.0

# Build and test
./build-all-platforms.sh all

# Create release
git tag v3.0.0
git push origin v3.0.0
```

### Custom Version

```bash
# Set to specific version
./version-manager.sh set 2.1.0-rc1
# ⚠️  Note: This will fail validation (pre-release format)

# Use valid format
./version-manager.sh set 2.1.0
```

## Files Updated by Version Manager

### 1. VERSION (Root)
```
2.0.0
```
The single source of truth for the project version.

### 2. snap/snapcraft.yaml
```yaml
version: '2.0.0'
```
Snap package metadata.

### 3. file_manager/pkg/version/version.go
```go
const Version = "2.0.0"
```
Go application version constant.

### 4. rust_ffi/Cargo.toml
```toml
version = "2.0.0"
```
Rust workspace version.

### 5. rust_ffi/crates/*/Cargo.toml
```toml
version = "2.0.0"
```
Individual Rust crate versions.

### 6. README.md
All references to version numbers:
```markdown
FileManager v2.0.0
```

## Build System Integration

All build scripts automatically read the version from the `VERSION` file:

```bash
# Reads VERSION file automatically
./build-all-platforms.sh all

# Creates packages with correct version
# - filemanager_2.0.0_amd64.deb
# - filemanager-2.0.0-linux-amd64.tar.gz
# - filemanager-2.0.0-windows-amd64.zip
```

No need to update version in build scripts!

## Troubleshooting

### Version file not found

```bash
# Create VERSION file with default
echo "2.0.0" > ../VERSION
```

### Invalid version format

```bash
# ❌ Invalid
./version-manager.sh set 2.1        # Missing patch
./version-manager.sh set v2.1.0     # Has 'v' prefix
./version-manager.sh set 2.1.0-beta # Pre-release format

# ✅ Valid
./version-manager.sh set 2.1.0
./version-manager.sh set 1.0.0
./version-manager.sh set 0.1.0
```

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

## Best Practices

### 1. Always Use Version Manager
Don't manually edit version numbers. Always use the version manager:

```bash
# ✅ Good
./version-manager.sh bump-minor

# ❌ Bad
# Manually editing VERSION file
# Manually editing snapcraft.yaml
# Manually editing Cargo.toml
```

### 2. Commit Version Changes
After updating version, commit the changes:

```bash
./version-manager.sh bump-minor
git add -A
git commit -m "Bump version to 2.1.0"
```

### 3. Tag Releases
Always tag releases with the version:

```bash
git tag v2.1.0
git push origin v2.1.0
```

### 4. Update Changelog
Update CHANGELOG.md with release notes:

```markdown
## [2.1.0] - 2024-11-23

### Added
- New feature X
- New feature Y

### Fixed
- Bug fix A
- Bug fix B
```

### 5. Test Before Release
Always build and test before releasing:

```bash
./version-manager.sh bump-minor
./build-all-platforms.sh all
# Test packages...
git tag v2.1.0
```

## Advanced Usage

### Scripting Version Updates

```bash
#!/bin/bash
# Update version and build

cd scripts
VERSION=$(./version-manager.sh get)
echo "Current version: $VERSION"

# Bump minor version
./version-manager.sh bump-minor
NEW_VERSION=$(./version-manager.sh get)
echo "New version: $NEW_VERSION"

# Build all platforms
./build-all-platforms.sh all

# Create git tag
git tag v$NEW_VERSION
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
    
- name: Create release
  run: |
    VERSION=$(cat VERSION | tr -d ' \n')
    git tag v$VERSION
    git push origin v$VERSION
```

## Summary

| Command | Effect | Use Case |
|---------|--------|----------|
| `get` | Show current version | Check version |
| `set X.Y.Z` | Set specific version | Custom version |
| `bump-patch` | 2.0.0 → 2.0.1 | Bug fixes |
| `bump-minor` | 2.0.0 → 2.1.0 | New features |
| `bump-major` | 2.0.0 → 3.0.0 | Breaking changes |
| `list` | Show all versions | Verify updates |
| `validate` | Check format | Validate version |

## Questions?

For more information, see:
- [CROSS_PLATFORM_BUILD.md](CROSS_PLATFORM_BUILD.md) - Build instructions
- [README.md](README.md) - Project overview
- [VERSION](VERSION) - Current version file
