// package main

// import (
// 	"bufio"
// 	"filemanager/internal/ffi"
// 	"filemanager/internal/handler"
// 	"filemanager/internal/service"
// 	"filemanager/pkg/version"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"sort"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// func main() {
// 	scanner := bufio.NewScanner(os.Stdin)

// 	// Handle command-line flags
// 	if len(os.Args) > 1 {
// 		switch os.Args[1] {
// 		case "--version", "-v":
// 			version.ShowVersion()
// 			return
// 		case "--update", "-u":
// 			version.CheckForUpdates()
// 			return
// 		case "--help", "-h":
// 			showHelp()
// 			return
// 		case "--web", "-w":
// 			// Start web server mode directly
// 			handler.StartWebServer()
// 			return
// 		}
// 	}

// 	version.ShowBanner()

// 	// Check for updates on startup (non-blocking)
// 	go func() {
// 		time.Sleep(500 * time.Millisecond)
// 		version.CheckForUpdates()
// 	}()

// 	for {
// 		displayMenu()

// 		fmt.Print("Enter your choice: ")
// 		if !scanner.Scan() {
// 			break
// 		}

// 		choice := strings.TrimSpace(scanner.Text())

// 		switch choice {
// 		case "0":
// 			fmt.Println("\n👋 Goodbye!")
// 			return
// 		case "1":
// 			handleCreateFolder(scanner)
// 		case "2":
// 			handleCreateFile(scanner)
// 		case "3":
// 			handleRename(scanner)
// 		case "4":
// 			handleDelete(scanner)
// 		case "5":
// 			handleChangePermissions(scanner)
// 		case "6":
// 			handleMove(scanner)
// 		case "7":
// 			handleCopy(scanner)
// 		case "8":
// 			handleCreateStructure(scanner)
// 		case "9":
// 			handleWebServerLaunch()
// 			return
// 		default:
// 			fmt.Println("❌ Invalid choice. Please try again.")
// 		}
// 	}
// }

// func showHelp() {
// 	fmt.Printf("%s v%s - Hybrid File Manager\n\n", version.AppName, version.Version)
// 	fmt.Println("Usage:")
// 	fmt.Println("  filemanager              Start interactive mode")
// 	fmt.Println("  filemanager --version    Show version information")
// 	fmt.Println("  filemanager --update     Check for updates")
// 	fmt.Println("  filemanager --help       Show this help message")
// 	fmt.Println("  filemanager --web        Start web interface")
// 	fmt.Println()
// 	fmt.Println("Features:")
// 	fmt.Println("  • Single & batch file/folder operations")
// 	fmt.Println("  • 12 project templates (Flask, Spring Boot, React, etc.)")
// 	fmt.Println("  • Interactive structure builder")
// 	fmt.Println("  • Auto parent directory creation")
// 	fmt.Println("  • Cross-platform support (Linux, macOS, Windows)")
// 	fmt.Println("  • Web interface for browser-based management")
// 	fmt.Println()
// }

// func displayMenu() {
// 	fmt.Println("┌────────────────────────────────────────┐")
// 	fmt.Println("│           Available Operations         │")
// 	fmt.Println("├────────────────────────────────────────┤")
// 	fmt.Println("│  1️⃣  Create Folder                     │")
// 	fmt.Println("│  2️⃣  Create File                       │")
// 	fmt.Println("│  3️⃣  Rename File/Folder                │")
// 	fmt.Println("│  4️⃣  Delete File/Folder                │")
// 	fmt.Println("│  5️⃣  Change Permissions                │")
// 	fmt.Println("│  6️⃣  Move File/Folder                  │")
// 	fmt.Println("│  7️⃣  Copy File/Folder                  │")
// 	fmt.Println("│  8️⃣  Create Structure (Multi-entity)   │")
// 	fmt.Println("│  9️⃣  Launch Web Interface              │")
// 	fmt.Println("│  0️⃣  Exit                              │")
// 	fmt.Println("└────────────────────────────────────────┘")
// 	fmt.Println()
// }

