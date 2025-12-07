#!/bin/bash

# Version Manager for FileManager
# This script manages version updates across the entire project

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION_FILE="$PROJECT_ROOT/VERSION"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display usage
usage() {
    cat << EOF
${BLUE}FileManager Version Manager${NC}

Usage: $(basename "$0") [COMMAND] [VERSION]

Commands:
  get                 Get current version
  set <VERSION>       Set new version (e.g., 2.1.0)
  bump-patch          Bump patch version (2.0.0 → 2.0.1)
  bump-minor          Bump minor version (2.0.0 → 2.1.0)
  bump-major          Bump major version (2.0.0 → 3.0.0)
  list                List all version files
  validate <VERSION>  Validate version format (MAJOR.MINOR.PATCH)
  help                Show this help message

Examples:
  $(basename "$0") get                    # Show current version
  $(basename "$0") set 2.1.0              # Set version to 2.1.0
  $(basename "$0") bump-minor             # Increment minor version
  $(basename "$0") validate 2.1.0         # Check if format is valid

EOF
}

# Function to validate version format
validate_version() {
    local version=$1
    if [[ $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        return 0
    else
        echo -e "${RED}❌ Invalid version format: $version${NC}"
        echo "Expected format: MAJOR.MINOR.PATCH (e.g., 2.1.0)"
        return 1
    fi
}

# Function to get current version
get_version() {
    if [ ! -f "$VERSION_FILE" ]; then
        echo -e "${RED}❌ VERSION file not found at $VERSION_FILE${NC}"
        return 1
    fi
    cat "$VERSION_FILE" | tr -d ' \n'
}

# Function to set version
set_version() {
    local new_version=$1
    
    if ! validate_version "$new_version"; then
        return 1
    fi
    
    local current_version=$(get_version)
    
    echo -e "${YELLOW}Updating version: $current_version → $new_version${NC}"
    
    # Update VERSION file
    echo "$new_version" > "$VERSION_FILE"
    echo -e "${GREEN}✅ Updated VERSION file${NC}"
    
    # Update snapcraft.yaml
    if [ -f "$PROJECT_ROOT/snap/snapcraft.yaml" ]; then
        sed -i "s/version: '[^']*'/version: '$new_version'/" "$PROJECT_ROOT/snap/snapcraft.yaml"
        echo -e "${GREEN}✅ Updated snap/snapcraft.yaml${NC}"
    fi
    
    # Update Go version.go if it exists
    if [ -f "$PROJECT_ROOT/file_manager/pkg/version/version.go" ]; then
        sed -i "s/Version = \"[^\"]*\"/Version = \"$new_version\"/" "$PROJECT_ROOT/file_manager/pkg/version/version.go"
        echo -e "${GREEN}✅ Updated file_manager/pkg/version/version.go${NC}"
    fi
    
    # Update Cargo.toml files
    find "$PROJECT_ROOT/rust_ffi" -name "Cargo.toml" -type f | while read -r cargo_file; do
        sed -i "s/version = \"[^\"]*\"/version = \"$new_version\"/" "$cargo_file"
    done
    echo -e "${GREEN}✅ Updated Cargo.toml files${NC}"
    
    # Update README.md version references
    if [ -f "$PROJECT_ROOT/README.md" ]; then
        sed -i "s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v$new_version/g" "$PROJECT_ROOT/README.md"
        echo -e "${GREEN}✅ Updated README.md${NC}"
    fi
    
    # Update debian/control version
    if [ -f "$PROJECT_ROOT/debian/control" ]; then
        sed -i "s/^Version: [0-9]\+\.[0-9]\+\.[0-9]\+/Version: $new_version/" "$PROJECT_ROOT/debian/control"
        echo -e "${GREEN}✅ Updated debian/control${NC}"
    fi
    
    # Update Flatpak manifest version
    if [ -f "$PROJECT_ROOT/flatpak/com.filemanager.FileManager.yaml" ]; then
        sed -i "s/version: '[^']*'/version: '$new_version'/" "$PROJECT_ROOT/flatpak/com.filemanager.FileManager.yaml"
        echo -e "${GREEN}✅ Updated flatpak/com.filemanager.FileManager.yaml${NC}"
    fi
    
    # Update Homebrew formula version
    if [ -f "$PROJECT_ROOT/homebrew/filemanager.rb" ]; then
        sed -i "s/url \"https:\/\/github\.com\/devonionrouting4Moses\/fileManager\/archive\/v[0-9]\+\.[0-9]\+\.[0-9]\+\.tar\.gz\"/url \"https:\/\/github.com\/devonionrouting4Moses\/fileManager\/archive\/v$new_version.tar.gz\"/" "$PROJECT_ROOT/homebrew/filemanager.rb"
        sed -i "s/version \"[0-9]\+\.[0-9]\+\.[0-9]\+\"/version \"$new_version\"/" "$PROJECT_ROOT/homebrew/filemanager.rb"
        echo -e "${GREEN}✅ Updated homebrew/filemanager.rb${NC}"
    fi
    
    # Update Windows NSIS installer version
    if [ -f "$PROJECT_ROOT/windows/installer.nsi" ]; then
        sed -i "s/Name \"FileManager v[0-9]\+\.[0-9]\+\.[0-9]\+\"/Name \"FileManager v$new_version\"/" "$PROJECT_ROOT/windows/installer.nsi"
        sed -i "s/OutFile \"FileManager-[0-9]\+\.[0-9]\+\.[0-9]\+-Setup\.exe\"/OutFile \"FileManager-$new_version-Setup.exe\"/" "$PROJECT_ROOT/windows/installer.nsi"
        sed -i "s/VIProductVersion \"[0-9]\+\.[0-9]\+\.[0-9]\+\.[0-9]\+\"/VIProductVersion \"$new_version.0\"/" "$PROJECT_ROOT/windows/installer.nsi"
        sed -i "s/VIAddVersionKey \"ProductVersion\" \"[0-9]\+\.[0-9]\+\.[0-9]\+\"/VIAddVersionKey \"ProductVersion\" \"$new_version\"/" "$PROJECT_ROOT/windows/installer.nsi"
        sed -i "s/VIAddVersionKey \"FileVersion\" \"[0-9]\+\.[0-9]\+\.[0-9]\+\"/VIAddVersionKey \"FileVersion\" \"$new_version\"/" "$PROJECT_ROOT/windows/installer.nsi"
        echo -e "${GREEN}✅ Updated windows/installer.nsi${NC}"
    fi
    
    # Update Arch PKGBUILD files
    if [ -f "$PROJECT_ROOT/arch/PKGBUILD" ]; then
        sed -i "s/pkgver=[0-9]\+\.[0-9]\+\.[0-9]\+/pkgver=$new_version/" "$PROJECT_ROOT/arch/PKGBUILD"
        echo -e "${GREEN}✅ Updated arch/PKGBUILD${NC}"
    fi
    
    if [ -f "$PROJECT_ROOT/arch-arm64/PKGBUILD" ]; then
        sed -i "s/pkgver=[0-9]\+\.[0-9]\+\.[0-9]\+/pkgver=$new_version/" "$PROJECT_ROOT/arch-arm64/PKGBUILD"
        echo -e "${GREEN}✅ Updated arch-arm64/PKGBUILD${NC}"
    fi
    
    # Update documentation files version
    if [ -f "$PROJECT_ROOT/IMPLEMENTATION_COMPLETE.md" ]; then
        sed -i "s/\*\*Version\*\*: [0-9]\+\.[0-9]\+\.[0-9]\+/\*\*Version\*\*: $new_version/" "$PROJECT_ROOT/IMPLEMENTATION_COMPLETE.md"
        echo -e "${GREEN}✅ Updated IMPLEMENTATION_COMPLETE.md${NC}"
    fi
    
    if [ -f "$PROJECT_ROOT/IMPLEMENTATION_STATUS.md" ]; then
        sed -i "s/\*\*Version\*\*: [0-9]\+\.[0-9]\+\.[0-9]\+/\*\*Version\*\*: $new_version/" "$PROJECT_ROOT/IMPLEMENTATION_STATUS.md"
        echo -e "${GREEN}✅ Updated IMPLEMENTATION_STATUS.md${NC}"
    fi
    
    if [ -f "$PROJECT_ROOT/VERIFICATION_CHECKLIST.md" ]; then
        sed -i "s/\*\*Version\*\*: [0-9]\+\.[0-9]\+\.[0-9]\+/\*\*Version\*\*: $new_version/" "$PROJECT_ROOT/VERIFICATION_CHECKLIST.md"
        echo -e "${GREEN}✅ Updated VERIFICATION_CHECKLIST.md${NC}"
    fi
    
    if [ -f "$PROJECT_ROOT/PACKAGING_QUICK_START.md" ]; then
        sed -i "s/\*\*Version\*\*: [0-9]\+\.[0-9]\+\.[0-9]\+/\*\*Version\*\*: $new_version/" "$PROJECT_ROOT/PACKAGING_QUICK_START.md"
        echo -e "${GREEN}✅ Updated PACKAGING_QUICK_START.md${NC}"
    fi
    
    if [ -f "$PROJECT_ROOT/FRONTEND_BUNDLING_INDEX.md" ]; then
        sed -i "s/\*\*Version\*\*: [0-9]\+\.[0-9]\+\.[0-9]\+/\*\*Version\*\*: $new_version/" "$PROJECT_ROOT/FRONTEND_BUNDLING_INDEX.md"
        echo -e "${GREEN}✅ Updated FRONTEND_BUNDLING_INDEX.md${NC}"
    fi
    
    # Note: Frontend version is now dynamically fetched from the API
    # The index.html no longer has a hardcoded version
    
    echo -e "${GREEN}🎉 Version updated to $new_version${NC}"
    echo ""
    echo "Files updated:"
    echo "  - VERSION"
    echo "  - snap/snapcraft.yaml"
    echo "  - file_manager/pkg/version/version.go"
    echo "  - rust_ffi/Cargo.toml files"
    echo "  - README.md"
    echo "  - debian/control"
    echo "  - flatpak/com.filemanager.FileManager.yaml"
    echo "  - homebrew/filemanager.rb"
    echo "  - windows/installer.nsi"
    echo "  - IMPLEMENTATION_COMPLETE.md"
    echo "  - IMPLEMENTATION_STATUS.md"
    echo "  - VERIFICATION_CHECKLIST.md"
    echo "  - PACKAGING_QUICK_START.md"
    echo "  - FRONTEND_BUNDLING_INDEX.md"
}

# Function to bump patch version
bump_patch() {
    local current=$(get_version)
    local major=$(echo $current | cut -d. -f1)
    local minor=$(echo $current | cut -d. -f2)
    local patch=$(echo $current | cut -d. -f3)
    
    patch=$((patch + 1))
    local new_version="$major.$minor.$patch"
    
    set_version "$new_version"
}

# Function to bump minor version
bump_minor() {
    local current=$(get_version)
    local major=$(echo $current | cut -d. -f1)
    local minor=$(echo $current | cut -d. -f2)
    
    minor=$((minor + 1))
    local new_version="$major.$minor.0"
    
    set_version "$new_version"
}

# Function to bump major version
bump_major() {
    local current=$(get_version)
    local major=$(echo $current | cut -d. -f1)
    
    major=$((major + 1))
    local new_version="$major.0.0"
    
    set_version "$new_version"
}

# Function to list all version files
list_versions() {
    echo -e "${BLUE}Version files in the project:${NC}"
    echo ""
    
    echo -e "${YELLOW}Core Version:${NC}"
    echo "  VERSION: $(get_version)"
    echo ""
    
    echo -e "${YELLOW}snapcraft.yaml:${NC}"
    grep "^version:" "$PROJECT_ROOT/snap/snapcraft.yaml" 2>/dev/null || echo "  Not found"
    echo ""
    
    echo -e "${YELLOW}Go version.go:${NC}"
    grep "Version = " "$PROJECT_ROOT/file_manager/pkg/version/version.go" 2>/dev/null || echo "  Not found"
    echo ""
    
    echo -e "${YELLOW}Cargo.toml files:${NC}"
    find "$PROJECT_ROOT/rust_ffi" -name "Cargo.toml" -type f -exec grep "^version = " {} + 2>/dev/null || echo "  Not found"
    echo ""
    
    echo -e "${YELLOW}debian/control:${NC}"
    grep "^Version:" "$PROJECT_ROOT/debian/control" 2>/dev/null || echo "  Not found"
    echo ""
    
    echo -e "${YELLOW}Flatpak manifest:${NC}"
    grep "^version:" "$PROJECT_ROOT/flatpak/com.filemanager.FileManager.yaml" 2>/dev/null || echo "  Not found"
    echo ""
    
    echo -e "${YELLOW}Homebrew formula:${NC}"
    grep "version \"" "$PROJECT_ROOT/homebrew/filemanager.rb" 2>/dev/null || echo "  Not found"
    echo ""
    
    echo -e "${YELLOW}Windows installer:${NC}"
    grep "VIProductVersion" "$PROJECT_ROOT/windows/installer.nsi" 2>/dev/null || echo "  Not found"
    echo ""
    
    echo -e "${YELLOW}Arch PKGBUILD:${NC}"
    grep "^pkgver=" "$PROJECT_ROOT/arch/PKGBUILD" 2>/dev/null || echo "  Not found"
    echo ""
    
    echo -e "${YELLOW}Arch ARM64 PKGBUILD:${NC}"
    grep "^pkgver=" "$PROJECT_ROOT/arch-arm64/PKGBUILD" 2>/dev/null || echo "  Not found"
}

# Main logic
case "${1:-help}" in
    get)
        get_version
        ;;
    set)
        if [ -z "$2" ]; then
            echo -e "${RED}❌ Version argument required${NC}"
            usage
            exit 1
        fi
        set_version "$2"
        ;;
    bump-patch)
        bump_patch
        ;;
    bump-minor)
        bump_minor
        ;;
    bump-major)
        bump_major
        ;;
    list)
        list_versions
        ;;
    validate)
        if [ -z "$2" ]; then
            echo -e "${RED}❌ Version argument required${NC}"
            usage
            exit 1
        fi
        if validate_version "$2"; then
            echo -e "${GREEN}✅ Version format is valid: $2${NC}"
        fi
        ;;
    help|--help|-h)
        usage
        ;;
    *)
        echo -e "${RED}❌ Unknown command: $1${NC}"
        usage
        exit 1
        ;;
esac
