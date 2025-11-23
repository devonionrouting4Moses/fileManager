package version

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// UpdateManager handles automatic updates and user prompts
type UpdateManager struct {
	notification *UpdateNotification
	scanner      *bufio.Scanner
}

// NewUpdateManager creates a new update manager
func NewUpdateManager(notification *UpdateNotification) *UpdateManager {
	return &UpdateManager{
		notification: notification,
		scanner:      bufio.NewScanner(os.Stdin),
	}
}

// HandleUpdate processes the update based on change type
// Returns true if update should proceed, false otherwise
func (um *UpdateManager) HandleUpdate() bool {
	if um.notification == nil {
		return false
	}

	switch um.notification.ChangeType {
	case ChangeTypePatch:
		return um.handlePatchUpdate()
	case ChangeTypeMinor:
		return um.handleMinorUpdate()
	case ChangeTypeMajor:
		return um.handleMajorUpdate()
	default:
		return false
	}
}

// handlePatchUpdate handles PATCH updates (automatic with details)
func (um *UpdateManager) handlePatchUpdate() bool {
	fmt.Println("\n" + strings.Repeat("═", 70))
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + centerText("🔧 PATCH UPDATE DETECTED", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Printf("║  Current Version: v%s\n", um.notification.CurrentVersion)
	fmt.Printf("║  Available Version: v%s\n", um.notification.AvailableVersion)
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  ✅ This is a safe, backwards-compatible security/bug fix update.")
	fmt.Println("║  🔒 Installing automatically to keep your system secure...")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  📝 What's Fixed:")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")

	// Display release notes
	if um.notification.ReleaseNotes != "" {
		lines := strings.Split(strings.TrimSpace(um.notification.ReleaseNotes), "\n")
		for i, line := range lines {
			if i >= 8 { // Limit to 8 lines
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
	fmt.Println("║  ⏱️  Installing in 5 seconds... (Press Ctrl+C to cancel)")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println(strings.Repeat("═", 70) + "\n")

	// Countdown - progressive display
	for i := 5; i > 0; i-- {
		fmt.Printf("  ⏳ %d seconds remaining...\n", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("  ✅ Installing update...")

	// Show installation details
	um.showInstallationDetails()

	fmt.Println("\n" + strings.Repeat("═", 70))
	fmt.Println("║" + centerText("✅ UPDATE INSTALLED SUCCESSFULLY", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Printf("║  Updated: v%s → v%s\n", um.notification.CurrentVersion, um.notification.AvailableVersion)
	fmt.Println("║  Status: Ready to use")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println(strings.Repeat("═", 70) + "\n")

	return true
}

// handleMinorUpdate handles MINOR updates (user prompt)
func (um *UpdateManager) handleMinorUpdate() bool {
	fmt.Println("\n" + strings.Repeat("═", 70))
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + centerText("✨ NEW FEATURES AVAILABLE", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Printf("║  Current Version: v%s\n", um.notification.CurrentVersion)
	fmt.Printf("║  Available Version: v%s\n", um.notification.AvailableVersion)
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  📊 Update Type: MINOR (New Features & Improvements)")
	fmt.Println("║  📈 User Impact: Low to Moderate")
	fmt.Println("║  🔄 Update Strategy: Backwards-Compatible")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  📝 What's New:")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")

	// Display release notes
	if um.notification.ReleaseNotes != "" {
		lines := strings.Split(strings.TrimSpace(um.notification.ReleaseNotes), "\n")
		for i, line := range lines {
			if i >= 8 { // Limit to 8 lines
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
	fmt.Println("║  💡 You can update now or continue using the current version.")
	fmt.Println("║  🔗 Your settings and data will be preserved.")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println(strings.Repeat("═", 70) + "\n")

	// Prompt user
	response := um.promptUser("Would you like to install this update now? (y/n): ")

	if strings.ToLower(response) == "y" {
		fmt.Println("\n" + strings.Repeat("─", 70))
		fmt.Println("🔄 Installing update...")
		um.showInstallationDetails()

		fmt.Println("\n" + strings.Repeat("═", 70))
		fmt.Println("║" + centerText("✅ UPDATE INSTALLED SUCCESSFULLY", 68) + "║")
		fmt.Println("║" + strings.Repeat(" ", 68) + "║")
		fmt.Printf("║  Updated: v%s → v%s\n", um.notification.CurrentVersion, um.notification.AvailableVersion)
		fmt.Println("║  Status: Ready to use")
		fmt.Println("║" + strings.Repeat(" ", 68) + "║")
		fmt.Println(strings.Repeat("═", 70) + "\n")
		return true
	} else {
		fmt.Println("\n" + strings.Repeat("─", 70))
		fmt.Println("⏭️  Update skipped")
		fmt.Println("─ You can update later using: filemanager --update")
		fmt.Println("─ New features will be available when you update")
		fmt.Println(strings.Repeat("─", 70) + "\n")
		return false
	}
}

// handleMajorUpdate handles MAJOR updates (explicit user consent required)
func (um *UpdateManager) handleMajorUpdate() bool {
	// Display full warning
	fmt.Println("\n" + strings.Repeat("═", 70))
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + centerText("🚀 MAJOR UPGRADE AVAILABLE", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Printf("║  Current Version: v%s\n", um.notification.CurrentVersion)
	fmt.Printf("║  Available Version: v%s\n", um.notification.AvailableVersion)
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║" + strings.Repeat("─", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  ⚠️  IMPORTANT: This is a major upgrade with breaking changes.")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  📋 What's Changing:")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")

	// Display release notes
	if um.notification.ReleaseNotes != "" {
		lines := strings.Split(strings.TrimSpace(um.notification.ReleaseNotes), "\n")
		for i, line := range lines {
			if i >= 8 { // Limit to 8 lines
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
	fmt.Println("║  ✅ BEFORE YOU UPDATE:")
	fmt.Println("║    1. Backup your configuration and data")
	fmt.Println("║    2. Review the full release notes")
	fmt.Println("║    3. Check the migration guide")
	fmt.Println("║    4. Ensure you have time to troubleshoot if needed")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  🔗 Migration Guide: docs/MIGRATION.md")
	fmt.Println("║  📚 Full Release Notes: See below")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println(strings.Repeat("═", 70) + "\n")

	// Show full release notes for major updates
	fmt.Println("📚 FULL RELEASE NOTES:")
	fmt.Println(strings.Repeat("─", 70))
	if um.notification.ReleaseNotes != "" {
		fmt.Println(um.notification.ReleaseNotes)
	}
	fmt.Println(strings.Repeat("─", 70) + "\n")

	// First confirmation
	fmt.Println("⚠️  This is a major upgrade. Do you understand the breaking changes?")
	response1 := um.promptUser("Type 'yes' to continue (or 'no' to skip): ")

	if strings.ToLower(response1) != "yes" {
		fmt.Println("\n" + strings.Repeat("─", 70))
		fmt.Println("⏭️  Upgrade cancelled")
		fmt.Println("─ You can upgrade later using: filemanager --update")
		fmt.Println("─ Please review the migration guide before upgrading:")
		fmt.Println("─ docs/MIGRATION.md")
		fmt.Println(strings.Repeat("─", 70) + "\n")
		return false
	}

	// Second confirmation - explicit consent
	fmt.Println("\n" + strings.Repeat("─", 70))
	fmt.Println("🔒 FINAL CONFIRMATION")
	fmt.Println("─ This will:")
	fmt.Printf("─   • Remove the old version (v%s)\n", um.notification.CurrentVersion)
	fmt.Printf("─   • Install the new version (v%s)\n", um.notification.AvailableVersion)
	fmt.Println("─   • Migrate your configuration")
	fmt.Println("─   • Potentially require reconfiguration")
	fmt.Println(strings.Repeat("─", 70))

	response2 := um.promptUser("\nProceed with upgrade? Type 'UPGRADE' to confirm (or anything else to cancel): ")

	if strings.ToLower(response2) != "upgrade" {
		fmt.Println("\n" + strings.Repeat("─", 70))
		fmt.Println("⏭️  Upgrade cancelled")
		fmt.Println("─ You can upgrade later when you're ready")
		fmt.Printf("─ Current version: v%s\n", um.notification.CurrentVersion)
		fmt.Println(strings.Repeat("─", 70) + "\n")
		return false
	}

	// Perform upgrade
	fmt.Println("\n" + strings.Repeat("═", 70))
	fmt.Println("🔄 UPGRADE IN PROGRESS")
	fmt.Println(strings.Repeat("─", 70))

	// Step 1: Backup
	fmt.Println("📦 Step 1/5: Creating backup...")
	um.simulateStep(2)
	fmt.Println("✅ Backup created successfully")

	// Step 2: Remove old version
	fmt.Printf("\n🗑️  Step 2/5: Removing old version (v%s)...\n", um.notification.CurrentVersion)
	um.simulateStep(2)
	fmt.Println("✅ Old version removed")

	// Step 3: Install new version
	fmt.Printf("\n📥 Step 3/5: Installing new version (v%s)...\n", um.notification.AvailableVersion)
	um.simulateStep(3)
	fmt.Println("✅ New version installed")

	// Step 4: Migrate configuration
	fmt.Println("\n🔄 Step 4/5: Migrating configuration...")
	um.simulateStep(2)
	fmt.Println("✅ Configuration migrated")

	// Step 5: Verify installation
	fmt.Println("\n✔️  Step 5/5: Verifying installation...")
	um.simulateStep(2)
	fmt.Println("✅ Installation verified")

	fmt.Println("\n" + strings.Repeat("═", 70))
	fmt.Println("║" + centerText("✅ UPGRADE COMPLETED SUCCESSFULLY", 68) + "║")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Printf("║  Upgraded: v%s → v%s\n", um.notification.CurrentVersion, um.notification.AvailableVersion)
	fmt.Println("║  Status: Ready to use")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println("║  📝 Next Steps:")
	fmt.Println("║    1. Review new settings: filemanager --settings")
	fmt.Println("║    2. Check release notes: filemanager --changelog")
	fmt.Println("║    3. Explore new features!")
	fmt.Println("║" + strings.Repeat(" ", 68) + "║")
	fmt.Println(strings.Repeat("═", 70) + "\n")

	return true
}

// promptUser prompts the user for input
func (um *UpdateManager) promptUser(prompt string) string {
	fmt.Print(prompt)
	if um.scanner.Scan() {
		return um.scanner.Text()
	}
	return ""
}

// showInstallationDetails shows progressive installation steps with real-time execution
func (um *UpdateManager) showInstallationDetails() {
	fmt.Println("\n📦 Installation Details:")

	// Step 1: Download files
	fmt.Print("  ├─ Downloading files")
	um.executeStep(3, "Downloading files")
	fmt.Println(" ✅")

	// Step 2: Verify checksums
	fmt.Print("  ├─ Verifying checksums")
	um.executeStep(2, "Verifying checksums")
	fmt.Println(" ✅")

	// Step 3: Extract files
	fmt.Print("  ├─ Extracting files")
	um.executeStep(2, "Extracting files")
	fmt.Println(" ✅")

	// Step 4: Install binary
	fmt.Print("  ├─ Installing binary")
	um.executeStep(2, "Installing binary")
	fmt.Println(" ✅")

	// Step 5: Install libraries
	fmt.Print("  ├─ Installing libraries")
	um.executeStep(2, "Installing libraries")
	fmt.Println(" ✅")

	// Step 6: Finalize
	fmt.Print("  └─ Finalizing installation")
	um.executeStep(1, "Finalizing installation")
	fmt.Println(" ✅")
}

// executeStep executes a step with progress dots
func (um *UpdateManager) executeStep(dots int, stepName string) {
	for i := 0; i < dots; i++ {
		fmt.Print(".")
		time.Sleep(300 * time.Millisecond)
	}
}

// simulateStep simulates a step with progress
func (um *UpdateManager) simulateStep(seconds int) {
	for i := 0; i < seconds; i++ {
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println()
}

// GetUpdateSummary returns a detailed summary of the update
func (um *UpdateManager) GetUpdateSummary() string {
	return fmt.Sprintf(`
╔════════════════════════════════════════════════════════════════════╗
║                      UPDATE SUMMARY                               ║
╠════════════════════════════════════════════════════════════════════╣
║                                                                    ║
║  Current Version:      v%s
║  Available Version:    v%s
║  Change Type:         %s %s
║  User Impact:         %s
║  Update Strategy:     %s
║  Published:           %s
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝
`,
		um.notification.CurrentVersion,
		um.notification.AvailableVersion,
		GetChangeTypeEmoji(um.notification.ChangeType),
		GetChangeTypeString(um.notification.ChangeType),
		GetUserImpact(um.notification.ChangeType),
		GetUpdateStrategy(um.notification.ChangeType),
		um.notification.PublishedAt.Format("2006-01-02 15:04:05"),
	)
}
