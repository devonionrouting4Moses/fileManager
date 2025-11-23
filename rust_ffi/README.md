# Rust FFI File Operations

A modular Rust library for file system operations with C FFI bindings and a CLI tool.

## Project Structure

```
rust_ffi/
├── Cargo.toml                 # Workspace configuration
├── README.md
├── benches/                   # Benchmarks (future)
├── examples/                  # Usage examples (future)
└── crates/
    ├── core/                  # Core library with FFI bindings
    │   ├── Cargo.toml
    │   └── src/
    │       ├── lib.rs         # Main library entry point
    │       ├── common/        # Shared types and utilities
    │       │   └── mod.rs
    │       ├── ffi/           # C FFI wrapper functions
    │       │   └── mod.rs
    │       └── operations/    # File system operations
    │           ├── mod.rs
    │           ├── copy.rs
    │           ├── create.rs
    │           ├── delete.rs
    │           ├── file_permissions.rs
    │           ├── move_ops.rs
    │           └── rename.rs
    └── cli/                   # Command-line interface
        ├── Cargo.toml
        └── src/
            └── main.rs
```

## Features

- **Create**: Create files and folders
- **Delete**: Remove files and directories
- **Rename**: Rename files and directories
- **Move**: Move files and directories
- **Copy**: Copy files and directories (including recursive)
- **Permissions**: Change file permissions (Unix only)

## Building

### Build the entire workspace
```bash
cargo build --release
```

### Build only the core library
```bash
cargo build -p fs-operations-core --release
```

### Build the CLI tool
```bash
cargo build -p fs-operations-cli --release
```

## Usage

### As a CLI Tool

```bash
# Create a folder
./target/release/fsops create-folder /path/to/folder

# Create a file
./target/release/fsops create-file /path/to/file.txt

# Copy a file or folder
./target/release/fsops copy /source /destination

# Move a file or folder
./target/release/fsops move /source /destination

# Rename
./target/release/fsops rename /old/path /new/path

# Delete
./target/release/fsops delete /path/to/delete

# Change permissions (Unix only, octal mode)
./target/release/fsops chmod /path/to/file 755
```

### As a Rust Library

```rust
use fs_operations_core::*;

fn main() -> FsResult<()> {
    // Create a folder
    create_folder("/tmp/my_folder")?;
    
    // Create a file
    create_file("/tmp/my_folder/file.txt")?;
    
    // Copy
    copy_operation("/tmp/my_folder", "/tmp/backup")?;
    
    // Change permissions (Unix)
    change_perms("/tmp/my_folder/file.txt", 0o644)?;
    
    Ok(())
}
```

### From C/C++ via FFI

```c
#include <stdio.h>

// Declare the FFI functions
typedef struct {
    int success;
    char* message;
} OperationResult;

extern OperationResult create_folder(const char* path);
extern void free_result(OperationResult result);

int main() {
    OperationResult result = create_folder("/tmp/test_folder");
    
    if (result.success) {
        printf("Success: %s\n", result.message);
    } else {
        printf("Error: %s\n", result.message);
    }
    
    free_result(result);
    return 0;
}
```

Compile with:
```bash
gcc -o example example.c -L./target/release -lfs_operations_core
```

## Testing

Run all tests:
```bash
cargo test
```

Run tests for a specific crate:
```bash
cargo test -p fs-operations-core
```

## Architecture

### Core Library (`crates/core`)

- **common/**: Shared types, error handling, and FFI utilities
- **operations/**: Pure Rust implementations of file operations
  - Each operation is in its own module for clarity
  - All functions return `FsResult<String>` for consistent error handling
- **ffi/**: C-compatible wrapper functions
  - Converts C strings to Rust strings
  - Calls internal operations
  - Returns `OperationResult` for C compatibility

### CLI (`crates/cli`)

- Built with `clap` for modern CLI parsing
- Provides user-friendly interface to all operations
- Uses the core library directly (not through FFI)

## Design Principles

1. **Separation of Concerns**: Operations, FFI, and CLI are separate
2. **Error Handling**: Custom error types with proper propagation
3. **Modularity**: Each operation in its own file
4. **Testability**: Unit tests for each module
5. **Safety**: Proper memory management in FFI layer
6. **Cross-platform**: Works on Unix and Windows (with platform-specific features)

## License

MIT OR Apache-2.0

## Contributing

Contributions are welcome! Please ensure all tests pass before submitting PRs.