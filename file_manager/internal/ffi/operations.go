package ffi

/*
#cgo LDFLAGS: -L../../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm
#include <stdlib.h>

typedef struct {
    int success;
    char* message;
} OperationResult;

// FFI function declarations
OperationResult create_folder(const char* path);
OperationResult create_file(const char* path);
OperationResult rename_path(const char* old_path, const char* new_path);
OperationResult delete_path(const char* path);
OperationResult change_permissions(const char* path, unsigned int mode);
OperationResult move_path(const char* src, const char* dst);
OperationResult copy_path(const char* src, const char* dst);
void free_result(OperationResult result);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Result represents the outcome of a file operation
type Result struct {
	Success bool
	Message string
}

// processResult converts a C OperationResult to a Go Result
// and properly frees the C memory
func processResult(cResult C.OperationResult) Result {
	defer C.free_result(cResult)

	result := Result{
		Success: cResult.success == 1,
		Message: C.GoString(cResult.message),
	}

	return result
}

// CreateFolder creates a new folder at the specified path
// Creates all parent directories if they don't exist
func CreateFolder(path string) Result {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	cResult := C.create_folder(cPath)
	return processResult(cResult)
}

// CreateFile creates a new empty file at the specified path
func CreateFile(path string) Result {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	cResult := C.create_file(cPath)
	return processResult(cResult)
}

// RenamePath renames a file or folder from oldPath to newPath
func RenamePath(oldPath, newPath string) Result {
	cOldPath := C.CString(oldPath)
	cNewPath := C.CString(newPath)
	defer C.free(unsafe.Pointer(cOldPath))
	defer C.free(unsafe.Pointer(cNewPath))

	cResult := C.rename_path(cOldPath, cNewPath)
	return processResult(cResult)
}

// DeletePath deletes a file or folder at the specified path
// Recursively deletes directories and their contents
func DeletePath(path string) Result {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	cResult := C.delete_path(cPath)
	return processResult(cResult)
}

// ChangePermissions changes file or directory permissions (Unix only)
// mode should be an octal value like 0755
func ChangePermissions(path string, mode uint32) Result {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	cResult := C.change_permissions(cPath, C.uint(mode))
	return processResult(cResult)
}

// MovePath moves a file or folder from src to dst
func MovePath(src, dst string) Result {
	cSrc := C.CString(src)
	cDst := C.CString(dst)
	defer C.free(unsafe.Pointer(cSrc))
	defer C.free(unsafe.Pointer(cDst))

	cResult := C.move_path(cSrc, cDst)
	return processResult(cResult)
}

// CopyPath copies a file or folder from src to dst
// Recursively copies directories and their contents
func CopyPath(src, dst string) Result {
	cSrc := C.CString(src)
	cDst := C.CString(dst)
	defer C.free(unsafe.Pointer(cSrc))
	defer C.free(unsafe.Pointer(cDst))

	cResult := C.copy_path(cSrc, cDst)
	return processResult(cResult)
}

// PrintResult prints a formatted result message
func PrintResult(result Result) {
	if result.Success {
		fmt.Printf("✅ %s\n", result.Message)
	} else {
		fmt.Printf("❌ %s\n", result.Message)
	}
}

// BatchOperation represents a batch of file operations
type BatchOperation struct {
	operations []func() Result
	results    []Result
}

// NewBatchOperation creates a new batch operation handler
func NewBatchOperation() *BatchOperation {
	return &BatchOperation{
		operations: make([]func() Result, 0),
		results:    make([]Result, 0),
	}
}

// AddCreateFolder adds a create folder operation to the batch
func (b *BatchOperation) AddCreateFolder(path string) {
	b.operations = append(b.operations, func() Result {
		return CreateFolder(path)
	})
}

// AddCreateFile adds a create file operation to the batch
func (b *BatchOperation) AddCreateFile(path string) {
	b.operations = append(b.operations, func() Result {
		return CreateFile(path)
	})
}

// Execute runs all batched operations and returns the results
func (b *BatchOperation) Execute() []Result {
	b.results = make([]Result, len(b.operations))
	for i, op := range b.operations {
		b.results[i] = op()
	}
	return b.results
}

// GetSummary returns success and failure counts
func (b *BatchOperation) GetSummary() (success int, failed int) {
	for _, result := range b.results {
		if result.Success {
			success++
		} else {
			failed++
		}
	}
	return
}