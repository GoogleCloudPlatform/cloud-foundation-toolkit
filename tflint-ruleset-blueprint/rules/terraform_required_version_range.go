package rules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformRequiredVersionRange checks if a module has a terraform required_version within valid range.
type TerraformRequiredVersionRange struct {
	tflint.DefaultRule
}

// TerraformRequiredVersionRangeConfig is a config of TerraformRequiredVersionRange
type TerraformRequiredVersionRangeConfig struct {
	MinVersion string `hclext:"min_version,optional"`
	MaxVersion string `hclext:"max_version,optional"`
}

// NewTerraformRequiredVersionRange returns a new rule.
func NewTerraformRequiredVersionRange() *TerraformRequiredVersionRange {
	return &TerraformRequiredVersionRange{}
}

// Name returns the rule name.
func (r *TerraformRequiredVersionRange) Name() string {
	return "terraform_required_version_range"
}

// Enabled returns whether the rule is enabled by default.
func (r *TerraformRequiredVersionRange) Enabled() bool {
	return false
}

// Severity returns the rule severity.
func (r *TerraformRequiredVersionRange) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *TerraformRequiredVersionRange) Link() string {
	return "https://googlecloudplatform.github.io/samples-style-guide/#language-specific"
}

const (
	minimumTerraformRequiredVersionRange = "1.3"
	maximumTerraformRequiredVersionRange = "1.5"
)

// Checks if a module has a terraform required_version within valid range.
func (r *TerraformRequiredVersionRange) Check(runner tflint.Runner) error {
	config := &TerraformRequiredVersionRangeConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	minVersion := minimumTerraformRequiredVersionRange
	if config.MinVersion != "" {
		if _, err := version.NewSemver(config.MinVersion); err != nil {
			return err
		}
		minVersion = config.MinVersion
	}

	maxVersion := maximumTerraformRequiredVersionRange
	if config.MaxVersion != "" {
		if _, err := version.NewSemver(config.MaxVersion); err != nil {
			return err
		}
		maxVersion = config.MaxVersion
	}

	logger.Info(fmt.Sprintf("Running with min_version: %q max_version: %q", minVersion, maxVersion))

	splitVersion := strings.Split(minVersion, ".")
	majorVersion, err := strconv.Atoi(splitVersion[0])
	if err != nil {
		return err
	}
	minorVersion, err := strconv.Atoi(splitVersion[1])
	if err != nil {
		return err
	}

	var terraform_below_minimum_required_version string
	if minorVersion > 0 {
		terraform_below_minimum_required_version = fmt.Sprintf(
			"v%d.%d.999",
			majorVersion,
			minorVersion - 1,
		 )
	} else {
		if majorVersion == 0 {
			return fmt.Errorf("Error: minimum version test constraint would be below zero: v%d.%d.999", majorVersion - 1, 999)
		}
		terraform_below_minimum_required_version = fmt.Sprintf(
			"v%d.%d.999",
			majorVersion - 1,
			999,
		 )
	}

	below_required_version, err := version.NewVersion(terraform_below_minimum_required_version)
	if err != nil {
		return err
	}

	minimum_required_version, err := version.NewVersion(minVersion)
	if err != nil {
		return err
	}

	maximum_required_version, err := version.NewVersion(maxVersion)
	if err != nil {
		return err
	}

	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}

	if !path.IsRoot() {
		return nil
	}

	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "terraform",
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{{Name: "required_version"}},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		requiredVersion, exists := block.Body.Attributes["required_version"]
		if !exists {
			logger.Info(fmt.Sprintf("terraform block does not contain required_version: %s", block.DefRange))
			continue
		}

		var raw_terraform_required_version string
		diags := gohcl.DecodeExpression(requiredVersion.Expr, nil, &raw_terraform_required_version)
		if diags.HasErrors() {
			return fmt.Errorf("failed to decode terraform block required_version: %v", diags.Error())
		}

		constraints, err := version.NewConstraint(raw_terraform_required_version)
		if err != nil {
			return err
		}

		if !((constraints.Check(minimum_required_version) || constraints.Check(maximum_required_version)) && !constraints.Check(below_required_version)) {
			//TODO: use EmitIssueWithFix()
			err := runner.EmitIssue(r, fmt.Sprintf("required_version is not inclusive of the the minimum %q and maximum %q terraform required_version: %q", minVersion, maxVersion, constraints.String()), block.DefRange)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
