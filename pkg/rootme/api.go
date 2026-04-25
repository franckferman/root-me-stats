// Package rootme provides a high-level API for Root-me statistics and badges
package rootme

import (
	"github.com/franckferman/root-me-stats/internal/cache"
	"github.com/franckferman/root-me-stats/internal/fetcher"
	"github.com/franckferman/root-me-stats/internal/generator"
	"github.com/franckferman/root-me-stats/internal/themes"
)

// Profile is an alias for the internal profile type
type Profile = fetcher.Profile

// Comparison is an alias for the internal comparison type
type Comparison = fetcher.Comparison

// BadgeOptions is an alias for the internal badge options type
type BadgeOptions = generator.BadgeOptions

// GetProfile fetches a Root-me profile (uses cache when available)
func GetProfile(nickname string) (*Profile, error) {
	return cache.GetOrFetch(nickname)
}

// CompareProfiles compares two Root-me profiles
func CompareProfiles(user1, user2 string) (*Comparison, error) {
	return fetcher.CompareProfiles(user1, user2)
}

// GenerateBadge creates an SVG badge for a profile
func GenerateBadge(profile *Profile, opts ...BadgeOptions) string {
	var options BadgeOptions
	if len(opts) > 0 {
		options = opts[0]
	} else {
		options = generator.DefaultBadgeOptions()
	}

	return generator.GenerateBadge(profile, options)
}

// GenerateComparisonBadge creates an SVG badge comparing two profiles
func GenerateComparisonBadge(comparison *Comparison, theme string, width int) string {
	return generator.GenerateComparisonBadge(comparison, theme, width)
}

// QuickBadge fetches a profile and generates a badge in one call
func QuickBadge(nickname string, opts ...BadgeOptions) (string, error) {
	profile, err := GetProfile(nickname)
	if err != nil {
		return "", err
	}

	return GenerateBadge(profile, opts...), nil
}

// QuickComparison compares profiles and generates a badge in one call
func QuickComparison(user1, user2, theme string, width int) (string, error) {
	comparison, err := CompareProfiles(user1, user2)
	if err != nil {
		return "", err
	}

	return GenerateComparisonBadge(comparison, theme, width), nil
}

// GetThemes returns all available theme names
func GetThemes() []string {
	return themes.GetThemeNames()
}

// ClearCache removes all cached data
func ClearCache() error {
	return cache.Clear()
}

// DefaultBadgeOptions returns sensible default options
func DefaultBadgeOptions() BadgeOptions {
	return generator.DefaultBadgeOptions()
}
