#!/bin/bash
set -e
echo "Installing FileManager..."

# Check if running with sudo
if [ "$EUID" -ne 0 ]; then
    echo "Please run with sudo: sudo ./install.sh"
    exit 1
fi

# Copy binary to /usr/local/bin
cp filemanager /usr/local/bin/
chmod +x /usr/local/bin/filemanager

# Copy library to /usr/local/lib
cp libfs_operations_core.dylib /usr/local/lib/

# Update library search path
if [ ! -f /etc/paths.d/filemanager ]; then
    echo "/usr/local/lib" > /etc/paths.d/filemanager
fi

echo "✅ FileManager installed successfully!"
echo "Run 'filemanager' from terminal to start"