// func handleWebServerLaunch() {
// 	fmt.Println("\n🌐 Launching Web Interface...")
// 	fmt.Println("════════════════════════════════════════")
// 	fmt.Println("Starting HTTP server on http://localhost:8080")
// 	fmt.Println("Press Ctrl+C to stop the server")
// 	fmt.Println()
// 	handler.StartWebServer()
// }

// func handleCreateFolder(scanner *bufio.Scanner) {
// 	fmt.Print("\n📁 Enter folder path(s) - space-separated for multiple folders: ")
// 	if !scanner.Scan() {
// 		return
// 	}

// 	input := strings.TrimSpace(scanner.Text())
// 	if input == "" {
// 		fmt.Println("❌ Path cannot be empty")
// 		return
// 	}

// 	paths := strings.Fields(input)
// 	if len(paths) == 0 {
// 		fmt.Println("❌ No valid paths provided")
// 		return
// 	}

// 	successCount := 0
// 	errorCount := 0

// 	for _, path := range paths {
// 		result := ffi.CreateFolder(path)
// 		if result.Success {
// 			successCount++
// 			fmt.Printf("  ✅ 📁 %s\n", path)
// 		} else {
// 			errorCount++
// 			fmt.Printf("  ❌ %s: %s\n", path, result.Message)
// 		}
// 	}

// 	if len(paths) > 1 {
// 		fmt.Printf("📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
// 	}
// 	fmt.Println()
// }

// func handleCreateFile(scanner *bufio.Scanner) {
// 	fmt.Print("📄 Enter file path(s) - space-separated for multiple files: ")
// 	if !scanner.Scan() {
// 		return
// 	}

// 	input := strings.TrimSpace(scanner.Text())
// 	if input == "" {
// 		fmt.Println("❌ Path cannot be empty")
// 		return
// 	}

// 	paths := strings.Fields(input)
// 	if len(paths) == 0 {
// 		fmt.Println("❌ No valid paths provided")
// 		return
// 	}

// 	successCount := 0
// 	errorCount := 0

// 	for _, path := range paths {
// 		// Validate the path
// 		if err := validateFilePath(path); err != nil {
// 			fmt.Printf("❌ %s: %s\n", path, err.Error())
// 			errorCount++
// 			continue
// 		}

// 		// Create parent directories if needed
// 		dir := filepath.Dir(path)
// 		if dir != "." && dir != path {
// 			result := ffi.CreateFolder(dir)
// 			if !result.Success {
// 				fmt.Printf("❌ Failed to create directory %s: %s\n", dir, result.Message)
// 				errorCount++
// 				continue
// 			}
// 			fmt.Printf("  📁 Created directory: %s\n", dir)
// 		}

// 		result := ffi.CreateFile(path)
// 		if result.Success {
// 			successCount++
// 			fmt.Printf("  ✅ 📄 %s\n", path)
// 		} else {
// 			errorCount++
// 			fmt.Printf("  ❌ %s: %s\n", path, result.Message)
// 		}
// 	}

// 	if len(paths) > 1 {
// 		fmt.Printf("📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
// 	}
// 	fmt.Println()
// }

// func validateFilePath(path string) error {
// 	parts := strings.Split(path, "/")
// 	lastPart := parts[len(parts)-1]

// 	dotCount := strings.Count(lastPart, ".")

// 	if dotCount > 1 {
// 		if !strings.HasPrefix(lastPart, ".") {
// 			extParts := strings.Split(lastPart, ".")
// 			if len(extParts) > 2 {
// 				validDoubleExt := []string{"tar.gz", "tar.bz2", "tar.xz"}
// 				extension := strings.Join(extParts[len(extParts)-2:], ".")
// 				isValid := false
// 				for _, validExt := range validDoubleExt {
// 					if extension == validExt {
// 						isValid = true
// 						break
// 					}
// 				}
// 				if !isValid {
// 					return fmt.Errorf("invalid path - multiple file extensions detected")
// 				}
// 			}
// 		}
// 	}

// 	if !strings.Contains(lastPart, ".") {
// 		return fmt.Errorf("invalid filename - no file extension detected")
// 	}

// 	return nil
// }

