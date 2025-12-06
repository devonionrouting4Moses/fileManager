# Version Management System Analysis

## 🎯 Overview

The FileManager version management system is a **three-layer architecture** that handles versioning, updates, and notifications:

1. **Version Tracking Layer** (`semver.go`) - Semantic versioning logic
2. **Version Management Layer** (`version.go`) - Update detection & notification
3. **Version Control Layer** (`version-manager.sh`) - Manual version updates

---

## 📊 System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    VERSION MANAGEMENT SYSTEM                │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  LAYER 1: Version Tracking (semver.go)              │   │
│  │  ├─ SemVer struct (Major.Minor.Patch)               │   │
│  │  ├─ ParseSemVer() - Parse version strings           │   │
│  │  ├─ Compare() - Compare two versions                │   │
│  │  ├─ DetermineChangeType() - Identify change impact  │   │
│  │  └─ GetUpdateStrategy() - Recommend notification    │   │
│  └──────────────────────────────────────────────────────┘   │
│                           ▲                                  │
│                           │ Uses                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  LAYER 2: Update Detection (version.go)             │   │
│  │  ├─ GetVersion() - Load current version             │   │
│  │  ├─ CheckForUpdates() - GitHub API polling          │   │
│  │  ├─ CreateUpdateNotification() - Build notification │   │
│  │  ├─ HandleUpdate() - Process updates                │   │
│  │  └─ Cache management (24-hour TTL)                  │   │
│  └──────────────────────────────────────────────────────┘   │
│                           ▲                                  │
│                           │ Calls                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  LAYER 3: Version Control (version-manager.sh)      │   │
│  │  ├─ get - Read current version                      │   │
│  │  ├─ set - Update all version files                  │   │
│  │  ├─ bump-patch/minor/major - Auto-increment         │   │
│  │  └─ list - Show all version references              │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔄 How It Works During Different Operations

### 1️⃣ PATCH UPDATE (Bug Fix / Security Patch)

**Scenario:** v2.0.0 → v2.0.1

#### Flow:
```
User launches FileManager
         ↓
CheckForUpdates() called
         ↓
GitHub API: fetch latest release
         ↓
Compare: v2.0.0 vs v2.0.1
         ↓
DetermineChangeType() → ChangeTypePatch (Z changed)
         ↓
GetUpdateStrategy() → "Silent/Direct Install"
         ↓
CreateUpdateNotification() builds notification
         ↓
DisplayNotification() shows minimal banner
         ↓
Auto-install on next restart (no user action needed)
```

#### Code Flow:
```go
// semver.go - Line 87-94
func (s SemVer) DetermineChangeType(other SemVer) ChangeType {
    if s.Major != other.Major {
        return ChangeTypeMajor
    }
    if s.Minor != other.Minor {
        return ChangeTypeMinor
    }
    return ChangeTypePatch  // ← Only patch changed
}

// semver.go - Line 126-136
func GetUpdateStrategy(changeType ChangeType) string {
    case ChangeTypePatch:
        return "Silent/Direct Install"  // ← Minimal disruption
}

// version.go - Line 274-310
func handleUpdateWithManager(release ReleaseInfo, showPrompts bool) {
    // Creates notification with appropriate strategy
    notification, err := CreateUpdateNotification(release)
    notification.DisplayNotification()  // ← Shows based on type
}
```

#### User Experience:
- ✅ Minimal notification (banner)
- ✅ Auto-installs on restart
- ✅ No user action required
- ✅ Zero disruption

---

### 2️⃣ MINOR UPDATE (New Features)

**Scenario:** v2.0.0 → v2.1.0

#### Flow:
```
User launches FileManager
         ↓
CheckForUpdates() called
         ↓
GitHub API: fetch latest release
         ↓
Compare: v2.0.0 vs v2.1.0
         ↓
DetermineChangeType() → ChangeTypeMinor (Y changed)
         ↓
GetUpdateStrategy() → "Subtle In-App Banner/Hotspot"
         ↓
CreateUpdateNotification() builds notification
         ↓
DisplayNotification() shows in-app banner
         ↓
User can update at convenience (optional)
```

#### Code Flow:
```go
// semver.go - Line 87-94
func (s SemVer) DetermineChangeType(other SemVer) ChangeType {
    if s.Major != other.Major {
        return ChangeTypeMajor
    }
    if s.Minor != other.Minor {
        return ChangeTypeMinor  // ← Minor changed
    }
    return ChangeTypePatch
}

// semver.go - Line 126-136
func GetUpdateStrategy(changeType ChangeType) string {
    case ChangeTypeMinor:
        return "Subtle In-App Banner/Hotspot"  // ← Non-intrusive
}
```

#### User Experience:
- 📢 Subtle in-app banner
- 🔗 Shows release notes
- ⏰ User decides when to update
- ✨ New features highlighted

---

### 3️⃣ MAJOR UPDATE (Breaking Changes)

**Scenario:** v1.5.0 → v2.0.0

