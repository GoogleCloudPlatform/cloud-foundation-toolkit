package util

import (
	"fmt"
	"io/fs"
	"path/filepath"
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
			return fmt.Errorf("failure in accessing the path %q: %v\n", path, err)
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
		return nil, fmt.Errorf("error walking the path %q: %v\n", topLevelPath, err)
	}

	return tfDirs, nil
}
