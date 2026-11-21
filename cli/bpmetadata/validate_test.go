package bpmetadata

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

const (
	yamlTestDirPath = "../testdata/bpmetadata/schema"
)

func TestValidateMetadata(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "empty metadata",
			path:    "empty-metadata.yaml",
			wantErr: true,
		},
		{
			name:    "valid metadata",
			path:    "valid-metadata.yaml",
			wantErr: false,
		},
		{
			name:    "valid metadata with connections",
			path:    "valid-metadata-connections.yaml",
			wantErr: false,
		},
		{
			name:    "valid display metadata with alternate defaults",
			path:    "valid-display-metadata-alternate-defaults.yaml",
			wantErr: false,
		},
		{
			name:    "invalid metadata - title missing",
			path:    "invalid-metadata.yaml",
			wantErr: true,
		},
		{
			name: "valid enums for QuotaType",
			path: "valid-metadata-w-enum.yaml",
		},
		{
			name:    "invalid enums for QuotaResourceType",
			path:    "invalid-metadata-w-enum.yaml",
			wantErr: true,
		},
	}

	// load schema from the binary
	s := gojsonschema.NewReferenceLoader("file://schema/gcp-blueprint-metadata.json")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMetadataYaml(path.Join(yamlTestDirPath, tt.path), s)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMetadataYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestValidateDisplayMetadataCompleteness(t *testing.T) {
	tests := []struct {
		name        string
		coreMeta    string
		dispMeta    string
		noCore      bool
		wantErr     bool
		errContains string
	}{
		{
			name: "display metadata has all variables matching core metadata",
			coreMeta: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintMetadata
metadata:
  name: test-module
spec:
  info:
    title: Test Module
  interfaces:
    variables:
      - name: var_one
        description: First variable
      - name: var_two
        description: Second variable
`,
			dispMeta: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintMetadata
metadata:
  name: test-module-display
spec:
  info:
    title: Test Module
  ui:
    input:
      variables:
        var_one:
          name: var_one
          title: Var One
        var_two:
          name: var_two
          title: Var Two
`,
			wantErr: false,
		},
		{
			name: "display metadata is missing a variable from core metadata",
			coreMeta: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintMetadata
metadata:
  name: test-module
spec:
  info:
    title: Test Module
  interfaces:
    variables:
      - name: var_one
        description: First variable
      - name: var_two
        description: Second variable
`,
			dispMeta: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintMetadata
metadata:
  name: test-module-display
spec:
  info:
    title: Test Module
  ui:
    input:
      variables:
        var_one:
          name: var_one
          title: Var One
`,
			wantErr:     true,
			errContains: "var_two",
		},
		{
			name: "display metadata exists without core metadata in directory",
			dispMeta: `apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintMetadata
metadata:
  name: test-module-display
spec:
  info:
    title: Test Module
`,
			noCore:  true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if !tt.noCore {
				err := os.WriteFile(path.Join(dir, metadataFileName), []byte(tt.coreMeta), 0644)
				if err != nil {
					t.Fatalf("failed to write core metadata: %v", err)
				}
			}
			dispPath := path.Join(dir, metadataDisplayFileName)
			err := os.WriteFile(dispPath, []byte(tt.dispMeta), 0644)
			if err != nil {
				t.Fatalf("failed to write display metadata: %v", err)
			}

			err = validateDisplayMetadataCompleteness(dispPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDisplayMetadataCompleteness() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("validateDisplayMetadataCompleteness() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}