#### Flow:
```
User launches FileManager
         ↓
CheckForUpdates() called
         ↓
GitHub API: fetch latest release
         ↓
Compare: v1.5.0 vs v2.0.0
         ↓
DetermineChangeType() → ChangeTypeMajor (X changed)
         ↓
GetUpdateStrategy() → "Modal Window/Full Screen Splash"
         ↓
CreateUpdateNotification() builds notification
         ↓
DisplayNotification() shows full-screen modal
         ↓
User must review breaking changes & consent
```

#### Code Flow:
```go
// semver.go - Line 87-94
func (s SemVer) DetermineChangeType(other SemVer) ChangeType {
    if s.Major != other.Major {
        return ChangeTypeMajor  // ← Major changed
    }
    // ...
}

// semver.go - Line 126-136
func GetUpdateStrategy(changeType ChangeType) string {
    case ChangeTypeMajor:
        return "Modal Window/Full Screen Splash"  // ← Prominent
}
```

#### User Experience:
- 🚀 Full-screen modal notification
- ⚠️ Breaking changes highlighted
- 📋 Migration guide provided
- ✅ Requires explicit consent

---

## 🛠️ Manual Version Management

### Updating Version (version-manager.sh)

#### Scenario: Release v2.1.0

```bash
# Step 1: Bump version
./version-manager.sh bump-minor
# Current: 2.0.0 → New: 2.1.0

# Step 2: What gets updated automatically
# ✅ VERSION file
# ✅ snap/snapcraft.yaml
# ✅ file_manager/pkg/version/version.go
# ✅ rust_ffi/Cargo.toml files
# ✅ README.md
```

#### Code Flow (version-manager.sh):
```bash
# Lines 128-137
bump_minor() {
    local current=$(get_version)
    local major=$(echo $current | cut -d. -f1)
    local minor=$(echo $current | cut -d. -f2)
    
    minor=$((minor + 1))
    local new_version="$major.$minor.0"  # ← Resets patch to 0
    
    set_version "$new_version"  # ← Updates all files
}

# Lines 66-113
set_version() {
    # Updates VERSION file
    echo "$new_version" > "$VERSION_FILE"
    
    # Updates snapcraft.yaml
    sed -i "s/version: '[^']*'/version: '$new_version'/" ...
    
    # Updates Go version
    sed -i "s/Version = \"[^\"]*\"/Version = \"$new_version\"/" ...
    
    # Updates Cargo.toml files
    find ... -name "Cargo.toml" ... sed -i ...
    
    # Updates README.md
    sed -i "s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v$new_version/g" ...
}
```

#### Files Updated:
| File | Before | After | Purpose |
|------|--------|-------|---------|
| VERSION | 2.0.0 | 2.1.0 | Source of truth |
| snap/snapcraft.yaml | version: '2.0.0' | version: '2.1.0' | Snap package |
| version.go | Version = "2.0.0" | Version = "2.1.0" | Go constant |
| Cargo.toml | version = "2.0.0" | version = "2.1.0" | Rust packages |
| README.md | v2.0.0 | v2.1.0 | Documentation |

---

## 📥 Version Loading Flow

### When App Starts (version.go):

```go
// Lines 31-36
func GetVersion() string {
    versionOnce.Do(func() {
        cachedVersion = loadVersionFromFile()
    })
    return cachedVersion
}

// Lines 39-63
func loadVersionFromFile() string {
    possiblePaths := []string{
        filepath.Join(getExecutableDir(), "..", "..", "VERSION"),  // From binary
        "VERSION",                                                   // From CWD
        filepath.Join(os.Getenv("HOME"), ".filemanager", "VERSION"), // From home
        "/etc/filemanager/VERSION",                                 // System config
    }
    
    for _, path := range possiblePaths {
        if data, err := os.ReadFile(path); err == nil {
            version := strings.TrimSpace(string(data))
            if version != "" {
                return version
            }
        }
    }
    
    return defaultVersion  // "2.0.0"
}
```

#### Priority Order:
1. **Binary directory** - `../../../VERSION` (relative to executable)
2. **Current working directory** - `./VERSION`
3. **Home directory** - `~/.filemanager/VERSION`
4. **System config** - `/etc/filemanager/VERSION`
5. **Default** - `"2.0.0"`

---

## 🔍 Update Detection Flow

### CheckForUpdates() Process (version.go):

```go
// Lines 197-272
func CheckForUpdatesWithPrompt(showPrompts bool) {
    // 1. Skip if dev build
    if isDevBuild() {
        return
    }
    
    // 2. Check cache first (24-hour TTL)
    if cached, ok := loadCache(); ok {
        if time.Since(cached.LastCheck) < 24*time.Hour {
            handleUpdateWithManager(cached.ReleaseInfo, false)
            return
        }
    }
    
    // 3. Fetch from GitHub API
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", releaseURL, nil)
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    
    // 4. Parse response
    var release ReleaseInfo
    json.Unmarshal(body, &release)
    
    // 5. Save to cache
    saveCache(UpdateCache{
        LastCheck:   time.Now(),
        ReleaseInfo: release,
    })
    
    // 6. Handle update
    handleUpdateWithManager(release, showPrompts)
}
```

#### Cache Behavior:
- **First check**: Fetches from GitHub API
- **Subsequent checks (24h)**: Uses cached data
- **After 24h**: Fetches fresh data from GitHub
- **Offline mode**: Falls back to cache

