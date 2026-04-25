// Package cache provides simple file-based caching for Root-me data
package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/franckferman/root-me-stats/internal/fetcher"
)

const (
	// CacheTTL is how long cached data is valid (24 hours)
	CacheTTL = 24 * time.Hour

	// CacheDir is the directory for cache files
	CacheDir = ".cache"
)

// CacheEntry represents a cached profile with timestamp
type CacheEntry struct {
	Profile   *fetcher.Profile `json:"profile"`
	Timestamp time.Time        `json:"timestamp"`
}

// Get retrieves a profile from cache if valid, nil if not found or expired
func Get(nickname string) (*fetcher.Profile, error) {
	cachePath := getCachePath(nickname)

	// Check if cache file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil, nil // Not found
	}

	// Read cache file
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, fmt.Errorf("reading cache file: %w", err)
	}

	// Parse cache entry
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		// Invalid cache file, delete it
		os.Remove(cachePath)
		return nil, nil
	}

	// Check if expired
	if time.Since(entry.Timestamp) > CacheTTL {
		// Expired, delete cache file
		os.Remove(cachePath)
		return nil, nil
	}

	return entry.Profile, nil
}

// Set stores a profile in cache
func Set(profile *fetcher.Profile) error {
	if profile == nil {
		return fmt.Errorf("profile cannot be nil")
	}

	// Ensure cache directory exists
	if err := ensureCacheDir(); err != nil {
		return fmt.Errorf("creating cache directory: %w", err)
	}

	cachePath := getCachePath(profile.Nickname)

	// Create cache entry
	entry := CacheEntry{
		Profile:   profile,
		Timestamp: time.Now(),
	}

	// Marshal to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshaling cache entry: %w", err)
	}

	// Write to file
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("writing cache file: %w", err)
	}

	return nil
}

// Clear removes all cache files
func Clear() error {
	return os.RemoveAll(CacheDir)
}

// GetOrFetch retrieves profile from cache or fetches if not available
func GetOrFetch(nickname string) (*fetcher.Profile, error) {
	// Try cache first
	if profile, err := Get(nickname); err != nil {
		// Cache error, continue to fetch
	} else if profile != nil {
		return profile, nil // Cache hit
	}

	// Cache miss or error, fetch from Root-me
	profile, err := fetcher.FetchProfile(nickname)
	if err != nil {
		return nil, err
	}

	// Store in cache (ignore errors)
	Set(profile)

	return profile, nil
}

// getCachePath returns the file path for a nickname's cache
func getCachePath(nickname string) string {
	// Use MD5 hash of nickname for filename (handles special characters)
	hash := md5.Sum([]byte(nickname))
	filename := fmt.Sprintf("%x.json", hash)
	return filepath.Join(CacheDir, filename)
}

// ensureCacheDir creates the cache directory if it doesn't exist
func ensureCacheDir() error {
	return os.MkdirAll(CacheDir, 0755)
}

// CleanExpired removes expired cache files
func CleanExpired() error {
	if err := ensureCacheDir(); err != nil {
		return err
	}

	entries, err := os.ReadDir(CacheDir)
	if err != nil {
		return fmt.Errorf("reading cache directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(CacheDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Remove files older than TTL
		if time.Since(info.ModTime()) > CacheTTL {
			os.Remove(filePath)
		}
	}

	return nil
}
