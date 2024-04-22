package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformDocSamplesRestrictedResources checks whether restricted resources are used.
type TerraformDocSamplesRestrictedResources struct {
	tflint.DefaultRule
}

// NewTerraformDocSamplesRestrictedResources returns a new rule.
func NewTerraformDocSamplesRestrictedResources() *TerraformDocSamplesRestrictedResources {
	return &TerraformDocSamplesRestrictedResources{}
}

// Name returns the rule name.
func (r *TerraformDocSamplesRestrictedResources) Name() string {
	return "terraform_doc_sample_restricted_resources"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformDocSamplesRestrictedResources) Enabled() bool {
	return false
}

// Severity returns the rule severity.
func (r *TerraformDocSamplesRestrictedResources) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *TerraformDocSamplesRestrictedResources) Link() string {
	return "https://googlecloudplatform.github.io/samples-style-guide/#language-specific"
}

const (
	nullResource = "null_resource"
)

var restrictedResources = []string{nullResource}

// Check checks whether config contains restricted resource types.
func (r *TerraformDocSamplesRestrictedResources) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// Each sample must be a root module.
		return nil
	}

	for _, restrictedResource := range restrictedResources {
		content, err := runner.GetResourceContent(restrictedResource, &hclext.BodySchema{}, nil)
		if err != nil {
			return err
		}
		for _, b := range content.Blocks {
			err := runner.EmitIssue(r, fmt.Sprintf("doc sample restricted resource type: %s", restrictedResource), b.DefRange)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
