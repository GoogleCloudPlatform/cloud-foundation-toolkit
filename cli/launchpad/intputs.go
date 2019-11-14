// Package launchpad file inputs.go contains all input processing logic.
package launchpad

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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

// outputFlavor defines launchpad's generated output language.
type outputFlavor int

const (
	DeploymentManager outputFlavor = iota
	Terraform
)

// String returns the string representation of an outputFlavor.
func (f outputFlavor) String() string {
	return []string{"DeploymentManager", "Terraform"}[f]
}

// newOutputFlavor parses string formatted output flavor and convert to internal format.
//
// Unsupported format given will terminate the application.
func newOutputFlavor(fStr string) outputFlavor {
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
