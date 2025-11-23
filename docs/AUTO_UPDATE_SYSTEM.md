# 🔄 Automatic Update System

FileManager v2 includes an intelligent automatic update system that respects user preferences while maintaining security and stability.

---

## 📋 Overview

The automatic update system uses **Semantic Versioning (SemVer)** to determine how to handle updates:

- **PATCH Updates** (v1.0.Z) → Automatic installation with countdown
- **MINOR Updates** (v1.Y.0) → User prompt (optional installation)
- **MAJOR Updates** (vX.0.0) → Explicit user consent required

---

## 🔧 PATCH Updates (Automatic)

### What Happens

When a PATCH update is available:

1. ✅ Notification is displayed
2. ✅ 5-second countdown starts
3. ✅ Installation proceeds automatically
4. ✅ User can press Ctrl+C to cancel
5. ✅ App continues after installation

### User Experience

```
══════════════════════════════════════════════════════════════════
║                  🔧 PATCH UPDATE DETECTED                      ║
║                                                                  ║
║  Current Version: v1.0.8
║  Available Version: v1.0.9
║                                                                  ║
║  ✅ This is a safe, backwards-compatible security/bug fix.
║  🔒 Installing automatically to keep your system secure...
║                                                                  ║
║  📝 What's Fixed:
║    • Fixed app crash when saving configuration
║    • Fixed memory leak in file operations
║    • Improved copy operation speed by 15%
║                                                                  ║
║  ⏱️  Installing in 5 seconds... (Press Ctrl+C to cancel)
║                                                                  ║
══════════════════════════════════════════════════════════════════

  ⏳ 5 seconds remaining...
  ⏳ 4 seconds remaining...
  ⏳ 3 seconds remaining...
  ⏳ 2 seconds remaining...
  ⏳ 1 second remaining...
  ✅ Installing update...

📦 Installation Details:
  ├─ Downloading files... ✅
  ├─ Verifying checksums... ✅
  ├─ Extracting files... ✅
  ├─ Installing binary... ✅
  ├─ Installing libraries... ✅
  └─ Finalizing installation... ✅

══════════════════════════════════════════════════════════════════
║              ✅ UPDATE INSTALLED SUCCESSFULLY                   ║
║                                                                  ║
║  Updated: v1.0.8 → v1.0.9
║  Status: Ready to use
║                                                                  ║
══════════════════════════════════════════════════════════════════

╔════════════════════════════════════════════════════════════════════╗
║                      UPDATE SUMMARY                               ║
╠════════════════════════════════════════════════════════════════════╣
║                                                                    ║
║  Current Version:      v1.0.8
║  Available Version:    v1.0.9
║  Change Type:         🔧 PATCH
║  User Impact:         Minimal
║  Update Strategy:     Silent/Direct Install
║  Published:           2025-11-23 14:30:45
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝
```

### Progressive Installation Process

The installation happens in real-time with each step executing sequentially:

1. **Countdown Phase** (5 seconds)
   - Each second is displayed as it passes
   - User can press Ctrl+C to cancel
   - Shows: "⏳ 5 seconds remaining..." → "⏳ 1 second remaining..."

2. **Installation Phase** (6 steps)
   - Each step executes in real-time
   - Progress dots appear as the step executes
   - Checkmark appears when step completes
   - Shows:
     ```
     ├─ Downloading files... ✅
     ├─ Verifying checksums... ✅
     ├─ Extracting files... ✅
     ├─ Installing binary... ✅
     ├─ Installing libraries... ✅
     └─ Finalizing installation... ✅
     ```

3. **Completion Phase**
   - Success message displayed
   - Update summary shown
   - App ready to use

### Key Features

- ✅ Non-disruptive
- ✅ Automatic installation
- ✅ Cancellable with Ctrl+C
- ✅ Progressive real-time display
- ✅ Each step executes sequentially
- ✅ Detailed progress feedback
- ✅ Security-focused

---

## ✨ MINOR Updates (User Prompt)

### What Happens

When a MINOR update is available:

1. 📢 Notification is displayed
2. ❓ User is prompted: "Would you like to install this update now? (y/n)"
3. 🔄 If yes → Installation proceeds
4. ⏭️ If no → User can update later

### User Experience

```
════════════════════════════════════════════════════════════════════
║              ✨ NEW FEATURES AVAILABLE                           ║
║                                                                  ║
║  Current Version: v1.2.5
║  Available Version: v1.3.0
║                                                                  ║
║  📊 Update Type: MINOR (New Features & Improvements)
║  📈 User Impact: Low to Moderate
║  🔄 Update Strategy: Backwards-Compatible
║                                                                  ║
║  📝 What's New:
║    • Added dark mode toggle
║    • Added batch file operations
║    • Added file search functionality
║    • Improved data synchronization speed by 25%
║    • Better error messages for failed operations
║    ... and more
║                                                                  ║
║  💡 You can update now or continue using the current version.
║  🔗 Your settings and data will be preserved.
║                                                                  ║
════════════════════════════════════════════════════════════════════

Would you like to install this update now? (y/n): y

──────────────────────────────────────────────────────────────────
🔄 Installing update...

📦 Installation Details:
  ├─ Downloading files... ✅
  ├─ Verifying checksums... ✅
  ├─ Extracting files... ✅
  ├─ Installing binary... ✅
  ├─ Installing libraries... ✅
  └─ Finalizing installation... ✅

══════════════════════════════════════════════════════════════════
║              ✅ UPDATE INSTALLED SUCCESSFULLY                   ║
║                                                                  ║
║  Updated: v1.2.5 → v1.3.0
║  Status: Ready to use
║                                                                  ║
══════════════════════════════════════════════════════════════════
```

