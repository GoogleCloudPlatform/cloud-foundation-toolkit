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
			ver, err := version.NewConstraint(conn.Version)
			if err != nil {
				return fmt.Errorf("invalid version: %w", err)
			}
			return nil
		}
	}
	return nil
}
