package bptest

import (
	"os"
	"path"
	"testing"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

const (
	cbTestDataDir   = "testdata/cb"
	kitchenTestData = "testdata/kitchen-tests"
)

func TestGetCFTCmd(t *testing.T) {
	tests := []struct {
		name       string
		kitchenCmd string
		want       string
		errMsg     string
	}{
		{
			name:       "simple",
			kitchenCmd: "kitchen_do create",
			want:       "cft test run all --stage init --verbose",
		},
		{
			name:       "explicit test",
			kitchenCmd: "kitchen_do converge foo",
			want:       "cft test run TestFoo --stage apply --verbose",
		},
		{
			name:       "not kitchen",
			kitchenCmd: "foo verify bar",
			errMsg:     "invalid kitchen command: foo verify bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got, err := getCFTCmd(tt.kitchenCmd)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestTransformBuild(t *testing.T) {
	tests := []struct {
		name   string
		fp     string
		wantFp string
		errMsg string
	}{
		{
			name:   "simple",
			fp:     path.Join(cbTestDataDir, "oldAll.yaml"),
			wantFp: path.Join(cbTestDataDir, "newAll.yaml"),
		},
		{
			name:   "targeted",
			fp:     path.Join(cbTestDataDir, "oldTarget.yaml"),
			wantFp: path.Join(cbTestDataDir, "newTarget.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			b, bf, err := getBuildFromFile(tt.fp)
			assert.NoError(err)
			gotBf, err := transformBuild(b, bf)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				_, wbf, err := getBuildFromFile(tt.wantFp)
				assert.NoError(err)
				assert.Equal(wbf, gotBf)
			}
		})
	}
}

func TestConvertTest(t *testing.T) {
	tests := []struct {
		name                  string
		dir                   string
		expectedFilesContents map[string]string
		errMsg                string
	}{
		{
			name: "simple",
			dir:  "simple-example",
			expectedFilesContents: map[string]string{"simple_example_test.go": `package simple_example

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestSimpleExample(t *testing.T) {
	bpt := tft.NewTFBlueprintTest(t)

	bpt.DefineVerify(func(assert *assert.Assertions) {
		bpt.DefaultVerify(assert)
		
		projectId := bpt.GetStringOutput("project_id")
		location := bpt.GetStringOutput("location")
		clusterName := bpt.GetStringOutput("cluster_name")
		masterKubernetesVersion := bpt.GetStringOutput("master_kubernetes_version")
		kubernetesEndpoint := bpt.GetStringOutput("kubernetes_endpoint")
		clientToken := bpt.GetStringOutput("client_token")
		serviceAccount := bpt.GetStringOutput("service_account")
		serviceAccount := bpt.GetStringOutput("service_account")
		databaseEncryptionKeyName := bpt.GetStringOutput("database_encryption_key_name")
		identityNamespace := bpt.GetStringOutput("identity_namespace")

		op := gcloud.Run(t,"")
		assert.Contains(op.Get("result").String(), "foo", "contains foo")
	})

	bpt.Test()
}
`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			tmpDir := path.Join(t.TempDir(), tt.dir)
			err := copy.Copy(path.Join(kitchenTestData, tt.dir), tmpDir)
			assert.NoError(err)
			err = convertTest(tmpDir)
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
