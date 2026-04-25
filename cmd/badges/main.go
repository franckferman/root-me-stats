// Simple badge generator - standalone tool
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/franckferman/root-me-stats/pkg/rootme"
)

func main() {
	var (
		nickname = flag.String("nickname", "", "Root-me username (required)")
		theme    = flag.String("theme", "dark", "Theme: dark, light, midnight, punk, weedy, astral")
		stats    = flag.Bool("stats", false, "Include global stats")
		output   = flag.String("output", "", "Output file (optional, prints to stdout if not specified)")
	)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Simple Root-me badge generator

Usage: %s [options]

Options:
`, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), `
Examples:
  %s --nickname=franckferman --theme=dark --stats
  %s --nickname=franckferman --output=badge.svg --theme=midnight
`, os.Args[0], os.Args[0])
	}

	flag.Parse()

	if *nickname == "" {
		fmt.Fprintln(os.Stderr, "Error: --nickname is required")
		flag.Usage()
		os.Exit(1)
	}

	// Generate badge
	opts := rootme.DefaultBadgeOptions()
	opts.Theme = *theme
	opts.ShowGlobalStats = *stats

	svg, err := rootme.QuickBadge(*nickname, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		// Save to file
		if err := os.WriteFile(*output, []byte(svg), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Badge saved to: %s\n", *output)
	} else {
		// Print to stdout
		fmt.Print(svg)
	}
}
