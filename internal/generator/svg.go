// Package generator handles SVG badge generation
package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/franckferman/root-me-stats/internal/fetcher"
	"github.com/franckferman/root-me-stats/internal/themes"
)

// BadgeOptions configures SVG badge generation
type BadgeOptions struct {
	Theme           string
	ShowGlobalStats bool
	Width           int
	MaxCategories   int
}

// DefaultBadgeOptions returns sensible defaults
func DefaultBadgeOptions() BadgeOptions {
	return BadgeOptions{
		Theme:           "dark",
		ShowGlobalStats: false,
		Width:           380,
		MaxCategories:   10,
	}
}

// GenerateBadge creates an SVG badge for a Root-me profile
func GenerateBadge(profile *fetcher.Profile, opts BadgeOptions) string {
	theme := themes.GetTheme(opts.Theme)
	if theme == nil {
		theme = themes.GetTheme("dark") // Fallback
	}

	categories := profile.Categories
	if len(categories) > opts.MaxCategories {
		categories = categories[:opts.MaxCategories]
	}

	// Calculate dimensions
	baseHeight := 120
	categoryHeight := len(categories) * 35
	statsHeight := 0
	if opts.ShowGlobalStats {
		statsHeight = 80
	}
	height := baseHeight + categoryHeight + statsHeight + 40

	var svg strings.Builder

	// SVG header
	svg.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		opts.Width, height, opts.Width, height))

	// CSS styles
	svg.WriteString(generateStyles(theme))

	// Background
	svg.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" rx="8" fill="%s" />`, theme.Background))

	// Header
	svg.WriteString(generateHeader(profile.Nickname, profile.Stats, theme))

	// Categories
	svg.WriteString(generateCategories(categories, theme))

	// Global stats footer
	if opts.ShowGlobalStats {
		svg.WriteString(generateFooter(profile.Stats, theme, height, opts.Width))
	}

	// Border
	svg.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" rx="8" fill="none" stroke="%s" stroke-width="1" />`, theme.Bar))

	svg.WriteString(`</svg>`)

	return svg.String()
}

// GenerateComparisonBadge creates an SVG for comparing two profiles
func GenerateComparisonBadge(comparison *fetcher.Comparison, themeName string, width int) string {
	if width == 0 {
		width = 500
	}

	theme := themes.GetTheme(themeName)
	if theme == nil {
		theme = themes.GetTheme("dark")
	}

	height := 200

	var svg strings.Builder

	svg.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		width, height, width, height))

	svg.WriteString(generateStyles(theme))

	// Background
	svg.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" rx="8" fill="%s" />`, theme.Background))
	svg.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" rx="8" fill="none" stroke="%s" stroke-width="1" />`, theme.Bar))

	// Header
	svg.WriteString(fmt.Sprintf(`<text x="20" y="30" fill="%s" font="600 16px 'Segoe UI', sans-serif">Root-me Comparison</text>`, theme.Title))

	// User 1
	svg.WriteString(generateUserComparison(comparison.User1, 20, theme, false))

	// VS
	svg.WriteString(fmt.Sprintf(`<text x="%d" y="90" fill="%s" font="600 18px 'Segoe UI', sans-serif" text-anchor="middle">VS</text>`,
		width/2, theme.Accent))

	// User 2
	svg.WriteString(generateUserComparison(comparison.User2, width-20, theme, true))

	// Diff stats
	svg.WriteString(fmt.Sprintf(`<text x="20" y="160" fill="%s" font="400 11px 'Segoe UI', sans-serif">Rank Diff: %s%d</text>`,
		theme.Text, formatDiff(comparison.Diff.Rank), comparison.Diff.Rank))
	svg.WriteString(fmt.Sprintf(`<text x="20" y="175" fill="%s" font="400 11px 'Segoe UI', sans-serif">Score Diff: %s%s</text>`,
		theme.Text, formatDiff(comparison.Diff.Score), formatNumber(comparison.Diff.Score)))

	svg.WriteString(`</svg>`)

	return svg.String()
}

// generateStyles creates the CSS styles for the SVG
func generateStyles(theme *themes.Theme) string {
	return fmt.Sprintf(`
  <defs>
    <style>
      .text-title { fill: %s; font: 600 18px 'Segoe UI', sans-serif; }
      .text-normal { fill: %s; font: 400 14px 'Segoe UI', sans-serif; }
      .text-small { fill: %s; font: 400 12px 'Segoe UI', sans-serif; }
      .icon { fill: %s; }

      @keyframes fadeIn {
        from { opacity: 0; }
        to { opacity: 1; }
      }

      @keyframes progressFill {
        from { width: 0; }
        to { width: var(--progress); }
      }

      .animated { animation: fadeIn 0.8s ease-in-out; }
      .progress-bar { animation: progressFill 1.2s ease-out; }
    </style>
  </defs>
`, theme.Title, theme.Text, theme.Text, theme.Accent)
}

