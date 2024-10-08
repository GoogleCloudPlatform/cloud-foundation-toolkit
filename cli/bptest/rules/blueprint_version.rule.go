package bptest

import (
	"fmt"
	"regexp"
)

// BlueprintVersionRule checks if the blueprint connection version is valid.
type BlueprintVersionRule struct{}

func (r *BlueprintVersionRule) Name() string {
	return "blueprint_version_rule"
}

func (r *BlueprintVersionRule) Enabled() bool {
	return true
}

func (r *BlueprintVersionRule) Link() string {
	return "https://example.com/blueprint-version-guidelines"
}

// Check validates the version format and ensures it exists.
func (r *BlueprintVersionRule) Check(ctx LintContext) error {
	for _, conn := range ctx.Metadata.Spec.Interfaces.Variables {
		if conn.Version != "" {
			if err := validateVersion(conn.Version, conn.Source); err != nil {
				return err
			}
		}
	}
	return nil
}

// // validateVersion checks for semantic versioning.
// func validateVersion(version string) error {
// matched, err := regexp.MatchString(`^v?(\d+\.\d+\.\d+)$`, version)
// if err != nil || !matched {
// return fmt.Errorf("invalid version format: %s", version)
// }
// return nil
// }
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
