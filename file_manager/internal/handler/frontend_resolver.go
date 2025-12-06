package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the YAML configuration structure
type Config struct {
	FrontendDir string `yaml:"frontend_dir"`
}

// GetFrontendDir returns the single, stable frontend directory location
// Priority: ENV → User Config → System Config → OS Default → Portable Fallback
func GetFrontendDir() (string, error) {
	// 🚨 SNAP SPECIAL HANDLING - Must come FIRST
	if snapDir := os.Getenv("SNAP"); snapDir != "" {
		// Check if bundled frontend exists (read-only but OK for serving)
		bundledPath := filepath.Join(snapDir, "usr/share/filemanager/frontend")
		if isValidDir(bundledPath) {
			// log.Printf("📍 Using bundled snap frontend: %s\n", bundledPath)
			return bundledPath, nil
		}

		// If bundled doesn't exist, use writable snap common directory
		// log.Printf("⚠️  Bundled frontend not found, using writable directory\n")
		snapUserCommon := os.Getenv("SNAP_USER_COMMON")
		if snapUserCommon == "" {
			home, _ := os.UserHomeDir()
			snapUserCommon = filepath.Join(home, "snap/filemanager/common")
		}
		frontendDir := filepath.Join(snapUserCommon, "frontend")
		// log.Printf("📍 Using snap user common: %s\n", frontendDir)
		return frontendDir, nil
	}

	// 1️⃣ Environment Variable Override (Highest Priority)
	if envDir := os.Getenv("FILEMANAGER_FRONTEND_DIR"); envDir != "" {
		absPath, err := filepath.Abs(envDir)
		if err == nil && isValidDir(absPath) {
			return absPath, nil
		}
		// ENV path exists but is not valid directory
	}

	// 2️⃣ User Config File
	if userConfigDir := getUserConfigFrontendDir(); userConfigDir != "" {
		if isValidDir(userConfigDir) {
			return userConfigDir, nil
		}
	}

	// 3️⃣ System Config File
	if systemConfigDir := getSystemConfigFrontendDir(); systemConfigDir != "" {
		if isValidDir(systemConfigDir) {
			return systemConfigDir, nil
		}
	}

	// 4️⃣ OS-Specific Default Directory
	defaultDir := getOSDefaultFrontendDir()
	if isValidDir(defaultDir) {
		return defaultDir, nil
	}

	// 5️⃣ Portable Fallback (directory next to executable)
	portableDir := getPortableFrontendDir()
	if isValidDir(portableDir) {
		return portableDir, nil
	}

	// If nothing exists, return the most appropriate default to create
	preferredDir := getPreferredCreationDir()
	return preferredDir, nil
}

// getPreferredCreationDir returns the best location to create frontend if it doesn't exist
func getPreferredCreationDir() string {
	// For user installs, prefer user directories
	// For system installs, prefer system directories

	// Check if running with elevated privileges (system install)
	if isElevated() {
		return getOSDefaultFrontendDir()
	}

	// For regular users, prefer user-writable locations
	switch runtime.GOOS {
	case "linux":
		// Check for Snap/Flatpak first
		if snapDir := os.Getenv("SNAP"); snapDir != "" {
			return filepath.Join(snapDir, "frontend")
		}
		if flatpakID := os.Getenv("FLATPAK_ID"); flatpakID != "" {
			home, _ := os.UserHomeDir()
			return filepath.Join(home, ".var/app", flatpakID, "data/frontend")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".local/share/filemanager/frontend")
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library/Application Support/FileManager/frontend")
	case "windows":
		appData := os.Getenv("LOCALAPPDATA")
		if appData != "" {
			return filepath.Join(appData, "FileManager", "frontend")
		}
		return getPortableFrontendDir()
	default:
		// HarmonyOS or unknown - use portable
		return getPortableFrontendDir()
	}
}

// isElevated checks if the process is running with elevated privileges
func isElevated() bool {
	switch runtime.GOOS {
	case "linux", "darwin":
		return os.Geteuid() == 0
	case "windows":
		// On Windows, check if we can write to Program Files
		testPath := filepath.Join(os.Getenv("ProgramFiles"), ".filemanager_write_test")
		if err := os.WriteFile(testPath, []byte("test"), 0644); err == nil {
			os.Remove(testPath)
			return true
		}
		return false
	default:
		return false
	}
}

// getUserConfigFrontendDir reads frontend_dir from user config file
func getUserConfigFrontendDir() string {
	configPath := getUserConfigPath()
	return parseFrontendDirFromConfig(configPath)
}

// getSystemConfigFrontendDir reads frontend_dir from system config file
func getSystemConfigFrontendDir() string {
	configPath := getSystemConfigPath()
	return parseFrontendDirFromConfig(configPath)
}

