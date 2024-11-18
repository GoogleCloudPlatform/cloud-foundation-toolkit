package bptest

import (
	"fmt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"os"
)

// lintRule defines the common interface for all metadata lint rules.
type lintRule interface {
	name() string            // Unique name of the rule
	enabled() bool           // Indicates if the rule is enabled by default
	check(lintContext) error // Main entrypoint for rule validation
}

// LintContext holds the metadata and other contextual information for a rule.
type lintContext struct {
	metadata *bpmetadata.BlueprintMetadata // Parsed metadata for the blueprint
	filePath string                        // Path of the metadata file being checked
}

// LintRunner is responsible for running all registered lint rules.
type lintRunner struct {
	rules []lintRule
}

// RegisterRule adds a new rule to the runner.
func (r *lintRunner) RegisterRule(rule lintRule) {
	r.rules = append(r.rules, rule)
}

// Run runs all the registered rules on the provided context.
func (r *lintRunner) Run(ctx lintContext) []error {
	var errs []error
	if os.Getenv("BLUEPRINT_LINT_DISABLE") == "1" {
		fmt.Println("BLUEPRINT_LINT_DISABLE is set to 1. Skipping lint checks.")
		return errs
	}

	for _, rule := range r.rules {
		if rule.enabled() {
			err := rule.check(ctx)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}
