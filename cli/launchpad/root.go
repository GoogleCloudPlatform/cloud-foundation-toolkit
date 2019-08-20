package launchpad

import (
	"gopkg.in/yaml.v2"
	"log"
)

//go:generate go run static/includestatic.go

func NewBootstrap() {
	// TODO (@rjerrems) Bootstrap entry point
}

// `$ cft launchpad generate` entry point
func NewGenerate(rawFilepath []string, outputFlavor string, outputDir string) {
	gState.outputDirectory = outputDir
	switch outputFlavor {
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

	fps, err := validateYAMLFilepath(rawFilepath)
	if err != nil {
		log.Fatalln(err)
	}
	if fps == nil || len(fps) == 0 {
		log.Fatalln("No valid YAML files given")
	}
	for _, conf := range fps { // Load all files into runtime
		// TODO multiple yaml documents in one file
		err := yaml.Unmarshal([]byte(loadFile(conf)), &configYAML{})
		if err != nil {
			log.Fatalln(err)
		}
	}
	generateOutput()
}
