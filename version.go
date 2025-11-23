package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	Version   = "0.1.2"
	AppName   = "FileManager"
	repoOwner = "devonionrouting4Moses"
	repoName  = "fileManager"
)

var (
	releaseURL   = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	downloadBase = fmt.Sprintf("https://github.com/%s/%s/releases/download", repoOwner, repoName)

	// Cache file for update checks
	cacheFile = filepath.Join(os.TempDir(), "filemanager_update_cache.json")
)

type ReleaseInfo struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type UpdateCache struct {
	LastCheck   time.Time   `json:"last_check"`
	ReleaseInfo ReleaseInfo `json:"release_info"`
}

// GetDownloadURL returns the download URL for the current version
func GetDownloadURL() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Map GOARCH to common architecture names
	var archName string
	switch arch {
	case "amd64":
		archName = "amd64"
	case "386":
		archName = "386"
	case "arm64":
		archName = "arm64"
	default:
		archName = arch
	}

	// Build target name to match GitHub Actions matrix
	target := fmt.Sprintf("%s-%s", osName, archName)

	var ext string
	if osName == "windows" {
		ext = "zip"
	} else {
		ext = "tar.gz"
	}

	return fmt.Sprintf("%s/v%s/filemanager-%s-%s.%s", downloadBase, Version, Version, target, ext)
}

// ShowVersion displays the current version and platform information
func ShowVersion() {
	fmt.Printf("%s v%s\n", AppName, Version)
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("Download: %s\n", GetDownloadURL())
}

// CheckForUpdates checks if a newer version is available on GitHub
func CheckForUpdates() {
	fmt.Println("\n🔍 Checking for updates...")

	// Skip update check if running in development
	if isDevBuild() {
		fmt.Println("Development build - skipping update check")
		return
	}

	// Check cache first (only check once per day)
	if cached, ok := loadCache(); ok {
		if time.Since(cached.LastCheck) < 24*time.Hour {
			fmt.Println("Using cached update information...")
			displayUpdateInfo(cached.ReleaseInfo)
			return
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", releaseURL, nil)
	if err != nil {
		fmt.Printf("⚠️  Could not create update request: %v\n", err)
		return
	}

	// Add GitHub API headers to avoid rate limiting
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("⚠️  Could not check for updates: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusForbidden {
			fmt.Println("⚠️  Rate limited by GitHub API. Set GITHUB_TOKEN environment variable to increase rate limit.")
		} else {
			fmt.Printf("⚠️  Failed to check for updates: %s\n", resp.Status)
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("⚠️  Could not read update response: %v\n", err)
		return
	}

	var release ReleaseInfo
	if err := json.Unmarshal(body, &release); err != nil {
		fmt.Printf("⚠️  Could not parse update information: %v\n", err)
		return
	}

	// Save to cache
	saveCache(UpdateCache{
		LastCheck:   time.Now(),
		ReleaseInfo: release,
	})

	displayUpdateInfo(release)
}

func displayUpdateInfo(release ReleaseInfo) {
	currentVersion := "v" + Version
	comparison := compareVersions(release.TagName, currentVersion)

	if comparison > 0 {
		fmt.Printf("\n🎉 New version available: %s (current: %s)\n", release.TagName, currentVersion)
		if release.Body != "" {
			// Limit the body preview to 200 characters
			bodyPreview := strings.TrimSpace(release.Body)
			if len(bodyPreview) > 200 {
				bodyPreview = bodyPreview[:197] + "..."
			}
			fmt.Printf("📝 %s\n\n", bodyPreview)
		}

		// Find the appropriate download URL for this platform
		downloadURL := findAssetURL(release)
		if downloadURL != "" {
			fmt.Printf("📦 Download from:\n   %s\n\n", downloadURL)

			// Ask user if they want to install
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Would you like to install this update now? (yes/no): ")
			if scanner.Scan() {
				response := strings.ToLower(strings.TrimSpace(scanner.Text()))
				if response == "yes" || response == "y" {
					installUpdate(release, downloadURL)
					return
				}
			}
		} else {
			releaseURL := fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", repoOwner, repoName, release.TagName)
			fmt.Printf("🌐 Visit: %s\n", releaseURL)
		}

		fmt.Println("\n💡 Manual installation instructions:")
		printInstallInstructions(release.TagName)
	} else if comparison < 0 {
		fmt.Printf("ℹ️  You're running a pre-release version (v%s, latest stable: %s)\n", Version, release.TagName)
	} else {
		fmt.Printf("✅ You're running the latest version (v%s)\n", Version)
	}
}

func findAssetURL(release ReleaseInfo) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Map GOARCH to common architecture names
	var archName string
	switch arch {
	case "amd64":
		archName = "amd64"
	case "386":
		archName = "386"
	case "arm64":
		archName = "arm64"
	default:
		archName = arch
	}

	// Build target name to match GitHub Actions matrix
	target := fmt.Sprintf("%s-%s", osName, archName)

	var ext string
	if osName == "windows" {
		ext = ".zip"
	} else {
		ext = ".tar.gz"
	}

	// Search through assets for matching filename
	for _, asset := range release.Assets {
		// Match pattern: filemanager-{version}-{target}.{ext}
		if strings.Contains(asset.Name, target) && strings.HasSuffix(asset.Name, ext) {
			return asset.DownloadURL
		}
	}

	return ""
}

func printInstallInstructions(version string) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Map GOARCH to common architecture names
	var archName string
	switch arch {
	case "amd64":
		archName = "amd64"
	case "386":
		archName = "386"
	case "arm64":
		archName = "arm64"
	default:
		archName = arch
	}

	target := fmt.Sprintf("%s-%s", osName, archName)
	downloadURL := fmt.Sprintf("%s/%s/filemanager-%s-%s", downloadBase, version, strings.TrimPrefix(version, "v"), target)

	switch osName {
	case "linux":
		fmt.Printf("curl -L %s.tar.gz | tar xz\n", downloadURL)
		fmt.Println("cd linux-amd64")
		fmt.Println("sudo ./install.sh")
	case "darwin":
		fmt.Printf("curl -L %s.tar.gz | tar xz\n", downloadURL)
		fmt.Println("cd darwin-amd64")
		fmt.Println("sudo ./install.sh")
	case "windows":
		fmt.Printf("Download: %s.zip\n", downloadURL)
		fmt.Println("Extract and run install.bat as Administrator")
	}
}

