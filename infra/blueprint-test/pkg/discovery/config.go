package discovery

import (
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	DefaultTestConfigFilename = "test.yaml"
	blueprintTestKind         = "BlueprintTest"
	blueprintTestAPIVersion   = "blueprints.cloud.google.com/v1alpha1"
)

type BlueprintTestConfig struct {
	yaml.ResourceMeta `json:",inline" yaml:",inline"`
	Spec              struct {
		Skip bool `json:"skip" yaml:"skip"`
	} `json:"spec" yaml:"spec"`
	Path string
}

// GetTestConfig returns BlueprintTestConfig if found
func GetTestConfig(path string) (BlueprintTestConfig, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// BlueprintTestConfig does not exist, so we return an equivalent empty config
		emptyCfg := BlueprintTestConfig{}
		emptyCfg.Spec.Skip = false
		return emptyCfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return BlueprintTestConfig{}, fmt.Errorf("error reading %s: %v", path, err)
	}
	var bpc BlueprintTestConfig
	err = yaml.Unmarshal(data, &bpc)
	if err != nil {
		return BlueprintTestConfig{}, fmt.Errorf("error unmarshalling %s: %v", data, err)
	}
	bpc.Path = path
	err = isValidTestConfig(bpc)
	if err != nil {
		return BlueprintTestConfig{}, fmt.Errorf("error validating testconfig in %s: %v", path, err)
	}
	return bpc, nil
}

// isValidTestConfig validates a given BlueprintTestConfig
func isValidTestConfig(b BlueprintTestConfig) error {
	if b.ResourceMeta.APIVersion != blueprintTestAPIVersion {
		return fmt.Errorf("invalid APIVersion %s expected %s", b.ResourceMeta.APIVersion, blueprintTestAPIVersion)
	}
	if b.ResourceMeta.Kind != blueprintTestKind {
		return fmt.Errorf("invalid Kind %s expected %s", b.ResourceMeta.Kind, blueprintTestKind)
	}
	return nil
}
