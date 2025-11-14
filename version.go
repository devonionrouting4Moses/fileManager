package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		fmt.Printf("\n✨ New version available: %s (current: %s)\n", release.TagName, currentVersion)
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
			fmt.Printf("⬇️  Download: %s\n", downloadURL)
		} else {
			releaseURL := fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", repoOwner, repoName, release.TagName)
			fmt.Printf("🌐 Visit: %s\n", releaseURL)
		}

		fmt.Println("\n💡 Installation instructions:")
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