package bptest

import (
	"errors"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"github.com/stretchr/testify/assert"
)

type MockLintRule struct {
	Name    string
	Enabled bool
	Err     error
}

func (m *MockLintRule) name() string {
	return m.Name
}

func (m *MockLintRule) enabled() bool {
	return m.Enabled
}

func (m *MockLintRule) check(ctx lintContext) error {
	return m.Err
}

func TestLintRunner(t *testing.T) {
	t.Run("register and run rules with lintRunner", func(t *testing.T) {
		mockRule1 := &MockLintRule{Name: "MockRule1", Enabled: true, Err: nil}
		mockRule2 := &MockLintRule{Name: "MockRule2", Enabled: true, Err: errors.New("lint error")}
		mockRule3 := &MockLintRule{Name: "MockRule3", Enabled: false, Err: nil}

		runner := lintRunner{}
		runner.RegisterRule(mockRule1)
		runner.RegisterRule(mockRule2)
		runner.RegisterRule(mockRule3)

		ctx := lintContext{
			metadata: &bpmetadata.BlueprintMetadata{ApiVersion: "v1", Kind: "Blueprint"},
			filePath: "/path/to/metadata/file.yaml",
		}

		errs := runner.Run(ctx)
		assert.Len(t, errs, 1, "Only one rule should return an error")
		assert.Equal(t, "lint error", errs[0].Error(), "Error message should match the expected lint error")
	})

	t.Run("run without registered rules", func(t *testing.T) {
		runner := lintRunner{}
		ctx := lintContext{
			metadata: &bpmetadata.BlueprintMetadata{ApiVersion: "v1", Kind: "Blueprint"},
			filePath: "/path/to/metadata/file.yaml",
		}

		errs := runner.Run(ctx)
		assert.Empty(t, errs, "No errors should be returned when no rules are registered")
	})
	t.Run("skip lint rules when BLUEPRINT_LINT_DISABLE is set", func(t *testing.T) {
		os.Setenv("BLUEPRINT_LINT_DISABLE", "1")
		defer os.Unsetenv("BLUEPRINT_LINT_DISABLE")

		mockRule1 := &MockLintRule{Name: "MockRule1", Enabled: true, Err: errors.New("lint error")}
		mockRule2 := &MockLintRule{Name: "MockRule2", Enabled: true, Err: errors.New("another lint error")}

		runner := lintRunner{}
		runner.RegisterRule(mockRule1)
		runner.RegisterRule(mockRule2)

		ctx := lintContext{
			metadata: &bpmetadata.BlueprintMetadata{ApiVersion: "v1", Kind: "Blueprint"},
			filePath: "/path/to/metadata/file.yaml",
		}

		errs := runner.Run(ctx)
		assert.Empty(t, errs, "No errors should be returned when BLUEPRINT_LINT_DISABLE is set")
	})
}

//func TestRunLint(t *testing.T) {
//t.Run("skip lint when BLUEPRINT_LINT_DISABLE is set", func(t *testing.T) {
//os.Setenv("BLUEPRINT_LINT_DISABLE", "1")
//defer os.Unsetenv("BLUEPRINT_LINT_DISABLE")

//err := Run()
//assert.NoError(t, err, "Lint should be skipped without error")
//})

//t.Run("run lint when BLUEPRINT_LINT_DISABLE is not set", func(t *testing.T) {
//os.Unsetenv("BLUEPRINT_LINT_DISABLE")

//err := Run()
//assert.NoError(t, err, "Lint should run without error")
//})
//}
