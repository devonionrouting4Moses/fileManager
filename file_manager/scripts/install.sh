#!/bin/bash
# FileManager Installation Script

set -e

APP_NAME="filemanager"
INSTALL_DIR="/usr/local/bin"
LIB_DIR="/usr/local/lib"

echo "📦 Installing FileManager..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  This script requires sudo privileges"
    echo "Please run: sudo ./install.sh"
    exit 1
fi

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Copy binary
echo "📋 Copying binary to ${INSTALL_DIR}..."
cp ${APP_NAME} ${INSTALL_DIR}/
chmod +x ${INSTALL_DIR}/${APP_NAME}

# Copy library
echo "📚 Copying library to ${LIB_DIR}..."
if [ "$OS" = "darwin" ]; then
    cp rust_ffi/target/release/libfs_operations_core.dylib ${LIB_DIR}/ 2>/dev/null || \
    cp rust_ffi/target/release/libfs_operations_core.so ${LIB_DIR}/libfs_operations_core.dylib
else
    cp rust_ffi/target/release/libfs_operations_core.so ${LIB_DIR}/
fi

# Update library cache (Linux only)
if [ "$OS" = "linux" ]; then
    ldconfig
fi

echo ""
echo "✅ FileManager installed successfully!"
echo ""
echo "Usage:"
echo "  filemanager              Start interactive mode"
echo "  filemanager --version    Show version"
echo "  filemanager --update     Check for updates"
echo "  filemanager --help       Show help"
echo "  filemanager --web        Start web interface"
echo ""