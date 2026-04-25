// Package fetcher handles Root-me data extraction
package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Profile represents a Root-me user profile
type Profile struct {
	Nickname   string      `json:"nickname"`
	Stats      Stats       `json:"stats"`
	Categories []Category  `json:"categories"`
	Challenges []Challenge `json:"challenges"`
	FetchedAt  time.Time   `json:"fetched_at"`
}

// Stats contains basic user statistics
type Stats struct {
	Rank           int `json:"rank"`
	Score          int `json:"score"`
	Challenges     int `json:"challenges"`
	Compromissions int `json:"compromissions"`
}

// Category represents a challenge category with progress
type Category struct {
	Name       string `json:"name"`
	Percentage int    `json:"percentage"`
}

// Challenge represents an individual challenge
type Challenge struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
	Completed bool   `json:"completed"`
	Category  string `json:"category"`
}

// Comparison holds data for comparing two profiles
type Comparison struct {
	User1                Profile                    `json:"user1"`
	User2                Profile                    `json:"user2"`
	Diff                 StatsDiff                  `json:"diff"`
	ChallengesByCategory map[string][]ChallengeDiff `json:"challenges_by_category"`
}

// StatsDiff represents the difference between two profiles
type StatsDiff struct {
	Rank       int                     `json:"rank"`
	Score      int                     `json:"score"`
	Challenges int                     `json:"challenges"`
	Categories map[string]CategoryDiff `json:"categories"`
}

// CategoryDiff represents the difference in a category
type CategoryDiff struct {
	User1 int `json:"user1"`
	User2 int `json:"user2"`
	Diff  int `json:"diff"`
}

// ChallengeDiff represents a challenge comparison
type ChallengeDiff struct {
	Name           string `json:"name"`
	Link           string `json:"link"`
	User1Completed bool   `json:"user1_completed"`
	User2Completed bool   `json:"user2_completed"`
}

// FetchProfile retrieves and parses a Root-me user profile
func FetchProfile(nickname string) (*Profile, error) {
	if nickname == "" {
		return nil, fmt.Errorf("nickname cannot be empty")
	}

	// Fetch the profile page
	url := fmt.Sprintf("https://www.root-me.org/%s?inc=score&lang=en", nickname)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; root-me-stats/2.0)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("profile not found: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	html := string(body)

	// Parse the HTML using regex (stdlib only - no external HTML parser)
	profile := &Profile{
		Nickname:  nickname,
		FetchedAt: time.Now(),
	}

	// Extract basic stats from h3 elements
	stats, err := parseStats(html)
	if err != nil {
		return nil, fmt.Errorf("parsing stats: %w", err)
	}
	profile.Stats = stats

	// Extract category progress from h4 elements
	categories := parseCategories(html)
	profile.Categories = categories

	// Extract challenge list
	challenges := parseChallenges(html)
	profile.Challenges = challenges

	return profile, nil
}

// CompareProfiles compares two Root-me profiles
func CompareProfiles(user1, user2 string) (*Comparison, error) {
	profile1, err := FetchProfile(user1)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", user1, err)
	}

	profile2, err := FetchProfile(user2)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", user2, err)
	}

	comparison := &Comparison{
		User1: *profile1,
		User2: *profile2,
		Diff: StatsDiff{
			Rank:       profile2.Stats.Rank - profile1.Stats.Rank,
			Score:      profile2.Stats.Score - profile1.Stats.Score,
			Challenges: profile2.Stats.Challenges - profile1.Stats.Challenges,
			Categories: make(map[string]CategoryDiff),
		},
		ChallengesByCategory: make(map[string][]ChallengeDiff),
	}

	// Compare categories
	for _, cat1 := range profile1.Categories {
		for _, cat2 := range profile2.Categories {
			if cat1.Name == cat2.Name {
				comparison.Diff.Categories[cat1.Name] = CategoryDiff{
					User1: cat1.Percentage,
					User2: cat2.Percentage,
					Diff:  cat2.Percentage - cat1.Percentage,
				}
				break
			}
		}
	}

	// Compare challenges by category
	challengeMap := make(map[string]map[string]*ChallengeDiff)

	for _, challenge := range profile1.Challenges {
		if challengeMap[challenge.Category] == nil {
			challengeMap[challenge.Category] = make(map[string]*ChallengeDiff)
		}
		challengeMap[challenge.Category][challenge.Name] = &ChallengeDiff{
			Name:           challenge.Name,
			Link:           challenge.Link,
			User1Completed: challenge.Completed,
			User2Completed: false,
		}
	}

	for _, challenge := range profile2.Challenges {
		if challengeMap[challenge.Category] == nil {
			challengeMap[challenge.Category] = make(map[string]*ChallengeDiff)
		}
		if existing := challengeMap[challenge.Category][challenge.Name]; existing != nil {
			existing.User2Completed = challenge.Completed
		} else {
			challengeMap[challenge.Category][challenge.Name] = &ChallengeDiff{
				Name:           challenge.Name,
				Link:           challenge.Link,
				User1Completed: false,
				User2Completed: challenge.Completed,
			}
		}
	}

	// Convert map to slice
	for category, challenges := range challengeMap {
		for _, challenge := range challenges {
			comparison.ChallengesByCategory[category] = append(
				comparison.ChallengesByCategory[category],
				*challenge,
			)
		}
	}

	return comparison, nil
}