// func handleRename(scanner *bufio.Scanner) {
// 	fmt.Print("🔄 Enter current path: ")
// 	if !scanner.Scan() {
// 		return
// 	}
// 	oldPath := strings.TrimSpace(scanner.Text())

// 	fmt.Print("Enter new path: ")
// 	if !scanner.Scan() {
// 		return
// 	}
// 	newPath := strings.TrimSpace(scanner.Text())

// 	if oldPath == "" || newPath == "" {
// 		fmt.Println("❌ Paths cannot be empty")
// 		return
// 	}

// 	result := ffi.RenamePath(oldPath, newPath)
// 	ffi.PrintResult(result)
// 	fmt.Println()
// }

// func handleDelete(scanner *bufio.Scanner) {
// 	fmt.Print("🗑️  Enter path to delete: ")
// 	if !scanner.Scan() {
// 		return
// 	}

// 	path := strings.TrimSpace(scanner.Text())
// 	if path == "" {
// 		fmt.Println("❌ Path cannot be empty")
// 		return
// 	}

// 	fmt.Printf("⚠️  Are you sure you want to delete '%s'? (yes/no): ", path)
// 	if !scanner.Scan() {
// 		return
// 	}

// 	confirmation := strings.ToLower(strings.TrimSpace(scanner.Text()))
// 	if confirmation != "yes" && confirmation != "y" {
// 		fmt.Println("❌ Deletion cancelled")
// 		return
// 	}

// 	result := ffi.DeletePath(path)
// 	ffi.PrintResult(result)
// 	fmt.Println()
// }

// func handleChangePermissions(scanner *bufio.Scanner) {
// 	fmt.Print("🔒 Enter path: ")
// 	if !scanner.Scan() {
// 		return
// 	}
// 	path := strings.TrimSpace(scanner.Text())

// 	fmt.Print("Enter permissions (octal, e.g., 755): ")
// 	if !scanner.Scan() {
// 		return
// 	}
// 	modeStr := strings.TrimSpace(scanner.Text())

// 	if path == "" || modeStr == "" {
// 		fmt.Println("❌ Path and permissions cannot be empty")
// 		return
// 	}

// 	mode, err := strconv.ParseUint(modeStr, 8, 32)
// 	if err != nil {
// 		fmt.Printf("❌ Invalid permission format: %v", err)
// 		return
// 	}

// 	result := ffi.ChangePermissions(path, uint32(mode))
// 	ffi.PrintResult(result)
// 	fmt.Println()
// }

// func handleMove(scanner *bufio.Scanner) {
// 	fmt.Print("➡️  Enter source path: ")
// 	if !scanner.Scan() {
// 		return
// 	}
// 	src := strings.TrimSpace(scanner.Text())

// 	fmt.Print("Enter destination path: ")
// 	if !scanner.Scan() {
// 		return
// 	}
// 	dst := strings.TrimSpace(scanner.Text())

// 	if src == "" || dst == "" {
// 		fmt.Println("❌ Paths cannot be empty")
// 		return
// 	}

// 	result := ffi.MovePath(src, dst)
// 	ffi.PrintResult(result)
// 	fmt.Println()
// }

// func handleCopy(scanner *bufio.Scanner) {
// 	fmt.Print("📋 Enter source path: ")
// 	if !scanner.Scan() {
// 		return
// 	}
// 	src := strings.TrimSpace(scanner.Text())

// 	fmt.Print("Enter destination path: ")
// 	if !scanner.Scan() {
// 		return
// 	}
// 	dst := strings.TrimSpace(scanner.Text())

// 	if src == "" || dst == "" {
// 		fmt.Println("❌ Paths cannot be empty")
// 		return
// 	}

// 	result := ffi.CopyPath(src, dst)
// 	ffi.PrintResult(result)
// 	fmt.Println()
// }

// func handleCreateStructure(scanner *bufio.Scanner) {
// 	for {
// 		fmt.Println("\n🏗️  Create Hierarchical Structure")
// 		fmt.Println("════════════════════════════════════════")
// 		fmt.Println()

// 		templates := service.GetAvailableTemplates()

