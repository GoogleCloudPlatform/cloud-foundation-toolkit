package launchpad

import "errors"

const (
	outDm = "dm"
	outTf = "tf"
)

var gState globalState

func init() {
	gState.evaluated.folders.YAMLs = make(map[string]*folderYAML)
}

type globalState struct {
	// A simple stack implementation
	//stack []stackFrame

	outputDirectory string
	outputFlavor    string
	evaluated       struct {
		folders folders
	}
}

// ==== Constructors for runtime state ====
func newFolder(f *folderYAML) error {
	if _, ok := gState.evaluated.folders.YAMLs[f.Spec.Id]; ok {
		return errors.New("duplicated definition of folder")
	}
	gState.evaluated.folders.YAMLs[f.Spec.Id] = f
	return nil
}
