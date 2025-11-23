package version

import (
	"fmt"
	"strings"
	"time"
)

// UpdateNotification represents a notification for an available update
type UpdateNotification struct {
	CurrentVersion   SemVer
	AvailableVersion SemVer
	ChangeType       ChangeType
	ReleaseNotes     string
	PublishedAt      time.Time
	DownloadURL      string
}

// NotificationStyle represents how to display the notification
type NotificationStyle int

const (
	NotificationStyleSilent NotificationStyle = iota
	NotificationStyleSubtle
	NotificationStyleModal
)

// GetNotificationStyle returns the appropriate notification style based on change type
func GetNotificationStyle(changeType ChangeType) NotificationStyle {
	switch changeType {
	case ChangeTypePatch:
		return NotificationStyleSilent
	case ChangeTypeMinor:
		return NotificationStyleSubtle
	case ChangeTypeMajor:
		return NotificationStyleModal
	default:
		return NotificationStyleSubtle
	}
}

// DisplayNotification displays the update notification based on the change type
func (n *UpdateNotification) DisplayNotification() {
	style := GetNotificationStyle(n.ChangeType)

	switch style {
	case NotificationStyleSilent:
		n.displaySilentNotification()
	case NotificationStyleSubtle:
		n.displaySubtleNotification()
	case NotificationStyleModal:
		n.displayModalNotification()
	}
}

// displaySilentNotification displays a minimal notification for patch updates
func (n *UpdateNotification) displaySilentNotification() {
	fmt.Println("\n" + strings.Repeat("─", 60))
	fmt.Printf("🔧 PATCH UPDATE AVAILABLE: v%s → v%s\n", n.CurrentVersion, n.AvailableVersion)
	fmt.Println("─ Security & Bug Fixes ─")
	fmt.Println("\n✅ This is a safe, backwards-compatible update.")
	fmt.Println("💡 It will be installed automatically on next restart.")
	fmt.Println(strings.Repeat("─", 60) + "\n")

	if n.ReleaseNotes != "" {
		fmt.Println("📝 What's Fixed:")
		n.displayReleaseNotes()
		fmt.Println()
	}
}

// displaySubtleNotification displays an in-app banner for minor updates
func (n *UpdateNotification) displaySubtleNotification() {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Printf("║ ✨ NEW FEATURES AVAILABLE: v%s → v%s\n", n.CurrentVersion, n.AvailableVersion)
	fmt.Println("║ " + strings.Repeat("─", 56))
	fmt.Println("║ 📊 Update Type: MINOR (New Features & Improvements)")
	fmt.Println("║ 📈 User Impact: Low to Moderate")
	fmt.Println("║ 🔄 Update Strategy: Subtle In-App Notification")
	fmt.Println("║ " + strings.Repeat("─", 56))
	fmt.Println("║")
	fmt.Println("║ 💡 Tip: Check the release notes to see what's new!")
	fmt.Println("║ 🔗 You can update at your convenience.")
	fmt.Println(strings.Repeat("═", 60) + "\n")

	if n.ReleaseNotes != "" {
		fmt.Println("📝 What's New:")
		n.displayReleaseNotes()
		fmt.Println()
	}

	fmt.Printf("📦 Download: %s\n\n", n.DownloadURL)
}

// displayModalNotification displays a full-screen splash for major updates
func (n *UpdateNotification) displayModalNotification() {
	// Clear screen effect
	fmt.Println("\n" + strings.Repeat("═", 70))
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + centerText("🚀 MAJOR UPGRADE AVAILABLE", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Printf("║  Current Version: v%s\n", n.CurrentVersion)
	fmt.Printf("║  Available Version: v%s\n", n.AvailableVersion)
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  ⚠️  IMPORTANT: This is a major upgrade with breaking changes.")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  📋 Key Changes:")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")

	if n.ReleaseNotes != "" {
		lines := strings.Split(strings.TrimSpace(n.ReleaseNotes), "\n")
		for i, line := range lines {
			if i >= 5 { // Limit to 5 lines in modal
				fmt.Println("║    ... and more")
				break
			}
			if len(line) > 60 {
				line = line[:57] + "..."
			}
			fmt.Printf("║    • %s\n", line)
		}
	}

	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  ✅ Action Required: Please review release notes before updating.")
	fmt.Println("║  🔗 You may need to reconfigure settings or migrate data.")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Printf("║  📦 Download: %s\n", truncateURL(n.DownloadURL, 60))
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println(strings.Repeat("═", 70) + "\n")

	fmt.Println("📚 Full Release Notes:")
	if n.ReleaseNotes != "" {
		n.displayReleaseNotes()
	}
	fmt.Println()
}

// displayReleaseNotes formats and displays release notes
func (n *UpdateNotification) displayReleaseNotes() {
	if n.ReleaseNotes == "" {
		fmt.Println("  No release notes available")
		return
	}

	notes := strings.TrimSpace(n.ReleaseNotes)

	// Parse release notes by category
	categories := map[string]string{
		"✨ New Features":       "✨",
		"🔧 Improvements":       "🔧",
		"🐛 Bug Fixes":          "🐛",
		"⚠️  Breaking Changes": "⚠️",
		"📚 Documentation":      "📚",
	}

	lines := strings.Split(notes, "\n")
	var currentCategory string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this line starts a new category
		isCategory := false
		for category, emoji := range categories {
			if strings.Contains(line, emoji) || strings.HasPrefix(line, category) {
				if currentCategory != category {
					currentCategory = category
					isCategory = true
				}
				break
			}
		}

		if !isCategory && currentCategory != "" {
			fmt.Printf("  %s\n", line)
		}
	}
}

// centerText centers text within a given width
func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}

// truncateURL truncates a URL to fit within a width
func truncateURL(url string, width int) string {
	if len(url) <= width {
		return url
	}
	return url[:width-3] + "..."
}

// CreateUpdateNotification creates an update notification from release info
func CreateUpdateNotification(release ReleaseInfo) (*UpdateNotification, error) {
	currentVer, err := ParseSemVer(Version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current version: %w", err)
	}

	availableVer, err := ParseSemVer(release.TagName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse available version: %w", err)
	}

	changeType := currentVer.DetermineChangeType(availableVer)

	// Parse published date
	publishedAt, _ := time.Parse(time.RFC3339, release.PublishedAt)

	return &UpdateNotification{
		CurrentVersion:   currentVer,
		AvailableVersion: availableVer,
		ChangeType:       changeType,
		ReleaseNotes:     release.Body,
		PublishedAt:      publishedAt,
		DownloadURL:      findAssetURL(release),
	}, nil
}

// GetUpdateSummary returns a summary of the update
func (n *UpdateNotification) GetUpdateSummary() string {
	return fmt.Sprintf(`
Update Summary:
  Current Version: v%s
  Available Version: v%s
  Change Type: %s %s
  User Impact: %s
  Update Strategy: %s
  Published: %s
`,
		n.CurrentVersion,
		n.AvailableVersion,
		GetChangeTypeEmoji(n.ChangeType),
		GetChangeTypeString(n.ChangeType),
		GetUserImpact(n.ChangeType),
		GetUpdateStrategy(n.ChangeType),
		n.PublishedAt.Format("2006-01-02 15:04:05"),
	)
}