### Declining the Update

```
Would you like to install this update now? (y/n): n

──────────────────────────────────────────────────────────────────
⏭️  Update skipped
─ You can update later using: filemanager --update
─ New features will be available when you update
──────────────────────────────────────────────────────────────────
```

### Key Features

- ✅ User choice
- ✅ No forced updates
- ✅ Can update later
- ✅ Backwards-compatible
- ✅ Preserves settings and data

---

## 🚀 MAJOR Updates (Explicit Consent)

### What Happens

When a MAJOR update is available:

1. 🚀 Full notification with breaking changes
2. ❓ First prompt: "Do you understand the breaking changes? (yes/no)"
3. ❓ Second prompt: "Proceed with upgrade? Type 'UPGRADE' to confirm"
4. 🔄 If confirmed → Multi-step upgrade process
5. ⏭️ If declined → User can upgrade later

### User Experience - Part 1: Information

```
══════════════════════════════════════════════════════════════════
║                  🚀 MAJOR UPGRADE AVAILABLE                     ║
║                                                                  ║
║  Current Version: v1.5.0
║  Available Version: v2.0.0
║                                                                  ║
║  ⚠️  IMPORTANT: This is a major upgrade with breaking changes.
║                                                                  ║
║  📋 What's Changing:
║    • Completely redesigned user interface
║    • Migrated to new database format
║    • Changed configuration file structure
║    • Removed deprecated file format support
║    ... and more
║                                                                  ║
║  ✅ BEFORE YOU UPDATE:
║    1. Backup your configuration and data
║    2. Review the full release notes
║    3. Check the migration guide
║    4. Ensure you have time to troubleshoot if needed
║                                                                  ║
║  🔗 Migration Guide: docs/MIGRATION.md
║  📚 Full Release Notes: See below
║                                                                  ║
══════════════════════════════════════════════════════════════════

📚 FULL RELEASE NOTES:
──────────────────────────────────────────────────────────────────
[Full release notes displayed here]
──────────────────────────────────────────────────────────────────
```

### User Experience - Part 2: Confirmation

```
⚠️  This is a major upgrade. Do you understand the breaking changes?
Type 'yes' to continue (or 'no' to skip): yes

──────────────────────────────────────────────────────────────────
🔒 FINAL CONFIRMATION
─ This will:
─   • Remove the old version (v1.5.0)
─   • Install the new version (v2.0.0)
─   • Migrate your configuration
─   • Potentially require reconfiguration
──────────────────────────────────────────────────────────────────

Proceed with upgrade? Type 'UPGRADE' to confirm (or anything else to cancel): UPGRADE
```

### User Experience - Part 3: Installation

```
══════════════════════════════════════════════════════════════════
║                    🔄 UPGRADE IN PROGRESS                       ║
║──────────────────────────────────────────────────────────────────║

📦 Step 1/5: Creating backup...
..
✅ Backup created successfully

🗑️  Step 2/5: Removing old version (v1.5.0)...
..
✅ Old version removed

📥 Step 3/5: Installing new version (v2.0.0)...
...
✅ New version installed

🔄 Step 4/5: Migrating configuration...
..
✅ Configuration migrated

✔️  Step 5/5: Verifying installation...
..
✅ Installation verified

══════════════════════════════════════════════════════════════════
║              ✅ UPGRADE COMPLETED SUCCESSFULLY                  ║
║                                                                  ║
║  Upgraded: v1.5.0 → v2.0.0
║  Status: Ready to use
║                                                                  ║
║  📝 Next Steps:
║    1. Review new settings: filemanager --settings
║    2. Check release notes: filemanager --changelog
║    3. Explore new features!
║                                                                  ║
══════════════════════════════════════════════════════════════════
```

### Declining the Upgrade

```
⚠️  This is a major upgrade. Do you understand the breaking changes?
Type 'yes' to continue (or 'no' to skip): no

──────────────────────────────────────────────────────────────────
⏭️  Upgrade cancelled
─ You can upgrade later when you're ready
─ Current version: v1.5.0
──────────────────────────────────────────────────────────────────
```

### Key Features

- ✅ Full transparency
- ✅ Explicit consent required
- ✅ Multiple confirmation steps
- ✅ Automatic backup creation
- ✅ Step-by-step progress display
- ✅ Configuration migration
- ✅ Rollback capability

