package testingx

import (
	"os"
	"path/filepath"
)

// FindProjectRoot walks up from the current directory to find the directory containing .env.
// Returns the absolute path to the project root, or an error if not found.
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached filesystem root
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

func ApplyProjectRootDir() error {
	rootDir, err := FindProjectRoot()
	if err != nil {
		return err
	}
	os.Chdir(rootDir)
	return nil
}
