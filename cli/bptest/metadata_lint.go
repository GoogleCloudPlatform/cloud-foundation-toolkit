package bptest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"gopkg.in/yaml.v2"
)

// BlueprintConnection structure to match metadata YAML
type BlueprintConnection struct {
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
}

// BlueprintVariable structure to match metadata YAML
type BlueprintVariable struct {
	Connections []BlueprintConnection `yaml:"connections"`
}

// BlueprintInterface structure to match metadata YAML
type BlueprintInterface struct {
	Variables []BlueprintVariable `yaml:"variables"`
}

// BlueprintMetadata structure for entire metadata YAML
type BlueprintMetadata struct {
	Spec struct {
		Interfaces []BlueprintInterface `yaml:"interfaces"`
	} `yaml:"spec"`
}

func metadataLint(metadataPath string) error {
	fmt.Printf("Start Linting metadata file: %s\n", metadataPath)
	if err := lintConnectionSourceAndVersion(metadataPath); err != nil {
		return err
	}
	fmt.Printf("Lint test all pass for metadata.yaml file: %s\n", metadataPath)
	return nil
}

func lintConnectionSourceAndVersion(metadataPath string) error {
	fmt.Printf("Linting: %s\n", metadataPath)

	// Read the metadata file
	content, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("error reading metadata file %s: %v", metadataPath, err)
	}

	// Unmarshal the YAML content into the BlueprintMetadata struct
	var metadata BlueprintMetadata
	if err := yaml.Unmarshal(content, &metadata); err != nil {
		return fmt.Errorf("error parsing metadata file %s: %v", metadataPath, err)
	}

	// Walk through the BlueprintInterface -> BlueprintVariable -> BlueprintConnection
	for _, blueprintInterface := range metadata.Spec.Interfaces {
		for _, blueprintVariable := range blueprintInterface.Variables {
			for _, connection := range blueprintVariable.Connections {
				// Validate blueprint source and version
				if err := validateBlueprintSource(connection.Source); err != nil {
					return err
				}
				if err := validateBlueprintVersion(connection.Version, connection.Source); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// validateBlueprintSource checks if the blueprint source exists
func validateBlueprintSource(source string) error {
	// Regex to validate source format
	sourceRegex := regexp.MustCompile(`^[a-zA-Z0-9\-\.]+/[a-zA-Z0-9\-\.]+$`)
	if !sourceRegex.MatchString(source) {
		return fmt.Errorf("Invalid blueprint source format: %s. Expected format: namespace/repository", source)
	}

	// Check if the blueprint source exists.
	resp, err := http.Get(fmt.Sprintf("https://github.com/%s", source))
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Blueprint source not found: %s", source)
	}

	fmt.Printf("Valid source: %s\n", source)
	return nil
}

// validateBlueprintVersion checks if the given version constraint is valid and if the specified version exists.
func validateBlueprintVersion(version string, source string) error {
	// Regex for Terraform version constraints
	versionConstraintRegex := regexp.MustCompile(`^((>=|<=|>|<|=|~>|\^)?\s?\d+(\.\d+){0,2})(,\s*(>=|<=|>|<|=|~>|\^)?\s?\d+(\.\d+){0,2})*$`)

	// Check if the version matches the version constraint syntax
	if !versionConstraintRegex.MatchString(version) {
		return fmt.Errorf("Invalid version constraint format: %s. Expected format: =1.0.0, >= 1.0.0, < 2.0.0, ~> 1.0", version)
	}

	// If the version constraint is exact (e.g., `=1.0.0`), validate that the version exists
	if exactVersion := extractExactVersion(version); exactVersion != "" {
		url := fmt.Sprintf("https://github.com/%s/releases/tag/%s", source, exactVersion)

		// Make an HTTP request to check if the version exists
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			return fmt.Errorf("Blueprint version not found: %s", exactVersion)
		}

		// Version exists
		fmt.Printf("Valid and existing version: %s\n", exactVersion)
	}

	return nil
}

// extractExactVersion checks if the version constraint is exact (e.g., `=1.0.0`) and extracts the version number.
func extractExactVersion(version string) string {
	// Regex to capture exact version numbers, e.g., =1.0.0
	exactVersionRegex := regexp.MustCompile(`^=([\d]+\.[\d]+\.[\d]+)$`)

	if matches := exactVersionRegex.FindStringSubmatch(version); len(matches) == 2 {
		return matches[1]
	}

	return ""
}
