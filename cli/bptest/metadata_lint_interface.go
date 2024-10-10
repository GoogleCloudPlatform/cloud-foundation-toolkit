package bptest

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
)

// LintRule defines the common interface for all metadata lint rules.
type LintRule interface {
	Name() string            // Unique name of the rule
	Enabled() bool           // Indicates if the rule is enabled by default
	Link() string            // A reference link for the rule
	Check(LintContext) error // Main entrypoint for rule validation
}

// LintContext holds the metadata and other contextual information for a rule.
type LintContext struct {
	Metadata *bpmetadata.BlueprintMetadata // Parsed metadata for the blueprint
	FilePath string                        // Path of the metadata file being checked
}

// LintRunner is responsible for running all registered lint rules.
type LintRunner struct {
	Rules []LintRule
}

// RegisterRule adds a new rule to the runner.
func (r *LintRunner) RegisterRule(rule LintRule) {
	r.Rules = append(r.Rules, rule)
}

// Run runs all the registered rules on the provided context.
func (r *LintRunner) Run(ctx LintContext) []error {
	var errs []error
	for _, rule := range r.Rules {
		if rule.Enabled() {
			err := rule.Check(ctx)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}
