package launchpad

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const ExtensionYAML string = ".yaml"

// validateYAMLFilepath returns .yaml suffix files based on filepath.Glob patterns.
func validateYAMLFilepath(raw []string) ([]string, error) {
	var fps []string
	for _, pattern := range raw {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			if strings.ToLower(filepath.Ext(m)) != ExtensionYAML {
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
func loadFile(fp string) string {
	_, err := os.Stat(fp)
	if err == nil { // file exist
		content, err := ioutil.ReadFile(fp)
		if err != nil {
			panic(err)
		}
		return string(content)
	} else if os.IsNotExist(err) { // file does not exist
		if content, ok := statics[fp]; ok {
			return content
		} else {
			fmt.Printf("Requested file does not exist in filesystem nor generated binary %s\n", fp)
			panic(errors.New("file not found"))
		}
	} else {
		panic(err)
	}
}
