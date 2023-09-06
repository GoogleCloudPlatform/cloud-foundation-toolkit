package launchpad

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// NewGenerate takes file patterns as input YAMLs and output Infrastructure as
// Code ready scripts based on specified output flavor.
//
// NewGenerate can be triggered by
//
//	$ cft launchpad generate *.yaml
//	$ cft lp g *.yaml
func NewGenerate(rawPaths []string, outFlavor OutputFlavor, outputDir string) {
	// attempt to load all configs with best effort
	log.Println("debug: output location", outputDir) // Remove after generate code is written
	log.Println("debug: output flavor", outFlavor)   // Remove after generate code is written
	resources := loadResources(rawPaths)
	log.Println(len(resources), "YAML documents loaded")

	assembled := assembleResourcesToOrg(resources)

	print(assembled.String()) // Place-holder for future trigger point of code generation
}

// OutputFlavor defines launchpad's generated output language.
type OutputFlavor int

const (
	DeploymentManager OutputFlavor = iota
	Terraform
)

// String returns the string representation of an OutputFlavor.
func (f OutputFlavor) String() string {
	return []string{"DeploymentManager", "Terraform"}[f]
}

// newOutputFlavor parses string formatted output flavor and convert to internal format.
//
// Unsupported format given will terminate the application.
func NewOutputFlavor(fStr string) OutputFlavor {
	switch strings.ToLower(fStr) {
	case "deploymentmanager", "dm":
		log.Println("Warning: Deployment Manager format not yet supported")
		return DeploymentManager
	case "terraform", "tf":
		return Terraform
	default:
		log.Fatalln("Unsupported output flavor", fStr)
	}
	return -1
}

const yamlDelimiter = "---\n"

// loadResources attempts to load YAMLs from all given path patterns in best effort.
//
// loadResources will silently ignore file I/O related errors, attempt to parse all
// files and extract resources if possible.
func loadResources(rawPaths []string) []resourceHandler {
	var buff []resourceHandler
	for _, pathPattern := range rawPaths {
		matches, err := filepath.Glob(pathPattern)
		if err != nil {
			log.Println("warning: Invalid file path pattern", pathPattern)
			continue
		}
		for _, fp := range matches {
			content, err := loadFile(fp)
			if err != nil {
				log.Println("warning: Unable to load requested file", fp)
				continue
			}
			// Multiple YAML doc can exist within one file
			docStrs := strings.Split(content, yamlDelimiter)
			for _, docStr := range docStrs {
				resource, err := loadYAML([]byte(docStr))
				if err != nil {
					log.Printf("warning: unable to parse YAML [%s]:\n%s", err.Error(), docStr)
					continue
				}
				if err = resource.validate(); err != nil {
					log.Printf("warning: YAML validation failed [%s]:\n%s", err.Error(), docStr)
					continue
				}
				buff = append(buff, resource)
			}
		}
	}
	return buff
}

// loadFile return the file content with the specified relative path to current location.
//
// loadFile will attempt to load from filesystem directly first, a not found will attempt
// to load from statics variable generated from `$ go generate`. A user using output binary
// can in theory place their own file in matching relative path and overwrite the binary
// default.
func loadFile(fp string) (string, error) {
	if content, err := os.ReadFile(fp); err == nil {
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

// loadYAML loads given byte slice as a CFT resource.
//
// loadYAML takes two pass to load YAML, first to determine the CRD kind,
// second to load the YAML into specific CFT resource.
func loadYAML(docStr []byte) (resourceHandler, error) {
	h := &headerYAML{}
	err := yaml.Unmarshal(docStr, h)
	if err != nil {
		log.Printf("Malformed YAML")
		return nil, err
	}
	kinds, ok := supportedVersion[h.APIVersion]
	if !ok {
		log.Printf("Not supported version")
		return nil, errors.New("unsupported version")
	}
	resourceFunc, ok := kinds[h.kind()]
	if !ok {
		log.Printf("Not supported kind")
		return nil, errors.New("unsupported custom resource kind for the version")
	}
	resource := resourceFunc()
	err = yaml.Unmarshal(docStr, resource)
	if err != nil {
		log.Printf("Malformed YAML")
		return nil, errors.New("unsupported custom resource kind")
	}
	return resource, nil
}