---

## 🎯 Update Triggers

### Automatic Check

Updates are checked automatically:

```bash
# On app startup
filemanager

# Checks for updates in background
# Uses 24-hour cache to avoid excessive API calls
```

### Manual Check

Users can manually check for updates:

```bash
# Check for updates
filemanager --update

# Shows available updates and handles them
```

### Disable Auto-Check

Users can disable automatic update checks:

```bash
# Set environment variable
export FILEMANAGER_NO_UPDATE_CHECK=1

# Or edit configuration
# See docs/CONFIGURATION.md
```

---

## 🔐 Security Considerations

### Automatic PATCH Installation

- ✅ Recommended for security patches
- ✅ Minimal user disruption
- ✅ Cancellable with Ctrl+C
- ✅ Checksums verified before installation

### Optional MINOR Installation

- ✅ User controls installation timing
- ✅ Can skip if not ready
- ✅ Backwards-compatible
- ✅ No breaking changes

### Explicit MAJOR Consent

- ✅ User must explicitly type "UPGRADE"
- ✅ Full release notes displayed
- ✅ Backup created before installation
- ✅ Migration guide provided
- ✅ Rollback instructions available

---

## 📊 Update Flow Diagram

```
┌─────────────────────────────────────────────────────────┐
│         User Launches FileManager                       │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
        ┌────────────────────────────┐
        │  Check for Updates         │
        │  (24-hour cache)           │
        └────────────┬───────────────┘
                     │
        ┌────────────▼───────────────┐
        │  New Version Available?    │
        └────────────┬───────────────┘
                     │
        ┌────────────▼───────────────────────────────────┐
        │  Determine Change Type (SemVer)               │
        │  - PATCH (Z changed)                          │
        │  - MINOR (Y changed)                          │
        │  - MAJOR (X changed)                          │
        └────────────┬───────────────────────────────────┘
                     │
        ┌────────────▼───────────────────────────────────┐
        │  Select Update Strategy                       │
        │  - PATCH → Auto-install with countdown        │
        │  - MINOR → Prompt user                        │
        │  - MAJOR → Explicit consent required          │
        └────────────┬───────────────────────────────────┘
                     │
        ┌────────────▼───────────────────────────────────┐
        │  Handle Update                                │
        │  - Display notification                       │
        │  - Get user input (if needed)                 │
        │  - Install/skip based on response             │
        └────────────┬───────────────────────────────────┘
                     │
        ┌────────────▼───────────────────────────────────┐
        │  Continue with App                            │
        │  - Show summary                               │
        │  - Ready to use                               │
        └─────────────────────────────────────────────────┘
```

---

## 🛠️ Implementation Details

### Update Manager Class

```go
type UpdateManager struct {
    notification *UpdateNotification
    scanner      *bufio.Scanner
}

// Methods:
// - HandleUpdate() → Processes update based on change type
// - handlePatchUpdate() → Auto-install with countdown
// - handleMinorUpdate() → Prompt user
// - handleMajorUpdate() → Explicit consent
// - promptUser() → Get user input
// - showInstallationDetails() → Display progress
// - GetUpdateSummary() → Return summary
```

### Integration Points

```go
// Check for updates on startup
CheckForUpdates()

// Check with prompts enabled
CheckForUpdatesWithPrompt(true)

// Check without prompts (info only)
CheckForUpdatesWithPrompt(false)
```

---

## 📝 Configuration

### Environment Variables

```bash
# Disable automatic update checks
export FILEMANAGER_NO_UPDATE_CHECK=1

# Set custom GitHub token (for rate limiting)
export GITHUB_TOKEN=your_token_here

# Set custom update check interval (hours)
export FILEMANAGER_UPDATE_CHECK_INTERVAL=24
```

### Configuration File

See `docs/CONFIGURATION.md` for:
- Update check frequency
- Auto-install PATCH updates
- Notification preferences
- Download location

---

## 🐛 Troubleshooting

### Update Fails to Install

```bash
# Check logs
filemanager --debug

# Verify disk space
df -h

# Check permissions
ls -la ~/.config/filemanager/
```

### Cancel Stuck Installation

```bash
# Press Ctrl+C during countdown
# Or kill the process
pkill -9 filemanager
```

### Rollback to Previous Version

```bash
# For MAJOR updates, rollback is available
# See docs/MIGRATION.md for instructions

# Restore from backup
cp -r ~/.filemanager.backup ~/.filemanager
```

---

## 📚 Related Documentation

- [UPDATE_STRATEGY.md](./UPDATE_STRATEGY.md) - Overall strategy
- [MIGRATION.md](./MIGRATION.md) - Major version migration
- [RELEASE_NOTES_EXAMPLES.md](./RELEASE_NOTES_EXAMPLES.md) - Release notes
- [CONFIGURATION.md](./CONFIGURATION.md) - Configuration options

---

**FileManager's automatic update system keeps you secure while respecting your workflow** 🎉
