package version

import (
	"fmt"
	"strconv"
	"strings"
)

// SemVer represents a semantic version (MAJOR.MINOR.PATCH)
type SemVer struct {
	Major int
	Minor int
	Patch int
}

// ChangeType represents the type of change in the update
type ChangeType int

const (
	ChangeTypePatch ChangeType = iota
	ChangeTypeMinor
	ChangeTypeMajor
)

// ParseSemVer parses a version string into SemVer
func ParseSemVer(versionStr string) (SemVer, error) {
	versionStr = strings.TrimPrefix(versionStr, "v")
	parts := strings.Split(versionStr, ".")

	if len(parts) < 3 {
		return SemVer{}, fmt.Errorf("invalid version format: %s", versionStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return SemVer{}, err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return SemVer{}, err
	}

	// Handle patch version with pre-release or metadata
	patchStr := strings.Split(parts[2], "-")[0]
	patch, err := strconv.Atoi(patchStr)
	if err != nil {
		return SemVer{}, err
	}

	return SemVer{Major: major, Minor: minor, Patch: patch}, nil
}

// String returns the string representation of SemVer
func (s SemVer) String() string {
	return fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
}

// Compare compares two semantic versions
// Returns: -1 if s < other, 0 if s == other, 1 if s > other
func (s SemVer) Compare(other SemVer) int {
	if s.Major != other.Major {
		if s.Major > other.Major {
			return 1
		}
		return -1
	}

	if s.Minor != other.Minor {
		if s.Minor > other.Minor {
			return 1
		}
		return -1
	}

	if s.Patch != other.Patch {
		if s.Patch > other.Patch {
			return 1
		}
		return -1
	}

	return 0
}

// DetermineChangeType determines what type of change occurred
func (s SemVer) DetermineChangeType(other SemVer) ChangeType {
	if s.Major != other.Major {
		return ChangeTypeMajor
	}
	if s.Minor != other.Minor {
		return ChangeTypeMinor
	}
	return ChangeTypePatch
}

// GetChangeTypeString returns a human-readable string for the change type
func GetChangeTypeString(changeType ChangeType) string {
	switch changeType {
	case ChangeTypePatch:
		return "PATCH"
	case ChangeTypeMinor:
		return "MINOR"
	case ChangeTypeMajor:
		return "MAJOR"
	default:
		return "UNKNOWN"
	}
}

// GetChangeTypeEmoji returns an emoji for the change type
func GetChangeTypeEmoji(changeType ChangeType) string {
	switch changeType {
	case ChangeTypePatch:
		return "🔧"
	case ChangeTypeMinor:
		return "✨"
	case ChangeTypeMajor:
		return "🚀"
	default:
		return "📦"
	}
}

// GetUpdateStrategy returns the recommended update strategy for a change type
func GetUpdateStrategy(changeType ChangeType) string {
	switch changeType {
	case ChangeTypePatch:
		return "Silent/Direct Install"
	case ChangeTypeMinor:
		return "Subtle In-App Banner/Hotspot"
	case ChangeTypeMajor:
		return "Modal Window/Full Screen Splash"
	default:
		return "Standard Notification"
	}
}

// GetUserImpact returns the user impact level for a change type
func GetUserImpact(changeType ChangeType) string {
	switch changeType {
	case ChangeTypePatch:
		return "Minimal"
	case ChangeTypeMinor:
		return "Low to Moderate"
	case ChangeTypeMajor:
		return "High"
	default:
		return "Unknown"
	}
}
