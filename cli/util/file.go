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

// walkTerraformDirs traverses a provided path to return a list of directories
// that hold terraform configs while skiping internal folders that have a
// .terraform.* prefix
func WalkTerraformDirs(topLevelPath string) ([]string, error) {
	var tfDirs []string
	err := filepath.Walk(topLevelPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("failure in accessing the path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && strings.HasPrefix(info.Name(), tfInternalDirPrefix) {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tf") {
			tfDirs = append(tfDirs, filepath.Dir(path))
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", topLevelPath, err)
		return nil, err
	}

	return tfDirs, nil
}
