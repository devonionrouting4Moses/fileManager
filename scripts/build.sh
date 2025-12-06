#!/bin/bash

# FileManager Local Build Script
# Builds Rust library and Go binary for local testing

set -e

VERSION="2.0.0"
APP_NAME="filemanager"
BUILD_DIR="dist"
RUST_DIR="../rust_ffi"
GO_DIR="../file_manager"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and architecture
UNAME_S=$(uname -s)
UNAME_M=$(uname -m)

case "$UNAME_S" in
    Linux*)
        OS="linux"
        LIB_EXT="so"
        LIB_NAME="libfs_operations_core.so"
        ;;
    Darwin*)
        OS="darwin"
        LIB_EXT="dylib"
        LIB_NAME="libfs_operations_core.dylib"
        ;;
    *)
        echo -e "${RED}❌ Unsupported OS: $UNAME_S${NC}"
        exit 1
        ;;
esac

case "$UNAME_M" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}❌ Unsupported architecture: $UNAME_M${NC}"
        exit 1
        ;;
esac

echo -e "${BLUE}🏗️  Building FileManager v${VERSION} for ${OS}/${ARCH}${NC}"
echo ""

# Step 1: Build Rust library
echo -e "${BLUE}📦 Step 1: Building Rust library...${NC}"
cd "$SCRIPT_DIR/$RUST_DIR"

if ! cargo build --release -p fs-operations-core; then
    echo -e "${RED}❌ Rust build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Rust library built${NC}"
echo ""

# Step 2: Build Go binary
echo -e "${BLUE}🐹 Step 2: Building Go binary...${NC}"
cd "$SCRIPT_DIR/$GO_DIR"

if ! env CGO_ENABLED=1 \
    CGO_LDFLAGS="-L$SCRIPT_DIR/$RUST_DIR/target/release -lfs_operations_core -ldl -lpthread -lm" \
    go build -ldflags="-s -w" -o "$SCRIPT_DIR/../${APP_NAME}.bin" ./cmd/app; then
    echo -e "${RED}❌ Go build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go binary built${NC}"
echo ""

# Step 3: Create local distribution
echo -e "${BLUE}📦 Step 3: Creating local distribution...${NC}"
cd "$SCRIPT_DIR"

mkdir -p "$BUILD_DIR"

# Copy binary
cp "../${APP_NAME}.bin" "$BUILD_DIR/"
chmod +x "$BUILD_DIR/${APP_NAME}.bin"

# Copy Rust library
cp "$RUST_DIR/target/release/$LIB_NAME" "$BUILD_DIR/"

# Copy wrapper script
if [ -f "../${APP_NAME}.sh" ]; then
    cp "../${APP_NAME}.sh" "$BUILD_DIR/"
    chmod +x "$BUILD_DIR/${APP_NAME}.sh"
fi

echo -e "${GREEN}✅ Local distribution created${NC}"
echo ""

# Step 4: Verify build
echo -e "${BLUE}🔍 Step 4: Verifying build...${NC}"

if [ ! -f "$BUILD_DIR/${APP_NAME}.bin" ]; then
    echo -e "${RED}❌ Binary not found: $BUILD_DIR/${APP_NAME}.bin${NC}"
    exit 1
fi

if [ ! -f "$BUILD_DIR/$LIB_NAME" ]; then
    echo -e "${RED}❌ Library not found: $BUILD_DIR/$LIB_NAME${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All artifacts verified${NC}"
echo ""

# Summary
echo -e "${GREEN}✨ Build complete!${NC}"
echo ""
echo -e "${YELLOW}📂 Build artifacts:${NC}"
ls -lh "$BUILD_DIR"/ | tail -n +2 | awk '{printf "   %-40s %8s\n", $9, $5}'
echo ""
echo -e "${YELLOW}🚀 To run the application:${NC}"
echo "   cd $BUILD_DIR"
echo "   LD_LIBRARY_PATH=. ./${APP_NAME}.bin"
echo ""
echo -e "${YELLOW}📝 Or use the wrapper script:${NC}"
echo "   ./${APP_NAME}.sh"
echo ""