// 		fmt.Println("📦 Available Templates:")
// 		for i, t := range templates {
// 			fmt.Printf("  %2d. %s\n", i+1, t.Description)
// 		}

// 		fmt.Printf("\n  %2d. Custom Structure (from definition)\n", len(templates)+1)
// 		fmt.Printf("  %2d. Parse Tree Structure (paste)\n", len(templates)+2)
// 		fmt.Printf("  %2d. Interactive Builder\n", len(templates)+3)
// 		fmt.Printf("  %2d. ← Back to Main Menu\n", len(templates)+4)
// 		fmt.Println()
// 		fmt.Print("Select option: ")

// 		if !scanner.Scan() {
// 			return
// 		}

// 		option := strings.TrimSpace(scanner.Text())
// 		optionNum, err := strconv.Atoi(option)

// 		if err != nil {
// 			fmt.Println("❌ Invalid input. Please enter a number.")
// 			continue
// 		}

// 		if optionNum == len(templates)+4 {
// 			return
// 		}

// 		if optionNum >= 1 && optionNum <= len(templates) {
// 			if !handleTemplateStructure(scanner, templates[optionNum-1]) {
// 				continue
// 			}
// 			return
// 		}

// 		switch optionNum {
// 		case len(templates) + 1:
// 			if !handleCustomStructure(scanner) {
// 				continue
// 			}
// 			return
// 		case len(templates) + 2:
// 			if !handleParseTreeStructure(scanner) {
// 				continue
// 			}
// 			return
// 		case len(templates) + 3:
// 			if !handleInteractiveBuilder(scanner) {
// 				continue
// 			}
// 			return
// 		default:
// 			fmt.Println("❌ Invalid option")
// 		}
// 	}
// }

// func handleTemplateStructure(scanner *bufio.Scanner, template service.StructureTemplate) bool {
// 	for {
// 		fmt.Printf("\n📋 Template: %s\n", template.Description)
// 		fmt.Println("────────────────────────────────────────")
// 		fmt.Print("\n📁 Enter root directory name (or 'back' to return): ")

// 		if !scanner.Scan() {
// 			return false
// 		}

// 		input := strings.TrimSpace(scanner.Text())

// 		if strings.ToLower(input) == "back" || strings.ToLower(input) == "b" {
// 			return false
// 		}

// 		if input == "" {
// 			fmt.Println("❌ Path cannot be empty")
// 			continue
// 		}

// 		fmt.Printf("\n🔨 Creating %s structure...\n", template.Name)

// 		successCount, errorCount := service.CreateFromTemplate(input, template)

// 		fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)

// 		if errorCount == 0 {
// 			fmt.Println("✨ Structure created successfully!")
// 		}

// 		return true
// 	}
// }

// // handleCustomStructure creates a custom structure from user-defined format
// // Format: d:path/to/dir or f:path/to/file on each line
// func handleCustomStructure(scanner *bufio.Scanner) bool {
// 	fmt.Println("\n📝 Custom Structure Definition")
// 	fmt.Println("════════════════════════════════════════")
// 	fmt.Println("Format: d:path/to/dir or f:path/to/file")
// 	fmt.Println("Example:")
// 	fmt.Println("  d:src")
// 	fmt.Println("  d:src/api")
// 	fmt.Println("  f:src/api/users.go")
// 	fmt.Println("  f:README.md")
// 	fmt.Println("Enter 'done' when finished (or 'back' to cancel)")
// 	fmt.Println()

// 	var lines []string
// 	for {
// 		fmt.Print("> ")
// 		if !scanner.Scan() {
// 			return false
// 		}

// 		input := strings.TrimSpace(scanner.Text())

// 		if strings.ToLower(input) == "done" {
// 			break
// 		}
// 		if strings.ToLower(input) == "back" || strings.ToLower(input) == "b" {
// 			return false
// 		}

// 		if input != "" {
// 			lines = append(lines, input)
// 		}
// 	}

// 	if len(lines) == 0 {
// 		fmt.Println("❌ No items defined")
// 		return false
// 	}

// 	successCount := 0
// 	errorCount := 0

// 	for _, line := range lines {
// 		line = strings.TrimSpace(line)
// 		if line == "" {
// 			continue
// 		}

