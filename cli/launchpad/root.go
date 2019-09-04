package launchpad

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
)

//go:generate go run static/includestatic.go

// CustomResourceDefinition Kind specifies the Kind key-value found in YAML config
type crdKind string

// Output Flavor supported
type outputFlavor string

// Supported CustomResourceDefinition Kind
const (
	KindCloudFoundation crdKind      = "CloudFoundation"
	KindFolder          crdKind      = "Folder"
	KindOrganization    crdKind      = "Organization"
	outDm               outputFlavor = "dm"
	outTf               outputFlavor = "tf"
)

// Global State facilitates evaluation and evaluated objects
var gState globalState

// init initialize tracking for evaluated objects
func init() {
	gState.evaluated.folders.YAMLs = make(map[string]*folderSpecYAML)
}

// NewGenerate takes file patterns as input YAMLs and output Infrastructure as
// Code ready scripts based on specified output flavor.
//
// NewGenerate can be triggered by
//   $ cft launchpad generate
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

// loadAllYAMLs parses input YAMLs and stores evaluated objects in gState
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
