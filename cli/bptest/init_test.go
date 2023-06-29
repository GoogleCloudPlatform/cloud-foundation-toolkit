package bptest

import (
	"os"
	"path"
	"testing"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

const (
	initTestDir = "testdata/init"
)

func TestInitTest(t *testing.T) {
	tests := []struct {
		name                  string
		bptName               string
		preProcessTestDir     func(t *testing.T, dir string)
		expectedFilesContents map[string]string
		errMsg                string
	}{
		{
			name:    "simple with mod",
			bptName: "foo",
			expectedFilesContents: map[string]string{
				"test/integration/foo/foo_test.go": `package foo

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	bpt := tft.NewTFBlueprintTest(t)

	bpt.DefineVerify(func(assert *assert.Assertions) {
		bpt.DefaultVerify(assert)
		
		foo := bpt.GetStringOutput("foo")

		op := gcloud.Run(t,"")
		assert.Contains(op.Get("result").String(), "foo", "contains foo")
	})

	bpt.Test()
}
`,
				"test/integration/go.mod": "", // we create an empty go.mod in preprocess so no generation is expected
			},
			preProcessTestDir: func(t *testing.T, dir string) {
				_, err := os.Create(path.Join(dir, intTestPath, "go.mod"))
				if err != nil {
					t.Fatalf("error creating go.mod: %v", err)
				}
			},
		},
		{
			name:    "simple without mod",
			bptName: "foo",
			expectedFilesContents: map[string]string{
				"test/integration/foo/foo_test.go": `package foo

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	bpt := tft.NewTFBlueprintTest(t)

	bpt.DefineVerify(func(assert *assert.Assertions) {
		bpt.DefaultVerify(assert)
		
		foo := bpt.GetStringOutput("foo")

		op := gcloud.Run(t,"")
		assert.Contains(op.Get("result").String(), "foo", "contains foo")
	})

	bpt.Test()
}
`,
				"test/integration/go.mod": `module github.com/terraform-google-modules/init/test/integration

go 1.16

require (
	github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test v0.4.0
	github.com/stretchr/testify v1.8.1
)
`,
			},
		},
		{
			name:    "invalid already exists",
			bptName: "bar",
			errMsg:  "test/integration/bar already exists",
		},
		{
			name:    "invalid no example",
			bptName: "baz",
			errMsg:  "unable to discover test configs for test/integration/baz: unable to find config in test/fixtures/baz nor examples/baz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			// copy fixture to tmpdir
			tmpDir := path.Join(t.TempDir(), initTestDir)
			err := copy.Copy(initTestDir, tmpDir)
			assert.NoError(err)
			// apply any pre processing before tests
			if tt.preProcessTestDir != nil {
				tt.preProcessTestDir(t, tmpDir)
			}
			// switch to tmp dir for test
			t.Cleanup(switchDir(t, tmpDir))

			err = initTest(tt.bptName)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				for name, expectedContent := range tt.expectedFilesContents {
					pth := path.Join(tmpDir, name)
					assert.FileExists(pth)
					gotContents, err := os.ReadFile(pth)
					assert.NoError(err)
					assert.Equal(expectedContent, string(gotContents))
				}
			}
		})
	}
}

func switchDir(t *testing.T, dir string) func() {
	assert := assert.New(t)
	currDir, err := os.Getwd()
	assert.NoError(err)
	err = os.Chdir(dir)
	assert.NoError(err)
	return func() {
		if err := os.Chdir(currDir); err != nil {
			assert.NoError(err)
		}
	}
}