// parseStats extracts basic statistics from HTML using regex
func parseStats(html string) (Stats, error) {
	var stats Stats

	// Match h3 elements containing stats
	h3Regex := regexp.MustCompile(`<h3[^>]*>([^<]+)</h3>`)
	matches := h3Regex.FindAllStringSubmatch(html, -1)

	if len(matches) >= 4 {
		stats.Rank = parseNumber(matches[0][1])
		stats.Score = parseNumber(matches[1][1])
		stats.Challenges = parseNumber(matches[2][1])
		stats.Compromissions = parseNumber(matches[3][1])
	}

	return stats, nil
}

// parseCategories extracts category progress from HTML
func parseCategories(html string) []Category {
	var categories []Category

	// Match h4 elements with percentages
	h4Regex := regexp.MustCompile(`<h4[^>]*>([^:]+)\s*:\s*(\d+)%</h4>`)
	matches := h4Regex.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			name := strings.TrimSpace(match[1])
			percentage := parseNumber(match[2])

			if name != "" {
				categories = append(categories, Category{
					Name:       name,
					Percentage: percentage,
				})
			}
		}
	}

	return categories
}

// parseChallenges extracts challenge list from HTML
func parseChallenges(html string) []Challenge {
	var challenges []Challenge

	// Match challenge links with completion status
	linkRegex := regexp.MustCompile(`<a[^>]*class="[^"]*(?:vert|rouge)[^"]*"[^>]*href="([^"]+)"[^>]*>([^<]+)</a>`)
	matches := linkRegex.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			link := match[1]
			name := strings.TrimSpace(match[2])

			// Determine completion status (vert = completed)
			completed := strings.Contains(match[0], "vert")

			// Extract category from URL
			category := extractCategory(link)

			if name != "" {
				// Make URL absolute if relative
				if strings.HasPrefix(link, "/") {
					link = "https://www.root-me.org" + link
				}

				challenges = append(challenges, Challenge{
					Name:      name,
					Link:      link,
					Completed: completed,
					Category:  category,
				})
			}
		}
	}

	return challenges
}

// parseNumber extracts integer from text
func parseNumber(text string) int {
	// Remove all non-digit characters
	numRegex := regexp.MustCompile(`\d+`)
	match := numRegex.FindString(strings.ReplaceAll(text, ",", ""))
	if match != "" {
		if num, err := strconv.Atoi(match); err == nil {
			return num
		}
	}
	return 0
}

// extractCategory determines category from challenge URL
func extractCategory(url string) string {
	categoryMap := map[string]string{
		"app-script":    "App-Script",
		"app-system":    "App-System",
		"cracking":      "Cracking",
		"cryptanalysis": "Cryptanalysis",
		"forensic":      "Forensic",
		"network":       "Network",
		"programming":   "Programming",
		"realist":       "Realist",
		"steganography": "Steganography",
		"web-client":    "Web-Client",
		"web-server":    "Web-Server",
	}

	for key, value := range categoryMap {
		if strings.Contains(url, key) {
			return value
		}
	}

	return "Unknown"
}
