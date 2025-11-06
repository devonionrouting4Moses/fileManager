package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	// Handle command-line flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			ShowVersion()
			return
		case "--update", "-u":
			CheckForUpdates()
			return
		case "--help", "-h":
			showHelp()
			return
		}
	}
	
	ShowBanner()
	
	// Check for updates on startup (non-blocking)
	go func() {
		time.Sleep(500 * time.Millisecond)
		CheckForUpdates()
	}()
	
	for {
		displayMenu()
		
		fmt.Print("Enter your choice: ")
		if !scanner.Scan() {
			break
		}
		
		choice := strings.TrimSpace(scanner.Text())
		
		switch choice {
		case "1":
			handleCreateFolder(scanner)
		case "2":
			handleCreateFile(scanner)
		case "3":
			handleRename(scanner)
		case "4":
			handleDelete(scanner)
		case "5":
			handleChangePermissions(scanner)
		case "6":
			handleMove(scanner)
		case "7":
			handleCopy(scanner)
		case "8":
			handleCreateStructure(scanner)
		case "9":
			fmt.Println("\n👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid choice. Please try again.\n")
		}
	}
}

func showHelp() {
	fmt.Printf("%s v%s - Hybrid File Manager\n\n", AppName, Version)
	fmt.Println("Usage:")
	fmt.Println("  filemanager              Start interactive mode")
	fmt.Println("  filemanager --version    Show version information")
	fmt.Println("  filemanager --update     Check for updates")
	fmt.Println("  filemanager --help       Show this help message")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  • Single & batch file/folder operations")
	fmt.Println("  • 12 project templates (Flask, Spring Boot, React, etc.)")
	fmt.Println("  • Interactive structure builder")
	fmt.Println("  • Auto parent directory creation")
	fmt.Println("  • Cross-platform support (Linux, macOS, Windows)")
}

func displayMenu() {
	fmt.Println("┌────────────────────────────────────────┐")
	fmt.Println("│           Available Operations         │")
	fmt.Println("├────────────────────────────────────────┤")
	fmt.Println("│  1️⃣  Create Folder                     │")
	fmt.Println("│  2️⃣  Create File                       │")
	fmt.Println("│  3️⃣  Rename File/Folder                │")
	fmt.Println("│  4️⃣  Delete File/Folder                │")
	fmt.Println("│  5️⃣  Change Permissions                │")
	fmt.Println("│  6️⃣  Move File/Folder                  │")
	fmt.Println("│  7️⃣  Copy File/Folder                  │")
	fmt.Println("│  8️⃣  Create Structure (Multi-entity)   │")
	fmt.Println("│  9️⃣  Exit                              │")
	fmt.Println("└────────────────────────────────────────┘")
	fmt.Println()
}

func handleCreateFolder(scanner *bufio.Scanner) {
	fmt.Print("\n📁 Enter folder path(s) - space-separated for multiple folders: ")
	if !scanner.Scan() {
		return
	}
	
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		fmt.Println("❌ Path cannot be empty\n")
		return
	}
	
	// Split by spaces to handle multiple folders
	paths := strings.Fields(input)
	
	if len(paths) == 0 {
		fmt.Println("❌ No valid paths provided\n")
		return
	}
	
	successCount := 0
	errorCount := 0
	
	for _, path := range paths {
		result := CreateFolder(path)
		if result.Success {
			successCount++
			fmt.Printf("  ✅ 📁 %s\n", path)
		} else {
			errorCount++
			fmt.Printf("  ❌ %s: %s\n", path, result.Message)
		}
	}
	
	if len(paths) > 1 {
		fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
	}
	fmt.Println()
}

func handleCreateFile(scanner *bufio.Scanner) {
	fmt.Print("\n📄 Enter file path(s) - space-separated for multiple files: ")
	if !scanner.Scan() {
		return
	}
	
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		fmt.Println("❌ Path cannot be empty\n")
		return
	}
	
	// Split by spaces to handle multiple files
	paths := strings.Fields(input)
	
	if len(paths) == 0 {
		fmt.Println("❌ No valid paths provided\n")
		return
	}
	
	successCount := 0
	errorCount := 0
	
	for _, path := range paths {
		// Validate the path
		if err := validateFilePath(path); err != nil {
			fmt.Printf("❌ %s: %s\n", path, err.Error())
			errorCount++
			continue
		}
		
		// Check if path contains directory structure
		dir := filepath.Dir(path)
		if dir != "." && dir != path {
			// Create parent directories if they don't exist
			result := CreateFolder(dir)
			if !result.Success {
				fmt.Printf("❌ Failed to create directory %s: %s\n", dir, result.Message)
				errorCount++
				continue
			}
			fmt.Printf("  📁 Created directory: %s\n", dir)
		}
		
		// Create the file
		result := CreateFile(path)
		if result.Success {
			successCount++
			fmt.Printf("  ✅ 📄 %s\n", path)
		} else {
			errorCount++
			fmt.Printf("  ❌ %s: %s\n", path, result.Message)
		}
	}
	
	if len(paths) > 1 {
		fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
	}
	fmt.Println()
}

