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

func (r *BlueprintVersionRule) Check(ctx LintContext) error {
	// Check if Spec or Interfaces is nil to avoid null pointer dereference
	if ctx.Metadata == nil || ctx.Metadata.Spec == nil || ctx.Metadata.Spec.Interfaces == nil {
		fmt.Println("metadata, spec, or interfaces are nil")
		return nil
	}

	for _, variable := range ctx.Metadata.Spec.Interfaces.Variables {
		if variable == nil {
			continue // Skip if variable is nil
		}

		for _, conn := range variable.Connections {
			if conn == nil || conn.Source == nil {
				continue // Skip if connection or source is nil
			}

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
