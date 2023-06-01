package discovery

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldSkipTest(t *testing.T) {
	tests := []struct {
		name       string
		testCfg    string
		shouldSkip bool
		errMsg     string
	}{
		{
			name: "simple",
			testCfg: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintTest
metadata:
  name: test
spec:
  skip: true
`,
			shouldSkip: true,
		},
		{
			name: "invalid",
			testCfg: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: foo
metadata:
  name: test
spec:
  skip: true
`,
			errMsg: "invalid Kind foo expected BlueprintTest",
		},
		{
			name:       "empty none skipped",
			shouldSkip: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			testCfgPath := ""
			if tt.testCfg != "" {
				testCfgPath = setupTestCfg(t, tt.testCfg)
			}
			defer os.RemoveAll(testCfgPath)
			bpTestCfg, err := GetTestConfig(testCfgPath)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.Equal(tt.shouldSkip, bpTestCfg.Spec.Skip)
			}
		})
	}
}

func setupTestCfg(t *testing.T, data string) string {
	t.Helper()
	assert := assert.New(t)
	baseDir, err := os.MkdirTemp("", "")
	assert.NoError(err)
	fPath := path.Join(baseDir, "test.yaml")
	err = os.WriteFile(fPath, []byte(data), 0644)
	assert.NoError(err)
	return fPath
}