func validateFilePath(path string) error {
	// Count file extensions (dots)
	parts := strings.Split(path, "/")
	lastPart := parts[len(parts)-1]
	
	// Count dots in the filename
	dotCount := strings.Count(lastPart, ".")
	
	// Check for multiple extensions (error case)
	if dotCount > 1 {
		// Check if it's a valid case like .tar.gz or hidden file
		if !strings.HasPrefix(lastPart, ".") {
			// Split by dots and check if multiple extensions exist
			extParts := strings.Split(lastPart, ".")
			if len(extParts) > 2 {
				// Valid cases: file.tar.gz, archive.tar.bz2
				validDoubleExt := []string{"tar.gz", "tar.bz2", "tar.xz"}
				extension := strings.Join(extParts[len(extParts)-2:], ".")
				isValid := false
				for _, validExt := range validDoubleExt {
					if extension == validExt {
						isValid = true
						break
					}
				}
				if !isValid {
					return fmt.Errorf("invalid path - multiple file extensions detected. Did you mean to create multiple files? Use spaces to separate file paths")
				}
			}
		}
	}
	
	// Check if filename has an extension
	if !strings.Contains(lastPart, ".") {
		return fmt.Errorf("invalid filename - no file extension detected. Files should have extensions (e.g., .txt, .html)")
	}
	
	return nil
}

func handleRename(scanner *bufio.Scanner) {
	fmt.Print("\n🔄 Enter current path: ")
	if !scanner.Scan() {
		return
	}
	oldPath := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Enter new path: ")
	if !scanner.Scan() {
		return
	}
	newPath := strings.TrimSpace(scanner.Text())
	
	if oldPath == "" || newPath == "" {
		fmt.Println("❌ Paths cannot be empty\n")
		return
	}
	
	result := RenamePath(oldPath, newPath)
	PrintResult(result)
	fmt.Println()
}

func handleDelete(scanner *bufio.Scanner) {
	fmt.Print("\n🗑️  Enter path to delete: ")
	if !scanner.Scan() {
		return
	}
	
	path := strings.TrimSpace(scanner.Text())
	if path == "" {
		fmt.Println("❌ Path cannot be empty\n")
		return
	}
	
	fmt.Printf("⚠️  Are you sure you want to delete '%s'? (yes/no): ", path)
	if !scanner.Scan() {
		return
	}
	
	confirmation := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirmation != "yes" && confirmation != "y" {
		fmt.Println("❌ Deletion cancelled\n")
		return
	}
	
	result := DeletePath(path)
	PrintResult(result)
	fmt.Println()
}

func handleChangePermissions(scanner *bufio.Scanner) {
	fmt.Print("\n🔒 Enter path: ")
	if !scanner.Scan() {
		return
	}
	path := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Enter permissions (octal, e.g., 755): ")
	if !scanner.Scan() {
		return
	}
	modeStr := strings.TrimSpace(scanner.Text())
	
	if path == "" || modeStr == "" {
		fmt.Println("❌ Path and permissions cannot be empty\n")
		return
	}
	
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		fmt.Printf("❌ Invalid permission format: %v\n\n", err)
		return
	}
	
	result := ChangePermissions(path, uint32(mode))
	PrintResult(result)
	fmt.Println()
}

func handleMove(scanner *bufio.Scanner) {
	fmt.Print("\n➡️  Enter source path: ")
	if !scanner.Scan() {
		return
	}
	src := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Enter destination path: ")
	if !scanner.Scan() {
		return
	}
	dst := strings.TrimSpace(scanner.Text())
	
	if src == "" || dst == "" {
		fmt.Println("❌ Paths cannot be empty\n")
		return
	}
	
	result := MovePath(src, dst)
	PrintResult(result)
	fmt.Println()
}

