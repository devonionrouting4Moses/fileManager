#!/bin/bash

# FileManager Wrapper Script
# Sets up the environment and runs the FileManager binary with proper library paths

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Get the project root (parent directory of scripts)
PROJECT_ROOT="$(dirname "${SCRIPT_DIR}")"

# Check for binary in different locations
BINARY_PATH=""
if [ -f "${SCRIPT_DIR}/dist/filemanager.bin" ]; then
    # Binary from local build (build.sh output)
    BINARY_PATH="${SCRIPT_DIR}/dist/filemanager.bin"
    LIBRARY_PATH="${SCRIPT_DIR}/dist"
elif [ -f "${PROJECT_ROOT}/filemanager.bin" ]; then
    # Binary in project root
    BINARY_PATH="${PROJECT_ROOT}/filemanager.bin"
    LIBRARY_PATH="${PROJECT_ROOT}/rust_ffi/target/release"
elif [ -f "${PROJECT_ROOT}/filemanager" ]; then
    # Fallback to old location
    BINARY_PATH="${PROJECT_ROOT}/filemanager"
    LIBRARY_PATH="${PROJECT_ROOT}/rust_ffi/target/release"
else
    echo "❌ Error: FileManager binary not found"
    echo ""
    echo "Please run the build script first:"
    echo "  cd ${SCRIPT_DIR}"
    echo "  bash build.sh"
    exit 1
fi

# Set library path for the Rust FFI library
export LD_LIBRARY_PATH="${LIBRARY_PATH}:${LD_LIBRARY_PATH}"
export DYLD_LIBRARY_PATH="${LIBRARY_PATH}:${DYLD_LIBRARY_PATH}"

# Run the FileManager binary with all arguments passed through
exec "${BINARY_PATH}" "$@"