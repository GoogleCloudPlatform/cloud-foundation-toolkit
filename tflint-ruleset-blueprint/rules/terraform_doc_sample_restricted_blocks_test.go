package rules

import (
	"path"
	"testing"
)

const (
	restrictedBlocksTestDir = "doc-sample-restricted-blocks"
)

func TestRestrictedBlocks(t *testing.T) {
	tests := []ruleTC{
		{
			dir: path.Join(restrictedBlocksTestDir, "valid"),
		},
		{
			dir: path.Join(restrictedBlocksTestDir, "single-restricted"),
		},
		{
			dir: path.Join(restrictedBlocksTestDir, "multiple-restricted"),
		},
	}

	rule := NewTerraformDocSamplesRestrictedBlocks()

	for _, tc := range tests {
		t.Run(tc.dir, func(t *testing.T) {
			ruleTest(t, rule, tc)
		})
	}
}
