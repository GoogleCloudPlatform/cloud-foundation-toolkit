package launchpad

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"path/filepath"
)

//go:generate go run static/includestatic.go

// gState is a global scoped state to facilitate evaluation and output generation.
var gState globalState

// init initialize tracking for evaluated objects
func init() { gState.init() }

// NewGenerate takes file patterns as input YAMLs and output Infrastructure as
// Code ready scripts based on specified output flavor.
//
// NewGenerate can be triggered by
//   $ cft launchpad generate
func NewGenerate(rawPaths []string, outputFlavorString string, outputDir string) {
	gState.outputDirectory = outputDir
	gState.outputFlavor = newOutputFlavor(outputFlavorString)

	// attempt to load all configs with best effort
	lenYAMLDocs := 0
	for _, pathPattern := range rawPaths {
		matches, err := filepath.Glob(pathPattern)
		if err != nil {
			log.Println("Warning: Invalid file path pattern", pathPattern)
			continue
		}
		for _, fp := range matches {
			content, err := loadFile(fp)
			if err != nil {
				log.Println("Warning: Unable to load requested file", fp)
				continue
			}
			// Multiple YAML doc can exist within one file
			decoder := yaml.NewDecoder(bytes.NewReader([]byte(content)))
			for err := decoder.Decode(&genericYAML{}); err != io.EOF; err = decoder.Decode(&genericYAML{}) {
				if err != nil { // sub document processing error
					log.Fatalln("Unable to process a YAML document within", fp, err.Error())
				}
				lenYAMLDocs += 1
			}
		}
	}
	log.Println(lenYAMLDocs, "YAML documents loaded")
	if err := gState.referenceMap.validate(); err != nil {
		log.Fatalln(err)
	}

	gState.dump()	// Place-holder for future trigger point of code generation
	//generateOutput()
}
