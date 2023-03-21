package kpt

import (
	"fmt"
	"os"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// UpsertSetters inserts or updates setters if apply-setters fn config is discovered.
func UpsertSetters(nodes []*yaml.RNode, setters map[string]string) error {
	kf, err := findKptfile(nodes)
	if err != nil {
		return err
	}
	// no pipeline defined for pkg
	if kf.Pipeline == nil {
		return nil
	}
	for _, fn := range kf.Pipeline.Mutators {
		if !strings.Contains(fn.Image, "apply-setters") {
			continue
		}
		// ignoring inlined configMap in Kptfile for now
		// all blueprint examples will have setters defined via configPath
		if fn.ConfigPath != "" {
			settersConfig, err := findSetterNode(nodes, fn.ConfigPath)
			if err != nil {
				return err
			}
			setterData := settersConfig.GetDataMap()
			for sKey, sVal := range setters {
				setterData[sKey] = sVal
			}
			settersConfig.SetDataMap(setterData)
		}
	}
	return nil
}

// findSetterNode finds setter node from a slice of nodes.
func findSetterNode(nodes []*yaml.RNode, path string) (*yaml.RNode, error) {
	for _, node := range nodes {
		np := node.GetAnnotations()[kioutil.PathAnnotation]
		if np == path {
			return node, nil
		}
	}
	return nil, fmt.Errorf(`file %s doesn't exist, please ensure the file specified in "configPath" exists and retry`, path)
}

// Generates setters from environment variables.
// Setter names are generated from variable name by lowercasing and replacing "_" to "-".
func GenerateSetterKVFromEnvVar(e string) (string, string, error) {
	sVal, found := os.LookupEnv(e)
	if !found {
		return "", "", fmt.Errorf("unable to find envvar %s", e)
	}
	sKey := strings.ReplaceAll(strings.ToLower(e), "_", "-")
	return sKey, sVal, nil
}

// MergeSetters merges two setter maps a and b.
// If duplicate key map b takes precedence.
func MergeSetters(a, b map[string]string) map[string]string {
	merged := make(map[string]string, len(a)+len(b))
	for k, v := range a {
		merged[k] = v
	}
	for k, v := range b {
		merged[k] = v
	}
	return merged
}

// ReadPkgResources returns a slice of resources from a dir.
func ReadPkgResources(dir string) ([]*yaml.RNode, error) {
	p := &kio.LocalPackageReader{
		PackagePath:        dir,
		PackageFileName:    "Kptfile",
		MatchFilesGlob:     append(kio.DefaultMatch, "Kptfile"),
		IncludeSubpackages: true,
	}
	return p.Read()
}

// WritePkgResources writes a slice of resources to a dir.
func WritePkgResources(dir string, rs []*yaml.RNode) error {
	p := &kio.LocalPackageWriter{
		PackagePath: dir,
	}
	return p.Write(rs)
}