func loadCache() (UpdateCache, bool) {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return UpdateCache{}, false
	}

	var cache UpdateCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return UpdateCache{}, false
	}

	return cache, true
}

func saveCache(cache UpdateCache) {
	data, err := json.Marshal(cache)
	if err != nil {
		return
	}
	os.WriteFile(cacheFile, data, 0644)
}

func isDevBuild() bool {
	// Check for common development environment indicators
	if strings.Contains(Version, "dev") ||
		strings.Contains(Version, "SNAPSHOT") ||
		Version == "0.0.0" {
		return true
	}

	// Check if running from source (go run)
	if len(os.Args) > 0 && strings.HasSuffix(os.Args[0], "go-build") {
		return true
	}

	return false
}

func compareVersions(v1, v2 string) int {
	// Remove 'v' prefix if present
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int

		if i < len(parts1) {
			// Handle pre-release versions (e.g., "1.0.0-beta")
			part := strings.Split(parts1[i], "-")[0]
			n1, _ = strconv.Atoi(part)
		}
		if i < len(parts2) {
			part := strings.Split(parts2[i], "-")[0]
			n2, _ = strconv.Atoi(part)
		}

		if n1 > n2 {
			return 1
		} else if n1 < n2 {
			return -1
		}
	}

	return 0
}

func ShowBanner() {
	// Dynamic banner that adjusts to version length
	versionStr := fmt.Sprintf("v%s", Version)
	bannerWidth := 40

	// Calculate padding
	contentStr := fmt.Sprintf("🗂️  %s %s (Rust+Go)", AppName, versionStr)
	// Remove emoji width (counts as more than 1 char)
	visualLen := len(contentStr) - 2 // Approximate emoji width adjustment
	padding := (bannerWidth - visualLen) / 2

	fmt.Println("╔" + strings.Repeat("═", bannerWidth) + "╗")
	fmt.Printf("║%s%s%s║\n",
		strings.Repeat(" ", padding),
		contentStr,
		strings.Repeat(" ", bannerWidth-visualLen-padding))
	fmt.Println("╚" + strings.Repeat("═", bannerWidth) + "╝")
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
}

func installUpdate(release ReleaseInfo, downloadURL string) {
	fmt.Println("\n⬇️  Detecting installation method...")

	// Detect installation method
	installMethod := detectInstallationMethod()

	switch installMethod {
	case "snap":
		installUpdateSnap(release)
	case "deb":
		installUpdateDeb(release)
	case "exe":
		installUpdateExe(release, downloadURL)
	default:
		// Direct binary installation (default for manual Linux/macOS)
		installUpdateDirect(release, downloadURL)
	}
}

func detectInstallationMethod() string {
	// Check if running from snap
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		return "snap"
	}

	currentExe, err := os.Executable()
	if err != nil {
		return "direct"
	}

	// Check if running from snap confinement (alternative check)
	if strings.Contains(currentExe, "/snap/") {
		return "snap"
	}

	// Check for deb/apt installation on Debian/Ubuntu systems
	if _, err := os.Stat("/var/lib/dpkg/status"); err == nil {
		// Debian/Ubuntu system with dpkg
		// Check if binary is in standard system locations
		if strings.HasPrefix(currentExe, "/usr/bin/") || strings.HasPrefix(currentExe, "/usr/local/bin/") {
			// Verify it's from deb by checking if package is installed
			if isDebPackageInstalled("filemanager") {
				return "deb"
			}
		}
	}

	// Check for Windows .exe installation
	if runtime.GOOS == "windows" {
		if strings.HasSuffix(currentExe, ".exe") {
			return "exe"
		}
	}

	return "direct"
}

func isDebPackageInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	err := cmd.Run()
	return err == nil
}

func installUpdateSnap(release ReleaseInfo) {
	fmt.Println("\n📦 Detected snap installation")
	fmt.Println("🔄 Updating via snap store...")

	cmd := exec.Command("snap", "refresh", "filemanager")
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to update via snap: %v\n", err)
		fmt.Println("💡 Try manually: snap refresh filemanager")
		return
	}

	fmt.Printf("\n✅ Update installed successfully via snap!\n")
	fmt.Printf("📝 New version: %s\n", release.TagName)
	fmt.Println("🔄 Snap will automatically restart the application")
}

func installUpdateDeb(release ReleaseInfo) {
	fmt.Println("\n📦 Detected deb package installation")
	fmt.Println("🔄 Updating via apt package manager...")

	// Update package list
	fmt.Println("📥 Updating package list...")
	updateCmd := exec.Command("sudo", "apt-get", "update")
	if err := updateCmd.Run(); err != nil {
		fmt.Printf("⚠️  Warning: Could not update package list: %v\n", err)
	}

	// Install/upgrade package
	fmt.Println("📦 Installing new version...")
	cmd := exec.Command("sudo", "apt-get", "install", "-y", "filemanager")
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to update via apt: %v\n", err)
		fmt.Println("💡 Try manually: sudo apt-get install filemanager")
		return
	}

	fmt.Printf("\n✅ Update installed successfully via deb!\n")
	fmt.Printf("📝 New version: %s\n", release.TagName)
	fmt.Println("🔄 Please restart the application to use the new version")
}

func installUpdateExe(release ReleaseInfo, downloadURL string) {
	fmt.Println("\n📦 Detected Windows .exe installation")
	fmt.Println("⬇️  Downloading Windows installer...")

	// Download the file
	resp, err := http.Get(downloadURL)
	if err != nil {
		fmt.Printf("❌ Failed to download update: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Download failed with status: %s\n", resp.Status)
		return
	}

	// Create temporary directory for extraction
	tempDir := filepath.Join(os.TempDir(), "filemanager_update")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create temp directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Save ZIP file
	zipPath := filepath.Join(tempDir, "update.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		fmt.Printf("❌ Failed to create temp file: %v\n", err)
		return
	}

	if _, err := io.Copy(zipFile, resp.Body); err != nil {
		zipFile.Close()
		fmt.Printf("❌ Failed to download file: %v\n", err)
		return
	}
	zipFile.Close()

	// Extract ZIP
	extractedBinary, err := extractZip(zipPath, tempDir)
	if err != nil {
		fmt.Printf("❌ Failed to extract update: %v\n", err)
		return
	}

	if extractedBinary == "" {
		fmt.Println("❌ Could not find binary in downloaded file")
		return
	}

	// For Windows, we need admin privileges to replace the running executable
	fmt.Println("\n⚠️  Windows update requires administrator privileges")
	fmt.Println("📝 Extracted installer to: " + tempDir)
	fmt.Println("📝 Please run the installer as Administrator:")
	fmt.Printf("   %s\\install.bat\n", tempDir)
	fmt.Println("\n💡 Or manually:")
	fmt.Printf("   1. Download: %s\n", downloadURL)
	fmt.Println("   2. Extract the ZIP file")
	fmt.Println("   3. Right-click install.bat → Run as Administrator")
	fmt.Printf("\n📝 New version: %s\n", release.TagName)
}

