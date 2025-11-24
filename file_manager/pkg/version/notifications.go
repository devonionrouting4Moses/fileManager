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
	// Midnight Purple color scheme
	primary := "\033[38;5;219m"   // Pink (#ffb7c5)
	accent := "\033[38;5;198m"    // Hot pink (#ff69b4)
	secondary := "\033[38;5;135m" // Violet (#9d4edd)
	reset := "\033[0m"
	bold := "\033[1m"
	boxWidth := 60

	fmt.Println()
	fmt.Printf("%s%s┌%s┐%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ 🔧 PATCH UPDATE AVAILABLE%s │%s\n", primary, bold, strings.Repeat(" ", 32), reset)
	fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ %sv%s → v%s%s │%s\n", accent, bold, secondary, n.CurrentVersion, n.AvailableVersion, strings.Repeat(" ", boxWidth-len(fmt.Sprintf("v%s → v%s", n.CurrentVersion, n.AvailableVersion))-4), reset)
	fmt.Printf("%s%s│ %s✅ Safe, backwards-compatible update%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 19), reset)
	fmt.Printf("%s%s│ %s💡 Auto-installs on next restart%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 23), reset)
	fmt.Printf("%s%s└%s┘%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Println()

	if n.ReleaseNotes != "" {
		fmt.Println("📝 What's Fixed:")
		n.displayReleaseNotes()
		fmt.Println()
	}
}

// displaySubtleNotification displays an in-app banner for minor updates
func (n *UpdateNotification) displaySubtleNotification() {
	// Midnight Purple color scheme
	primary := "\033[38;5;219m"   // Pink (#ffb7c5)
	accent := "\033[38;5;198m"    // Hot pink (#ff69b4)
	secondary := "\033[38;5;135m" // Violet (#9d4edd)
	reset := "\033[0m"
	bold := "\033[1m"
	boxWidth := 60

	fmt.Println()
	fmt.Printf("%s%s┌%s┐%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ ✨ NEW FEATURES AVAILABLE%s │%s\n", primary, bold, strings.Repeat(" ", 32), reset)
	fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ %sv%s → v%s%s │%s\n", accent, bold, secondary, n.CurrentVersion, n.AvailableVersion, strings.Repeat(" ", boxWidth-len(fmt.Sprintf("v%s → v%s", n.CurrentVersion, n.AvailableVersion))-4), reset)
	fmt.Printf("%s%s│ %s📊 Update Type: MINOR (New Features)%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 19), reset)
	fmt.Printf("%s%s│ %s📈 User Impact: Low to Moderate%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 24), reset)
	fmt.Printf("%s%s│ %s💡 Update at your convenience%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 27), reset)
	fmt.Printf("%s%s└%s┘%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Println()

	if n.ReleaseNotes != "" {
		fmt.Println("📝 What's New:")
		n.displayReleaseNotes()
		fmt.Println()
	}

	fmt.Printf("📦 Download: %s\n\n", n.DownloadURL)
}

// displayModalNotification displays a full-screen splash for major updates
func (n *UpdateNotification) displayModalNotification() {
	// Midnight Purple color scheme
	primary := "\033[38;5;219m"   // Pink (#ffb7c5)
	accent := "\033[38;5;198m"    // Hot pink (#ff69b4)
	secondary := "\033[38;5;135m" // Violet (#9d4edd)
	reset := "\033[0m"
	bold := "\033[1m"
	boxWidth := 70

	fmt.Println()
	fmt.Printf("%s%s┌%s┐%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ 🚀 MAJOR UPGRADE AVAILABLE%s │%s\n", primary, bold, strings.Repeat(" ", 40), reset)
	fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ %sCurrent Version: v%s%s │%s\n", accent, bold, secondary, n.CurrentVersion, strings.Repeat(" ", boxWidth-len(fmt.Sprintf("Current Version: v%s", n.CurrentVersion))-4), reset)
	fmt.Printf("%s%s│ %sAvailable Version: v%s%s │%s\n", accent, bold, secondary, n.AvailableVersion, strings.Repeat(" ", boxWidth-len(fmt.Sprintf("Available Version: v%s", n.AvailableVersion))-4), reset)
	fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ %s⚠️  IMPORTANT: Major upgrade with breaking changes%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 15), reset)
	fmt.Printf("%s%s│ %s✅ Action Required: Review release notes first%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 18), reset)
	fmt.Printf("%s%s│ %s🔗 May need to reconfigure or migrate data%s │%s\n", accent, bold, secondary, strings.Repeat(" ", 22), reset)

	if n.ReleaseNotes != "" {
		fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
		lines := strings.Split(strings.TrimSpace(n.ReleaseNotes), "\n")
		for i, line := range lines {
			if i >= 5 { // Limit to 5 lines in modal
				fmt.Printf("%s%s│ %s... and more%s │%s\n", accent, bold, secondary, strings.Repeat(" ", boxWidth-18), reset)
				break
			}
			if len(line) > boxWidth-8 {
				line = line[:boxWidth-11] + "..."
			}
			fmt.Printf("%s%s│ %s• %s%s │%s\n", accent, bold, secondary, line, strings.Repeat(" ", boxWidth-len(line)-6), reset)
		}
	}

	fmt.Printf("%s%s├%s┤%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Printf("%s%s│ %s📦 Download: %s%s │%s\n", accent, bold, secondary, truncateURL(n.DownloadURL, 40), strings.Repeat(" ", boxWidth-len(truncateURL(n.DownloadURL, 40))-14), reset)
	fmt.Printf("%s%s└%s┘%s\n", primary, bold, strings.Repeat("─", boxWidth-2), reset)
	fmt.Println()

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
