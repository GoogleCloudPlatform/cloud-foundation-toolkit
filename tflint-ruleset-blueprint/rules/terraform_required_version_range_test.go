package rules

import (
	"path"
	"testing"
)

const (
	TerraformRequiredVersionRangeTestDir = "terraform_required_version_range"
)

func TestTerraformMinimumRequiredVersionRange(t *testing.T) {
	tests := []ruleTC{
		{
			dir: path.Join(TerraformRequiredVersionRangeTestDir, "multiple-valid"),
		},
		{
			dir: path.Join(TerraformRequiredVersionRangeTestDir, "multiple-invalid"),
		},
		{
			dir: path.Join(TerraformRequiredVersionRangeTestDir, "multiple-valid-config"),
		},
		{
			dir: path.Join(TerraformRequiredVersionRangeTestDir, "multiple-invalid-config"),
		},
		{
			dir: path.Join(TerraformRequiredVersionRangeTestDir, "multiple-valid-config-single"),
		},
	}

	rule := NewTerraformRequiredVersionRange()

	for _, tc := range tests {
		t.Run(tc.dir, func(t *testing.T) {
			ruleTest(t, rule, tc)
		})
	}
}
