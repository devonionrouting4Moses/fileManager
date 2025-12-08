# Quick Setup Instructions

## Step 1: Replace Your Build Scripts

Replace these 3 files in your `scripts/` directory:

### 1. scripts/build-macos.sh
Copy the content from the "Improved build-macos.sh" artifact above.

### 2. scripts/build-macos-all.sh  
Copy the content from the "build-macos-all.sh" artifact above.

### 3. scripts/generate-homebrew-formula.sh
Copy the content from the "generate-homebrew-formula.sh" artifact above (document #4).

## Step 2: Make Scripts Executable

```bash
cd ~/Documents/MOSESE/fileManager
chmod +x scripts/build-macos.sh
chmod +x scripts/build-macos-all.sh
chmod +x scripts/generate-homebrew-formula.sh
```

## Step 3: Clean Up Old Builds (Optional)

```bash
# Remove old dist folder if it exists
rm -rf dist/

# Remove old homebrew folder if it exists
rm -rf homebrew/

# Remove old binary files at root if they exist
rm -f filemanager-macos-*
```

## Step 4: Build Everything

```bash
bash scripts/build-macos-all.sh
```

This will create the following structure:

```
fileManager/
├── macos/
│   ├── filemanager-macos-amd64          # Intel binary
│   ├── filemanager-macos-arm64          # ARM binary
│   ├── dist/
│   │   ├── amd64/
│   │   │   ├── FileManager.app/
│   │   │   ├── FileManager-2.0.2-macos-amd64.dmg
│   │   │   └── filemanager-2.0.2-macos-amd64.tar.gz
│   │   ├── arm64/
│   │   │   ├── FileManager.app/
│   │   │   ├── FileManager-2.0.2-macos-arm64.dmg
│   │   │   └── filemanager-2.0.2-macos-arm64.tar.gz
│   │   ├── checksums-amd64.txt
│   │   └── checksums-arm64.txt
│   └── homebrew/
│       ├── filemanager.rb
│       ├── Casks/
│       │   └── filemanager.rb
│       ├── CHECKSUMS.txt
│       └── README.md
```

## Step 5: Test Your Builds

```bash
# Test Intel build
open macos/dist/amd64/FileManager.app

# Test ARM build (if you're on Apple Silicon)
open macos/dist/arm64/FileManager.app
```

## Step 6: View Checksums

```bash
cat macos/homebrew/CHECKSUMS.txt
```

You should see all the SHA256 checksums for your DMG and TAR.GZ files.

## Troubleshooting

### If build-macos.sh still overwrites files:
Make sure you've replaced the file with the new version from the artifact. The key changes are:
- Line 53-57: Creates `macos/` directory structure
- Line 70-72: Saves binary to `macos/` directory
- Line 77: Creates app bundle in `macos/dist/{arch}/`
- All subsequent paths use the new directory structure

### If checksums are missing:
The checksums are now saved to:
- `macos/dist/checksums-amd64.txt`
- `macos/dist/checksums-arm64.txt`
- `macos/homebrew/CHECKSUMS.txt` (combined, after running generate-homebrew-formula.sh)

### If Homebrew formulas have wrong paths:
Make sure you run `bash scripts/generate-homebrew-formula.sh` AFTER building both architectures.

## What's Different Now?

### Before (Old Structure):
```
dist/
├── FileManager.app          # Gets overwritten each build
├── FileManager-*-amd64.dmg  # Gets overwritten
└── FileManager-*-arm64.dmg  # Gets overwritten
homebrew/
└── filemanager.rb           # Gets overwritten each build
```

### After (New Structure):
```
macos/
├── dist/
│   ├── amd64/               # Intel-specific files
│   └── arm64/               # ARM-specific files
└── homebrew/                # Generated ONCE after both builds
```

Each architecture has its own directory, and the Homebrew formula is only generated after both builds complete!
