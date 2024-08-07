package rules

import (
	"path"
	"testing"
)

const (
	terraformRequiredVersionTestDir = "terraform_required_version"
)

func TestTerraformMinimumRequiredVersion(t *testing.T) {
	tests := []ruleTC{
		{
			dir: path.Join(terraformRequiredVersionTestDir, "multiple-valid"),
		},
		{
			dir: path.Join(terraformRequiredVersionTestDir, "multiple-invalid"),
		},
	}

	rule := NewTerraformRequiredVersion()

	for _, tc := range tests {
		t.Run(tc.dir, func(t *testing.T) {
			ruleTest(t, rule, tc)
		})
	}
}