// 		var result ffi.Result
// 		if strings.HasPrefix(line, "d:") {
// 			path := strings.TrimPrefix(line, "d:")
// 			path = strings.TrimSpace(path)
// 			result = ffi.CreateFolder(path)
// 			if result.Success {
// 				successCount++
// 				fmt.Printf("  ✅ 📁 %s\n", path)
// 			} else {
// 				errorCount++
// 				fmt.Printf("  ❌ %s: %s\n", path, result.Message)
// 			}
// 		} else if strings.HasPrefix(line, "f:") {
// 			path := strings.TrimPrefix(line, "f:")
// 			path = strings.TrimSpace(path)
// 			result = ffi.CreateFile(path)
// 			if result.Success {
// 				successCount++
// 				fmt.Printf("  ✅ 📄 %s\n", path)
// 			} else {
// 				errorCount++
// 				fmt.Printf("  ❌ %s: %s\n", path, result.Message)
// 			}
// 		}
// 	}

// 	fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
// 	return true
// }

// // handleParseTreeStructure creates a structure from pasted tree format
// func handleParseTreeStructure(scanner *bufio.Scanner) bool {
// 	fmt.Println("\n🌳 Parse Tree Structure")
// 	fmt.Println("════════════════════════════════════════")
// 	fmt.Println("Paste your tree structure (supports ├──, └──, │ characters)")
// 	fmt.Println("Example:")
// 	fmt.Println("  myapp/")
// 	fmt.Println("  ├── src/")
// 	fmt.Println("  │   ├── main.go")
// 	fmt.Println("  │   └── utils.go")
// 	fmt.Println("  ├── tests/")
// 	fmt.Println("  └── README.md")
// 	fmt.Println("Enter 'done' when finished (or 'back' to cancel)")
// 	fmt.Println()

// 	var lines []string
// 	for {
// 		fmt.Print("> ")
// 		if !scanner.Scan() {
// 			return false
// 		}

// 		input := scanner.Text()

// 		if strings.ToLower(strings.TrimSpace(input)) == "done" {
// 			break
// 		}
// 		if strings.ToLower(strings.TrimSpace(input)) == "back" || strings.ToLower(strings.TrimSpace(input)) == "b" {
// 			return false
// 		}

// 		lines = append(lines, input)
// 	}

// 	if len(lines) == 0 {
// 		fmt.Println("❌ No structure provided")
// 		return false
// 	}

// 	treeInput := strings.Join(lines, "\n")
// 	dirs, files, err := service.ParseTreeStructure(treeInput)
// 	if err != nil {
// 		fmt.Printf("❌ Error parsing structure: %v\n", err)
// 		return false
// 	}

// 	successCount := 0
// 	errorCount := 0

// 	// Sort directories by depth to create parent dirs first
// 	sort.Slice(dirs, func(i, j int) bool {
// 		depthI := strings.Count(dirs[i], string(filepath.Separator))
// 		depthJ := strings.Count(dirs[j], string(filepath.Separator))
// 		if depthI != depthJ {
// 			return depthI < depthJ
// 		}
// 		return dirs[i] < dirs[j]
// 	})

// 	// Create directories
// 	for _, dir := range dirs {
// 		result := ffi.CreateFolder(dir)
// 		if result.Success {
// 			successCount++
// 			fmt.Printf("  ✅ 📁 %s\n", dir)
// 		} else {
// 			errorCount++
// 			fmt.Printf("  ❌ %s: %s\n", dir, result.Message)
// 		}
// 	}

// 	// Create files
// 	for filePath := range files {
// 		result := ffi.CreateFile(filePath)
// 		if result.Success {
// 			successCount++
// 			fmt.Printf("  ✅ 📄 %s\n", filePath)
// 		} else {
// 			errorCount++
// 			fmt.Printf("  ❌ %s: %s\n", filePath, result.Message)
// 		}
// 	}

// 	fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
// 	return true
// }

