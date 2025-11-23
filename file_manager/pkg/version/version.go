package version

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
	cacheFile    = filepath.Join(os.TempDir(), "filemanager_update_cache.json")
)

// ReleaseInfo represents GitHub release information
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

// UpdateCache stores cached update information
type UpdateCache struct {
	LastCheck   time.Time   `json:"last_check"`
	ReleaseInfo ReleaseInfo `json:"release_info"`
}

// ShowVersion displays version information
func ShowVersion() {
	fmt.Printf("%s v%s\n", AppName, Version)
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("Download: %s\n", GetDownloadURL())
}

// ShowBanner displays the application banner
func ShowBanner() {
	versionStr := fmt.Sprintf("v%s", Version)
	bannerWidth := 40

	contentStr := fmt.Sprintf("🗂️  %s %s (Rust+Go)", AppName, versionStr)
	visualLen := len(contentStr) - 2
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

// GetDownloadURL returns the download URL for the current platform
func GetDownloadURL() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

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

	var ext string
	if osName == "windows" {
		ext = "zip"
	} else {
		ext = "tar.gz"
	}

	return fmt.Sprintf("%s/v%s/filemanager-%s-%s.%s", downloadBase, Version, Version, target, ext)
}

// CheckForUpdates checks for available updates and handles them automatically
func CheckForUpdates() {
	CheckForUpdatesWithPrompt(true)
}

// CheckForUpdatesWithPrompt checks for updates with optional user prompts
func CheckForUpdatesWithPrompt(showPrompts bool) {
	fmt.Println("\n🔍 Checking for updates...")

	if isDevBuild() {
		fmt.Println("Development build - skipping update check")
		return
	}

	// Check cache first
	if cached, ok := loadCache(); ok {
		if time.Since(cached.LastCheck) < 24*time.Hour {
			fmt.Println("Using cached update information...")
			handleUpdateWithManager(cached.ReleaseInfo, showPrompts)
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

	handleUpdateWithManager(release, showPrompts)
}

// handleUpdateWithManager handles updates using the update manager
func handleUpdateWithManager(release ReleaseInfo, showPrompts bool) {
	currentVersion := "v" + Version
	comparison := compareVersions(release.TagName, currentVersion)

	if comparison > 0 {
		// Create notification
		notification, err := CreateUpdateNotification(release)
		if err != nil {
			fmt.Printf("⚠️  Could not process update: %v\n", err)
			return
		}

		if !showPrompts {
			// Just display info without prompts
			displayUpdateInfo(release)
			return
		}

		// Create update manager and handle update
		manager := NewUpdateManager(notification)
		updateApplied := manager.HandleUpdate()

		if updateApplied {
			fmt.Println(manager.GetUpdateSummary())
		}
	} else if comparison < 0 {
		fmt.Printf("ℹ️  You're running a pre-release version (v%s, latest stable: %s)\n", Version, release.TagName)
	} else {
		fmt.Printf("✅ You're running the latest version (v%s)\n", Version)
	}
}

func displayUpdateInfo(release ReleaseInfo) {
	currentVersion := "v" + Version
	comparison := compareVersions(release.TagName, currentVersion)

	if comparison > 0 {
		// Use new semantic versioning notification system
		notification, err := CreateUpdateNotification(release)
		if err != nil {
			// Fallback to old display if parsing fails
			fmt.Printf("\n🎉 New version available: %s (current: %s)\n", release.TagName, currentVersion)
			if release.Body != "" {
				bodyPreview := strings.TrimSpace(release.Body)
				if len(bodyPreview) > 200 {
					bodyPreview = bodyPreview[:197] + "..."
				}
				fmt.Printf("📝 %s\n\n", bodyPreview)
			}
			downloadURL := findAssetURL(release)
			if downloadURL != "" {
				fmt.Printf("📦 Download from:\n   %s\n\n", downloadURL)
			}
			return
		}

		// Display notification based on change type
		notification.DisplayNotification()

		// Show installation instructions
		fmt.Println("\n💡 Installation Instructions:")
		printInstallInstructions(release.TagName)

		// Show update summary
		fmt.Println(notification.GetUpdateSummary())
	} else if comparison < 0 {
		fmt.Printf("ℹ️  You're running a pre-release version (v%s, latest stable: %s)\n", Version, release.TagName)
	} else {
		fmt.Printf("✅ You're running the latest version (v%s)\n", Version)
	}
}

func findAssetURL(release ReleaseInfo) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

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

	var ext string
	if osName == "windows" {
		ext = ".zip"
	} else {
		ext = ".tar.gz"
	}

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, target) && strings.HasSuffix(asset.Name, ext) {
			return asset.DownloadURL
		}
	}

	return ""
}

func printInstallInstructions(version string) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

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
	if strings.Contains(Version, "dev") ||
		strings.Contains(Version, "SNAPSHOT") ||
		Version == "0.0.0" {
		return true
	}

	if len(os.Args) > 0 && strings.HasSuffix(os.Args[0], "go-build") {
		return true
	}

	return false
}

func compareVersions(v1, v2 string) int {
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