---

## 🎯 Notification Strategy Selection

### Decision Tree (semver.go):

```
New Version Available?
    ↓
Compare Current vs Available
    ↓
┌───────────────────────────────────────────┐
│                                           │
├─ MAJOR changed (X.Y.Z)                   │
│  └─ Strategy: Modal/Full-screen          │
│     Impact: High                          │
│     Action: Requires consent              │
│                                           │
├─ MINOR changed (X.Y.Z)                   │
│  └─ Strategy: Subtle banner               │
│     Impact: Low-Moderate                  │
│     Action: Optional                      │
│                                           │
└─ PATCH changed (X.Y.Z)                   │
   └─ Strategy: Silent/Direct               │
      Impact: Minimal                       │
      Action: Auto-install                  │
```

---

## 📋 Comparison Matrix

| Aspect | PATCH | MINOR | MAJOR |
|--------|-------|-------|-------|
| **Version Change** | Z only | Y only | X only |
| **Example** | 2.0.0 → 2.0.1 | 2.0.0 → 2.1.0 | 2.0.0 → 3.0.0 |
| **Change Type** | ChangeTypePatch | ChangeTypeMinor | ChangeTypeMajor |
| **Strategy** | Silent/Direct | Subtle Banner | Modal/Splash |
| **User Impact** | Minimal | Low-Moderate | High |
| **User Action** | None | Optional | Required |
| **Display** | Minimal banner | In-app notification | Full-screen modal |
| **Emoji** | 🔧 | ✨ | 🚀 |
| **Auto-install** | Yes | No | No |
| **Requires Consent** | No | No | Yes |

---

## 🔐 Security Considerations

### Patch Updates:
- ✅ Automatic installation recommended
- ✅ Silent background download
- ✅ Always verify checksums
- ✅ Critical for security patches

### Minor Updates:
- ✅ Optional (user can skip)
- ✅ Backwards compatible
- ✅ Safe to defer
- ✅ Reminder on next app open

### Major Updates:
- ✅ Explicit user consent required
- ✅ Backup before migration
- ✅ Migration guide provided
- ✅ Rollback option documented

---

## 🚀 Release Workflow

### Complete Release Process:

```bash
# 1. Update version
cd scripts
./version-manager.sh bump-minor
# 2.0.0 → 2.1.0

# 2. Verify all files updated
./version-manager.sh list

# 3. Build all platforms
./build-all-platforms.sh all

# 4. Test packages
# ... manual testing ...

# 5. Commit changes
git add -A
git commit -m "Bump version to 2.1.0"

# 6. Create git tag
git tag v2.1.0

# 7. Push to GitHub
git push origin v2.1.0

# 8. GitHub Actions creates release
# 9. Users see update notification
# 10. Update strategy applied based on version change
```

---

## 🔄 Integration Points

### Build System Integration:
```bash
# build-all-platforms.sh automatically reads VERSION
VERSION=$(cat ../VERSION | tr -d ' \n')

# Creates packages with correct version
# - filemanager_2.1.0_amd64.deb
# - filemanager-2.1.0-linux-amd64.tar.gz
# - filemanager-2.1.0-windows-amd64.zip
```

### CI/CD Integration:
```yaml
# GitHub Actions can use version-manager.sh
- name: Bump version
  run: |
    cd scripts
    ./version-manager.sh bump-patch
    
- name: Get version
  run: |
    VERSION=$(cd scripts && ./version-manager.sh get)
    echo "Building version: $VERSION"
```

---

## ✅ Verification Checklist

### Before Release:
- [ ] Version bumped using `version-manager.sh`
- [ ] All files updated (run `list` command)
- [ ] Build succeeds for all platforms
- [ ] Tests pass
- [ ] Release notes prepared
- [ ] Git tag created with version
- [ ] GitHub release published

### After Release:
- [ ] Users see appropriate notification
- [ ] Update strategy matches version change
- [ ] Download links work
- [ ] Installation succeeds
- [ ] No duplicate version files

---

## 📚 Key Files Summary

| File | Purpose | Key Functions |
|------|---------|----------------|
| **semver.go** | Version comparison logic | ParseSemVer, Compare, DetermineChangeType |
| **version.go** | Update detection & notification | GetVersion, CheckForUpdates, CreateUpdateNotification |
| **version-manager.sh** | Manual version control | bump-patch/minor/major, set, list |
| **UPDATE_STRATEGY.md** | User-facing strategy | Notification types, user impact |
| **VERSION_MANAGEMENT.md** | Developer guide | How to use version-manager.sh |

---

## 🎯 Summary

The FileManager version management system provides:

1. **Semantic Versioning** - Clear communication of change impact
2. **Automatic Detection** - GitHub API polling with 24-hour cache
3. **Smart Notifications** - Appropriate UI based on change type
4. **Centralized Control** - Single `version-manager.sh` updates all files
5. **Security** - Patch updates auto-install, major updates require consent
6. **Build Integration** - Automatic version injection into packages

**Result:** Users always know what's changing, and the system respects their time while keeping them secure.
