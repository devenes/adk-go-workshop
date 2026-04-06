// Package main scans files for potential API key leaks.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	patterns = map[string]*regexp.Regexp{
		"Google API Key":    regexp.MustCompile(`AIza[0-9A-Za-z\-_]{35}`),
		"OpenAI Secret":     regexp.MustCompile(`sk-[0-9A-Za-z]{20,}`),
		"GitHub PAT":        regexp.MustCompile(`ghp_[0-9A-Za-z]{36,}`),
		"AWS Access Key":    regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		"Generic Secret":    regexp.MustCompile(`(password|secret|token|api_key)\s*[:=]\s*["']?[^\s"']+["']?`),
	}

	skipDirs = map[string]bool{
		".git":    true,
		".venv":   true,
		"node_modules": true,
		"vendor":   true,
		".env":    true,
	}

	allowlistedPatterns = map[string]bool{
		"__GEMINI_FLASH_LITE_MODEL_ID__": true,
		"placeholder":                     true,
		"example":                         true,
	}
)

func shouldSkip(path string) bool {
	for skipDir := range skipDirs {
		if strings.Contains(path, skipDir) {
			return true
		}
	}

	// Allow .env.example
	if filepath.Base(path) == ".env.example" {
		return true
	}

	return false
}

func scanFile(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	var findings []string

	for name, pattern := range patterns {
		matches := pattern.FindAllString(content, -1)
		for _, match := range matches {
			// Check if it's in allowlist
			isAllowed := false
			for allowed := range allowlistedPatterns {
				if match == allowed {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				findings = append(findings, fmt.Sprintf("%s: %s", name, match))
			}
		}
	}

	return findings, nil
}

func main() {
	rootDir := "."
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}

	var hasLeaks bool

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if shouldSkip(path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		// Check text files
		if !isTextFile(path) {
			return nil
		}

		findings, err := scanFile(path)
		if err != nil {
			return err
		}

		if len(findings) > 0 {
			hasLeaks = true
			fmt.Printf("⚠️  %s\n", path)
			for _, finding := range findings {
				fmt.Printf("   - %s\n", finding)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	if hasLeaks {
		fmt.Println("\n❌ Potential API key leaks detected!")
		os.Exit(1)
	}

	fmt.Println("✅ No API key leaks detected")
	os.Exit(0)
}

func isTextFile(path string) bool {
	ext := filepath.Ext(path)
	textExtensions := map[string]bool{
		".go":    true,
		".py":    true,
		".js":    true,
		".ts":    true,
		".txt":   true,
		".md":    true,
		".yaml":  true,
		".yml":   true,
		".json":  true,
		".env":   true,
		".sh":    true,
		".bash":  true,
		".zsh":   true,
		".mod":   true,
		".sum":   true,
	}
	return textExtensions[ext]
}
