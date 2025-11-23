# 🔄 Migration Guides for Major Updates

This document provides migration guides for upgrading between major versions of FileManager.

---

## 📋 Table of Contents

- [v1.x → v2.0 Migration](#v1x--v20-migration)
- [General Migration Best Practices](#general-migration-best-practices)
- [Troubleshooting Migration Issues](#troubleshooting-migration-issues)

---

## v1.x → v2.0 Migration

### Overview

**v2.0.0** is a major upgrade introducing:
- ✨ Modern web interface alongside terminal mode
- 🎨 Complete UI redesign
- 📊 New database schema
- 🔐 Enhanced security features
- ⚡ Performance improvements

### What's Changing

| Component | v1.x | v2.0 | Impact |
|-----------|------|------|--------|
| **UI** | Terminal only | Terminal + Web | New interface available |
| **Database** | SQLite v3.x | SQLite v3.4+ | Auto-migrated |
| **Config** | `.filemanager.conf` | `.config/filemanager/` | Restructured |
| **Permissions** | Unix only | Unix + Windows | Enhanced |
| **API** | N/A | REST API | New capability |

### Pre-Migration Checklist

- [ ] Backup your configuration files
- [ ] Note any custom settings
- [ ] Close all FileManager instances
- [ ] Ensure 100MB free disk space
- [ ] Check system requirements (see below)

### System Requirements for v2.0

**Linux:**
- Ubuntu 20.04+ / Debian 11+ / Fedora 35+
- 512 MB RAM (1 GB recommended)
- 100 MB disk space

**macOS:**
- macOS 10.15+ (Catalina or later)
- Intel or Apple Silicon (M1/M2)
- 512 MB RAM (1 GB recommended)

**Windows:**
- Windows 10 version 1909+ or Windows 11
- 512 MB RAM (1 GB recommended)
- 100 MB disk space

### Migration Steps

#### Step 1: Backup Current Configuration

```bash
# Linux/macOS
cp -r ~/.filemanager ~/.filemanager.backup
cp -r ~/.config/filemanager ~/.config/filemanager.backup

# Windows
xcopy "%APPDATA%\FileManager" "%APPDATA%\FileManager.backup" /E /I
```

#### Step 2: Uninstall v1.x

```bash
# Linux/macOS
sudo make uninstall

# Or manually:
sudo rm /usr/local/bin/filemanager
sudo rm /usr/local/lib/libfilemanager.so

# Windows
# Use Control Panel → Programs → Uninstall a program
# Or run uninstall.bat as Administrator
```

#### Step 3: Install v2.0

```bash
# Download v2.0.0
curl -L https://github.com/DevChigarlicMoses/FileManager/releases/download/v2.0.0/filemanager-2.0.0-linux-amd64.tar.gz | tar xz

# Install
cd linux-amd64
sudo ./install.sh
```

#### Step 4: Automatic Migration

When you first run v2.0:

```bash
filemanager
```

The app will:
1. ✅ Detect v1.x configuration
2. ✅ Create backup at `~/.filemanager.v1.backup`
3. ✅ Migrate settings automatically
4. ✅ Create new config structure
5. ✅ Display migration summary

**Expected Output:**
```
🔄 Migrating from FileManager v1.x...

📋 Migration Steps:
  ✅ Detected v1.x configuration
  ✅ Creating backup: ~/.filemanager.v1.backup
  ✅ Migrating settings...
  ✅ Creating new config structure
  ✅ Updating database schema

✨ Migration Complete!
  • Old config backed up: ~/.filemanager.v1.backup
  • New config location: ~/.config/filemanager/
  • Database updated: ~/.config/filemanager/filemanager.db

📝 Next Steps:
  1. Review new settings: filemanager --settings
  2. Try web mode: filemanager --web
  3. Check release notes: filemanager --changelog
```

#### Step 5: Verify Installation

```bash
# Check version
filemanager --version
# Output: FileManager v2.0.0

# Test terminal mode
filemanager
# Should show new menu

# Test web mode
filemanager --web
# Should open browser to http://localhost:8080
```

### Configuration Migration

#### Old Config Location (v1.x)
```
~/.filemanager/
├── config.json
├── history.json
└── templates.json
```

#### New Config Location (v2.0)
```
~/.config/filemanager/
├── config.json
├── history.json
├── templates.json
├── web-config.json
└── filemanager.db
```

#### Migrated Settings

| Setting | v1.x | v2.0 | Notes |
|---------|------|------|-------|
| Default path | `~` | `~` | Preserved |
| History | ✅ | ✅ | Automatically migrated |
| Templates | ✅ | ✅ | Automatically migrated |
| Permissions | ✅ | ✅ Enhanced | New permission model |
| Web port | N/A | 8080 | New in v2.0 |

### New Features in v2.0

After migration, you can use:

1. **Web Interface**
   ```bash
   filemanager --web
   # Opens modern web UI
   ```

2. **REST API**
   ```bash
   curl http://localhost:8080/api/health
   ```

3. **Advanced Operations**
   - Batch file operations
   - File search
   - Archive support
   - Cloud integration (coming soon)

### Troubleshooting Migration

#### Issue: "Migration Failed"

```bash
# Check error logs
filemanager --debug

# Restore from backup
cp -r ~/.filemanager.v1.backup ~/.filemanager
```

#### Issue: "Config Not Found"

```bash
# Manually migrate config
filemanager --migrate-config

# Or reset to defaults
filemanager --reset-config
```

#### Issue: "Database Locked"

```bash
# Close all FileManager instances
pkill filemanager

# Remove lock file
rm ~/.config/filemanager/.lock

# Try again
filemanager
```

#### Issue: "Permission Denied"

```bash
# Fix permissions
chmod -R 755 ~/.config/filemanager/
chmod -R 644 ~/.config/filemanager/*.*
```

### Rollback to v1.x (If Needed)

If you need to revert to v1.x:

```bash
# Uninstall v2.0
sudo make uninstall

# Restore backup
cp -r ~/.filemanager.v1.backup ~/.filemanager

# Reinstall v1.x
# Download and install v1.x release
```

---

## General Migration Best Practices

### Before Any Major Update

1. **Backup Everything**
   ```bash
   tar -czf filemanager-backup-$(date +%Y%m%d).tar.gz ~/.config/filemanager/
   ```

2. **Read Release Notes**
   - Understand breaking changes
   - Check system requirements
   - Review migration guide

3. **Test in Safe Environment**
   - Use test account if possible
   - Don't update production immediately
   - Wait for patch releases (v2.0.1, v2.0.2)

4. **Document Custom Configurations**
   - Export settings
   - Note any customizations
   - Screenshot important configs

### During Migration

1. **Follow Official Guide**
   - Use provided migration steps
   - Don't skip steps
   - Keep backups until verified

2. **Monitor Process**
   - Watch for error messages
   - Check logs for warnings
   - Verify each step completes

3. **Verify Functionality**
   - Test core features
   - Check custom configurations
   - Verify data integrity

### After Migration

1. **Validate Data**
   ```bash
   filemanager --verify-database
   filemanager --check-config
   ```

2. **Test All Features**
   - Terminal mode
   - Web mode (if applicable)
   - Custom operations
   - Integrations

3. **Monitor for Issues**
   - Check logs regularly
   - Watch for crashes
   - Report any problems

4. **Clean Up**
   - Remove old backups (after 1 month)
   - Delete temporary files
   - Archive old configs

---

## Troubleshooting Migration Issues

### Common Issues and Solutions

#### 1. "Old version still running"

**Problem:** Can't install new version because old version is active

**Solution:**
```bash
# Kill all instances
pkill -9 filemanager

# Check if killed
ps aux | grep filemanager

# Then proceed with installation
```

#### 2. "Permission denied" during migration

**Problem:** Migration fails due to file permissions

**Solution:**
```bash
# Fix permissions
sudo chown -R $USER:$USER ~/.config/filemanager/
sudo chmod -R 755 ~/.config/filemanager/

# Retry migration
filemanager
```

#### 3. "Database corruption"

**Problem:** Old database can't be migrated

**Solution:**
```bash
# Backup corrupted database
mv ~/.config/filemanager/filemanager.db ~/.config/filemanager/filemanager.db.corrupted

# Let app create new database
filemanager

# Restore from backup if needed
cp ~/.filemanager.v1.backup/history.json ~/.config/filemanager/
```

#### 4. "Config file not found"

**Problem:** Migration can't find old configuration

**Solution:**
```bash
# Check if backup exists
ls -la ~/.filemanager.v1.backup/

# Manually restore
cp -r ~/.filemanager.v1.backup/* ~/.config/filemanager/

# Run migration again
filemanager --migrate-config
```

#### 5. "Out of disk space"

**Problem:** Not enough space for migration

**Solution:**
```bash
# Check available space
df -h

# Free up space
# Remove old backups, temporary files, etc.

# Retry migration
filemanager
```

### Getting Help

If migration fails:

1. **Check Logs**
   ```bash
   filemanager --debug > migration.log 2>&1
   cat migration.log
   ```

2. **Report Issue**
   - Include error message
   - Attach migration.log
   - Describe steps taken
   - Provide system info

3. **Contact Support**
   - GitHub Issues: [Report Issue](https://github.com/DevChigarlicMoses/FileManager/issues)
   - Email: moses.muranja@strathmore.edu
   - Wiki: [Migration Help](https://github.com/DevChigarlicMoses/FileManager/wiki/Migration)

---

## Version-Specific Guides

### v0.x → v1.0

Coming soon - will be added when v1.0 is released.

### v1.0 → v1.5

No migration needed - fully backwards compatible.

### v1.5 → v2.0

See [v1.x → v2.0 Migration](#v1x--v20-migration) above.

---

**Successful migration ensures smooth transition to new features while preserving your data!** 🎉
