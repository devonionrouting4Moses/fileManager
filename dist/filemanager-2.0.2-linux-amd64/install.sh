#!/bin/bash
set -e
echo "Installing FileManager..."
sudo cp filemanager /usr/local/bin/
sudo chmod +x /usr/local/bin/filemanager
sudo cp libfs_operations_core.so /usr/local/lib/
sudo ldconfig
echo "✅ FileManager installed successfully!"
echo "Run 'filemanager' to start"
