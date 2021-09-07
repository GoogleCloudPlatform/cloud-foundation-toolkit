package discovery

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldSkipTest(t *testing.T) {
	tests := []struct {
		name         string
		testCfg      string
		checkTest    string
		shouldSkip   bool
		skippedTests []string
		errMsg       string
	}{
		{
			name: "simple",
			testCfg: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintTest
metadata:
  name: test
spec:
  skip:
    - foo
`,
			checkTest:    "bar",
			shouldSkip:   false,
			skippedTests: []string{"foo"},
		},
		{
			name: "simple skipped",
			testCfg: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintTest
metadata:
  name: test
spec:
  skip:
    - foo
`,
			checkTest:    "foo",
			shouldSkip:   true,
			skippedTests: []string{"foo"},
		},
		{
			name: "multiple skip with match",
			testCfg: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintTest
metadata:
  name: test
spec:
  skip:
    - foo
    - bar
    - baz
`,
			checkTest:    "bar",
			shouldSkip:   true,
			skippedTests: []string{"foo", "bar", "baz"},
		},
		{
			name: "multiple skip with no match",
			testCfg: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintTest
metadata:
  name: test
spec:
  skip:
    - foo
    - bar
    - baz
`,
			checkTest:    "quux",
			shouldSkip:   false,
			skippedTests: []string{"foo", "bar", "baz"},
		},
		{
			name: "fixture skipped",
			testCfg: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintTest
metadata:
  name: test
spec:
  skip:
    - with-bar
`,
			checkTest:    "fixtures/with-bar",
			shouldSkip:   true,
			skippedTests: []string{"with-bar"},
		},
		{
			name:         "empty none skipped",
			checkTest:    "fixtures/with-bar",
			shouldSkip:   false,
			skippedTests: []string{},
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
			bpTestCfg, err := getTestConfig(testCfgPath)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Equal(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.Equal(tt.skippedTests, bpTestCfg.Spec.Skip)
				assert.Equal(tt.shouldSkip, bpTestCfg.ShouldSkipTest(tt.checkTest))
			}
		})
	}
}

func setupTestCfg(t *testing.T, data string) string {
	t.Helper()
	assert := assert.New(t)
	baseDir, err := ioutil.TempDir("", "")
	assert.NoError(err)
	fPath := path.Join(baseDir, "test.yaml")
	err = ioutil.WriteFile(fPath, []byte(data), 0644)
	assert.NoError(err)
	return fPath
}
