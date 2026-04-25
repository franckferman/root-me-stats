// Root-me Stats CLI Tool
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/franckferman/root-me-stats/pkg/rootme"
)

const (
	version = "2.0.0"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "badge":
		runBadgeCommand()
	case "compare":
		runCompareCommand()
	case "profile":
		runProfileCommand()
	case "themes":
		runThemesCommand()
	case "version":
		fmt.Printf("root-me-stats v%s\n", version)
	case "help", "--help", "-h":
		showUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		showUsage()
		os.Exit(1)
	}
}

// runBadgeCommand generates SVG badges
func runBadgeCommand() {
	fs := flag.NewFlagSet("badge", flag.ExitOnError)

	nickname := fs.String("nickname", "", "Root-me username (required)")
	output := fs.String("output", "", "Output file path (required)")
	theme := fs.String("theme", "dark", "Theme name")
	showStats := fs.Bool("stats", false, "Include global stats")
	width := fs.Int("width", 380, "Badge width")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s badge [options]\n\n", os.Args[0])
		fmt.Fprintf(fs.Output(), "Generate SVG badge for a Root-me profile\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nExample:\n")
		fmt.Fprintf(fs.Output(), "  %s badge --nickname=franckferman --output=./badge.svg --theme=dark --stats\n", os.Args[0])
	}

	fs.Parse(os.Args[2:])

	if *nickname == "" {
		fmt.Fprintln(os.Stderr, "Error: --nickname is required")
		fs.Usage()
		os.Exit(1)
	}

	if *output == "" {
		fmt.Fprintln(os.Stderr, "Error: --output is required")
		fs.Usage()
		os.Exit(1)
	}

	// Validate theme
	themes := rootme.GetThemes()
	themeValid := false
	for _, validTheme := range themes {
		if *theme == validTheme {
			themeValid = true
			break
		}
	}

	if !themeValid {
		fmt.Fprintf(os.Stderr, "Error: Invalid theme '%s'. Available: %s\n", *theme, strings.Join(themes, ", "))
		os.Exit(1)
	}

	fmt.Printf("🔍 Fetching profile: %s\n", *nickname)

	// Generate badge
	opts := rootme.DefaultBadgeOptions()
	opts.Theme = *theme
	opts.ShowGlobalStats = *showStats
	opts.Width = *width

	svg, err := rootme.QuickBadge(*nickname, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Ensure output directory exists
	if err := ensureDir(filepath.Dir(*output)); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Write SVG file
	if err := os.WriteFile(*output, []byte(svg), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error writing file: %v\n", err)
		os.Exit(1)
	}

	// Get profile for stats display
	profile, _ := rootme.GetProfile(*nickname)
	if profile != nil {
		fmt.Printf("✨ Generating badge (theme: %s)\n", *theme)
		fmt.Printf("✅ Badge saved to: %s\n", *output)
		fmt.Printf("📊 Profile stats: Rank %s • Score %s • %s challenges\n",
			formatNumber(profile.Stats.Rank),
			formatNumber(profile.Stats.Score),
			formatNumber(profile.Stats.Challenges))
	} else {
		fmt.Printf("✅ Badge saved to: %s\n", *output)
	}
}

