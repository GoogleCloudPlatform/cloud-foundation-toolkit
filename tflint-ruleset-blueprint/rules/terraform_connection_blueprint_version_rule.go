package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformConnectionBlueprintVersionRule checks blueprint connection source reference and version.
type TerraformConnectionBlueprintVersionRule struct {
	tflint.DefaultRule
}

// NewTerraformDocSamplesRestrictedBlocks returns a new rule.
func NewTerraformConnectionBlueprintVersionRule() *TerraformConnectionBlueprintVersionRule {
	return &TerraformConnectionBlueprintVersionRule{}
}

// Name returns the rule name.
func (r *TerraformConnectionBlueprintVersionRule) Name() string {
	return "terraform_connection_blueprint_version_rule"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformConnectionBlueprintVersionRule) Enabled() bool {
	return false
}

// Severity returns the rule severity.
func (r *TerraformConnectionBlueprintVersionRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *TerraformConnectionBlueprintVersionRule) Link() string {
	// TODO: Update the rule reference link
	return "https://cloud.google.com/docs/blueprints"
}

// Check checks the reference of connection and version.
func (r *TerraformConnectionBlueprintVersionRule) Check(runner tflint.Runner) error {
	return runner.WalkAttributes(func(attribute *terraform.Attribute) error {
		// check source reference in connection.
		if attribute.Name == "source" {
			value, diags := attribute.Value(nil)
			if diags.HasErrors() {
				return nil // Skip if there's an error getting the value
			}

			sourceStr := value.AsString()
			if !r.regex.MatchString(sourceStr) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("Invalid blueprint reference in 'source': %s", sourceStr),
					attribute.Expr.Range(),
				)
			}
		}
		return nil
	})
}
