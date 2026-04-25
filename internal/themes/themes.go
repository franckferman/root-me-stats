// Package themes defines visual themes and icons for badges
package themes

// Theme defines colors for a visual theme
type Theme struct {
	Background string
	Bar        string
	Accent     string
	Text       string
	Title      string
}

// Available themes
var themes = map[string]*Theme{
	"light": {
		Background: "#ffffff",
		Bar:        "#f1f1f1",
		Accent:     "#91B302",
		Text:       "#333333",
		Title:      "#000000",
	},
	"dark": {
		Background: "#000000",
		Bar:        "#2A3609",
		Accent:     "#91B302",
		Text:       "#ffffff",
		Title:      "#ffffff",
	},
	"midnight": {
		Background: "#1e1b4b",
		Bar:        "#2e1065",
		Accent:     "#6d28d9",
		Text:       "#e2e8f0",
		Title:      "#f8fafc",
	},
	"punk": {
		Background: "#500724",
		Bar:        "#881337",
		Accent:     "#db2777",
		Text:       "#fce7f3",
		Title:      "#fdf2f8",
	},
	"weedy": {
		Background: "#022c22",
		Bar:        "#365314",
		Accent:     "#65a30d",
		Text:       "#dcfce7",
		Title:      "#f0fdf4",
	},
	"astral": {
		Background: "#000000",
		Bar:        "#340000",
		Accent:     "#DF0000",
		Text:       "#fecaca",
		Title:      "#fef2f2",
	},
}