// // handleInteractiveBuilder provides an interactive directory navigation builder
// func handleInteractiveBuilder(scanner *bufio.Scanner) bool {
// 	fmt.Println("\n🏗️  Interactive Structure Builder")
// 	fmt.Println("════════════════════════════════════════")
// 	fmt.Print("\n📁 Enter root directory name: ")

// 	if !scanner.Scan() {
// 		return false
// 	}

// 	rootDir := strings.TrimSpace(scanner.Text())
// 	if rootDir == "" {
// 		fmt.Println("❌ Root directory name cannot be empty")
// 		return false
// 	}

// 	result := ffi.CreateFolder(rootDir)
// 	if !result.Success {
// 		fmt.Printf("❌ Failed to create root: %s\n", result.Message)
// 		return false
// 	}

// 	fmt.Printf("✅ Root directory created: %s\n", rootDir)

// 	successCount := 1
// 	errorCount := 0

// 	// Interactive builder loop
// 	currentPath := rootDir
// 	for {
// 		fmt.Printf("\n📂 Current: %s\n", currentPath)
// 		fmt.Println("────────────────────────────────────────")
// 		fmt.Println("Commands:")
// 		fmt.Println("  mkdir <name>  - Create directory")
// 		fmt.Println("  touch <name>  - Create file")
// 		fmt.Println("  cd <name>     - Enter subdirectory")
// 		fmt.Println("  back          - Go to parent directory")
// 		fmt.Println("  list          - List subdirectories")
// 		fmt.Println("  done          - Finish building")
// 		fmt.Println()

// 		fmt.Print("> ")
// 		if !scanner.Scan() {
// 			break
// 		}

// 		input := strings.TrimSpace(scanner.Text())
// 		if input == "" {
// 			continue
// 		}

// 		parts := strings.Fields(input)
// 		command := parts[0]

// 		switch command {
// 		case "mkdir":
// 			if len(parts) < 2 {
// 				fmt.Println("❌ Usage: mkdir <name>")
// 				continue
// 			}
// 			dirName := strings.Join(parts[1:], " ")
// 			newPath := filepath.Join(currentPath, dirName)
// 			result := ffi.CreateFolder(newPath)
// 			if result.Success {
// 				successCount++
// 				fmt.Printf("✅ 📁 %s\n", newPath)
// 			} else {
// 				errorCount++
// 				fmt.Printf("❌ %s\n", result.Message)
// 			}

// 		case "touch":
// 			if len(parts) < 2 {
// 				fmt.Println("❌ Usage: touch <name>")
// 				continue
// 			}
// 			fileName := strings.Join(parts[1:], " ")
// 			newPath := filepath.Join(currentPath, fileName)
// 			result := ffi.CreateFile(newPath)
// 			if result.Success {
// 				successCount++
// 				fmt.Printf("✅ 📄 %s\n", newPath)
// 			} else {
// 				errorCount++
// 				fmt.Printf("❌ %s\n", result.Message)
// 			}

// 		case "cd":
// 			if len(parts) < 2 {
// 				fmt.Println("❌ Usage: cd <name>")
// 				continue
// 			}
// 			dirName := parts[1]
// 			if dirName == ".." {
// 				parent := filepath.Dir(currentPath)
// 				if parent != currentPath {
// 					currentPath = parent
// 					fmt.Printf("↩️  Moved to: %s\n", currentPath)
// 				} else {
// 					fmt.Println("❌ Already at root")
// 				}
// 			} else {
// 				newPath := filepath.Join(currentPath, dirName)
// 				currentPath = newPath
// 				fmt.Printf("📂 Entering: %s\n", currentPath)
// 			}

// 		case "back":
// 			parent := filepath.Dir(currentPath)
// 			if parent != currentPath {
// 				currentPath = parent
// 				fmt.Printf("↩️  Moved to: %s\n", currentPath)
// 			} else {
// 				fmt.Println("❌ Already at root")
// 			}

// 		case "list":
// 			fmt.Println("📁 Current directory structure:")
// 			fmt.Printf("  %s/\n", currentPath)

// 		case "done":
// 			fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
// 			return true

// 		default:
// 			fmt.Println("❌ Unknown command. Type 'done' to finish.")
// 		}
// 	}

// 	fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
// 	return true
// }