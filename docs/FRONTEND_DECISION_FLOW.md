# 🔄 Frontend Directory Decision Flow

## Application Startup Flow

```
┌─────────────────────────────────────────────────────────────┐
│  Application Starts                                         │
│  StartWebServer()                                           │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│  GetFrontendDir()                                           │
│  Determine frontend location using priority:               │
│  1. ENV var (FILEMANAGER_FRONTEND_DIR)                     │
│  2. User config (~/.config/filemanager/config.yaml)        │
│  3. System config (/etc/filemanager/config.yaml)           │
│  4. OS default (platform-specific)                         │
│  5. Portable fallback (next to executable)                 │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│  isSystemLevelInstall(frontendDir)                          │
│  Check if path is in system location                        │
└────────────────────┬────────────────────────────────────────┘
                     │
        ┌────────────┴────────────┐
        │                         │
        ▼                         ▼
    YES (System)             NO (User)
        │                         │
        ▼                         ▼
┌──────────────────┐      ┌──────────────────┐
│ System Install   │      │ User Install     │
│                  │      │                  │
│ ✓ Verify exists  │      │ ✓ Create if      │
│ ✓ Check index    │      │   missing        │
│ ✓ Serve files    │      │ ✓ Create basic   │
│ ✓ Read-only OK   │      │   structure      │
│ ✗ No creation    │      │ ✓ User writable  │
└────────┬─────────┘      └────────┬─────────┘
         │                         │
         └────────────┬────────────┘
                      │
                      ▼
         ┌────────────────────────┐
         │ Serve Frontend Files   │
         │ http://localhost:8080  │
         └────────────────────────┘
```

## Frontend Location Resolution Priority

```
┌─────────────────────────────────────────────────────────────┐
│  GetFrontendDir() Priority Order                            │
└─────────────────────────────────────────────────────────────┘

1️⃣  ENVIRONMENT VARIABLE
    │
    ├─ FILEMANAGER_FRONTEND_DIR=/custom/path
    │
    └─ If set and valid → USE THIS ✓

2️⃣  USER CONFIG FILE
    │
    ├─ Linux:   ~/.config/filemanager/config.yaml
    ├─ macOS:   ~/Library/Preferences/filemanager/config.yaml
    ├─ Windows: %APPDATA%\filemanager\config.yaml
    │
    └─ If exists and valid → USE THIS ✓

3️⃣  SYSTEM CONFIG FILE
    │
    ├─ Linux:   /etc/filemanager/config.yaml
    ├─ macOS:   /Library/Application Support/FileManager/config.yaml
    ├─ Windows: %ProgramFiles%\FileManager\config.yaml
    │
    └─ If exists and valid → USE THIS ✓

4️⃣  OS-SPECIFIC DEFAULT
    │
    ├─ Snap:     $SNAP/usr/share/filemanager/frontend
    ├─ Flatpak:  /app/share/filemanager/frontend
    ├─ Linux:    /usr/share/filemanager/frontend
    ├─ macOS:    /usr/local/opt/filemanager/share/frontend
    ├─ Windows:  %ProgramFiles%\FileManager\frontend
    │
    └─ If exists → USE THIS ✓

5️⃣  PORTABLE FALLBACK
    │
    ├─ Next to executable
    │
    └─ If exists → USE THIS ✓

6️⃣  PREFERRED CREATION LOCATION
    │
    ├─ If nothing exists, create at:
    │
    ├─ Linux:   ~/.local/share/filemanager/frontend
    ├─ macOS:   ~/Library/Application Support/FileManager/frontend
    ├─ Windows: %LOCALAPPDATA%\FileManager\frontend
    │
    └─ CREATE HERE ✓
```

## Platform-Specific Decision Tree

### Linux (Snap)
```
Is $SNAP set?
├─ YES
│  ├─ Check $SNAP/usr/share/filemanager/frontend
│  │  ├─ EXISTS → SYSTEM INSTALL (read-only)
│  │  └─ NOT EXISTS → Use $SNAP_USER_COMMON/frontend (user writable)
│  │
│  └─ Return path
│
└─ NO → Continue to next check
```

### Linux (Flatpak)
```
Is $FLATPAK_ID set?
├─ YES
│  ├─ Check /app/share/filemanager/frontend
│  │  ├─ EXISTS → SYSTEM INSTALL (read-only)
│  │  └─ NOT EXISTS → Use ~/.var/app/$FLATPAK_ID/data/frontend
│  │
│  └─ Return path
│
└─ NO → Continue to next check
```

### Linux (System/Manual)
```
Check /usr/share/filemanager/frontend
├─ EXISTS → SYSTEM INSTALL (read-only)
└─ NOT EXISTS → Use ~/.local/share/filemanager/frontend (user writable)
```

### macOS (Homebrew)
```
Check /usr/local/opt/filemanager/share/frontend
├─ EXISTS → SYSTEM INSTALL (read-only)
└─ NOT EXISTS → Check /opt/homebrew/opt/filemanager/share/frontend
   ├─ EXISTS → SYSTEM INSTALL (read-only)
   └─ NOT EXISTS → Use ~/Library/Application Support/FileManager/frontend
```

### Windows (System)
```
Check %ProgramFiles%\FileManager\frontend
├─ EXISTS → SYSTEM INSTALL (admin-only write)
└─ NOT EXISTS → Use %LOCALAPPDATA%\FileManager\frontend (user writable)
```

## System vs User Install Detection