// runCompareCommand generates comparison badges
func runCompareCommand() {
	fs := flag.NewFlagSet("compare", flag.ExitOnError)

	user1 := fs.String("user1", "", "First username (required)")
	user2 := fs.String("user2", "", "Second username (required)")
	output := fs.String("output", "", "Output file path (required)")
	theme := fs.String("theme", "dark", "Theme name")
	width := fs.Int("width", 500, "Badge width")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s compare [options]\n\n", os.Args[0])
		fmt.Fprintf(fs.Output(), "Generate comparison badge between two Root-me profiles\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nExample:\n")
		fmt.Fprintf(fs.Output(), "  %s compare --user1=user1 --user2=user2 --output=./compare.svg --theme=midnight\n", os.Args[0])
	}

	fs.Parse(os.Args[2:])

	if *user1 == "" || *user2 == "" {
		fmt.Fprintln(os.Stderr, "Error: Both --user1 and --user2 are required")
		fs.Usage()
		os.Exit(1)
	}

	if *output == "" {
		fmt.Fprintln(os.Stderr, "Error: --output is required")
		fs.Usage()
		os.Exit(1)
	}

	fmt.Printf("🔍 Fetching profiles: %s vs %s\n", *user1, *user2)

	svg, err := rootme.QuickComparison(*user1, *user2, *theme, *width)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Ensure output directory exists
	if err := ensureDir(filepath.Dir(*output)); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Write SVG file
	if err := os.WriteFile(*output, []byte(svg), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error writing file: %v\n", err)
		os.Exit(1)
	}

	// Get comparison for stats display
	comparison, _ := rootme.CompareProfiles(*user1, *user2)
	if comparison != nil {
		fmt.Printf("✨ Generating comparison (theme: %s)\n", *theme)
		fmt.Printf("✅ Comparison saved to: %s\n", *output)
		fmt.Printf("📊 Rank diff: %d • Score diff: %s\n",
			comparison.Diff.Rank,
			formatNumber(comparison.Diff.Score))
	} else {
		fmt.Printf("✅ Comparison saved to: %s\n", *output)
	}
}

// runProfileCommand fetches profile data
func runProfileCommand() {
	fs := flag.NewFlagSet("profile", flag.ExitOnError)

	nickname := fs.String("nickname", "", "Root-me username (required)")
	output := fs.String("output", "", "Output file path (optional, prints to stdout if not specified)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s profile [options]\n\n", os.Args[0])
		fmt.Fprintf(fs.Output(), "Fetch Root-me profile data as JSON\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nExamples:\n")
		fmt.Fprintf(fs.Output(), "  %s profile --nickname=franckferman --output=./profile.json\n", os.Args[0])
		fmt.Fprintf(fs.Output(), "  %s profile --nickname=franckferman  # Print to stdout\n", os.Args[0])
	}

	fs.Parse(os.Args[2:])

	if *nickname == "" {
		fmt.Fprintln(os.Stderr, "Error: --nickname is required")
		fs.Usage()
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "🔍 Fetching profile: %s\n", *nickname)

	profile, err := rootme.GetProfile(*nickname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		// Save to file
		if err := ensureDir(filepath.Dir(*output)); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error creating directory: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(*output, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error writing file: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "✅ Profile data saved to: %s\n", *output)
	} else {
		// Print to stdout
		fmt.Print(string(data))
	}
}

// runThemesCommand lists available themes
func runThemesCommand() {
	themes := rootme.GetThemes()

	fmt.Println("Available themes:")
	for _, theme := range themes {
		fmt.Printf("  • %s\n", theme)
	}
}

// showUsage displays help information
func showUsage() {
	fmt.Printf(`root-me-stats v%s

USAGE:
  %s <command> [options]

COMMANDS:
  badge     Generate SVG badge for a Root-me profile
  compare   Generate SVG comparison between two profiles
  profile   Fetch profile data as JSON
  themes    List available themes
  version   Show version information
  help      Show this help message

EXAMPLES:
  # Generate badge for GitHub README
  %s badge --nickname=franckferman --output=./assets/rootme-badge.svg --theme=dark --stats

  # Generate comparison badge
  %s compare --user1=user1 --user2=user2 --output=./assets/compare.svg --theme=midnight

  # Fetch profile data for processing
  %s profile --nickname=franckferman --output=./data/profile.json

  # Output profile to stdout (for GitHub Actions)
  %s profile --nickname=franckferman

GITHUB ACTIONS INTEGRATION:
  Add this to your workflow:

    - name: Generate Root-me Badge
      run: |
        ./rootme-stats badge \
          --nickname=${{ github.repository_owner }} \
          --output=./assets/rootme-badge.svg \
          --theme=dark \
          --stats

    - name: Commit Badge
      run: |
        git add assets/rootme-badge.svg
        git commit -m "Update Root-me badge" || exit 0
        git push

For more information: https://github.com/franckferman/root-me-stats
`, version, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

// ensureDir creates directory if it doesn't exist
func ensureDir(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

// formatNumber adds commas to numbers
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
