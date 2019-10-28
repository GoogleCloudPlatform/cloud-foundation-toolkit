package launchpad

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
)

//go:generate go run static/includestatic.go

// crdKind is the CustomResourceDefinition (CRD) which is indicated YAML's Kind value.
type crdKind string

// outputFlavor defines launchpad's generated output language.
type outputFlavor string

// Supported crdKind and outputFlavor.
const (
	KindCloudFoundation crdKind      = "CloudFoundation"
	KindFolder          crdKind      = "Folder"
	KindOrganization    crdKind      = "Organization"
	outDm               outputFlavor = "dm"
	outTf               outputFlavor = "tf"
)

// gState is a global scoped state to facilitate evaluation and output generation.
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

	if err := loadAllYAMLs(rawFilepath); err != nil {
		log.Fatalln(err)
	}
	generateOutput()
}

// loadAllYAMLs parses input YAMLs and stores evaluated objects in gState.
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
		if content, err := loadFile(conf); err != nil {
			return err
		} else {
			if err := yaml.Unmarshal([]byte(content), &configYAML{}); err != nil {
				return errors.New(fmt.Sprintf("%s %s", conf, err.Error()))
			}
		}
	}
	return nil
}
