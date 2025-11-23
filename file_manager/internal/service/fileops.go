package service

import (
	"filemanager/internal/ffi"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CreateFromTemplate creates a project structure from a template
func CreateFromTemplate(rootPath string, template StructureTemplate) (int, int) {
	successCount := 0
	errorCount := 0

	result := ffi.CreateFolder(rootPath)
	if !result.Success {
		fmt.Printf("❌ Failed to create root: %s\n", result.Message)
		return 0, 1
	}
	successCount++

	for _, dir := range template.Directories {
		fullPath := filepath.Join(rootPath, dir)
		result := ffi.CreateFolder(fullPath)
		if result.Success {
			successCount++
			fmt.Printf("  ✅ 📁 %s\n", fullPath)
		} else {
			errorCount++
			fmt.Printf("  ❌ %s: %s\n", fullPath, result.Message)
		}
	}

	for filePath, content := range template.Files {
		fullPath := filepath.Join(rootPath, filePath)
		result := ffi.CreateFile(fullPath)
		if result.Success {
			if err := os.WriteFile(fullPath, []byte(content), 0644); err == nil {
				successCount++
				fmt.Printf("  ✅ 📄 %s\n", fullPath)
			} else {
				errorCount++
				fmt.Printf("  ❌ %s: %s\n", fullPath, err.Error())
			}
		} else {
			errorCount++
			fmt.Printf("  ❌ %s: %s\n", fullPath, result.Message)
		}
	}

	return successCount, errorCount
}

// ParseTreeStructure parses a tree-like structure from pasted text
func ParseTreeStructure(input string) ([]string, map[string]string, error) {
	lines := strings.Split(input, "\n")
	var dirs []string
	files := make(map[string]string)

	pathStack := []string{}
	rootSet := false

	for _, rawLine := range lines {
		if strings.TrimSpace(rawLine) == "" {
			continue
		}

		indent := 0
		for _, ch := range rawLine {
			if ch == ' ' {
				indent++
			} else if ch == '\t' {
				indent += 4
			} else {
				break
			}
		}

		trimmed := strings.TrimLeft(rawLine, " \t")

		treeDepth := 0
		temp := trimmed
		for {
			if strings.HasPrefix(temp, "│") {
				treeDepth++
				temp = strings.TrimPrefix(temp, "│")
				temp = strings.TrimLeft(temp, " ")
			} else {
				break
			}
		}

		cleaned := trimmed
		for {
			old := cleaned
			cleaned = strings.TrimPrefix(cleaned, "│")
			cleaned = strings.TrimPrefix(cleaned, "├──")
			cleaned = strings.TrimPrefix(cleaned, "├─")
			cleaned = strings.TrimPrefix(cleaned, "└──")
			cleaned = strings.TrimPrefix(cleaned, "└─")
			cleaned = strings.TrimLeft(cleaned, " \t")
			if old == cleaned {
				break
			}
		}

		if cleaned == "" {
			continue
		}

		name := cleaned
		if idx := strings.Index(cleaned, "#"); idx != -1 {
			name = strings.TrimSpace(cleaned[:idx])
		}

		name = strings.TrimRight(name, "/")
		name = strings.TrimSpace(name)

		if name == "" {
			continue
		}

		var depth int
		if !rootSet {
			depth = 0
			rootSet = true
		} else {
			if strings.HasPrefix(trimmed, "├") || strings.HasPrefix(trimmed, "└") {
				depth = 1
			} else {
				depth = treeDepth + 1
			}
		}

		if depth < len(pathStack) {
			pathStack = pathStack[:depth]
		}

		var fullPath string
		if len(pathStack) == 0 {
			fullPath = name
		} else {
			fullPath = filepath.Join(append(pathStack, name)...)
		}

		isFile := hasFileExtension(name)

		if isFile {
			files[fullPath] = ""
		} else {
			dirs = append(dirs, fullPath)
			if depth >= len(pathStack) {
				pathStack = append(pathStack, name)
			} else {
				pathStack = append(pathStack[:depth], name)
			}
		}
	}

	return dirs, files, nil
}

// hasFileExtension checks if a name has a file extension
func hasFileExtension(name string) bool {
	fileExtensions := []string{
		".java", ".go", ".py", ".js", ".ts", ".jsx", ".tsx",
		".c", ".cpp", ".h", ".hpp", ".cs", ".rb", ".php",
		".html", ".css", ".scss", ".sass", ".json", ".xml",
		".yaml", ".yml", ".md", ".txt", ".sh", ".bat", ".ps1",
		".sql", ".rs", ".kt", ".swift", ".m", ".mm", ".r",
		".pl", ".lua", ".dart", ".vue", ".svelte", ".class",
		".exe", ".dll", ".so", ".jar", ".war", ".properties",
		".toml", ".ini", ".conf", ".lock", ".log",
	}

	nameLower := strings.ToLower(name)

	for _, ext := range fileExtensions {
		if strings.HasSuffix(nameLower, ext) {
			return true
		}
	}

	if idx := strings.LastIndex(name, "."); idx > 0 && idx < len(name)-1 {
		ext := name[idx+1:]
		if len(ext) >= 2 && len(ext) <= 10 {
			isAlphaNum := true
			for _, ch := range ext {
				if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
					isAlphaNum = false
					break
				}
			}
			return isAlphaNum
		}
	}

	return false
}

// BatchCreateFiles creates multiple files in batch
func BatchCreateFiles(paths []string) (successCount int, errorCount int) {
	for _, path := range paths {
		// Create parent directories if needed
		dir := filepath.Dir(path)
		if dir != "." && dir != path {
			ffi.CreateFolder(dir)
		}

		result := ffi.CreateFile(path)
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
	}
	return
}

// BatchCreateFolders creates multiple folders in batch
func BatchCreateFolders(paths []string) (successCount int, errorCount int) {
	// Sort by depth to ensure parent directories are created first
	sort.Slice(paths, func(i, j int) bool {
		depthI := strings.Count(paths[i], string(os.PathSeparator))
		depthJ := strings.Count(paths[j], string(os.PathSeparator))
		if depthI != depthJ {
			return depthI < depthJ
		}
		return paths[i] < paths[j]
	})

	for _, path := range paths {
		result := ffi.CreateFolder(path)
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
	}
	return
}