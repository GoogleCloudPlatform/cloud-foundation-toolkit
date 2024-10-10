package bptest

import (
	"fmt"
	"github.com/hashicorp/go-version"
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
	for _, variable := range ctx.Metadata.Spec.Interfaces.Variables {
		for _, conn := range variable.Connections {
			if conn.Source.Version != "" {
				_, err := version.NewConstraint(conn.Source.Version)
				if err != nil {
					return fmt.Errorf("invalid version: %w", err)
				}
				return nil
			}
		}
	}
	return nil
}
