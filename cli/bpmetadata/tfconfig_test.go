package bpmetadata

import (
	"path"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
)

const (
	tfTestdataPath = "../testdata/bpmetadata/tf"
	interfaces     = "sample-module"
)

func TestTFInterfaces(t *testing.T) {
	varTests := []struct {
		name            string
		varName         string
		wantDescription string
		wantVarType     string
		wantDefault     interface{}
		wantRequired    bool
	}{
		{
			name:            "just name and description",
			varName:         "project_id",
			wantDescription: "The project ID to host the cluster in",
			wantRequired:    true,
		},
		{
			name:            "with type and string default",
			varName:         "description",
			wantDescription: "The description of the cluster",
			wantVarType:     "string",
			wantDefault:     "some description",
			wantRequired:    false,
		},
		{
			name:            "with required as fasle",
			varName:         "regional",
			wantDescription: "Whether is a regional cluster",
			wantVarType:     "bool",
			wantDefault:     true,
			wantRequired:    false,
		},
	}

	outTests := []struct {
		name            string
		outName         string
		wantDescription string
	}{
		{
			name:            "just name and description",
			outName:         "cluster_id",
			wantDescription: "Cluster ID",
		},
		{
			name:            "more than just name and description",
			outName:         "endpoint",
			wantDescription: "Cluster endpoint",
		},
	}

	got, _ := getBlueprintInterfaces(path.Join(tfTestdataPath, interfaces))

	for _, tt := range varTests {
		t.Run(tt.name, func(t *testing.T) {
			for _, gotV := range got.Variables {
				if gotV.Name != tt.varName {
					continue
				}

				if gotV.Description != tt.wantDescription {
					t.Errorf("getBlueprintVariable() Description  = %v, want %v", gotV.Description, tt.wantDescription)
				}

				if gotV.VarType != tt.wantVarType {
					t.Errorf("getBlueprintVariable() VarType = %v, want %v", gotV.VarType, tt.wantVarType)
				}

				if gotV.Default != tt.wantDefault {
					t.Errorf("getBlueprintVariable() Default = %v, want %v", gotV.Default, tt.wantDefault)
				}

				if gotV.Required != tt.wantRequired {
					t.Errorf("getBlueprintVariable() Required = %v, want %v", gotV.Required, tt.wantRequired)
				}

				break
			}
		})
	}

	for _, tt := range outTests {
		t.Run(tt.name, func(t *testing.T) {
			for _, gotO := range got.Outputs {
				if gotO.Name != tt.name {
					continue
				}

				if gotO.Description != tt.wantDescription {
					t.Errorf("getBlueprintOutput() Description  = %v, want %v", gotO.Description, tt.wantDescription)
				}

				break
			}
		})
	}
}

func TestTFVersions(t *testing.T) {
	tests := []struct {
		name                string
		configName          string
		wantRequiredVersion string
		wantModuleVersion   string
	}{
		{
			name:                "core version only",
			configName:          "versions-core.tf",
			wantRequiredVersion: ">= 0.13.0",
		},
		{
			name:              "module version only",
			configName:        "versions-module.tf",
			wantModuleVersion: "23.1.0",
		},
		{
			name:                "bad module version good core version",
			configName:          "versions-bad-module.tf",
			wantRequiredVersion: ">= 0.13.0",
			wantModuleVersion:   "",
		},
		{
			name:                "bad core version good module version",
			configName:          "versions-bad-core.tf",
			wantRequiredVersion: "",
			wantModuleVersion:   "23.1.0",
		},
		{
			name:                "all bad",
			configName:          "versions-bad-all.tf",
			wantRequiredVersion: "",
			wantModuleVersion:   "",
		},
		{
			name:                "both versions",
			configName:          "versions.tf",
			wantRequiredVersion: ">= 0.13.0",
			wantModuleVersion:   "23.1.0",
		},
		{
			name:                "both versions with beta",
			configName:          "versions-beta.tf",
			wantRequiredVersion: ">= 0.13.0",
			wantModuleVersion:   "23.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBlueprintVersion(path.Join(tfTestdataPath, tt.configName))

			if got != nil {
				if got.requiredTfVersion != tt.wantRequiredVersion {
					t.Errorf("getBlueprintVersion() = %v, want %v", got.requiredTfVersion, tt.wantRequiredVersion)
					return
				}

				if got.moduleVersion != tt.wantModuleVersion {
					t.Errorf("getBlueprintVersion() = %v, want %v", got.moduleVersion, tt.wantModuleVersion)
					return
				}
			} else {
				if tt.wantModuleVersion != "" && tt.wantRequiredVersion != "" {
					t.Errorf("getBlueprintVersion() = returned nil when we want core: %v and bpVersion: %v", tt.wantRequiredVersion, tt.wantModuleVersion)
				}
			}
		})
	}
}

func TestTFServices(t *testing.T) {
	tests := []struct {
		name         string
		configName   string
		wantServices []string
	}{
		{
			name:       "simple list of apis",
			configName: "main.tf",
			wantServices: []string{
				"cloudkms.googleapis.com",
				"cloudresourcemanager.googleapis.com",
				"container.googleapis.com",
				"pubsub.googleapis.com",
				"serviceusage.googleapis.com",
				"storage-api.googleapis.com",
				"anthos.googleapis.com",
				"anthosconfigmanagement.googleapis.com",
				"logging.googleapis.com",
				"meshca.googleapis.com",
				"meshtelemetry.googleapis.com",
				"meshconfig.googleapis.com",
				"cloudresourcemanager.googleapis.com",
				"monitoring.googleapis.com",
				"stackdriver.googleapis.com",
				"cloudtrace.googleapis.com",
				"meshca.googleapis.com",
				"iamcredentials.googleapis.com",
				"gkeconnect.googleapis.com",
				"privateca.googleapis.com",
				"gkehub.googleapis.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := hclparse.NewParser()
			content, _ := p.ParseHCLFile(path.Join(tfTestdataPath, tt.configName))
			got := parseBlueprintServices(content)

			assert.Equal(t, got, tt.wantServices)
		})
	}
}
