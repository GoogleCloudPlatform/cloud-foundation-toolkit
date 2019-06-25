// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// configs helps with loading and parsing configuration files
package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/smallfish/simpleyaml"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type yamlFile struct {
	source       string // helpful information to rediscover this data
	yaml         *simpleyaml.Yaml
	fileContents []byte
}

// UnclassifiedConfig stores loosely parsed information not specific to constraints or templates.
type UnclassifiedConfig struct {
	Group        string
	MetadataName string
	Kind         string
	Yaml         *simpleyaml.Yaml
	// keep the file path to help debug logging
	FilePath string
	// Preserve the raw user data to forward into rego
	// This prevents any data loss issues from going though parsing libraries.
	RawFile string
}

// ConstraintTemplate stores parsed information including the raw data.
type ConstraintTemplate struct {
	Confg *UnclassifiedConfig
	// This is the kind that this template generates.
	GeneratedKind string
	Rego          string
}

// Constraint stores parsed information including the raw data.
type Constraint struct {
	Confg *UnclassifiedConfig
}

const (
	validTemplateGroup   = "templates.gatekeeper.sh/v1alpha1"
	validConstraintGroup = "constraints.gatekeeper.sh/v1alpha1"
	expectedTarget       = "validation.gcp.forsetisecurity.org"
)

// AsInterface returns the the config data as a structured golang object. This uses yaml.Unmarshal to create this object.
func (c *UnclassifiedConfig) AsInterface() (interface{}, error) {
	// Use yaml.Unmarshal to create a proper golang object that maintains the same structure
	var f interface{}
	if err := yaml.Unmarshal([]byte(c.RawFile), &f); err != nil {
		return nil, errors.Wrap(err, "converting from yaml")
	}
	return f, nil
}

// asConstraint attempts to convert to constraint
// Returns:
//   *Constraint: only set if valid constraint
//   bool: (always set) if this is a constraint
func asConstraint(data *UnclassifiedConfig) (*Constraint, error) {
	// There is no validation matching this constraint to the template here that happens after
	// basic parsing has happened when we have more context.
	if data.Group != validConstraintGroup {
		return nil, fmt.Errorf("group expected to be %s not %s", validConstraintGroup, data.Group)
	}
	if data.Kind == "ConstraintTemplate" {
		return nil, fmt.Errorf("kind should not be ConstraintTemplate")
	}
	return &Constraint{
		Confg: data,
	}, nil
}

// asConstraintTemplate attempts to convert to template
// Returns:
//   *ConstraintTemplate: only set if valid template
//   bool: (always set) if this is a template
func asConstraintTemplate(data *UnclassifiedConfig) (*ConstraintTemplate, error) {
	if data.Group != validTemplateGroup {
		return nil, fmt.Errorf("group expected to be %s not %s", validTemplateGroup, data.Group)
	}
	if data.Kind != "ConstraintTemplate" {
		return nil, fmt.Errorf("kind expected to be ConstraintTemplate not %s", data.Kind)
	}
	generatedKind, err := data.Yaml.GetPath("spec", "crd", "spec", "names", "kind").String()
	if err != nil {
		return nil, err // field expected to exist
	}
	rego, err := extractRego(data.Yaml)
	if err != nil {
		return nil, err // field expected to exist
	}
	return &ConstraintTemplate{
		Confg:         data,
		GeneratedKind: generatedKind,
		Rego:          rego,
	}, nil
}

func extractRego(yaml *simpleyaml.Yaml) (string, error) {
	targets := yaml.GetPath("spec", "targets")
	if !targets.IsArray() {
		// Old format looks like the following
		// targets:
		//   validation.gcp.forsetisecurity.org:
		//     rego:
		return targets.GetPath(expectedTarget, "rego").String()
	}
	// New format looks like the following
	// targets:
	//  - target: validation.gcp.forsetisecurity.org
	//    rego:
	size, err := targets.GetArraySize()
	if err != nil {
		return "", err
	}
	for i := 0; i < size; i++ {
		target := targets.GetIndex(i)
		targetString, err := target.Get("target").String()
		if err != nil {
			return "", err
		}
		if targetString == expectedTarget {
			return target.Get("rego").String()
		}
	}

	return "", status.Error(codes.InvalidArgument, "Unable to locate rego field in constraint template")
}

func arrayFilterSuffix(arr []string, suffix string) []string {
	filteredList := []string{}
	for _, s := range arr {
		if strings.HasSuffix(strings.ToLower(s), strings.ToLower(suffix)) {
			filteredList = append(filteredList, s)
		}
	}
	return filteredList
}

// listFiles returns a list of files under a dir. Errors will be grpc errors.
func listFiles(dir string) ([]string, error) {
	files := []string{}

	visit := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "error visiting path %s", path)
		}
		if !f.IsDir() {
			files = append(files, path)
		}
		return nil
	}

	err := filepath.Walk(dir, visit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return files, nil
}

// ListYAMLFiles returns a list of YAML files under a dir. Errors will be grpc errors.
func ListYAMLFiles(dir string) ([]string, error) {
	files, err := listFiles(dir)
	if err != nil {
		return nil, err
	}
	return arrayFilterSuffix(files, ".yaml"), nil
}

// ListRegoFiles returns a list of rego files under a dir. Errors will be grpc errors.
func ListRegoFiles(dir string) ([]string, error) {
	files, err := listFiles(dir)
	if err != nil {
		return nil, err
	}
	return arrayFilterSuffix(files, ".rego"), nil
}

// convertYAMLToUnclassifiedConfig converts yaml file to an unclassified config, if expected fields don't exist, a log message is printed and the config is skipped.
func convertYAMLToUnclassifiedConfig(config *yamlFile) (*UnclassifiedConfig, error) {
	kind, err := config.yaml.Get("kind").String()
	if err != nil {
		return nil, fmt.Errorf("error in converting %s: %v", config.source, err)
	}
	group, err := config.yaml.Get("apiVersion").String()
	if err != nil {
		return nil, fmt.Errorf("error in converting %s: %v", config.source, err)
	}
	metadataName, err := config.yaml.GetPath("metadata", "name").String()
	if err != nil {
		return nil, fmt.Errorf("error in converting %s: %v", config.source, err)
	}
	convertedConfig := &UnclassifiedConfig{
		Group:        group,
		MetadataName: metadataName,
		Kind:         kind,
		Yaml:         config.yaml,
		FilePath:     config.source,
		RawFile:      string(config.fileContents),
	}
	return convertedConfig, nil
}

// Returns either a *ConstraintTemplate or a *Constraint or an error
// dataSource should be helpful documentation to help rediscover the source of this information.
func CategorizeYAMLFile(data []byte, dataSource string) (interface{}, error) {
	y, err := simpleyaml.NewYaml(data)
	if err != nil {
		return nil, err
	}
	unclassified, err := convertYAMLToUnclassifiedConfig(&yamlFile{
		yaml:         y,
		fileContents: data,
		source:       dataSource,
	})
	if err != nil {
		return nil, err
	}
	switch unclassified.Group {
	case validTemplateGroup:
		return asConstraintTemplate(unclassified)
	case validConstraintGroup:
		return asConstraint(unclassified)
	}
	return nil, fmt.Errorf("unable to determine configuration type for data %s", dataSource)
}
