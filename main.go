package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	dir := flag.String("dir", ".", "Directory to sanitize")
	flag.Parse()

	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		dirPath := filepath.Dir(path)
		oldName := info.Name()
		newName := sanitizeName(oldName)

		if oldName != newName {
			oldPath := filepath.Join(dirPath, oldName)
			newPath := filepath.Join(dirPath, newName)
			fmt.Printf("Renaming: %s → %s\n", oldName, newName)
			if err := os.Rename(oldPath, newPath); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to rename %s: %v\n", oldName, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the path: %v\n", err)
		os.Exit(1)
	}
}

func sanitizeName(name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	// Replace all problematic characters with "-"
	base = strings.ReplaceAll(base, " ", "-")
	re := regexp.MustCompile(`[^\w\.-]+`)
	base = re.ReplaceAllString(base, "-")

	// Collapse multiple dashes
	reDash := regexp.MustCompile(`-+`)
	base = reDash.ReplaceAllString(base, "-")

	// Trim trailing/leading dashes from base name ONLY
	base = strings.Trim(base, "-")

	return base + ext
}