func handleCopy(scanner *bufio.Scanner) {
	fmt.Print("\n📋 Enter source path: ")
	if !scanner.Scan() {
		return
	}
	src := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Enter destination path: ")
	if !scanner.Scan() {
		return
	}
	dst := strings.TrimSpace(scanner.Text())
	
	if src == "" || dst == "" {
		fmt.Println("❌ Paths cannot be empty\n")
		return
	}
	
	result := CopyPath(src, dst)
	PrintResult(result)
	fmt.Println()
}

func handleCreateStructure(scanner *bufio.Scanner) {
	for {
		fmt.Println("\n🏗️  Create Hierarchical Structure")
		fmt.Println("════════════════════════════════════════")
		fmt.Println()
		
		templates := GetAvailableTemplates()
		
		fmt.Println("📦 Available Templates:")
		for i, t := range templates {
			fmt.Printf("  %2d. %s\n", i+1, t.Description)
		}
		
		fmt.Printf("\n  %2d. Custom Structure (from definition)\n", len(templates)+1)
		fmt.Printf("  %2d. Parse Tree Structure (paste)\n", len(templates)+2)
		fmt.Printf("  %2d. Interactive Builder\n", len(templates)+3)
		fmt.Printf("  %2d. ← Back to Main Menu\n", len(templates)+4)
		fmt.Println()
		fmt.Print("Select option: ")
		
		if !scanner.Scan() {
			return
		}
		
		option := strings.TrimSpace(scanner.Text())
		optionNum, err := strconv.Atoi(option)
		
		if err != nil {
			fmt.Println("❌ Invalid input. Please enter a number.\n")
			continue
		}
		
		// Back to main menu
		if optionNum == len(templates)+4 {
			return
		}
		
		// Template selection
		if optionNum >= 1 && optionNum <= len(templates) {
			if !handleTemplateStructure(scanner, templates[optionNum-1]) {
				continue // Go back to structure menu if user wants to go back
			}
			return
		}
		
		// Custom/Parse/Interactive options
		switch optionNum {
		case len(templates) + 1:
			if !handleCustomStructure(scanner) {
				continue
			}
			return
		case len(templates) + 2:
			if !handleParseTreeStructure(scanner) {
				continue
			}
			return
		case len(templates) + 3:
			if !handleInteractiveBuilder(scanner) {
				continue
			}
			return
		default:
			fmt.Println("❌ Invalid option\n")
		}
	}
}

func handleTemplateStructure(scanner *bufio.Scanner, template StructureTemplate) bool {
	for {
		fmt.Printf("\n📋 Template: %s\n", template.Description)
		fmt.Println("────────────────────────────────────────")
		fmt.Print("\n📁 Enter root directory name (or 'back' to return): ")
		
		if !scanner.Scan() {
			return false
		}
		
		input := strings.TrimSpace(scanner.Text())
		
		if strings.ToLower(input) == "back" || strings.ToLower(input) == "b" {
			return false
		}
		
		if input == "" {
			fmt.Println("❌ Path cannot be empty")
			continue
		}
		
		fmt.Printf("\n🔨 Creating %s structure...\n", template.Name)
		
		successCount, errorCount := CreateFromTemplate(input, template)
		
		fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
		
		if errorCount == 0 {
			fmt.Println("✨ Structure created successfully!\n")
		}
		
		return true
	}
}

func handleCustomStructure(scanner *bufio.Scanner) bool {
	for {
		fmt.Println("\n📝 Custom Structure Builder")
		fmt.Println("════════════════════════════════════════")
		fmt.Println("Enter your structure definition (one path per line)")
		fmt.Println("Prefix with 'f:' for files, 'd:' for directories")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  d:myproject")
		fmt.Println("  d:myproject/src")
		fmt.Println("  f:myproject/src/main.go")
		fmt.Println()
		fmt.Println("Commands: 'done' to finish, 'back' to cancel")
		fmt.Println()
		
		var paths []struct {
			path   string
			isFile bool
		}
		
		lineCount := 0
		for {
			fmt.Printf("> ")
			if !scanner.Scan() {
				return false
			}
			
			line := strings.TrimSpace(scanner.Text())
			
			if line == "done" || line == "exit" {
				break
			}
			
			if strings.ToLower(line) == "back" || strings.ToLower(line) == "cancel" {
				return false
			}
			
			if line == "" {
				continue
			}
			
			if strings.HasPrefix(line, "d:") {
				paths = append(paths, struct {
					path   string
					isFile bool
				}{strings.TrimPrefix(line, "d:"), false})
				lineCount++
			} else if strings.HasPrefix(line, "f:") {
				paths = append(paths, struct {
					path   string
					isFile bool
				}{strings.TrimPrefix(line, "f:"), true})
				lineCount++
			} else {
				fmt.Println("⚠️  Invalid format. Use 'd:' or 'f:' prefix")
			}
		}
		
		if len(paths) == 0 {
			fmt.Println("❌ No paths defined")
			continue
		}
		
		fmt.Printf("\n🔨 Creating %d entities...\n", len(paths))
		
		successCount := 0
		errorCount := 0
		
		for _, item := range paths {
			var result Result
			if item.isFile {
				result = CreateFile(item.path)
			} else {
				result = CreateFolder(item.path)
			}
			
			if result.Success {
				successCount++
				icon := "📁"
				if item.isFile {
					icon = "📄"
				}
				fmt.Printf("  ✅ %s %s\n", icon, item.path)
			} else {
				errorCount++
				fmt.Printf("  ❌ %s: %s\n", item.path, result.Message)
			}
		}
		
		fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n\n", successCount, errorCount)
		return true
	}
}

