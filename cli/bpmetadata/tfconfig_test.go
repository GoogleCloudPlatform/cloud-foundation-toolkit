package bpmetadata

import (
	"path"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	got, err := getBlueprintInterfaces(path.Join(tfTestdataPath, interfaces))
	require.NoError(t, err)
	for _, tt := range varTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, got.Variables, BlueprintVariable{
				Name:        tt.varName,
				Description: tt.wantDescription,
				Default:     tt.wantDefault,
				Required:    tt.wantRequired,
				VarType:     tt.wantVarType,
			})
		})
	}

	for _, tt := range outTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, got.Outputs, BlueprintOutput{
				Name:        tt.outName,
				Description: tt.wantDescription,
			})
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
			got, _ := getBlueprintVersion(path.Join(tfTestdataPath, tt.configName))

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
			got, err := parseBlueprintServices(content)
			require.NoError(t, err)
			assert.Equal(t, got, tt.wantServices)
		})
	}
}

func TestTFRoles(t *testing.T) {
	tests := []struct {
		name       string
		configName string
		wantRoles  []BlueprintRoles
	}{
		{
			name:       "simple list of roles",
			configName: "iam.tf",
			wantRoles: []BlueprintRoles{
				{
					Level: "Project",
					Roles: []string{
						"roles/cloudsql.admin",
						"roles/compute.networkAdmin",
						"roles/iam.serviceAccountAdmin",
						"roles/resourcemanager.projectIamAdmin",
						"roles/storage.admin",
						"roles/workflows.admin",
						"roles/cloudscheduler.admin",
						"roles/iam.serviceAccountUser",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := hclparse.NewParser()
			content, _ := p.ParseHCLFile(path.Join(tfTestdataPath, tt.configName))
			got, err := parseBlueprintRoles(content)
			require.NoError(t, err)
			assert.Equal(t, got, tt.wantRoles)
		})
	}
}