// parseFrontendDirFromConfig parses YAML config and returns frontend_dir
func parseFrontendDirFromConfig(configPath string) string {
	if configPath == "" {
		return ""
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return ""
	}

	if config.FrontendDir != "" {
		absPath, err := filepath.Abs(config.FrontendDir)
		if err == nil {
			return absPath
		}
	}

	return ""
}

// getUserConfigPath returns the user config file path based on OS
func getUserConfigPath() string {
	switch runtime.GOOS {
	case "linux":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config/filemanager/config.yaml")
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library/Preferences/filemanager/config.yaml")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "filemanager", "config.yaml")
		}
		return ""
	default:
		// HarmonyOS - requires UID
		uid := os.Getuid()
		return fmt.Sprintf("/data/accounts/%d/appdata/filemanager/config.yaml", uid)
	}
}

// getSystemConfigPath returns the system config file path based on OS
func getSystemConfigPath() string {
	switch runtime.GOOS {
	case "linux":
		return "/etc/filemanager/config.yaml"
	case "darwin":
		return "/Library/Application Support/FileManager/config.yaml"
	case "windows":
		programFiles := os.Getenv("ProgramFiles")
		if programFiles != "" {
			return filepath.Join(programFiles, "FileManager", "config.yaml")
		}
		return ""
	default:
		// HarmonyOS
		return "/system/app/FileManager/config.yaml"
	}
}

// getOSDefaultFrontendDir returns the OS-specific default frontend directory
func getOSDefaultFrontendDir() string {
	switch runtime.GOOS {
	case "linux":
		return getLinuxDefaultFrontendDir()
	case "darwin":
		return getMacOSDefaultFrontendDir()
	case "windows":
		return getWindowsDefaultFrontendDir()
	default:
		return getHarmonyOSDefaultFrontendDir()
	}
}

// getLinuxDefaultFrontendDir returns Linux-specific frontend directory
func getLinuxDefaultFrontendDir() string {
	// Check for Snap (bundled in read-only location)
	if snapDir := os.Getenv("SNAP"); snapDir != "" {
		bundledPath := filepath.Join(snapDir, "usr/share/filemanager/frontend")
		if isValidDir(bundledPath) {
			return bundledPath
		}
		// Fallback to writable snap directory
		snapUserCommon := os.Getenv("SNAP_USER_COMMON")
		if snapUserCommon == "" {
			home, _ := os.UserHomeDir()
			snapUserCommon = filepath.Join(home, "snap/filemanager/common")
		}
		return filepath.Join(snapUserCommon, "frontend")
	}

	// Check for Flatpak (bundled in read-only location)
	if flatpakID := os.Getenv("FLATPAK_ID"); flatpakID != "" {
		bundledPath := filepath.Join("/app/share/filemanager/frontend")
		if isValidDir(bundledPath) {
			return bundledPath
		}
		// Fallback to writable flatpak directory
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".var/app", flatpakID, "data/frontend")
	}

	// Check system installation (APT/DEB/RPM)
	systemDir := "/usr/share/filemanager/frontend"
	if isValidDir(systemDir) {
		return systemDir
	}

	// Fallback to user directory
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local/share/filemanager/frontend")
}

// getMacOSDefaultFrontendDir returns macOS-specific frontend directory
func getMacOSDefaultFrontendDir() string {
	// Check Homebrew installation
	brewPaths := []string{
		"/usr/local/opt/filemanager/share/frontend",
		"/opt/homebrew/opt/filemanager/share/frontend", // Apple Silicon
	}
	for _, brewDir := range brewPaths {
		if isValidDir(brewDir) {
			return brewDir
		}
	}

	// Check system installation
	systemDir := "/Library/Application Support/FileManager/frontend"
	if isValidDir(systemDir) {
		return systemDir
	}

	// Fallback to user directory
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library/Application Support/FileManager/frontend")
}

// getWindowsDefaultFrontendDir returns Windows-specific frontend directory
func getWindowsDefaultFrontendDir() string {
	// Check Program Files installation
	programFiles := os.Getenv("ProgramFiles")
	if programFiles != "" {
		systemDir := filepath.Join(programFiles, "FileManager", "frontend")
		if isValidDir(systemDir) {
			return systemDir
		}
	}

	// Fallback to AppData
	appData := os.Getenv("LOCALAPPDATA")
	if appData != "" {
		return filepath.Join(appData, "FileManager", "frontend")
	}

	// Last resort: portable mode
	return getPortableFrontendDir()
}