// generateHeader creates the profile header
func generateHeader(nickname string, stats fetcher.Stats, theme *themes.Theme) string {
	return fmt.Sprintf(`
  <g class="animated">
    <text x="20" y="35" class="text-title">Root-me Stats</text>
    <text x="20" y="55" class="text-normal">@%s</text>
    <text x="20" y="75" class="text-small">Rank: %s • Score: %s</text>
  </g>
`, nickname, formatNumber(stats.Rank), formatNumber(stats.Score))
}

// generateCategories creates the category progress bars
func generateCategories(categories []fetcher.Category, theme *themes.Theme) string {
	var result strings.Builder

	for i, category := range categories {
		y := 110 + (i * 35)
		iconPath := themes.GetCategoryIcon(category.Name)
		progressWidth := float64(category.Percentage) / 100.0 * 150

		result.WriteString(fmt.Sprintf(`
    <g class="animated" style="animation-delay: %fs">
      <!-- Category Icon -->
      <svg x="20" y="%d" width="16" height="16" viewBox="0 0 24 24">
        <path d="%s" class="icon" />
      </svg>

      <!-- Category Name -->
      <text x="45" y="%d" class="text-small">%s</text>

      <!-- Progress Bar Background -->
      <rect x="200" y="%d" width="150" height="8" rx="4" fill="%s" />

      <!-- Progress Bar Fill -->
      <rect x="200" y="%d" width="%.0f" height="8" rx="4" fill="%s" class="progress-bar" />

      <!-- Percentage -->
      <text x="360" y="%d" class="text-small">%d%%</text>
    </g>
`, float64(i)*0.1, y-12, iconPath, y, category.Name, y-8, theme.Bar, y-8, progressWidth, theme.Accent, y, category.Percentage))
	}

	return result.String()
}

// generateFooter creates the global stats footer
func generateFooter(stats fetcher.Stats, theme *themes.Theme, height, width int) string {
	return fmt.Sprintf(`
  <g class="animated" style="animation-delay: 1s">
    <line x1="20" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1" />
    <text x="20" y="%d" class="text-small">Challenges: %s</text>
    <text x="20" y="%d" class="text-small">Compromissions: %s</text>
    <text x="%d" y="%d" class="text-small" text-anchor="end">Updated: %s</text>
  </g>
`, height-65, width-20, height-65, theme.Bar,
		height-40, formatNumber(stats.Challenges),
		height-25, formatNumber(stats.Compromissions),
		width-20, height-15, time.Now().Format("2006-01-02"))
}

// generateUserComparison creates comparison text for a user
func generateUserComparison(profile fetcher.Profile, x int, theme *themes.Theme, rightAlign bool) string {
	anchor := ""
	if rightAlign {
		anchor = ` text-anchor="end"`
	}

	return fmt.Sprintf(`
  <g>
    <text x="%d" y="60" fill="%s" font="400 13px 'Segoe UI', sans-serif"%s>@%s</text>
    <text x="%d" y="80" fill="%s" font="400 11px 'Segoe UI', sans-serif"%s>Rank: %s</text>
    <text x="%d" y="100" fill="%s" font="400 11px 'Segoe UI', sans-serif"%s>Score: %s</text>
    <text x="%d" y="120" fill="%s" font="400 11px 'Segoe UI', sans-serif"%s>Challenges: %s</text>
  </g>
`, x, theme.Text, anchor, profile.Nickname,
		x, theme.Text, anchor, formatNumber(profile.Stats.Rank),
		x, theme.Text, anchor, formatNumber(profile.Stats.Score),
		x, theme.Text, anchor, formatNumber(profile.Stats.Challenges))
}

// formatNumber adds commas to numbers for display
func formatNumber(n int) string {
	if n == 0 {
		return "0"
	}

	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}

	var result []string
	for i := len(str); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{str[start:i]}, result...)
	}

	return strings.Join(result, ",")
}

// formatDiff formats a difference with + or - sign
func formatDiff(n int) string {
	if n > 0 {
		return "+"
	}
	return ""
}