```
┌──────────────────────────────────────────────────────────┐
│  isSystemLevelInstall(frontendDir)                       │
└──────────────────────────────────────────────────────────┘

Check if path matches system patterns:

SNAP:
├─ $SNAP env var set?
└─ Path contains "/snap/"?
   └─ YES → SYSTEM INSTALL ✓

FLATPAK:
├─ $FLATPAK_ID env var set?
└─ Path starts with "/app/"?
   └─ YES → SYSTEM INSTALL ✓

LINUX SYSTEM:
├─ Path starts with "/usr/share/"?
├─ Path starts with "/usr/local/share/"?
└─ Path starts with "/opt/"?
   └─ YES → SYSTEM INSTALL ✓

MACOS SYSTEM:
├─ Path starts with "/usr/local/opt/"?
├─ Path starts with "/opt/homebrew/"?
└─ Path starts with "/Library/Application Support/"?
   └─ YES → SYSTEM INSTALL ✓

WINDOWS SYSTEM:
├─ Path starts with %ProgramFiles%?
   └─ YES → SYSTEM INSTALL ✓

DEFAULT:
└─ NOT A SYSTEM INSTALL → USER INSTALL ✓
```

## Action Based on Install Type

```
┌─────────────────────────────────────────────────────────┐
│  SYSTEM INSTALL (Read-Only)                             │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Verify directory exists                            │
│     └─ If not → ERROR (packaging issue)                │
│                                                         │
│  2. Verify index.html exists                           │
│     └─ If not → ERROR (incomplete package)             │
│                                                         │
│  3. Serve files (read-only is OK)                      │
│     └─ Users can read, not write                       │
│                                                         │
│  4. Log: "Using system-bundled frontend (read-only)"   │
│                                                         │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  USER INSTALL (Writable)                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Check if directory exists                          │
│     ├─ YES → Verify index.html exists                 │
│     │        └─ If yes → Use existing                 │
│     │        └─ If no → Create basic structure        │
│     │                                                  │
│     └─ NO → Create directory and basic structure      │
│                                                         │
│  2. Create basic frontend structure                    │
│     ├─ Create directories (css, js, images)           │
│     ├─ Create index.html                              │
│     ├─ Create style.css                               │
│     └─ Create main.js                                 │
│                                                         │
│  3. Serve files (user can read and write)             │
│     └─ Users can modify files                         │
│                                                         │
│  4. Log: "Using user frontend (writable)"             │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Error Handling Flow

```
┌─────────────────────────────────────────────────────────┐
│  Frontend Setup Error                                   │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
            Is it a system install?
                     │
        ┌────────────┴────────────┐
        │                         │
        ▼                         ▼
      YES                        NO
        │                         │
        ▼                         ▼
   ┌──────────────┐        ┌──────────────┐
   │ FATAL ERROR  │        │ Try to fix   │
   │              │        │              │
   │ "System      │        │ Create dir   │
   │  install     │        │ Create files │
   │  error:      │        │              │
   │  frontend    │        │ If still     │
   │  not found"  │        │ fails →      │
   │              │        │ FATAL ERROR  │
   └──────────────┘        └──────────────┘
        │                         │
        └────────────┬────────────┘
                     │
                     ▼
         ┌──────────────────────┐
         │ Stop server          │
         │ Log error message    │
         │ Exit with error code │
         └──────────────────────┘
```

## Example Scenarios

### Scenario 1: Snap Installation
```
1. User installs: snap install filemanager
2. App starts
3. GetFrontendDir() checks:
   - $SNAP = /snap/filemanager/3
   - Checks $SNAP/usr/share/filemanager/frontend
   - FOUND ✓
4. isSystemLevelInstall() checks:
   - $SNAP env var set? YES
   - Path contains /snap/? YES
   - Result: SYSTEM INSTALL
5. Action:
   - Verify directory exists ✓
   - Verify index.html exists ✓
   - Serve files (read-only)
6. Result: ✅ Frontend loads successfully
```

### Scenario 2: Manual Linux Installation
```
1. User downloads binary
2. User runs: ./filemanager
3. GetFrontendDir() checks:
   - ENV var? NO
   - User config? NO
   - System config? NO
   - OS default? NO
   - Portable? NO
   - Use preferred: ~/.local/share/filemanager/frontend
4. isSystemLevelInstall() checks:
   - Path starts with /usr/share/? NO
   - Path starts with /opt/? NO
   - Result: USER INSTALL
5. Action:
   - Directory doesn't exist
   - Create directory ✓
   - Create basic structure ✓
   - Serve files (writable)
6. Result: ✅ Frontend created and loads successfully
```

### Scenario 3: DEB Installation
```
1. User installs: sudo apt install filemanager
2. DEB postinst copies frontend to /usr/share/filemanager/frontend
3. App starts
4. GetFrontendDir() checks:
   - Finds /usr/share/filemanager/frontend ✓
5. isSystemLevelInstall() checks:
   - Path starts with /usr/share/? YES
   - Result: SYSTEM INSTALL
6. Action:
   - Verify directory exists ✓
   - Verify index.html exists ✓
   - Serve files (read-only)
7. Result: ✅ Frontend loads successfully
```

---

## Key Decision Points

| Decision | System Install | User Install |
|----------|---|---|
| **Frontend Location** | System path | User directory |
| **Writable?** | No (read-only) | Yes (writable) |
| **Create at Runtime?** | No | Yes |
| **Bundled During Build?** | Yes | No |
| **Error if Missing?** | Yes (fatal) | No (create) |
| **Permission Setup?** | By installer | By app |

---

**This flow ensures reliable frontend loading across all platforms and installation types.**
