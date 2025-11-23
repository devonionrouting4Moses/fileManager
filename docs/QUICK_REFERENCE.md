# FileManager v1 - Quick Reference Card

## 🚀 Most Common Commands

### Check Version
```bash
cd scripts
./version-manager.sh get
```

### Update Version
```bash
# Bump minor (2.0.0 → 2.1.0)
./version-manager.sh bump-minor

# Or set specific version
./version-manager.sh set 2.1.0
```

### Build Everything
```bash
./build-all-platforms.sh all
```

### Build Specific Platform
```bash
./build-all-platforms.sh linux-amd64      # Linux 64-bit
./build-all-platforms.sh linux-arm64      # Linux ARM
./build-all-platforms.sh windows-amd64    # Windows
```

## 📋 Version Manager Cheat Sheet

```bash
./version-manager.sh get              # Show current version
./version-manager.sh set 2.1.0        # Set to 2.1.0
./version-manager.sh bump-patch       # 2.0.0 → 2.0.1
./version-manager.sh bump-minor       # 2.0.0 → 2.1.0
./version-manager.sh bump-major       # 2.0.0 → 3.0.0
./version-manager.sh list             # Show all versions
./version-manager.sh validate 2.1.0   # Check format
./version-manager.sh help             # Show help
```

## 🏗️ Build Cheat Sheet

```bash
./build-all-platforms.sh all                    # All platforms
./build-all-platforms.sh linux-amd64            # Linux 64-bit
./build-all-platforms.sh linux-arm64            # Linux ARM
./build-all-platforms.sh windows-amd64          # Windows
./build-all-platforms.sh linux-amd64 windows    # Multiple
```

## 📦 Output Files

### Linux amd64
```
filemanager_2.0.0_amd64.deb
filemanager-2.0.0-linux-amd64.tar.gz
filemanager-apt_2.0.0_amd64.deb
arch/PKGBUILD
```

### Linux arm64
```
filemanager_2.0.0_arm64.deb
filemanager-2.0.0-linux-arm64.tar.gz
filemanager-apt_2.0.0_arm64.deb
arch-arm64/PKGBUILD
```

### Windows
```
filemanager.exe
filemanager-2.0.0-windows-amd64.zip
dist/filemanager-2.0.0-windows-amd64/install.bat
```

## 🔧 Setup Requirements

### Linux amd64
✅ Already installed

### Linux arm64
```bash
sudo apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu binutils-aarch64-linux-gnu
```

### Windows
```bash
sudo apt-get install -y mingw-w64
```

## 📖 Documentation

| File | Purpose |
|------|---------|
| `VERSION_MANAGEMENT.md` | Complete version guide |
| `DYNAMIC_VERSION_SYSTEM.md` | Version system overview |
| `ARM64_BUILD_SETUP.md` | ARM64 setup guide |
| `CROSS_PLATFORM_BUILD.md` | Build guide |
| `BUILD_SYSTEM_SUMMARY.md` | System summary |
| `QUICK_REFERENCE.md` | This file |

## 🎯 Common Workflows

### Release Bug Fix
```bash
cd scripts
./version-manager.sh bump-patch
./build-all-platforms.sh all
git tag v2.0.1
git push origin v2.0.1
```

### Release New Features
```bash
cd scripts
./version-manager.sh bump-minor
./build-all-platforms.sh all
git tag v2.1.0
git push origin v2.1.0
```

### Release Major Version
```bash
cd scripts
./version-manager.sh bump-major
./build-all-platforms.sh all
git tag v3.0.0
git push origin v3.0.0
```

## 🐛 Troubleshooting

### ARM64 build fails
```bash
sudo apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu binutils-aarch64-linux-gnu
```

### Windows build fails
```bash
sudo apt-get install -y mingw-w64
```

### Version not updating
```bash
cat VERSION                    # Check VERSION file
./version-manager.sh list      # Check all versions
```

## 📊 Version Format

```
MAJOR.MINOR.PATCH
  2    .  0    .  0

Breaking changes → Major (2.0.0 → 3.0.0)
New features → Minor (2.0.0 → 2.1.0)
Bug fixes → Patch (2.0.0 → 2.0.1)
```

## ✨ Key Features

✅ Single VERSION file
✅ Automatic updates everywhere
✅ No hardcoded versions
✅ Semantic versioning
✅ Cross-platform builds
✅ Multiple package formats
✅ Error handling
✅ Comprehensive docs

## 🎓 Pro Tips

1. **Always use version-manager.sh** - Never edit VERSION manually
2. **Build after version update** - Ensures correct version in packages
3. **Tag releases** - Use `git tag v2.1.0` for releases
4. **Check version** - Use `./version-manager.sh get` to verify
5. **Read docs** - See `VERSION_MANAGEMENT.md` for details

## 🚀 One-Command Release

```bash
cd scripts && \
./version-manager.sh bump-minor && \
./build-all-platforms.sh all && \
cd .. && \
git add -A && \
git commit -m "Release v2.1.0" && \
git tag v2.1.0 && \
git push origin main v2.1.0
```

## 📞 Need Help?

- Version management: `./version-manager.sh help`
- Build help: `./build-all-platforms.sh`
- Read docs: See `*.md` files in project root
- Check logs: Build output shows what's happening

---

**FileManager v1 - Production Ready Build System** 🎉
