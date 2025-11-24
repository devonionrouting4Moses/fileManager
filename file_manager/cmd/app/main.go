package main

import (
	"bufio"
	"filemanager/internal/ffi"
	"filemanager/internal/handler"
	"filemanager/internal/service"
	"filemanager/pkg/version"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Color scheme for each operation
type OperationStyle struct {
	Color       string
	Icon        string
	Title       string
	Description string
}

var operationStyles = map[int]OperationStyle{
	1: {Color: "\033[35m", Icon: "📁", Title: "CREATE FOLDER", Description: "Enter folder path(s) - space-separated for multiple folders"},
	2: {Color: "\033[35m", Icon: "📄", Title: "CREATE FILE", Description: "Enter file(s) - space-separated for multiple"},
	3: {Color: "\033[35m", Icon: "🔄", Title: "RENAME FILE/FOLDER", Description: "Enter old path and new name"},
	4: {Color: "\033[35m", Icon: "🗑️", Title: "DELETE FILE/FOLDER", Description: "Enter path(s) to delete - space-separated"},
	5: {Color: "\033[35m", Icon: "🔐", Title: "CHANGE PERMISSIONS", Description: "Enter path and permissions"},
	6: {Color: "\033[35m", Icon: "➡️", Title: "MOVE FILE/FOLDER", Description: "Enter source and destination paths"},
	7: {Color: "\033[35m", Icon: "📋", Title: "COPY FILE/FOLDER", Description: "Enter source and destination paths"},
}

// displayOperationProgress shows styled progress output for operations
func displayOperationProgress(operation int, status string, success bool) {
	style, exists := operationStyles[operation]
	if !exists {
		return
	}

	reset := "\033[0m"
	bold := "\033[1m"

	if success {
		fmt.Printf("%s%s✅ %s%s\n", style.Color, bold, status, reset)
	} else {
		fmt.Printf("%s%s❌ %s%s\n", style.Color, bold, status, reset)
	}
}

// displayInputBox shows a styled input box with Midnight Purple theme
func displayInputBox(prompt string) {
	// Midnight Purple color scheme
	purple := "\033[35m" // Light purple for border
	cyan := "\033[36m"   // Cyan for text
	reset := "\033[0m"
	bold := "\033[1m"
	boxWidth := 50

	fmt.Printf("%s%s┌%s┐%s\n", purple, bold, strings.Repeat("─", boxWidth-2), reset)

	// Prompt line with safe padding
	promptLen := len(prompt)
	promptPadding := boxWidth - promptLen - 2
	if promptPadding < 0 {
		promptPadding = 0
	}
	fmt.Printf("%s%s│ %s%s%s │%s\n", purple, bold, cyan, prompt, strings.Repeat(" ", promptPadding), reset)
	fmt.Printf("%s%s└%s┘%s ", purple, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Print("> ")
}

// displayStyledPrompt shows a styled input prompt for operations
func displayStyledPrompt(operation int, prompt string) {
	style, exists := operationStyles[operation]
	if !exists {
		fmt.Print(prompt)
		return
	}

	reset := "\033[0m"
	bold := "\033[1m"
	boxWidth := 50

	// Header
	fmt.Printf("%s%s┌%s┐%s\n", style.Color, bold, strings.Repeat("─", boxWidth-2), reset)

	// Title line with safe padding
	titleLen := len(style.Icon) + len(style.Title) + 3
	titlePadding := boxWidth - titleLen - 2
	if titlePadding < 0 {
		titlePadding = 0
	}
	fmt.Printf("%s%s│ %s %s%s │%s\n", style.Color, bold, style.Icon, style.Title, strings.Repeat(" ", titlePadding), reset)

	fmt.Printf("%s%s├%s┤%s\n", style.Color, bold, strings.Repeat("─", boxWidth-2), reset)

	// Description line with safe padding
	descLen := len(style.Description)
	descPadding := boxWidth - descLen - 2
	if descPadding < 0 {
		descPadding = 0
	}
	fmt.Printf("%s%s│ %s%s │%s\n", style.Color, bold, style.Description, strings.Repeat(" ", descPadding), reset)

	fmt.Printf("%s%s└%s┘%s\n", style.Color, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Print("> ")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Handle command-line flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			version.ShowVersion()
			return
		case "--update", "-u":
			version.CheckForUpdates()
			return
		case "--help", "-h":
			showHelp()
			return
		case "--web", "-w":
			// Start web server mode directly
			if err := handler.StartWebServer(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Server failed to start: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	version.ShowBanner()

	// Check for updates on startup (non-blocking)
	go func() {
		time.Sleep(500 * time.Millisecond)
		version.CheckForUpdates()
	}()

	for {
		displayMenu()

		fmt.Println()
		displayInputBox("Enter your choice (0-9)")
		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "0":
			fmt.Println("\n👋 Goodbye!")
			return
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
			handleWebServerLaunch()
			return
		default:
			fmt.Println("❌ Invalid choice. Please try again.")
		}
	}
}

func showHelp() {
	fmt.Printf("%s v%s - Hybrid File Manager\n\n", version.AppName, version.Version)
	fmt.Println("Usage:")
	fmt.Println("  filemanager              Start interactive mode")
	fmt.Println("  filemanager --version    Show version information")
	fmt.Println("  filemanager --update     Check for updates")
	fmt.Println("  filemanager --help       Show this help message")
	fmt.Println("  filemanager --web        Start web interface")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  • Single & batch file/folder operations")
	fmt.Println("  • 12 project templates (Flask, Spring Boot, React, etc.)")
	fmt.Println("  • Interactive structure builder")
	fmt.Println("  • Auto parent directory creation")
	fmt.Println("  • Cross-platform support (Linux, macOS, Windows)")
	fmt.Println("  • Web interface for browser-based management")
	fmt.Println()
}

func displayMenu() {
	// ANSI color codes
	green := "\033[32m"
	reset := "\033[0m"
	bold := "\033[1m"

	// Menu header with green border
	fmt.Printf("%s%s┌%s┐%s\n", green, bold, strings.Repeat("─", 48), reset)
	fmt.Printf("%s%s│%s AVAILABLE OPERATIONS %s│%s\n", green, bold, strings.Repeat(" ", 14), strings.Repeat(" ", 14), reset)
	fmt.Printf("%s%s├%s┤%s\n", green, bold, strings.Repeat("─", 48), reset)

	// Menu items
	menuItems := []string{
		"1️⃣  Create Folder",
		"2️⃣  Create File",
		"3️⃣  Rename File/Folder",
		"4️⃣  Delete File/Folder",
		"5️⃣  Change Permissions",
		"6️⃣  Move File/Folder",
		"7️⃣  Copy File/Folder",
		"8️⃣  Create Structure (Multi-entity)",
		"9️⃣  Launch Web Interface",
		"0️⃣  Exit",
	}

	for _, item := range menuItems {
		padding := 48 - len(item) - 2
		fmt.Printf("%s%s│ %s%s │%s\n", green, bold, item, strings.Repeat(" ", padding), reset)
	}

	fmt.Printf("%s%s└%s┘%s\n", green, bold, strings.Repeat("─", 48), reset)
	fmt.Println()
}

func handleWebServerLaunch() {
	// Midnight Purple color scheme
	primary := "\033[38;5;219m"   // Pink (#ffb7c5)
	accent := "\033[38;5;198m"    // Hot pink (#ff69b4)
	secondary := "\033[38;5;135m" // Violet (#9d4edd)
	reset := "\033[0m"
	bold := "\033[1m"
	boxWidth := 60

	fmt.Println()

	// Unified box - Launching Web Interface
	fmt.Printf("%s%s┌%s┐%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ 🌐 Launching Web Interface%s │%s\n", primary, bold, strings.Repeat(" ", 32), reset)
	fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ %sStarting HTTP server...%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 34), reset)
	fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ 📡 Connection Details%s │%s\n", primary, bold, strings.Repeat(" ", 37), reset)
	fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ %sURL: http://localhost:8080%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 28), reset)
	fmt.Printf("%s%s│ %sPress Ctrl+C to stop the server%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 27), reset)
	fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ %s✅ Server started successfully!%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 25), reset)
	fmt.Printf("%s%s│ %s🌐 Open your browser to begin%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 26), reset)
	fmt.Printf("%s%s└%s┘%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Println()

	if err := handler.StartWebServer(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Server failed to start: %v\n", err)
	}
}

func handleCreateFolder(scanner *bufio.Scanner) {
	fmt.Println()
	displayStyledPrompt(1, "")
	if !scanner.Scan() {
		return
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		fmt.Println("❌ Path cannot be empty")
		return
	}

	paths := strings.Fields(input)
	if len(paths) == 0 {
		fmt.Println("❌ No valid paths provided")
		return
	}

	successCount := 0
	errorCount := 0

	fmt.Println()
	for _, path := range paths {
		result := ffi.CreateFolder(path)
		if result.Success {
			successCount++
			displayOperationProgress(1, fmt.Sprintf("Created folder: %s", path), true)
		} else {
			errorCount++
			displayOperationProgress(1, fmt.Sprintf("Failed to create %s: %s", path, result.Message), false)
		}
	}

	if len(paths) > 1 {
		fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
	}
	fmt.Println()
}

func handleCreateFile(scanner *bufio.Scanner) {
	fmt.Println()
	displayStyledPrompt(2, "")
	if !scanner.Scan() {
		return
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		fmt.Println("❌ Path cannot be empty")
		return
	}

	paths := strings.Fields(input)
	if len(paths) == 0 {
		fmt.Println("❌ No valid paths provided")
		return
	}

	successCount := 0
	errorCount := 0

	fmt.Println()
	for _, path := range paths {
		// Validate the path
		if err := validateFilePath(path); err != nil {
			displayOperationProgress(2, fmt.Sprintf("Invalid path %s: %s", path, err.Error()), false)
			errorCount++
			continue
		}

		// Create parent directories if needed
		dir := filepath.Dir(path)
		if dir != "." && dir != path {
			result := ffi.CreateFolder(dir)
			if !result.Success {
				displayOperationProgress(2, fmt.Sprintf("Failed to create directory %s: %s", dir, result.Message), false)
				errorCount++
				continue
			}
			displayOperationProgress(2, fmt.Sprintf("Created directory: %s", dir), true)
		}

		result := ffi.CreateFile(path)
		if result.Success {
			successCount++
			displayOperationProgress(2, fmt.Sprintf("Created file: %s", path), true)
		} else {
			errorCount++
			displayOperationProgress(2, fmt.Sprintf("Failed to create %s: %s", path, result.Message), false)
		}
	}

	if len(paths) > 1 {
		fmt.Printf("📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
	}
	fmt.Println()
}

func validateFilePath(path string) error {
	parts := strings.Split(path, "/")
	lastPart := parts[len(parts)-1]

	dotCount := strings.Count(lastPart, ".")

	if dotCount > 1 {
		if !strings.HasPrefix(lastPart, ".") {
			extParts := strings.Split(lastPart, ".")
			if len(extParts) > 2 {
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
					return fmt.Errorf("invalid path - multiple file extensions detected")
				}
			}
		}
	}

	if !strings.Contains(lastPart, ".") {
		return fmt.Errorf("invalid filename - no file extension detected")
	}

	return nil
}

func handleRename(scanner *bufio.Scanner) {
	fmt.Println()
	displayStyledPrompt(3, "")
	if !scanner.Scan() {
		return
	}
	oldPath := strings.TrimSpace(scanner.Text())

	fmt.Println()
	displayInputBox("Enter new path")
	if !scanner.Scan() {
		return
	}
	newPath := strings.TrimSpace(scanner.Text())

	if oldPath == "" || newPath == "" {
		fmt.Println("❌ Paths cannot be empty")
		return
	}

	if oldPath == newPath {
		fmt.Println("❌ Source and destination paths are the same")
		return
	}

	fmt.Println()
	result := ffi.RenamePath(oldPath, newPath)
	if !result.Success {
		displayOperationProgress(3, fmt.Sprintf("Failed to rename %s: %s", oldPath, result.Message), false)
	} else {
		displayOperationProgress(3, fmt.Sprintf("Renamed %s to %s", oldPath, newPath), true)
	}
	fmt.Println()
}

func handleDelete(scanner *bufio.Scanner) {
	fmt.Println()
	displayStyledPrompt(4, "")
	if !scanner.Scan() {
		return
	}

	path := strings.TrimSpace(scanner.Text())
	if path == "" {
		fmt.Println("❌ Path cannot be empty")
		return
	}

	// Display confirmation in a styled box
	red := "\033[31m"
	reset := "\033[0m"
	bold := "\033[1m"
	boxWidth := 60

	fmt.Println()
	fmt.Printf("%s%s┌%s┐%s\n", red, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ ⚠️  DELETE CONFIRMATION%s │%s\n", red, bold, strings.Repeat(" ", 35), reset)
	fmt.Printf("%s%s├%s┤%s\n", red, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ Path: %s%s │%s\n", red, bold, path, strings.Repeat(" ", boxWidth-len(path)-8), reset)
	fmt.Printf("%s%s├%s┤%s\n", red, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ Are you sure? (yes/no)%s │%s\n", red, bold, strings.Repeat(" ", 34), reset)
	fmt.Printf("%s%s└%s┘%s\n", red, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Print("> ")

	if !scanner.Scan() {
		return
	}

	confirmation := strings.ToLower(strings.TrimSpace(scanner.Text()))
	fmt.Println()
	if confirmation != "yes" && confirmation != "y" {
		fmt.Println("❌ Deletion cancelled")
		return
	}

	fmt.Println()
	result := ffi.DeletePath(path)
	if result.Success {
		displayOperationProgress(4, fmt.Sprintf("Deleted: %s", path), true)
	} else {
		displayOperationProgress(4, fmt.Sprintf("Failed to delete %s: %s", path, result.Message), false)
	}
	fmt.Println()
}

func handleChangePermissions(scanner *bufio.Scanner) {
	fmt.Println()
	displayStyledPrompt(5, "")
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
		fmt.Println("❌ Path and permissions cannot be empty")
		return
	}

	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		fmt.Printf("❌ Invalid permission format. Please enter a valid octal value (e.g., 755, 644)\n")
		fmt.Printf("   Tip: 755 = rwxr-xr-x, 644 = rw-r--r--\n")
		return
	}

	fmt.Println()
	result := ffi.ChangePermissions(path, uint32(mode))
	if !result.Success {
		displayOperationProgress(5, fmt.Sprintf("Failed to change permissions: %s", result.Message), false)
		return
	}

	displayOperationProgress(5, fmt.Sprintf("Changed permissions to %s for %s", modeStr, path), true)
	fmt.Println()
	displayPermissionInfo(modeStr)
	fmt.Println()
}

// displayPermissionInfo shows detailed explanation of the permission mode
func displayPermissionInfo(modeStr string) {
	// Cyan color for box styling
	cyan := "\033[36m"
	reset := "\033[0m"
	bold := "\033[1m"
	boxWidth := 60

	if len(modeStr) != 3 {
		fmt.Println()
		fmt.Printf("%s%s┌%s┐%s\n", cyan, bold, strings.Repeat("─", boxWidth-2), reset)
		fmt.Printf("%s%s│ 📋 Permission Breakdown%s │%s\n", cyan, bold, strings.Repeat(" ", 34), reset)
		fmt.Printf("%s%s├%s┤%s\n", cyan, bold, strings.Repeat("─", boxWidth-2), reset)
		fmt.Printf("%s%s│ Invalid format%s │%s\n", cyan, bold, strings.Repeat(" ", 43), reset)
		fmt.Printf("%s%s└%s┘%s\n", cyan, bold, strings.Repeat("─", boxWidth-2), reset)
		fmt.Println()
		return
	}

	owner := string(modeStr[0])
	group := string(modeStr[1])
	others := string(modeStr[2])

	fmt.Println()
	fmt.Printf("%s%s┌%s┐%s\n", cyan, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ 📋 Permission Breakdown: %s%s │%s\n", cyan, bold, modeStr, strings.Repeat(" ", 29), reset)
	fmt.Printf("%s%s├%s┤%s\n", cyan, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ Owner:  %s (%s)%s │%s\n", cyan, bold, owner, decodePermission(owner), strings.Repeat(" ", boxWidth-len(fmt.Sprintf("Owner:  %s (%s)", owner, decodePermission(owner)))-4), reset)
	fmt.Printf("%s%s│ Group:  %s (%s)%s │%s\n", cyan, bold, group, decodePermission(group), strings.Repeat(" ", boxWidth-len(fmt.Sprintf("Group:  %s (%s)", group, decodePermission(group)))-4), reset)
	fmt.Printf("%s%s│ Others: %s (%s)%s │%s\n", cyan, bold, others, decodePermission(others), strings.Repeat(" ", boxWidth-len(fmt.Sprintf("Others: %s (%s)", others, decodePermission(others)))-4), reset)
	fmt.Printf("%s%s├%s┤%s\n", cyan, bold, strings.Repeat("─", boxWidth-2), reset)

	// Show common permission patterns
	switch modeStr {
	case "755":
		fmt.Printf("%s%s│ 📌 Directories & executable scripts%s │%s\n", cyan, bold, strings.Repeat(" ", 21), reset)
		fmt.Printf("%s%s│ ✓ Owner: read, write, execute%s │%s\n", cyan, bold, strings.Repeat(" ", 25), reset)
		fmt.Printf("%s%s│ ✓ Group: read, execute%s │%s\n", cyan, bold, strings.Repeat(" ", 32), reset)
		fmt.Printf("%s%s│ ✓ Others: read, execute%s │%s\n", cyan, bold, strings.Repeat(" ", 31), reset)

	case "644":
		fmt.Printf("%s%s│ 📌 Regular files & documents%s │%s\n", cyan, bold, strings.Repeat(" ", 28), reset)
		fmt.Printf("%s%s│ ✓ Owner: read, write%s │%s\n", cyan, bold, strings.Repeat(" ", 36), reset)
		fmt.Printf("%s%s│ ✓ Group: read only%s │%s\n", cyan, bold, strings.Repeat(" ", 37), reset)
		fmt.Printf("%s%s│ ✓ Others: read only%s │%s\n", cyan, bold, strings.Repeat(" ", 36), reset)

	case "700":
		fmt.Printf("%s%s│ 📌 Private files & directories%s │%s\n", cyan, bold, strings.Repeat(" ", 26), reset)
		fmt.Printf("%s%s│ ✓ Owner: read, write, execute%s │%s\n", cyan, bold, strings.Repeat(" ", 25), reset)
		fmt.Printf("%s%s│ ✗ Group: no access%s │%s\n", cyan, bold, strings.Repeat(" ", 37), reset)
		fmt.Printf("%s%s│ ✗ Others: no access%s │%s\n", cyan, bold, strings.Repeat(" ", 36), reset)

	case "777":
		fmt.Printf("%s%s│ 📌 Temporary files (NOT recommended)%s │%s\n", cyan, bold, strings.Repeat(" ", 19), reset)
		fmt.Printf("%s%s│ ✓ Owner: read, write, execute%s │%s\n", cyan, bold, strings.Repeat(" ", 25), reset)
		fmt.Printf("%s%s│ ✓ Group: read, write, execute%s │%s\n", cyan, bold, strings.Repeat(" ", 25), reset)
		fmt.Printf("%s%s│ ✓ Others: read, write, execute%s │%s\n", cyan, bold, strings.Repeat(" ", 24), reset)

	case "600":
		fmt.Printf("%s%s│ 📌 Sensitive files (configs, keys)%s │%s\n", cyan, bold, strings.Repeat(" ", 21), reset)
		fmt.Printf("%s%s│ ✓ Owner: read, write%s │%s\n", cyan, bold, strings.Repeat(" ", 36), reset)
		fmt.Printf("%s%s│ ✗ Group: no access%s │%s\n", cyan, bold, strings.Repeat(" ", 37), reset)
		fmt.Printf("%s%s│ ✗ Others: no access%s │%s\n", cyan, bold, strings.Repeat(" ", 36), reset)

	case "750":
		fmt.Printf("%s%s│ 📌 Project directories%s │%s\n", cyan, bold, strings.Repeat(" ", 34), reset)
		fmt.Printf("%s%s│ ✓ Owner: read, write, execute%s │%s\n", cyan, bold, strings.Repeat(" ", 25), reset)
		fmt.Printf("%s%s│ ✓ Group: read, execute%s │%s\n", cyan, bold, strings.Repeat(" ", 32), reset)
		fmt.Printf("%s%s│ ✗ Others: no access%s │%s\n", cyan, bold, strings.Repeat(" ", 36), reset)

	default:
		fmt.Printf("%s%s│ 📌 Custom Permission: %s%s │%s\n", cyan, bold, modeStr, strings.Repeat(" ", boxWidth-len(fmt.Sprintf("📌 Custom Permission: %s", modeStr))-4), reset)
		fmt.Printf("%s%s│ Owner: %s%s │%s\n", cyan, bold, decodePermission(owner), strings.Repeat(" ", boxWidth-len(fmt.Sprintf("Owner: %s", decodePermission(owner)))-4), reset)
		fmt.Printf("%s%s│ Group: %s%s │%s\n", cyan, bold, decodePermission(group), strings.Repeat(" ", boxWidth-len(fmt.Sprintf("Group: %s", decodePermission(group)))-4), reset)
		fmt.Printf("%s%s│ Others: %s%s │%s\n", cyan, bold, decodePermission(others), strings.Repeat(" ", boxWidth-len(fmt.Sprintf("Others: %s", decodePermission(others)))-4), reset)
	}

	fmt.Printf("%s%s└%s┘%s\n", cyan, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Println()
}

// decodePermission converts a single octal digit to readable format
func decodePermission(digit string) string {
	switch digit {
	case "0":
		return "---"
	case "1":
		return "--x (execute only)"
	case "2":
		return "-w- (write only)"
	case "3":
		return "-wx (write + execute)"
	case "4":
		return "r-- (read only)"
	case "5":
		return "r-x (read + execute)"
	case "6":
		return "rw- (read + write)"
	case "7":
		return "rwx (read + write + execute)"
	default:
		return "unknown"
	}
}

func handleMove(scanner *bufio.Scanner) {
	fmt.Println()
	displayStyledPrompt(6, "")
	if !scanner.Scan() {
		return
	}
	src := strings.TrimSpace(scanner.Text())

	fmt.Println()
	displayInputBox("Enter destination path")
	if !scanner.Scan() {
		return
	}
	dst := strings.TrimSpace(scanner.Text())

	if src == "" || dst == "" {
		fmt.Println("❌ Paths cannot be empty")
		return
	}

	if src == dst {
		fmt.Println("❌ Source and destination paths are the same")
		return
	}

	fmt.Println()
	result := ffi.MovePath(src, dst)
	if !result.Success {
		displayOperationProgress(6, fmt.Sprintf("Failed to move %s: %s", src, result.Message), false)
	} else {
		displayOperationProgress(6, fmt.Sprintf("Moved %s to %s", src, dst), true)
	}
	fmt.Println()
}

func handleCopy(scanner *bufio.Scanner) {
	fmt.Println()
	displayStyledPrompt(7, "")
	if !scanner.Scan() {
		return
	}
	src := strings.TrimSpace(scanner.Text())

	fmt.Println()
	displayInputBox("Enter destination path")
	if !scanner.Scan() {
		return
	}
	dst := strings.TrimSpace(scanner.Text())

	if src == "" || dst == "" {
		fmt.Println("❌ Paths cannot be empty")
		return
	}

	fmt.Println()
	result := ffi.CopyPath(src, dst)
	if result.Success {
		displayOperationProgress(7, fmt.Sprintf("Copied %s to %s", src, dst), true)
	} else {
		displayOperationProgress(7, fmt.Sprintf("Failed to copy %s: %s", src, result.Message), false)
	}
	fmt.Println()
}

func handleCreateStructure(scanner *bufio.Scanner) {
	for {
		// ANSI color codes
		cyan := "\033[36m"
		green := "\033[32m"
		yellow := "\033[33m"
		magenta := "\033[35m"
		reset := "\033[0m"
		bold := "\033[1m"

		// Header with cyan border
		fmt.Printf("%s%s┌%s┐%s\n", cyan, bold, strings.Repeat("─", 48), reset)
		fmt.Printf("%s%s│ 🏗️  CREATE HIERARCHICAL STRUCTURE%s │%s\n", cyan, bold, strings.Repeat(" ", 13), reset)
		fmt.Printf("%s%s└%s┘%s\n", cyan, bold, strings.Repeat("─", 48), reset)
		fmt.Println()

		templates := service.GetAvailableTemplates()

		// Templates section with green border
		fmt.Printf("%s%s┌%s┐%s\n", green, bold, strings.Repeat("─", 48), reset)
		fmt.Printf("%s%s│ 📦 AVAILABLE TEMPLATES%s │%s\n", green, bold, strings.Repeat(" ", 24), reset)
		fmt.Printf("%s%s├%s┤%s\n", green, bold, strings.Repeat("─", 48), reset)

		// Display templates in two columns for better layout
		for i, t := range templates {
			padding := 48 - len(fmt.Sprintf("%2d. %s", i+1, t.Description)) - 2
			fmt.Printf("%s%s│ %2d. %s%s │%s\n", green, bold, i+1, t.Description, strings.Repeat(" ", padding), reset)
		}

		fmt.Printf("%s%s├%s┤%s\n", green, bold, strings.Repeat("─", 48), reset)

		// Special options section with yellow accent
		specialOptions := []string{
			fmt.Sprintf("%2d. Custom Structure (from definition)", len(templates)+1),
			fmt.Sprintf("%2d. Parse Tree Structure (paste)", len(templates)+2),
			fmt.Sprintf("%2d. Interactive Builder", len(templates)+3),
		}

		for _, opt := range specialOptions {
			padding := 48 - len(opt) - 2
			fmt.Printf("%s%s│ %s%s │%s\n", yellow, bold, opt, strings.Repeat(" ", padding), reset)
		}

		fmt.Printf("%s%s└%s┘%s\n", green, bold, strings.Repeat("─", 48), reset)
		fmt.Println()

		// Back option with magenta border at bottom
		fmt.Printf("%s%s┌%s┐%s\n", magenta, bold, strings.Repeat("─", 48), reset)
		backOpt := fmt.Sprintf("%2d. ← Back to Main Menu", len(templates)+4)
		padding := 48 - len(backOpt) - 2
		fmt.Printf("%s%s│ %s%s │%s\n", magenta, bold, backOpt, strings.Repeat(" ", padding), reset)
		fmt.Printf("%s%s└%s┘%s\n", magenta, bold, strings.Repeat("─", 48), reset)
		fmt.Println()

		// Input prompt
		fmt.Printf("%s%s┌%s┐%s\n", magenta, bold, strings.Repeat("─", 48), reset)
		fmt.Printf("%s%s│ Select option (1-%d):%s │%s\n", magenta, bold, len(templates)+4, strings.Repeat(" ", 20), reset)
		fmt.Printf("%s%s└%s┘%s\n", magenta, bold, strings.Repeat("─", 48), reset)
		fmt.Print("> ")

		if !scanner.Scan() {
			return
		}

		option := strings.TrimSpace(scanner.Text())
		optionNum, err := strconv.Atoi(option)

		if err != nil {
			fmt.Println("❌ Invalid input. Please enter a number.")
			continue
		}

		if optionNum == len(templates)+4 {
			return
		}

		if optionNum >= 1 && optionNum <= len(templates) {
			if !handleTemplateStructure(scanner, templates[optionNum-1]) {
				continue
			}
			return
		}

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
			fmt.Println("❌ Invalid option")
		}
	}
}

func handleTemplateStructure(scanner *bufio.Scanner, template service.StructureTemplate) bool {
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

		successCount, errorCount := service.CreateFromTemplate(input, template)

		fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)

		if errorCount == 0 {
			fmt.Println("✨ Structure created successfully!")
		}

		return true
	}
}

// handleCustomStructure creates a custom structure from user-defined format
// Format: d:path/to/dir or f:path/to/file on each line
func handleCustomStructure(scanner *bufio.Scanner) bool {
	fmt.Println("\n📝 Custom Structure Definition")
	fmt.Println("════════════════════════════════════════")
	fmt.Println("Format: d:path/to/dir or f:path/to/file")
	fmt.Println("Example:")
	fmt.Println("  d:src")
	fmt.Println("  d:src/api")
	fmt.Println("  f:src/api/users.go")
	fmt.Println("  f:README.md")
	fmt.Println("Enter 'done' when finished (or 'back' to cancel)")
	fmt.Println()

	var lines []string
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			return false
		}

		input := strings.TrimSpace(scanner.Text())

		if strings.ToLower(input) == "done" {
			break
		}
		if strings.ToLower(input) == "back" || strings.ToLower(input) == "b" {
			return false
		}

		if input != "" {
			lines = append(lines, input)
		}
	}

	if len(lines) == 0 {
		fmt.Println("❌ No items defined")
		return false
	}

	successCount := 0
	errorCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result ffi.Result
		if strings.HasPrefix(line, "d:") {
			path := strings.TrimPrefix(line, "d:")
			path = strings.TrimSpace(path)
			result = ffi.CreateFolder(path)
			if result.Success {
				successCount++
				fmt.Printf("  ✅ 📁 %s\n", path)
			} else {
				errorCount++
				fmt.Printf("  ❌ %s: %s\n", path, result.Message)
			}
		} else if strings.HasPrefix(line, "f:") {
			path := strings.TrimPrefix(line, "f:")
			path = strings.TrimSpace(path)
			result = ffi.CreateFile(path)
			if result.Success {
				successCount++
				fmt.Printf("  ✅ 📄 %s\n", path)
			} else {
				errorCount++
				fmt.Printf("  ❌ %s: %s\n", path, result.Message)
			}
		}
	}

	fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
	return true
}

// handleParseTreeStructure creates a structure from pasted tree format
func handleParseTreeStructure(scanner *bufio.Scanner) bool {
	// Midnight Purple color scheme
	purple := "\033[35m" // Light purple for border
	cyan := "\033[36m"   // Cyan for text
	reset := "\033[0m"
	bold := "\033[1m"
	boxWidth := 60

	fmt.Println()
	fmt.Printf("%s%s┌%s┐%s\n", purple, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ 🌳 Parse Tree Structure%s │%s\n", purple, bold, strings.Repeat(" ", 35), reset)
	fmt.Printf("%s%s├%s┤%s\n", purple, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ %sPaste your tree structure (supports ├──, └──, │)%s │%s\n", purple, bold, cyan, strings.Repeat(" ", 7), reset)
	fmt.Printf("%s%s│ %sExample:%s │%s\n", purple, bold, cyan, strings.Repeat(" ", 50), reset)
	fmt.Printf("%s%s│ %s  myapp/%s │%s\n", purple, bold, cyan, strings.Repeat(" ", 50), reset)
	fmt.Printf("%s%s│ %s  ├── src/%s │%s\n", purple, bold, cyan, strings.Repeat(" ", 48), reset)
	fmt.Printf("%s%s│ %s  │   ├── main.go%s │%s\n", purple, bold, cyan, strings.Repeat(" ", 41), reset)
	fmt.Printf("%s%s│ %s  │   └── utils.go%s │%s\n", purple, bold, cyan, strings.Repeat(" ", 40), reset)
	fmt.Printf("%s%s│ %s  ├── tests/%s │%s\n", purple, bold, cyan, strings.Repeat(" ", 48), reset)
	fmt.Printf("%s%s│ %s  └── README.md%s │%s\n", purple, bold, cyan, strings.Repeat(" ", 42), reset)
	fmt.Printf("%s%s│ %sEnter 'done' when finished (or 'back' to cancel)%s │%s\n", purple, bold, cyan, strings.Repeat(" ", 8), reset)
	fmt.Printf("%s%s└%s┘%s\n", purple, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Println()

	var lines []string
	for {
		fmt.Printf("%s%s┌%s┐%s\n", purple, bold, strings.Repeat("─", boxWidth-2), reset)
		fmt.Printf("%s%s│ Input line:%s │%s\n", purple, bold, strings.Repeat(" ", boxWidth-14), reset)
		fmt.Printf("%s%s└%s┘%s ", purple, bold, strings.Repeat("─", boxWidth-2), reset)
		fmt.Print("> ")
		if !scanner.Scan() {
			return false
		}

		input := scanner.Text()

		if strings.ToLower(strings.TrimSpace(input)) == "done" {
			break
		}
		if strings.ToLower(strings.TrimSpace(input)) == "back" || strings.ToLower(strings.TrimSpace(input)) == "b" {
			return false
		}

		lines = append(lines, input)
	}

	if len(lines) == 0 {
		fmt.Println("❌ No structure provided")
		return false
	}

	treeInput := strings.Join(lines, "\n")
	dirs, files, err := service.ParseTreeStructure(treeInput)
	if err != nil {
		fmt.Printf("❌ Error parsing structure: %v\n", err)
		return false
	}

	successCount := 0
	errorCount := 0

	// Sort directories by depth to create parent dirs first
	sort.Slice(dirs, func(i, j int) bool {
		depthI := strings.Count(dirs[i], string(filepath.Separator))
		depthJ := strings.Count(dirs[j], string(filepath.Separator))
		if depthI != depthJ {
			return depthI < depthJ
		}
		return dirs[i] < dirs[j]
	})

	// Create directories
	for _, dir := range dirs {
		result := ffi.CreateFolder(dir)
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
		result := ffi.CreateFile(filePath)
		if result.Success {
			successCount++
			fmt.Printf("  ✅ 📄 %s\n", filePath)
		} else {
			errorCount++
			fmt.Printf("  ❌ %s: %s\n", filePath, result.Message)
		}
	}

	fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
	return true
}

// handleInteractiveBuilder provides an interactive directory navigation builder
func handleInteractiveBuilder(scanner *bufio.Scanner) bool {
	fmt.Println("\n🏗️  Interactive Structure Builder")
	fmt.Println("════════════════════════════════════════")
	fmt.Print("\n📁 Enter root directory name: ")

	if !scanner.Scan() {
		return false
	}

	rootDir := strings.TrimSpace(scanner.Text())
	if rootDir == "" {
		fmt.Println("❌ Root directory name cannot be empty")
		return false
	}

	result := ffi.CreateFolder(rootDir)
	if !result.Success {
		fmt.Printf("❌ Failed to create root: %s\n", result.Message)
		return false
	}

	fmt.Printf("✅ Root directory created: %s\n", rootDir)

	successCount := 1
	errorCount := 0

	// Interactive builder loop
	currentPath := rootDir
	for {
		fmt.Printf("\n📂 Current: %s\n", currentPath)
		fmt.Println("────────────────────────────────────────")
		fmt.Println("Commands:")
		fmt.Println("  mkdir <name>  - Create directory")
		fmt.Println("  touch <name>  - Create file")
		fmt.Println("  cd <name>     - Enter subdirectory")
		fmt.Println("  back          - Go to parent directory")
		fmt.Println("  list          - List subdirectories")
		fmt.Println("  done          - Finish building")
		fmt.Println()

		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "mkdir":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: mkdir <name>")
				continue
			}
			dirName := strings.Join(parts[1:], " ")
			newPath := filepath.Join(currentPath, dirName)
			result := ffi.CreateFolder(newPath)
			if result.Success {
				successCount++
				fmt.Printf("✅ 📁 %s\n", newPath)
			} else {
				errorCount++
				fmt.Printf("❌ %s\n", result.Message)
			}

		case "touch":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: touch <name>")
				continue
			}
			fileName := strings.Join(parts[1:], " ")
			newPath := filepath.Join(currentPath, fileName)
			result := ffi.CreateFile(newPath)
			if result.Success {
				successCount++
				fmt.Printf("✅ 📄 %s\n", newPath)
			} else {
				errorCount++
				fmt.Printf("❌ %s\n", result.Message)
			}

		case "cd":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: cd <name>")
				continue
			}
			dirName := parts[1]
			if dirName == ".." {
				parent := filepath.Dir(currentPath)
				if parent != currentPath {
					currentPath = parent
					fmt.Printf("↩️  Moved to: %s\n", currentPath)
				} else {
					fmt.Println("❌ Already at root")
				}
			} else {
				newPath := filepath.Join(currentPath, dirName)
				currentPath = newPath
				fmt.Printf("📂 Entering: %s\n", currentPath)
			}

		case "back":
			parent := filepath.Dir(currentPath)
			if parent != currentPath {
				currentPath = parent
				fmt.Printf("↩️  Moved to: %s\n", currentPath)
			} else {
				fmt.Println("❌ Already at root")
			}

		case "list":
			fmt.Println("📁 Current directory structure:")
			fmt.Printf("  %s/\n", currentPath)

		case "done":
			fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
			return true

		default:
			fmt.Println("❌ Unknown command. Type 'done' to finish.")
		}
	}

	fmt.Printf("\n📊 Summary: %d succeeded, %d failed\n", successCount, errorCount)
	return true
}