// Category icons (SVG paths)
var categoryIcons = map[string]string{
	"App-Script":    "M12 2L2 7V17L12 22L22 17V7L12 2Z",
	"App-System":    "M4 4H20V6H4V4M4 8H20V10H4V8M4 12H20V14H4V12M4 16H20V18H4V16M4 20H20V22H4V20Z",
	"Cracking":      "M12 1L21.5 6.5V17.5L12 23L2.5 17.5V6.5L12 1M12 7A2 2 0 0 0 10 9A2 2 0 0 0 12 11A2 2 0 0 0 14 9A2 2 0 0 0 12 7M7 18A2 2 0 0 1 9 16H15A2 2 0 0 1 17 18A2 2 0 0 1 15 20H9A2 2 0 0 1 7 18Z",
	"Cryptanalysis": "M12.65 10C11.7 7.31 8.9 5.5 5.77 6.12C4.68 6.34 4 7.28 4 8.39V9A4 4 0 0 0 8 13H10V19A2 2 0 0 0 12 21A2 2 0 0 0 14 19V13H16A4 4 0 0 0 20 9V8.39C20 7.28 19.32 6.34 18.23 6.12C16.06 5.71 14.15 6.9 13.09 8.77L12.65 10M8 11A1 1 0 0 1 7 10A1 1 0 0 1 8 9A1 1 0 0 1 9 10A1 1 0 0 1 8 11M16 11A1 1 0 0 1 15 10A1 1 0 0 1 16 9A1 1 0 0 1 17 10A1 1 0 0 1 16 11Z",
	"Forensic":      "M15.5 14L20.5 19L19 20.5L14 15.5V14.71L13.73 14.44C12.59 15.41 11.11 16 9.5 16A6.5 6.5 0 0 1 3 9.5A6.5 6.5 0 0 1 9.5 3A6.5 6.5 0 0 1 16 9.5C16 11.11 15.41 12.59 14.44 13.73L14.71 14H15.5M9.5 14C12 14 14 12 14 9.5S12 5 9.5 5S5 7 5 9.5S7 14 9.5 14Z",
	"Network":       "M1 9V15H4L5 19H7L6 15H9V9H1M11 9V15H14L15 19H17L16 15H19V9H11M21 9V15H24V9H21Z",
	"Programming":   "M8 3A2 2 0 0 0 6 5V9A2 2 0 0 1 4 11H3V13H4A2 2 0 0 1 6 15V19A2 2 0 0 0 8 21H10V19H8V14A2 2 0 0 0 6 12A2 2 0 0 0 8 10V5H10V3H8M16 3A2 2 0 0 1 18 5V9A2 2 0 0 0 20 11H21V13H20A2 2 0 0 0 18 15V19A2 2 0 0 1 16 21H14V19H16V14A2 2 0 0 1 18 12A2 2 0 0 1 16 10V5H14V3H16Z",
	"Realist":       "M12 2A10 10 0 0 0 2 12A10 10 0 0 0 12 22A10 10 0 0 0 22 12A10 10 0 0 0 12 2M12 4A8 8 0 0 1 20 12A8 8 0 0 1 12 20A8 8 0 0 1 4 12A8 8 0 0 1 12 4Z",
	"Steganography": "M9 11H7V9H9V11M13 11H11V9H13V11M17 11H15V9H17V11M19 9H17V7H19V9M17 13H15V11H17V13M19 11H21V13H19V11M13 7H11V5H13V7M13 15H11V13H13V15M5 11H3V9H5V11M7 7H5V5H7V7M9 15H7V13H9V15M5 13H3V15H5V13M9 19H7V17H9V19M5 17H3V19H5V17M13 19H11V17H13V19M17 17H15V19H17V17Z",
	"Web-Client":    "M16.36 14C16.44 13.3 16.5 12.66 16.5 12S16.44 10.7 16.36 10H19.74C19.9 10.64 20 11.31 20 12S19.9 13.36 19.74 14H16.36M14.59 19.56C15.19 18.45 15.65 17.25 15.97 16H18.92C17.96 17.65 16.43 18.93 14.59 19.56M14.34 14H9.66C9.56 13.34 9.5 12.68 9.5 12S9.56 10.66 9.66 10H14.34C14.44 10.66 14.5 11.32 14.5 12S14.44 13.34 14.34 14M12 19.96C11.17 18.76 10.5 17.43 10.09 16H13.91C13.5 17.43 12.83 18.76 12 19.96M8 8H5.08C6.03 6.34 7.57 5.06 9.4 4.44C8.8 5.55 8.35 6.75 8 8M5.08 16H8C8.35 17.25 8.8 18.45 9.4 19.56C7.57 18.93 6.03 17.65 5.08 16M4.26 14C4.1 13.36 4 12.69 4 12S4.1 10.64 4.26 10H7.64C7.56 10.7 7.5 11.34 7.5 12S7.56 13.3 7.64 14H4.26M12 4.03C12.83 5.23 13.5 6.57 13.91 8H10.09C10.5 6.57 11.17 5.23 12 4.03M18.92 8H15.97C15.65 6.75 15.19 5.55 14.59 4.44C16.43 5.07 17.96 6.34 18.92 8M12 2C6.48 2 2 6.48 2 12S6.48 22 12 22S22 17.52 22 12S17.52 2 12 2Z",
	"Web-Server":    "M4 1C2.89 1 2 1.89 2 3V7C2 8.11 2.89 9 4 9H20C21.11 9 22 8.11 22 7V3C22 1.89 21.11 1 20 1H4M4 11C2.89 11 2 11.89 2 13V17C2 18.11 2.89 19 4 19H20C21.11 19 22 18.11 22 17V13C22 11.89 21.11 11 20 11H4M4 21C2.89 21 2 21.89 2 23V27C2 28.11 2.89 29 4 29H20C21.11 29 22 28.11 22 27V23C22 21.89 21.11 21 20 21H4Z",
}

// GetTheme returns a theme by name, nil if not found
func GetTheme(name string) *Theme {
	return themes[name]
}

// GetThemeNames returns all available theme names
func GetThemeNames() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	return names
}

// GetCategoryIcon returns SVG path for a category icon
func GetCategoryIcon(category string) string {
	if icon, exists := categoryIcons[category]; exists {
		return icon
	}
	return categoryIcons["Programming"] // Default fallback
}
