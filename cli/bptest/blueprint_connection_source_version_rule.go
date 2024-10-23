package bptest

import (
	"fmt"
	"github.com/hashicorp/go-version"
)

type BlueprintConnectionSourceVersionRule struct{}

func (r *BlueprintConnectionSourceVersionRule) name() string {
	return "blueprint_connection_source_version_rule"
}

func (r *BlueprintConnectionSourceVersionRule) enabled() bool {
	return true
}

func (r *BlueprintConnectionSourceVersionRule) check(ctx lintContext) error {
	// Check if Spec or Interfaces is nil to avoid null pointer dereference
	if ctx.metadata == nil || ctx.metadata.Spec == nil || ctx.metadata.Spec.Interfaces == nil {
		fmt.Println("metadata, spec, or interfaces are nil")
		return nil
	}

	for _, variable := range ctx.metadata.Spec.Interfaces.Variables {
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
