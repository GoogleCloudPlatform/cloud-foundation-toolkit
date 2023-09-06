package bpmetadata

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	regexExamples = ".*/(examples/.*)"
	regexModules  = ".*/(modules/.*)"
)

var (
	reExamples = regexp.MustCompile(regexExamples)
	reModules  = regexp.MustCompile(regexModules)
)

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("unable to read file at the provided path: %w", err)
	}

	if info.IsDir() {
		return false, fmt.Errorf("provided path is a directory, need a valid file path.")
	}

	return true, nil
}

func getExamples(configPath string) ([]*BlueprintMiscContent, error) {
	return getDirPaths(configPath, reExamples)
}

func getModules(configPath string) ([]*BlueprintMiscContent, error) {
	return getDirPaths(configPath, reModules)
}

// getDirPaths traverses a given path and looks for directories
// with TF configs while ignoring the .terraform* directories created and
// used internally by the Terraform CLI
func getDirPaths(configPath string, re *regexp.Regexp) ([]*BlueprintMiscContent, error) {
	paths := []*BlueprintMiscContent{}
	err := filepath.Walk(configPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing examples in the path %q: %w", configPath, err)
		}

		// skip if this is a .terraform dir
		if info.IsDir() && strings.HasPrefix(info.Name(), ".terraform") {
			return filepath.SkipDir
		}

		// only interested if it has a TF config
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tf") {
			d := filepath.Dir(path)
			if l := trimPath(d, re); l != "" {
				dirPath := &BlueprintMiscContent{
					Name:     filepath.Base(d),
					Location: l,
				}

				paths = append(paths, dirPath)
			}
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error accessing examples in the path %q: %w", configPath, err)
	}

	// Sort by configPath name before returning
	sort.SliceStable(paths, func(i, j int) bool { return paths[i].Name < paths[j].Name })
	return paths, nil
}

func trimPath(assetPath string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(assetPath)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
