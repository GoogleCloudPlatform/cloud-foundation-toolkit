// Package launchpad file generate.go contains all output generation logic
//
// A component is a set of related scripts that generally resides under the same
// folder beneath outputDirectory root. A functionality is a particular action
// that can be applied to a component to achieve some purpose.
//
// Output generation depends on evaluated gState, and looping through components
// in specified order to apply functionality in sequence to generate output
// based on defined outputFlavor.
package launchpad

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

// component interface allows implementer to be put onto the generate processing loop
type component interface {
	componentName() string
}

// generateOutput loops through components and applies functionality in sequence
func generateOutput() {
	activeComponents := []component{
		newOutputDirectory(), // Create Top Level Output Directory for all Launchpad configs
		newFolders(),         // GCP Folder Generation
	}

	for _, c := range activeComponents {
		// Apply Functionality to each component
		withDirectory(c)
		withFiles(c)
	}

	// re-indent with terraform fmt
	if gState.outputFlavor == outTf {
		_, err := exec.Command("terraform", "fmt", gState.outputDirectory).Output()
		if err != nil {
			// Only warning user since output terraform files are technically able to execute, just not indented properly
			log.Printf("Failed to format terraform output")
		}
	}
}

// ==== Core Components ===

// outputDirectory serves as the top level output directory
type outputDirectory struct{}

// directoryProperty to implement directoryOwner
func (l *outputDirectory) directoryProperty() *directoryProperty {
	return newDirectoryProperty(
		gState.outputDirectory,
		directoryPropertyBackup(false))
}
func newOutputDirectory() *outputDirectory       { return &outputDirectory{} }
func (l *outputDirectory) componentName() string { return "outputDirectory" }

// ==== Components ====

// folders component allows generation of sub-directory under outputDirectory for GCP Folder scripts
type folders struct {
	YAMLs       map[string]*folderSpecYAML
	dirname     string
	dirProperty *directoryProperty
}

func newFolders() *folders               { return &gState.evaluated.folders }
func (f *folders) componentName() string { return "folders" }
func (f *folders) directoryProperty() *directoryProperty {
	if f.dirProperty == nil {
		f.dirProperty = newDirectoryProperty(f.componentName(), directoryPropertyBackup(false), directoryPropertyDirname(gState.outputDirectory))
	}
	return f.dirProperty
}
func (f *folders) files() (fs []file) {
	dir := f.dirProperty.path()
	switch gState.outputFlavor {
	case outTf:
		var outputCons, varCons []tfConstruct
		mainCons := []tfConstruct{newTfGoogleProvider()}
		for _, y := range f.YAMLs {
			mainCons = append(mainCons, newTfGoogleFolder(y.Id, y.DisplayName, &y.ParentRef))
			outputCons = append(outputCons, newTfOutput(y.Id, fmt.Sprintf("${google_folder.%s.name}", y.Id)))
		}
		varCons = append(
			varCons,
			newTfVariable("organization_id", "GCP Organization ID", ""),
			newTfVariable("credentials_file_path", "Service account key path", "credentials.json"),
		)

		return []file{
			newTfFile("main", dir, mainCons),
			newTfFile("output", dir, outputCons),
			newTfFile("variables", dir, varCons),
		}
	default:
		panic(errors.New("output format not yet implemented"))
	}
}