// getHarmonyOSDefaultFrontendDir returns HarmonyOS-specific frontend directory
func getHarmonyOSDefaultFrontendDir() string {
	// Check system installation
	systemDir := "/system/app/FileManager/frontend"
	if isValidDir(systemDir) {
		return systemDir
	}

	// Fallback to user directory
	uid := os.Getuid()
	return fmt.Sprintf("/data/accounts/%d/applications/filemanager/frontend", uid)
}

// getPortableFrontendDir returns the portable frontend directory (next to executable)
func getPortableFrontendDir() string {
	exePath, err := os.Executable()
	if err == nil {
		// Resolve symlinks
		realPath, err := filepath.EvalSymlinks(exePath)
		if err == nil {
			exePath = realPath
		}
		return filepath.Join(filepath.Dir(exePath), "frontend")
	}
	// Absolute fallback
	return "./frontend"
}

// isValidDir checks if a path exists and is a directory
func isValidDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ensureFrontendExists ensures the frontend directory exists
// If it doesn't exist, creates it with the basic frontend structure
func ensureFrontendExists(frontendDir string) error {
	if isValidDir(frontendDir) {
		// Verify critical files exist
		indexPath := filepath.Join(frontendDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			return nil
		}
	}

	// Create parent directories if needed
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return fmt.Errorf("failed to create frontend directory: %w", err)
	}

	// Create the basic frontend structure
	createBasicFrontend(frontendDir)
	return nil
}

// CreateDefaultUserConfig creates a default user config file with the provided frontend directory
func CreateDefaultUserConfig(frontendDir string) error {
	configPath := getUserConfigPath()
	if configPath == "" {
		return fmt.Errorf("unable to determine user config path")
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Don't overwrite existing config
	if _, err := os.Stat(configPath); err == nil {
		return nil
	}

	// Create config with frontend_dir
	config := Config{
		FrontendDir: frontendDir,
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add helpful comments
	configContent := fmt.Sprintf(`# FileManager Configuration
# This file controls where the web interface frontend files are located
# You can override this by setting the FILEMANAGER_FRONTEND_DIR environment variable

# Frontend directory path
%s

# Example configurations:
# frontend_dir: /usr/share/filemanager/frontend
# frontend_dir: ~/.local/share/filemanager/frontend
# frontend_dir: /custom/path/to/frontend
`, string(data))

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// PrintFrontendInfo prints diagnostic information about frontend location
func PrintFrontendInfo() {
	fmt.Println("🔍 Frontend Directory Resolution Info:")
	fmt.Println(strings.Repeat("─", 50))

	// 1. Environment variable
	if envDir := os.Getenv("FILEMANAGER_FRONTEND_DIR"); envDir != "" {
		fmt.Printf("1️⃣  ENV FILEMANAGER_FRONTEND_DIR: %s [%s]\n", envDir, getPathStatus(envDir))
	} else {
		fmt.Println("1️⃣  ENV FILEMANAGER_FRONTEND_DIR: (not set)")
	}

	// 2. User config
	userConfigPath := getUserConfigPath()
	userConfigDir := getUserConfigFrontendDir()
	if userConfigDir != "" {
		fmt.Printf("2️⃣  User Config (%s): %s [%s]\n", userConfigPath, userConfigDir, getPathStatus(userConfigDir))
	} else {
		fmt.Printf("2️⃣  User Config (%s): (not configured)\n", userConfigPath)
	}

	// 3. System config
	systemConfigPath := getSystemConfigPath()
	systemConfigDir := getSystemConfigFrontendDir()
	if systemConfigDir != "" {
		fmt.Printf("3️⃣  System Config (%s): %s [%s]\n", systemConfigPath, systemConfigDir, getPathStatus(systemConfigDir))
	} else {
		fmt.Printf("3️⃣  System Config (%s): (not configured)\n", systemConfigPath)
	}

	// 4. OS default
	defaultDir := getOSDefaultFrontendDir()
	fmt.Printf("4️⃣  OS Default: %s [%s]\n", defaultDir, getPathStatus(defaultDir))

	// 5. Portable fallback
	portableDir := getPortableFrontendDir()
	fmt.Printf("5️⃣  Portable Fallback: %s [%s]\n", portableDir, getPathStatus(portableDir))

	fmt.Println(strings.Repeat("─", 50))

	// Show which one will be used
	frontendDir, _ := GetFrontendDir()
	fmt.Printf("✅ Active Frontend: %s\n", frontendDir)
}

func getPathStatus(path string) string {
	if isValidDir(path) {
		return "EXISTS"
	}
	if _, err := os.Stat(path); err == nil {
		return "EXISTS (not a directory)"
	}
	return "DOES NOT EXIST"
}
