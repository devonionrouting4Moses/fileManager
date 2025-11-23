#!/bin/bash
# FileManager Uninstallation Script

set -e

APP_NAME="filemanager"
INSTALL_DIR="/usr/local/bin"
LIB_DIR="/usr/local/lib"

echo "🗑️  Uninstalling FileManager..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  This script requires sudo privileges"
    echo "Please run: sudo ./uninstall.sh"
    exit 1
fi

# Remove binary
if [ -f "${INSTALL_DIR}/${APP_NAME}" ]; then
    rm -f ${INSTALL_DIR}/${APP_NAME}
    echo "✅ Removed binary"
fi

# Remove library
rm -f ${LIB_DIR}/libfs_operations_core.so ${LIB_DIR}/libfs_operations_core.dylib
echo "✅ Removed library"

# Update library cache (Linux only)
if [ "$(uname -s | tr '[:upper:]' '[:lower:]')" = "linux" ]; then
    ldconfig
fi

echo ""
echo "✅ FileManager uninstalled successfully!"
echo ""