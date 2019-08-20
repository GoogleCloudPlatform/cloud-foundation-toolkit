package launchpad

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

// Whomever implements component can be processed by the functionality loop
type component interface {
	componentName() string
}

var activeComponents []component

func init() {
	activeComponents = []component{
		newOutputDirectory(), // Create Top Level Output Directory for all Launchpad configs
		newFolders(),         // GCP Folder Generation
	}
}

// Entry point for output generation
// Treating each component as separate entity, and each component is processed via a loop of functionality in sequence
// If a component wishes to be processed by a functionality, it has too implement required interface
func generateOutput() {
	for _, c := range activeComponents {
		// Apply Functionality to each component
		withDirectory(c)
		withFiles(c)
	}

	// Run terraform fmt to re-indent
	if gState.outputFlavor == outTf {
		_, err := exec.Command("terraform", "fmt", gState.outputDirectory).Output()
		if err != nil {
			// Only warning user since output terraform files are technically able to execute, just not indented properly
			log.Printf("Failed to format terraform output")
		}
	}
}

// ==== Core Components ===
type outputDirectory struct{}

// implement directoryOwner to generate top level directory
func (l *outputDirectory) directoryProperty() *directoryProperty {
	return newDirectoryProperty(
		gState.outputDirectory,
		directoryPropertyBackup(false))
}
func newOutputDirectory() *outputDirectory       { return &outputDirectory{} }
func (l *outputDirectory) componentName() string { return "outputDirectory" }

// ==== Components ====
type folders struct {
	YAMLs   map[string]*folderYAML
	dirname string
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
			mainCons = append(mainCons, newTfGoogleFolder(y.Spec.Id, y.Spec.DisplayName, y.Spec.ParentRef.ParentId))
			outputCons = append(outputCons, newTfOutput(y.Spec.Id, fmt.Sprintf("${google_folder.%s.name}", y.Spec.Id)))
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
