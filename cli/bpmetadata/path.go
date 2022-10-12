package bpmetadata

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

const (
	regexExamples = ".*/(examples/.*)"
	regexModules  = ".*/(modules/.*)"
)

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("unable to read file at the provided path: %v", err)
	}

	if info.IsDir() {
		return false, fmt.Errorf("provided path is a directory, need a valid file path.")
	}

	return true, nil
}

func getExamples(configPath string) ([]BlueprintMiscContent, error) {
	re := regexp.MustCompile(regexExamples)
	return getDirPaths(configPath, re)
}

func getModules(configPath string) ([]BlueprintMiscContent, error) {
	re := regexp.MustCompile(regexModules)
	return getDirPaths(configPath, re)
}

// getDirPaths traverses a given path and looks for directories
// with TF configs while ignoring the .terraform* directories created and
// used internally by the Terraform CLI
func getDirPaths(configPath string, re *regexp.Regexp) ([]BlueprintMiscContent, error) {
	dirContent, err := os.ReadDir(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir contents: %v", err)
	}

	var paths []BlueprintMiscContent
	for _, c := range dirContent {
		if !c.IsDir() {
			continue
		}

		currDirPath := path.Join(configPath, c.Name())
		// ignore .terraform directories
		if strings.HasPrefix(c.Name(), ".terraform") {
			continue
		}

		targetDirContent, err := os.ReadDir(currDirPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read dir contents: %v", err)
		}

		for _, t := range targetDirContent {
			// only evaluate dirs with tf configs
			if !strings.HasSuffix(t.Name(), ".tf") {
				continue
			}

			if l := trimPath(currDirPath, re); l != "" {
				dirPath := BlueprintMiscContent{
					Name:     c.Name(),
					Location: l,
				}

				paths = append(paths, dirPath)
				break
			}
		}

		rPaths, _ := getDirPaths(currDirPath, re)
		if len(rPaths) > 0 {
			paths = append(paths, rPaths...)
			}
		}
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
