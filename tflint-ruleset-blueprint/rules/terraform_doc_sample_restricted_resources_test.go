package rules

import (
	"path"
	"testing"
)

const (
	restrictedResourcesTestDir = "doc-sample-restricted-resources"
)

func TestRestrictedResources(t *testing.T) {
	tests := []ruleTC{
		{
			dir: path.Join(restrictedResourcesTestDir, "valid"),
		},
		{
			dir: path.Join(restrictedResourcesTestDir, "single-invalid"),
		},
		{
			dir: path.Join(restrictedResourcesTestDir, "multiple-invalid"),
		},
	}

	rule := NewTerraformDocSamplesRestrictedResources()

	for _, tc := range tests {
		t.Run(tc.dir, func(t *testing.T) {
			ruleTest(t, rule, tc)
		})
	}
}
