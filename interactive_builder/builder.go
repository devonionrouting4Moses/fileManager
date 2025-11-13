package interactive_builder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// InteractiveBuilder handles the interactive structure building
func InteractiveBuilder(scanner *bufio.Scanner, currentPath string, successCount, errorCount *int) (int, int) {
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
			return *successCount, *errorCount
		}

		line := strings.TrimSpace(scanner.Text())

		if line == "exit" {
			return *successCount, *errorCount
		}

		if line == "done" {
			subdirs, err := getSubdirectories(currentPath)
			if err != nil || len(subdirs) == 0 {
				return *successCount, *errorCount
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

			sort.Strings(subdirs)

			var selectedDir string
			if len(subdirs) == 1 {
				selectedDir = subdirs[0]
				fmt.Printf("📂 Entering: %s\n", selectedDir)
			} else {
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

			newPath := filepath.Join(currentPath, selectedDir)
			success, errorC := InteractiveBuilder(scanner, newPath, successCount, errorCount)
			*successCount = success
			*errorCount = errorC
			continue
		}

		if line == "move out" {
			fmt.Println("↩️  Moving to parent directory...")
			return *successCount, *errorCount
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

		fullPath := filepath.Join(currentPath, name)

		var result Result
		switch command {
		case "mkdir":
			result = CreateFolder(fullPath)
			if result.Success {
				*successCount++
				fmt.Printf("  ✅ 📁 %s\n", name)
			} else {
				*errorCount++
				fmt.Printf("  ❌ %s\n", result.Message)
			}
		case "touch":
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

			result = CreateFile(fullPath)
			if result.Success {
				*successCount++
				fmt.Printf("  ✅ 📄 %s\n", name)
			} else {
				*errorCount++
				fmt.Printf("  ❌ %s\n", result.Message)
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

// Result represents the result of a file operation
type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateFolder creates a new directory
func CreateFolder(path string) Result {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return Result{Success: false, Message: err.Error()}
	}
	return Result{Success: true, Message: "Directory created successfully"}
}

// CreateFile creates a new file
func CreateFile(path string) Result {
	file, err := os.Create(path)
	if err != nil {
		return Result{Success: false, Message: err.Error()}
	}
	file.Close()
	return Result{Success: true, Message: "File created successfully"}
}
