#!/bin/bash

# Get the project root (parent of scripts directory)
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Read version from VERSION file
if [ -f "$PROJECT_ROOT/VERSION" ]; then
    VERSION=$(cat "$PROJECT_ROOT/VERSION" | tr -d ' \n')
else
    VERSION="2.0.0"
    echo "⚠️  VERSION file not found, using default: $VERSION"
fi

echo "🔨 FileManager v${VERSION} - Multi-Platform Build"
echo "=================================================="
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check what to build
if [ $# -eq 0 ]; then
    echo "Usage: $0 [linux-amd64|amd64] [linux-arm64|arm64] [windows-amd64] [harmonyos] [all]"
    echo ""
    echo "Examples:"
    echo "  $0 linux-amd64          # Build for Linux amd64 only"
    echo "  $0 amd64                # Build for Linux amd64 only (shorthand)"
    echo "  $0 arm64                # Build for Linux arm64 only (shorthand)"
    echo "  $0 macos-amd64          # Build for macOS Intel"
    echo "  $0 macos-arm64          # Build for macOS Apple Silicon"
    echo "  $0 linux-amd64 windows-amd64  # Build for Linux and Windows"
    echo "  $0 all                  # Build for all platforms"
    echo ""
    exit 1
fi

BUILD_LINUX_AMD64=0
BUILD_LINUX_ARM64=0
BUILD_WINDOWS=0
BUILD_MACOS_AMD64=0
BUILD_MACOS_ARM64=0
BUILD_HARMONYOS=0

for arg in "$@"; do
    case $arg in
        linux-amd64|amd64)
            BUILD_LINUX_AMD64=1
            ;;
        linux-arm64|arm64)
            BUILD_LINUX_ARM64=1
            ;;
        windows-amd64)
            BUILD_WINDOWS=1
            ;;
        macos-amd64)
            BUILD_MACOS_AMD64=1
            ;;
        macos-arm64)
            BUILD_MACOS_ARM64=1
            ;;
        macos)
            BUILD_MACOS_AMD64=1
            BUILD_MACOS_ARM64=1
            ;;
        harmonyos)
            BUILD_HARMONYOS=1
            ;;
        all)
            BUILD_LINUX_AMD64=1
            BUILD_LINUX_ARM64=1
            BUILD_WINDOWS=1
            BUILD_MACOS_AMD64=1
            BUILD_MACOS_ARM64=1
            BUILD_HARMONYOS=1
            ;;
        *)
            echo "Unknown platform: $arg"
            exit 1
            ;;
    esac
done

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Linux amd64
if [ $BUILD_LINUX_AMD64 -eq 1 ]; then
    echo -e "${BLUE}Building for Linux amd64...${NC}"
    if bash "${SCRIPT_DIR}/build-linux-amd64.sh"; then
        echo -e "${GREEN}✅ Linux amd64 build completed${NC}"
    else
        echo -e "${RED}❌ Linux amd64 build failed${NC}"
        exit 1
    fi
    echo ""
fi

# Linux arm64
if [ $BUILD_LINUX_ARM64 -eq 1 ]; then
    echo -e "${BLUE}Building for Linux arm64...${NC}"
    if bash "${SCRIPT_DIR}/build-linux-arm64.sh"; then
        echo -e "${GREEN}✅ Linux arm64 build completed${NC}"
    else
        echo -e "${RED}❌ Linux arm64 build failed${NC}"
        exit 1
    fi
    echo ""
fi

# Windows amd64
if [ $BUILD_WINDOWS -eq 1 ]; then
    echo -e "${BLUE}Building for Windows amd64...${NC}"
    if bash "${SCRIPT_DIR}/build-windows-amd64.sh"; then
        echo -e "${GREEN}✅ Windows amd64 build completed${NC}"
    else
        echo -e "${RED}❌ Windows amd64 build failed${NC}"
        exit 1
    fi
    echo ""
fi

# macOS amd64
if [ $BUILD_MACOS_AMD64 -eq 1 ]; then
    echo -e "${BLUE}Building for macOS amd64 (Intel)...${NC}"
    if bash "${SCRIPT_DIR}/build-macos.sh" amd64; then
        echo -e "${GREEN}✅ macOS amd64 build completed${NC}"
    else
        echo -e "${YELLOW}⚠️  macOS amd64 build failed or requires macOS system${NC}"
    fi
    echo ""
fi

# macOS arm64
if [ $BUILD_MACOS_ARM64 -eq 1 ]; then
    echo -e "${BLUE}Building for macOS arm64 (Apple Silicon)...${NC}"
    if bash "${SCRIPT_DIR}/build-macos.sh" arm64; then
        echo -e "${GREEN}✅ macOS arm64 build completed${NC}"
    else
        echo -e "${YELLOW}⚠️  macOS arm64 build failed or requires macOS system${NC}"
    fi
    echo ""
fi

# HarmonyOS
if [ $BUILD_HARMONYOS -eq 1 ]; then
    echo -e "${BLUE}Building for HarmonyOS...${NC}"
    if bash "${SCRIPT_DIR}/build-harmonyos.sh"; then
        echo -e "${GREEN}✅ HarmonyOS build completed${NC}"
    else
        echo -e "${YELLOW}⚠️  HarmonyOS build skipped (SDK not installed)${NC}"
    fi
    echo ""
fi

echo -e "${GREEN}🎉 All requested builds completed!${NC}"
echo ""
echo "Build artifacts:"
echo "  - Linux amd64: filemanager_${VERSION}_amd64.deb, filemanager-${VERSION}-linux-amd64.tar.gz"
echo "  - Linux arm64: filemanager_${VERSION}_arm64.deb, filemanager-${VERSION}-linux-arm64.tar.gz"
echo "  - Windows: filemanager.exe, filemanager-${VERSION}-windows-amd64.zip"
echo "  - macOS amd64: FileManager.app, FileManager-${VERSION}-macos-amd64.dmg"
echo "  - macOS arm64: FileManager.app, FileManager-${VERSION}-macos-arm64.dmg"
echo "  - HarmonyOS: harmonyos/entry/build/outputs/hap/"