func handleParseTreeStructure(scanner *bufio.Scanner) bool {
	for {
		fmt.Println("\n🌳 Parse Tree Structure")
		fmt.Println("════════════════════════════════════════")
		fmt.Println("Paste your folder structure (tree format)")
		fmt.Println("Example:")
		fmt.Println("  myproject/")
		fmt.Println("  ├── src/")
		fmt.Println("  │   ├── main.go")
		fmt.Println("  │   └── utils.go")
		fmt.Println("  └── README.md")
		fmt.Println()
		fmt.Println("Enter 'END' on a new line when done, 'back' to cancel")
		fmt.Println()
		
		var lines []string
		for {
			if !scanner.Scan() {
				return false
			}
			
			line := scanner.Text()
			
			if strings.TrimSpace(line) == "END" {
				break
			}
			
			if strings.ToLower(strings.TrimSpace(line)) == "back" {
				return false
			}
			
			lines = append(lines, line)
		}
		
		if len(lines) == 0 {
			fmt.Println("❌ No structure provided")
			continue
		}
		
		input := strings.Join(lines, "\n")
		dirs, files, err := ParseTreeStructure(input)
		
		if err != nil {
			fmt.Printf("❌ Error parsing structure: %v\n", err)
			continue
		}
		
		fmt.Printf("\n📊 Parsed structure: %d directories, %d files\n", len(dirs), len(files))
		fmt.Print("\nProceed with creation? (yes/no): ")
		
		if !scanner.Scan() {
			return false
		}
		
		confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if confirm != "yes" && confirm != "y" {
			fmt.Println("❌ Creation cancelled")
			continue
		}
		
		fmt.Println("\n🔨 Creating structure...")
		
		successCount := 0
		errorCount := 0
		
		// Create directories
		for _, dir := range dirs {
			result := CreateFolder(dir)
			if result.Success {
				successCount++
				fmt.Printf("  ✅ 📁 %s\n", dir)
			} else {
				errorCount++
				fmt.Printf("  ❌ %s: %s\n", dir, result.Message)
			}
		}
		
		// Create files
		for filePath := range files {
			result := CreateFile(filePath)
			if result.Success {
				successCount++
				fmt.Printf("  ✅ 📄 %s\n", filePath)
			} else {
				errorCount++
				fmt.Printf("  ❌ %s: %s\n", filePath, result.Message)
			}
		}
		
		fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n\n", successCount, errorCount)
		return true
	}
}

func handleInteractiveBuilder(scanner *bufio.Scanner) bool {
	for {
		fmt.Println("\n🛠️  Interactive Structure Builder")
		fmt.Println("════════════════════════════════════════")
		fmt.Print("\n📁 Enter root directory (or 'back' to cancel): ")
		
		if !scanner.Scan() {
			return false
		}
		
		rootPath := strings.TrimSpace(scanner.Text())
		
		if strings.ToLower(rootPath) == "back" || strings.ToLower(rootPath) == "cancel" {
			return false
		}
		
		if rootPath == "" {
			fmt.Println("❌ Path cannot be empty")
			continue
		}
		
		// Create root
		result := CreateFolder(rootPath)
		if !result.Success {
			fmt.Printf("❌ Failed to create root: %s\n", result.Message)
			continue
		}
		
		fmt.Println("✅ Root directory created")
		fmt.Println()
		
		// Start interactive session in root
		successCount := 1 // Root already created
		errorCount := 0
		counts := interactiveBuildSession(scanner, rootPath, &successCount, &errorCount)
		successCount = counts[0]
		errorCount = counts[1]
		
		fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n\n", successCount, errorCount)
		return true
	}
}

