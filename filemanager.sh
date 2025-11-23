#!/bin/bash

# FileManager Wrapper Script
# Sets up the environment and runs the FileManager binary with proper library paths

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Set library path for the Rust FFI library
export LD_LIBRARY_PATH="${SCRIPT_DIR}/rust_ffi/target/release:${LD_LIBRARY_PATH}"
export DYLD_LIBRARY_PATH="${SCRIPT_DIR}/rust_ffi/target/release:${DYLD_LIBRARY_PATH}"

# Run the FileManager binary with all arguments passed through
exec "${SCRIPT_DIR}/filemanager" "$@"