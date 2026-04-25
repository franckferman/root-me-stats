// Root-me Stats API Server
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/franckferman/root-me-stats/pkg/rootme"
)

const (
	defaultPort = "3000"
	defaultHost = "0.0.0.0"
	version     = "2.0.0"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = defaultHost
	}

	mux := http.NewServeMux()

	// Badge endpoints (compatible with original)
	mux.HandleFunc("/rm-gh", handleBadge)
	mux.HandleFunc("/badge", handleBadge)

	// Comparison endpoint
	mux.HandleFunc("/compare", handleComparison)

	// JSON API endpoints
	mux.HandleFunc("/api/profile", handleProfileAPI)
	mux.HandleFunc("/api/compare", handleCompareAPI)

	// Utility endpoints
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/themes", handleThemes)
	mux.HandleFunc("/", handleRoot)

	// Add CORS middleware
	handler := corsMiddleware(mux)

	addr := host + ":" + port

	log.Printf("🚀 Root-me Stats API v%s starting on http://%s", version, addr)
	log.Printf("📋 Available endpoints:")
	log.Printf("  GET /rm-gh?nickname=USER&style=THEME&gstats=show")
	log.Printf("  GET /badge?nickname=USER&style=THEME&gstats=show")
	log.Printf("  GET /compare?user1=USER1&user2=USER2&style=THEME")
	log.Printf("  GET /api/profile?nickname=USER")
	log.Printf("  GET /api/compare?user1=USER1&user2=USER2")
	log.Printf("  GET /themes")
	log.Printf("  GET /health")

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// handleBadge generates SVG badges
func handleBadge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	nickname := query.Get("nickname")

	if nickname == "" {
		http.Error(w, "Missing required parameter: nickname", http.StatusBadRequest)
		return
	}

	theme := query.Get("style")
	if theme == "" {
		theme = "dark"
	}

	showStats := query.Get("gstats") == "show"

	// Validate theme
	validThemes := rootme.GetThemes()
	themeValid := false
	for _, validTheme := range validThemes {
		if theme == validTheme {
			themeValid = true
			break
		}
	}

	if !themeValid {
		http.Error(w, fmt.Sprintf("Invalid theme: %s. Available: %v", theme, validThemes), http.StatusBadRequest)
		return
	}

	// Generate badge
	opts := rootme.DefaultBadgeOptions()
	opts.Theme = theme
	opts.ShowGlobalStats = showStats

	svg, err := rootme.QuickBadge(nickname, opts)
	if err != nil {
		log.Printf("Badge generation error: %v", err)
		http.Error(w, fmt.Sprintf("Profile not found: %s", nickname), http.StatusNotFound)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("ETag", fmt.Sprintf("\"%s-%s-%d\"", nickname, theme, len(svg)))

	w.Write([]byte(svg))
}

// handleComparison generates comparison badges
func handleComparison(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	user1 := query.Get("user1")
	user2 := query.Get("user2")

	if user1 == "" || user2 == "" {
		http.Error(w, "Missing required parameters: user1, user2", http.StatusBadRequest)
		return
	}

	theme := query.Get("style")
	if theme == "" {
		theme = "dark"
	}

	widthStr := query.Get("width")
	width := 500
	if widthStr != "" {
		if w, err := strconv.Atoi(widthStr); err == nil && w > 0 {
			width = w
		}
	}

	svg, err := rootme.QuickComparison(user1, user2, theme, width)
	if err != nil {
		log.Printf("Comparison error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to compare profiles: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	w.Write([]byte(svg))
}

// handleProfileAPI returns profile data as JSON
func handleProfileAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	nickname := query.Get("nickname")

	if nickname == "" {
		writeJSONError(w, "Missing required parameter: nickname", http.StatusBadRequest)
		return
	}

	profile, err := rootme.GetProfile(nickname)
	if err != nil {
		log.Printf("Profile API error: %v", err)
		writeJSONError(w, fmt.Sprintf("Profile not found: %s", nickname), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	json.NewEncoder(w).Encode(profile)
}

// handleCompareAPI returns comparison data as JSON
func handleCompareAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	user1 := query.Get("user1")
	user2 := query.Get("user2")

	if user1 == "" || user2 == "" {
		writeJSONError(w, "Missing required parameters: user1, user2", http.StatusBadRequest)
		return
	}

	comparison, err := rootme.CompareProfiles(user1, user2)
	if err != nil {
		log.Printf("Compare API error: %v", err)
		writeJSONError(w, fmt.Sprintf("Comparison failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	json.NewEncoder(w).Encode(comparison)
}

// handleHealth returns service health status
func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "ok",
		"service": "root-me-stats",
		"version": version,
		"endpoints": []string{
			"/rm-gh", "/badge", "/compare",
			"/api/profile", "/api/compare",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleThemes returns available themes
func handleThemes(w http.ResponseWriter, r *http.Request) {
	themes := rootme.GetThemes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(themes)
}

// handleRoot handles root path
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Redirect to health endpoint
	http.Redirect(w, r, "/health", http.StatusMovedPermanently)
}

// writeJSONError writes an error response in JSON format
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]string{
		"error":   http.StatusText(statusCode),
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
