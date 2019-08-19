package launchpad

import (
	"log"
	"path/filepath"
	"strings"
)

const ExtensionYAML string = ".yaml"

func NewBootstrap() {
	// TODO (@rjerrems) Bootstrap entry point
}

// `$ cft launchpad generate` entry point
func NewGenerate(rawFilepath []string) {
	fps, err := validateYAMLFilepath(rawFilepath)
	if err != nil {
		log.Fatalln(err)
	}
	if fps == nil || len(fps) == 0 {
		log.Fatalln("No valid YAML files given")
	}
	for _, conf := range fps {
		// TODO (@wengm) Model loading .. etc
		println(conf)
	}
}

// Validates raw strings, including Glob patterns, and returns a validated list of
// yaml filepath ready for consumption
func validateYAMLFilepath(raw []string) ([]string, error) {
	var fps []string
	for _, pattern := range raw {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, m := range matches { // Glob will return exist files
			if strings.ToLower(filepath.Ext(m)) != ExtensionYAML {
				continue
			}
			fps = append(fps, m)
		}
	}
	return fps, nil
}