func installUpdateDirect(release ReleaseInfo, downloadURL string) {
	// Download the file
	resp, err := http.Get(downloadURL)
	if err != nil {
		fmt.Printf("❌ Failed to download update: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Download failed with status: %s\n", resp.Status)
		return
	}

	// Create temporary directory for extraction
	tempDir := filepath.Join(os.TempDir(), "filemanager_update")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create temp directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Determine file extension and extract
	osName := runtime.GOOS
	var extractedBinary string

	if osName == "windows" {
		// Handle ZIP file
		zipPath := filepath.Join(tempDir, "update.zip")
		zipFile, err := os.Create(zipPath)
		if err != nil {
			fmt.Printf("❌ Failed to create temp file: %v\n", err)
			return
		}

		if _, err := io.Copy(zipFile, resp.Body); err != nil {
			zipFile.Close()
			fmt.Printf("❌ Failed to download file: %v\n", err)
			return
		}
		zipFile.Close()

		// Extract ZIP
		extractedBinary, err = extractZip(zipPath, tempDir)
		if err != nil {
			fmt.Printf("❌ Failed to extract update: %v\n", err)
			return
		}
	} else {
		// Handle TAR.GZ file
		tarPath := filepath.Join(tempDir, "update.tar.gz")
		tarFile, err := os.Create(tarPath)
		if err != nil {
			fmt.Printf("❌ Failed to create temp file: %v\n", err)
			return
		}

		if _, err := io.Copy(tarFile, resp.Body); err != nil {
			tarFile.Close()
			fmt.Printf("❌ Failed to download file: %v\n", err)
			return
		}
		tarFile.Close()

		// Extract TAR.GZ
		extractedBinary, err = extractTarGz(tarPath, tempDir)
		if err != nil {
			fmt.Printf("❌ Failed to extract update: %v\n", err)
			return
		}
	}

	if extractedBinary == "" {
		fmt.Println("❌ Could not find binary in downloaded file")
		return
	}

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		fmt.Printf("❌ Failed to get current executable path: %v\n", err)
		return
	}

	// For Windows, we need admin privileges
	if osName == "windows" {
		fmt.Println("\n⚠️  Windows update requires administrator privileges")
		fmt.Println("📝 Please run the installer manually or restart with admin rights")
		return
	}

	// For Linux/macOS, check if we need sudo
	if !isWritableDirectory(filepath.Dir(currentExe)) {
		fmt.Println("\n⚠️  Installation directory requires elevated privileges")
		fmt.Println("📝 Please run: sudo filemanager --update")
		return
	}

	// Backup current executable
	backupPath := currentExe + ".backup"
	if err := os.Rename(currentExe, backupPath); err != nil {
		fmt.Printf("❌ Failed to backup current executable: %v\n", err)
		return
	}

	// Copy new binary to current location
	newBinary, err := os.Open(extractedBinary)
	if err != nil {
		os.Rename(backupPath, currentExe) // Restore backup
		fmt.Printf("❌ Failed to open new binary: %v\n", err)
		return
	}
	defer newBinary.Close()

	newExe, err := os.Create(currentExe)
	if err != nil {
		os.Rename(backupPath, currentExe) // Restore backup
		fmt.Printf("❌ Failed to create new executable: %v\n", err)
		return
	}
	defer newExe.Close()

	if _, err := io.Copy(newExe, newBinary); err != nil {
		os.Rename(backupPath, currentExe) // Restore backup
		fmt.Printf("❌ Failed to copy new binary: %v\n", err)
		return
	}

	// Make executable
	if err := os.Chmod(currentExe, 0755); err != nil {
		os.Rename(backupPath, currentExe) // Restore backup
		fmt.Printf("❌ Failed to set permissions: %v\n", err)
		return
	}

	// Remove backup
	os.Remove(backupPath)

	fmt.Printf("\n✅ Update installed successfully!\n")
	fmt.Printf("📝 New version: %s\n", release.TagName)
	fmt.Println("🔄 Please restart the application to use the new version")
}

func isWritableDirectory(dir string) bool {
	testFile := filepath.Join(dir, ".filemanager_write_test")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err == nil {
		os.Remove(testFile)
		return true
	}
	return false
}

func extractZip(zipPath, destDir string) (string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var binaryPath string

	for _, file := range reader.File {
		filePath := filepath.Join(destDir, file.Name)

		if strings.HasSuffix(file.Name, "/") {
			os.MkdirAll(filePath, 0755)
			continue
		}

		// Create parent directories
		os.MkdirAll(filepath.Dir(filePath), 0755)

		// Extract file
		srcFile, err := file.Open()
		if err != nil {
			return "", err
		}

		dstFile, err := os.Create(filePath)
		if err != nil {
			srcFile.Close()
			return "", err
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			srcFile.Close()
			dstFile.Close()
			return "", err
		}

		srcFile.Close()
		dstFile.Close()

		// Look for the binary
		if strings.HasSuffix(file.Name, "filemanager.exe") || strings.HasSuffix(file.Name, "filemanager") {
			binaryPath = filePath
		}
	}

	return binaryPath, nil
}

func extractTarGz(tarPath, destDir string) (string, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	var binaryPath string

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		filePath := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(filePath, 0755)
		case tar.TypeReg:
			// Create parent directories
			os.MkdirAll(filepath.Dir(filePath), 0755)

			// Extract file
			outFile, err := os.Create(filePath)
			if err != nil {
				return "", err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return "", err
			}

			outFile.Close()

			// Look for the binary
			if strings.HasSuffix(header.Name, "filemanager") && !strings.Contains(header.Name, "/") {
				binaryPath = filePath
			}
		}
	}

	return binaryPath, nil
}
