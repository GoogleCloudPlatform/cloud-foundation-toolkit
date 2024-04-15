package main

import (
	"github.com/cloud-foundation-toolkit/tflint-ruleset-blueprint/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "blueprint",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewTerraformDocSamplesRestrictedBlocks(),
				rules.NewTerraformDocSamplesRestrictedResources(),
			},
		},
	})
}
