package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	Version     = "0.1.0"
	AppName     = "FileManager"
	ReleaseURL  = "https://api.github.com/repos/devonionrouting4Moses/fileManager/releases/latest"
	DownloadURL = "https://github.com/devonionrouting4Moses/fileManager/releases/latest"
)

type ReleaseInfo struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	PublishedAt string `json:"published_at"`
	Assets     []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func ShowVersion() {
	fmt.Printf("%s v%s\n", AppName, Version)
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Go version: %s\n", runtime.Version())
}

func CheckForUpdates() {
	fmt.Println("\n🔍 Checking for updates...")
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(ReleaseURL)
	if err != nil {
		fmt.Printf("⚠️  Could not check for updates: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		fmt.Printf("⚠️  Could not check for updates (HTTP %d)\n", resp.StatusCode)
		return
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("⚠️  Could not read update info: %v\n", err)
		return
	}
	
	var release ReleaseInfo
	if err := json.Unmarshal(body, &release); err != nil {
		fmt.Printf("⚠️  Could not parse update info: %v\n", err)
		return
	}
	
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := Version
	
	if compareVersions(latestVersion, currentVersion) > 0 {
		fmt.Printf("\n🎉 New version available: v%s (current: v%s)\n", latestVersion, currentVersion)
		fmt.Printf("📝 Release: %s\n", release.Name)
		fmt.Printf("📅 Published: %s\n", release.PublishedAt)
		fmt.Println("\n📦 Download from:")
		fmt.Printf("   %s\n", DownloadURL)
		
		// Show platform-specific download
		platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
		for _, asset := range release.Assets {
			if strings.Contains(asset.Name, platform) {
				fmt.Printf("\n💾 Direct download for your platform:\n")
				fmt.Printf("   %s\n", asset.BrowserDownloadURL)
				break
			}
		}
	} else {
		fmt.Printf("✅ You're running the latest version (v%s)\n", Version)
	}
}

func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	
	for i := 0; i < 3; i++ {
		var n1, n2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
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
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Printf("║   🗂️  %s v%-6s (Rust+Go)     ║\n", AppName, Version)
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
}