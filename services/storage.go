package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func Storage(cCtx *cli.Context) error {
	tempDir := os.TempDir()

	fmt.Printf("Scanning temp directory: %s\n", tempDir)

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Size() == 0 {
			fmt.Printf("Deleting empty file: %s\n", path)
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to delete empty file %s: %v", path, err)
			}
		}

		if isCacheFile(info.Name()) {
			fmt.Printf("Deleting cache file: %s\n", path)
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to delete cache file %s: %v", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning temp directory: %v\n", err)
	} else {
		fmt.Println("Temporary files cleanup complete.")
	}

	return nil
}

func isCacheFile(filename string) bool {
	cacheExtensions := []string{
		".cache",
		".tmp",
		".swp",
		".bak",
	}

	for _, ext := range cacheExtensions {
		if filepath.Ext(filename) == ext {
			return true
		}
	}

	return containsCacheKeyword(filename)
}

func containsCacheKeyword(filename string) bool {
	keywords := []string{"cache", "temp", "swap", "backup"}
	for _, keyword := range keywords {
		if filepath.Base(filename) == keyword {
			return true
		}
	}
	return false
}
