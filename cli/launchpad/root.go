package launchpad

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
)

//go:generate go run static/includestatic.go

func NewBootstrap() {
	// TODO (@rjerrems) Bootstrap entry point
}

// `$ cft launchpad generate` entry point
func NewGenerate(rawFilepath []string, outFlavor string, outputDir string) {
	gState.outputDirectory = outputDir

	switch outputFlavor(outFlavor) {
	case outTf:
		gState.outputFlavor = outTf
	case outDm:
		gState.outputFlavor = outDm
		log.Println("Deployment Manager format not yet supported")
		return
	default:
		log.Fatalln("Unrecognized output format")
		return
	}

	err := loadAllYAMLs(rawFilepath)
	if err != nil {
		log.Fatalln(err)
	}
	generateOutput()
}

func loadAllYAMLs(rawFilepath []string) error {
	fps, err := validateYAMLFilepath(rawFilepath)
	if err != nil {
		return err
	}
	if fps == nil || len(fps) == 0 {
		return errors.New("no valid YAML files given")
	}
	for _, conf := range fps { // Load all files into runtime
		// TODO multiple yaml documents in one file
		err := yaml.Unmarshal([]byte(loadFile(conf)), &configYAML{})
		if err != nil {
			return errors.New(fmt.Sprintf("%s %s", conf, err.Error()))
		}
	}
	return nil
}