func interactiveBuildSession(scanner *bufio.Scanner, currentPath string, successCount, errorCount *int) [2]int {
	for {
		fmt.Printf("\n📂 Current: %s\n", currentPath)
		fmt.Println("────────────────────────────────────────")
		fmt.Println("Commands:")
		fmt.Println("  mkdir <name>  - Create directory")
		fmt.Println("  touch <name>  - Create file")
		fmt.Println("  move in       - Enter a subdirectory")
		fmt.Println("  move out      - Go to parent directory")
		fmt.Println("  done          - Finish in this directory")
		fmt.Println("  exit          - Finish building completely")
		fmt.Println()
		
		fmt.Print("> ")
		if !scanner.Scan() {
			return [2]int{*successCount, *errorCount}
		}
		
		line := strings.TrimSpace(scanner.Text())
		
		if line == "exit" {
			return [2]int{*successCount, *errorCount}
		}
		
		if line == "done" {
			// Check if there are subdirectories to navigate into
			subdirs, err := getSubdirectories(currentPath)
			if err != nil || len(subdirs) == 0 {
				return [2]int{*successCount, *errorCount}
			}
			
			fmt.Printf("\n📁 Found %d subdirectories. Ready to navigate?\n", len(subdirs))
			fmt.Println("Type 'move in' to enter a directory, or 'exit' to finish")
			continue
		}
		
		if line == "move in" {
			subdirs, err := getSubdirectories(currentPath)
			if err != nil {
				fmt.Printf("❌ Error reading directories: %v\n", err)
				continue
			}
			
			if len(subdirs) == 0 {
				fmt.Println("ℹ️  No subdirectories to enter")
				continue
			}
			
			// Sort alphabetically
			sort.Strings(subdirs)
			
			var selectedDir string
			if len(subdirs) == 1 {
				// Auto-select if only one directory
				selectedDir = subdirs[0]
				fmt.Printf("📂 Entering: %s\n", selectedDir)
			} else {
				// Let user choose
				fmt.Println("\n📁 Available directories:")
				for i, dir := range subdirs {
					fmt.Printf("  %d. %s\n", i+1, dir)
				}
				fmt.Print("\nSelect directory (number): ")
				
				if !scanner.Scan() {
					continue
				}
				
				choice := strings.TrimSpace(scanner.Text())
				choiceNum, err := strconv.Atoi(choice)
				
				if err != nil || choiceNum < 1 || choiceNum > len(subdirs) {
					fmt.Println("❌ Invalid selection")
					continue
				}
				
				selectedDir = subdirs[choiceNum-1]
			}
			
			// Recursively enter the subdirectory
			newPath := filepath.Join(currentPath, selectedDir)
			counts := interactiveBuildSession(scanner, newPath, successCount, errorCount)
			*successCount = counts[0]
			*errorCount = counts[1]
			continue
		}
		
		if line == "move out" {
			fmt.Println("↩️  Moving to parent directory...")
			return [2]int{*successCount, *errorCount}
		}
		
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			fmt.Println("⚠️  Usage: mkdir <name> or touch <name>")
			continue
		}
		
		command := parts[0]
		name := parts[1]
		
		// Create path relative to current directory
		fullPath := filepath.Join(currentPath, name)
		
		var cmdResult Result
		switch command {
		case "mkdir":
			cmdResult = CreateFolder(fullPath)
			if cmdResult.Success {
				*successCount++
				fmt.Printf("  ✅ 📁 %s\n", name)
			} else {
				*errorCount++
				fmt.Printf("  ❌ %s\n", cmdResult.Message)
			}
		case "touch":
			// Validate file has extension
			if !strings.Contains(name, ".") {
				fmt.Printf("  ⚠️  Warning: '%s' has no file extension. Continue? (y/n): ", name)
				if scanner.Scan() {
					confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
					if confirm != "y" && confirm != "yes" {
						fmt.Println("  ❌ File creation cancelled")
						continue
					}
				}
			}
			
			cmdResult = CreateFile(fullPath)
			if cmdResult.Success {
				*successCount++
				fmt.Printf("  ✅ 📄 %s\n", name)
			} else {
				*errorCount++
				fmt.Printf("  ❌ %s\n", cmdResult.Message)
			}
		default:
			fmt.Println("⚠️  Unknown command. Use: mkdir or touch")
		}
	}
}

func getSubdirectories(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	
	return dirs, nil
}