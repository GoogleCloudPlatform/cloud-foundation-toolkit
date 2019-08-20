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

// Validates raw strings, including Glob patterns, and returns a validated list of
// yaml filepath ready for consumption
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

// Given a filepath relative to `pwd`, attempt to load from filesystem. If failed, load from statics variable
// populated from `$ go generate`. We can embed static file into the golang binary
// and use this method to load efficiently.
//
// Filesystem as priority enables efficient development time, and also easy to overwrite if user do so choose.
//
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
