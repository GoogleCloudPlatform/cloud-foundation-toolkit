package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	tfInternalDirPrefix = ".terraform"
)

// skipDiscoverDirs are directories that are skipped when discovering test cases.
var skipDiscoverDirs = map[string]bool{
	"test":  true,
	"build": true,
	".git":  true,
}

// walkTerraformDirs traverses a provided path to return a list of directories
// that hold terraform configs while skiping internal folders that have a
// .terraform.* prefix
func WalkTerraformDirs(topLevelPath string) ([]string, error) {
	var tfDirs []string
	err := filepath.Walk(topLevelPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failure in accessing the path %q: %w\n", path, err)
		}
		if info.IsDir() && (strings.HasPrefix(info.Name(), tfInternalDirPrefix) || skipDiscoverDirs[info.Name()]) {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tf") {
			tfDirs = append(tfDirs, filepath.Dir(path))
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the path %q: %w\n", topLevelPath, err)
	}

	return tfDirs, nil
}

func FindFilesWithPattern(dir string, pattern string, skipPaths []string) ([]string, error) {
	f, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("no such dir: %w", err)
	}
	if !f.IsDir() {
		return nil, fmt.Errorf("expected dir %s: got file", dir)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %w", err)
	}

	filePaths := []string{}

	err = filepath.WalkDir(dir, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !re.MatchString(filepath.Base(path)) {
			return nil
		}

		for _, p := range skipPaths {
			if strings.Contains(path, p) {
				return nil
			}
		}

		if !file.IsDir() {
			filePaths = append(filePaths, path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error accessing the path: %q. error: %v\n", dir, err)
		return nil, err
	}

	return filePaths, nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("error checking if %s exists: %w", path, err)
}
