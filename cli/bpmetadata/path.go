package bpmetadata

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	regexEx  = ".*/(examples/.*)"
	regexMod = ".*/(modules/.*)"
)

func isPathValid(path string) (bool, error) {
	_, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("unable to read file at the provided path: %v", err)
	}

	return true, nil
}

func getExamples(configPath string) ([]BlueprintMiscContent, error) {
	re := regexp.MustCompile(regexEx)
	return getDirPaths(configPath, re)
}

func getModules(configPath string) ([]BlueprintMiscContent, error) {
	re := regexp.MustCompile(regexMod)
	return getDirPaths(configPath, re)
}

func getDirPaths(configPath string, re *regexp.Regexp) ([]BlueprintMiscContent, error) {
	dirContent, err := os.ReadDir(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir contents: %v", err)
	}

	var paths []BlueprintMiscContent
	for _, c := range dirContent {
		if c.IsDir() {
			currDirPath := path.Join(configPath, c.Name())
			// ignore .terraform directories
			if strings.HasPrefix(c.Name(), ".terraform") {
				continue
			}

			targetDirContent, _ := os.ReadDir(currDirPath)
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

	return paths, nil
}

func trimPath(assetPath string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(assetPath)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
