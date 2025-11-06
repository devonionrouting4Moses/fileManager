package main

/*
#cgo LDFLAGS: -L./target/release -lfilemanager -ldl -lpthread -lm
#include <stdlib.h>

typedef struct {
    int success;
    char* message;
} OperationResult;

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

type Result struct {
	Success bool
	Message string
}

func processResult(cResult C.OperationResult) Result {
	defer C.free_result(cResult)
	
	result := Result{
		Success: cResult.success == 1,
		Message: C.GoString(cResult.message),
	}
	
	return result
}

func CreateFolder(path string) Result {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	
	cResult := C.create_folder(cPath)
	return processResult(cResult)
}

func CreateFile(path string) Result {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	
	cResult := C.create_file(cPath)
	return processResult(cResult)
}

func RenamePath(oldPath, newPath string) Result {
	cOldPath := C.CString(oldPath)
	cNewPath := C.CString(newPath)
	defer C.free(unsafe.Pointer(cOldPath))
	defer C.free(unsafe.Pointer(cNewPath))
	
	cResult := C.rename_path(cOldPath, cNewPath)
	return processResult(cResult)
}

func DeletePath(path string) Result {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	
	cResult := C.delete_path(cPath)
	return processResult(cResult)
}

func ChangePermissions(path string, mode uint32) Result {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	
	cResult := C.change_permissions(cPath, C.uint(mode))
	return processResult(cResult)
}

func MovePath(src, dst string) Result {
	cSrc := C.CString(src)
	cDst := C.CString(dst)
	defer C.free(unsafe.Pointer(cSrc))
	defer C.free(unsafe.Pointer(cDst))
	
	cResult := C.move_path(cSrc, cDst)
	return processResult(cResult)
}

func CopyPath(src, dst string) Result {
	cSrc := C.CString(src)
	cDst := C.CString(dst)
	defer C.free(unsafe.Pointer(cSrc))
	defer C.free(unsafe.Pointer(cDst))
	
	cResult := C.copy_path(cSrc, cDst)
	return processResult(cResult)
}

func PrintResult(result Result) {
	if result.Success {
		fmt.Printf("✅ %s\n", result.Message)
	} else {
		fmt.Printf("❌ %s\n", result.Message)
	}
}