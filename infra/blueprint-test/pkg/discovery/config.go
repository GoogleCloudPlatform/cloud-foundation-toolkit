package discovery

import (
	"fmt"
	"os"
	"path"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	defaultTestConfigFilename = "test.yaml"
	blueprintTestKind         = "BlueprintTest"
	blueprintTestAPIVersion   = "blueprints.cloud.google.com/v1alpha1"
)

type BlueprintTestConfig struct {
	yaml.ResourceMeta `json:",inline" yaml:",inline"`
	Spec              struct {
		Skip []string `json:"skip" yaml:"skip"`
	} `json:"spec" yaml:"spec"`
}

// ShouldSkipTest checks if a given test should be skipped
func (b BlueprintTestConfig) ShouldSkipTest(dir string) bool {
	testDir := path.Base(dir)
	for _, skip := range b.Spec.Skip {
		if skip == testDir {
			return true
		}
	}
	return false
}

// getTestConfig returns BlueprintTestConfig if found
func getTestConfig(path string) (*BlueprintTestConfig, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// BlueprintTestConfig does not exist, so we return an equivalent empty config
		emptyCfg := BlueprintTestConfig{}
		emptyCfg.Spec.Skip = []string{}
		return &emptyCfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", path, err)
	}
	var bpc BlueprintTestConfig
	err = yaml.Unmarshal(data, &bpc)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling %s: %v", data, err)
	}
	err = isValidTestConfig(bpc)
	if err != nil {
		return nil, fmt.Errorf("error validating testconfig in %s: %v", path, err)
	}
	return &bpc, nil
}

// isValidTestConfig validates a given BlueprintTestConfig
func isValidTestConfig(b BlueprintTestConfig) error {
	if b.APIVersion != blueprintTestAPIVersion {
		return fmt.Errorf("invalid APIVersion %s expected %s", b.APIVersion, blueprintTestAPIVersion)
	}
	if b.Kind != blueprintTestKind {
		return fmt.Errorf("invalid Kind %s expected %s", b.Kind, blueprintTestKind)
	}
	return nil
}
