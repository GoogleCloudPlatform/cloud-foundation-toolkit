package bpmetadata

import (
	"path"
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
