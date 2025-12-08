# FileManager Homebrew Formulas

This directory contains Homebrew formulas for installing FileManager on macOS.

## Installation Options

### Option 1: GUI App (Cask - Recommended for most users)

Installs the full macOS application bundle:

```bash
# From your tap
brew tap devonionrouting4Moses/tap
brew install --cask filemanager

# Or directly (once published to official Homebrew)
brew install --cask filemanager
```

### Option 2: CLI Tool (Formula)

Installs the command-line binary:

```bash
# From your tap
brew tap devonionrouting4Moses/tap
brew install filemanager

# Or directly (once published to official Homebrew)
brew install filemanager
```

## Usage

### GUI App
- Launch from Applications folder
- Or run: `open -a FileManager`

### CLI Tool
```bash
# Start the file manager
filemanager

# Access web interface at http://localhost:8080
```

## Uninstall

### GUI App
```bash
brew uninstall --cask filemanager
```

### CLI Tool
```bash
brew uninstall filemanager
```

## Development

### Testing Formulas Locally

```bash
# Test the Cask
brew install --cask --debug ./Casks/filemanager.rb

# Test the Formula
brew install --formula --debug ./filemanager.rb
```

### Audit Before Submission

```bash
# Audit the Cask
brew audit --cask --strict ./Casks/filemanager.rb

# Audit the Formula
brew audit --strict ./filemanager.rb
```

## Publishing

### To Your Own Tap

1. Create a repository: `homebrew-tap`
2. Copy formulas:
   - `filemanager.rb` → `Formula/filemanager.rb` (for CLI)
   - `Casks/filemanager.rb` → `Casks/filemanager.rb` (for GUI)
3. Users install via: `brew tap yourusername/tap`

### To Official Homebrew

#### For Cask (GUI):
1. Fork `homebrew/homebrew-cask`
2. Add your cask to `Casks/f/filemanager.rb`
3. Submit pull request

#### For Formula (CLI):
1. Fork `homebrew/homebrew-core`
2. Add your formula to `Formula/f/filemanager.rb`
3. Submit pull request

## Checksums

Checksums are automatically generated during build and embedded in these formulas.

To verify:
```bash
shasum -a 256 FileManager-*.dmg
shasum -a 256 filemanager-*-macos-*.tar.gz
```
