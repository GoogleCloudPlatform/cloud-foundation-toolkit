package bptest

import (
	"fmt"
	"strings"
)

// BlueprintSourceRule checks if the blueprint connection source is valid.
type BlueprintSourceRule struct{}

func (r *BlueprintSourceRule) Name() string {
	return "blueprint_source_rule"
}

func (r *BlueprintSourceRule) Enabled() bool {
	return true
}

func (r *BlueprintSourceRule) Link() string {
	return "https://example.com/blueprint-source-guidelines"
}

// Check validates the source format and ensures it exists.
func (r *BlueprintSourceRule) Check(ctx LintContext) error {
	for _, conn := range ctx.Metadata.Spec.Interfaces.Variables {
		if conn.Source != "" {
			if err := validateBlueprintSource(conn.Source); err != nil {
				return err
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
