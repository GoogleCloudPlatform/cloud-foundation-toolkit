package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformDocSamplesRestrictedBlocks checks whether restricted blocks are used.
type TerraformDocSamplesRestrictedBlocks struct {
	tflint.DefaultRule
}

// NewTerraformDocSamplesRestrictedBlocks returns a new rule.
func NewTerraformDocSamplesRestrictedBlocks() *TerraformDocSamplesRestrictedBlocks {
	return &TerraformDocSamplesRestrictedBlocks{}
}

// Name returns the rule name.
func (r *TerraformDocSamplesRestrictedBlocks) Name() string {
	return "terraform_doc_sample_restricted_blocks"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformDocSamplesRestrictedBlocks) Enabled() bool {
	return false
}

// Severity returns the rule severity.
func (r *TerraformDocSamplesRestrictedBlocks) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *TerraformDocSamplesRestrictedBlocks) Link() string {
	return "https://googlecloudplatform.github.io/samples-style-guide/#language-specific"
}

const (
	moduleBlockType   = "module"
	variableBlockType = "variable"
)

var restrictedBlocks = []string{moduleBlockType, variableBlockType}

// Check checks whether config contains restricted block types.
func (r *TerraformDocSamplesRestrictedBlocks) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// Each sample must be a root module.
		return nil
	}

	// Extract restricted blocks if any from config.
	restrictedBlocksSchema := make([]hclext.BlockSchema, 0, len(restrictedBlocks))
	for _, rb := range restrictedBlocks {
		rs := hclext.BlockSchema{
			Type:       rb,
			LabelNames: []string{"name"},
			Body:       &hclext.BodySchema{},
		}
		restrictedBlocksSchema = append(restrictedBlocksSchema, rs)
	}
	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: restrictedBlocksSchema,
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	// Emit issues if extracted blocks are found.
	blocks := body.Blocks.ByType()
	for _, rBlockType := range restrictedBlocks {
		rBlocks, ok := blocks[rBlockType]
		if ok {
			for _, rBlock := range rBlocks {
				err := runner.EmitIssue(
					r,
					fmt.Sprintf("doc sample restricted block type %s", rBlockType),
					rBlock.DefRange,
				)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
