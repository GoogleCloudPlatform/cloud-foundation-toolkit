package launchpad

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const extensionYAML = ".yaml"

// validateYAMLFilepath returns .yaml suffix files based on filepath.Glob patterns.
func validateYAMLFilepath(raw []string) ([]string, error) {
	var fps []string
	for _, pattern := range raw {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			if strings.ToLower(filepath.Ext(m)) != extensionYAML {
				continue
			}
			fps = append(fps, m)
		}
	}
	return fps, nil
}

// loadFile return the file content with the specified relative path to current location.
//
// loadFile will attempt to load from filesystem directly first, a not found will attempt
// to load from statics variable generated from `$ go generate`. A user using output binary
// can in theory place their own file in matching relative path and overwrite the binary
// default.
func loadFile(fp string) (string, error) {
	if content, err := ioutil.ReadFile(fp); err == nil {
		return string(content), nil
	} else {
		if !os.IsNotExist(err) {
			fmt.Printf("Request file %s exists but cannot be read\n", fp)
			return "", err
		}
		if content, ok := statics[fp]; ok { // attempt to load from binary statics
			return content, nil
		}
		fmt.Printf("Requested file does not exist in filesystem nor generated binary %s\n", fp)
		return "", os.ErrNotExist
	}
}
