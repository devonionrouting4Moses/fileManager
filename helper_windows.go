// +build windows

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Result represents the result of an operation
type Result struct {
	Success bool
	Message string
}

func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}

// CreateFolder creates a new directory and all necessary parent directories
func CreateFolder(path string) Result {
	err := os.MkdirAll(filepath.Clean(path), 0755)
	if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}
	return Result{
		Success: true,
		Message: "Folder created successfully",
	}
}

// PrintResult prints the result of an operation
func PrintResult(result Result) {
	if result.Success {
		fmt.Println("✅", result.Message)
	} else {
		fmt.Println("❌", result.Message)
	}
}

// CreateFile creates a new empty file at the specified path
func CreateFile(path string) Result {
	// Create parent directories if they don't exist
	dir := filepath.Dir(path)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return Result{
				Success: false,
				Message: fmt.Sprintf("Failed to create parent directories for '%s': %v", path, err),
			}
		}
	}

	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to create file '%s': %v", path, err),
		}
	}
	file.Close()

	return Result{
		Success: true,
		Message: fmt.Sprintf("Successfully created file '%s'", path),
	}
}

// DeletePath deletes a file or directory at the specified path
func DeletePath(path string) Result {
	// First, check if the path exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Result{
				Success: false,
				Message: fmt.Sprintf("Path '%s' does not exist", path),
			}
		}
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to access path '%s': %v", path, err),
		}
	}

	// Handle directory deletion
	if fileInfo.IsDir() {
		err = os.RemoveAll(path)
		if err != nil {
			return Result{
				Success: false,
				Message: fmt.Sprintf("Failed to delete directory '%s': %v", path, err),
			}
		}
		return Result{
			Success: true,
			Message: fmt.Sprintf("Successfully deleted directory '%s'", path),
		}
	}

	// Handle file deletion
	err = os.Remove(path)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to delete file '%s': %v", path, err),
		}
	}

	return Result{
		Success: true,
		Message: fmt.Sprintf("Successfully deleted file '%s'", path),
	}
}

// ChangePermissions changes the permissions of a file or directory
// Note: On Windows, only the read-only flag is supported
func ChangePermissions(path string, mode uint32) Result {
	// On Windows, we can only really set the read-only flag
	// The mode parameter is mostly ignored except for the read-only bit
	readOnly := (mode & 0222) == 0 // If write bits are not set, it's read-only

	err := os.Chmod(path, os.FileMode(mode))
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to change permissions for '%s': %v", path, err),
		}
	}

	// On Windows, we also need to set the read-only attribute separately
	err = setReadOnly(path, readOnly)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to set read-only attribute for '%s': %v", path, err),
		}
	}

	return Result{
		Success: true,
		Message: fmt.Sprintf("Successfully changed permissions for '%s' (read-only: %v)", path, readOnly),
	}
}

// MovePath moves a file or directory from src to dst
func MovePath(src, dst string) Result {
	// First check if source exists
	_, err := os.Stat(src)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Source path '%s' does not exist: %v", src, err),
		}
	}

	// Create parent directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err = os.MkdirAll(dstDir, 0755)
		if err != nil {
			return Result{
				Success: false,
				Message: fmt.Sprintf("Failed to create destination directory '%s': %v", dstDir, err),
			}
		}
	}

	err = os.Rename(src, dst)
	if err != nil {
		// If rename fails (e.g., across different volumes), try copy+delete
		copyResult := CopyPath(src, dst)
		if !copyResult.Success {
			return copyResult
		}
		
		// If copy succeeded, delete the original
		deleteResult := DeletePath(src)
		if !deleteResult.Success {
			return Result{
				Success: false,
				Message: fmt.Sprintf("Moved file but failed to remove original: %s", deleteResult.Message),
			}
		}
		
		return Result{
			Success: true,
			Message: fmt.Sprintf("Successfully moved '%s' to '%s' (using copy+delete)", src, dst),
		}
	}

	return Result{
		Success: true,
		Message: fmt.Sprintf("Successfully moved '%s' to '%s'", src, dst),
	}
}

// CopyPath copies a file or directory from src to dst
func CopyPath(src, dst string) Result {
	// Get source info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to access source '%s': %v", src, err),
		}
	}

	// If source is a directory, copy it recursively
	if srcInfo.IsDir() {
		return copyDirectory(src, dst)
	}

	// Handle file copy
	return copyFile(src, dst)
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) Result {
	srcFile, err := os.Open(src)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to open source file '%s': %v", src, err),
		}
	}
	defer srcFile.Close()

	// Create parent directories if they don't exist
	dstDir := filepath.Dir(dst)
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err = os.MkdirAll(dstDir, 0755)
		if err != nil {
			return Result{
				Success: false,
				Message: fmt.Sprintf("Failed to create directory '%s': %v", dstDir, err),
			}
		}
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to create destination file '%s': %v", dst, err),
		}
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to copy data from '%s' to '%s': %v", src, dst, err),
		}
	}

	// Preserve file mode
	srcInfo, err := os.Stat(src)
	if err == nil {
		os.Chmod(dst, srcInfo.Mode())
	}

	return Result{
		Success: true,
		Message: fmt.Sprintf("Successfully copied file '%s' to '%s'", src, dst),
	}
}

// copyDirectory copies a directory recursively
func copyDirectory(src, dst string) Result {
	// Create the destination directory
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to create destination directory '%s': %v", dst, err),
		}
	}

	// Read the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to read source directory '%s': %v", src, err),
		}
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			result := copyDirectory(srcPath, dstPath)
			if !result.Success {
				return result
			}
		} else {
			result := copyFile(srcPath, dstPath)
			if !result.Success {
				return result
			}
		}
	}

	return Result{
		Success: true,
		Message: fmt.Sprintf("Successfully copied directory '%s' to '%s'", src, dst),
	}
}

// setReadOnly sets the read-only attribute on a file or directory
func setReadOnly(path string, readOnly bool) error {
	// Convert path to UTF16 for Windows API
	path16, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	// Get current attributes
	attrs, err := syscall.GetFileAttributes(path16)
	if err != nil {
		return err
	}

	// Set or clear the read-only bit
	if readOnly {
		attrs |= syscall.FILE_ATTRIBUTE_READONLY
	} else {
		attrs &^= syscall.FILE_ATTRIBUTE_READONLY
	}

	// Set the new attributes
	return syscall.SetFileAttributes(path16, attrs)
}

// RenamePath renames a file or directory from oldPath to newPath
func RenamePath(oldPath, newPath string) Result {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to rename '%s' to '%s': %v", oldPath, newPath, err),
		}
	}
	return Result{
		Success: true,
		Message: fmt.Sprintf("Successfully renamed '%s' to '%s'", oldPath, newPath),
	}
}