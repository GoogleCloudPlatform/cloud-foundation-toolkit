package bpmetadata

import (
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"path"
	"testing"
)

const (
	tfTestdataPath = "../testdata/bpmetadata/tf"
	interfaces     = "interfaces"
)

func TestTFVariables(t *testing.T) {
	varTests := []struct {
		name            string
		configDirPath   string
		varName         string
		wantDescription string
		wantVarType     string
		wantDefault     interface{}
		wantRequired    bool
	}{
		{
			name:            "just name and description",
			configDirPath:   path.Join(tfTestdataPath, interfaces),
			varName:         "project_id",
			wantDescription: "The project ID to host the cluster in",
			wantRequired:    true,
		},
		{
			name:            "with type and string default",
			configDirPath:   path.Join(tfTestdataPath, interfaces),
			varName:         "description",
			wantDescription: "The description of the cluster",
			wantVarType:     "string",
			wantDefault:     "some description",
			wantRequired:    false,
		},
		{
			name:            "with required as fasle",
			configDirPath:   path.Join(tfTestdataPath, interfaces),
			varName:         "regional",
			wantDescription: "Whether is a regional cluster",
			wantVarType:     "bool",
			wantDefault:     true,
			wantRequired:    false,
		},
	}

	for _, tt := range varTests {
		t.Run(tt.name, func(t *testing.T) {
			mod, _ := tfconfig.LoadModule(tt.configDirPath)
			got := getBlueprintVariable(mod.Variables[tt.varName])

			if got.Description != tt.wantDescription {
				t.Errorf("getBlueprintVariable() Description  = %v, want %v", got.Description, tt.wantDescription)
			}

			if got.VarType != tt.wantVarType {
				t.Errorf("getBlueprintVariable() VarType = %v, want %v", got.VarType, tt.wantVarType)
			}

			if got.Default != tt.wantDefault {
				t.Errorf("getBlueprintVariable() Default = %v, want %v", got.Default, tt.wantDefault)
			}

			if got.Required != tt.wantRequired {
				t.Errorf("getBlueprintVariable() Required = %v, want %v", got.Required, tt.wantRequired)
			}
		})
	}
}

func TestTFOutputs(t *testing.T) {
	outTests := []struct {
		name            string
		configDirPath   string
		outName         string
		wantDescription string
	}{
		{
			name:            "just name and description",
			configDirPath:   path.Join(tfTestdataPath, interfaces),
			outName:         "cluster_id",
			wantDescription: "Cluster ID",
		},
		{
			name:            "more than just name and description",
			configDirPath:   path.Join(tfTestdataPath, interfaces),
			outName:         "endpoint",
			wantDescription: "Cluster endpoint",
		},
	}

	for _, tt := range outTests {
		t.Run(tt.name, func(t *testing.T) {
			mod, _ := tfconfig.LoadModule(tt.configDirPath)
			got := getBlueprintOutput(mod.Outputs[tt.outName])

			if got.Description != tt.wantDescription {
				t.Errorf("getBlueprintOutput() Description  = %v, want %v", got.Description, tt.wantDescription)
			}
		})
	}
}

func TestTFVersions(t *testing.T) {
	tests := []struct {
		name                string
		configPath          string
		wantRequiredVersion string
		wantModuleVersion   string
	}{
		{
			name:                "core version only",
			configPath:          path.Join(tfTestdataPath, "versions-core.tf"),
			wantRequiredVersion: ">= 0.13.0",
		},
		{
			name:              "module version only",
			configPath:        path.Join(tfTestdataPath, "versions-module.tf"),
			wantModuleVersion: "23.1.0",
		},
		{
			name:                "bad module version good core version",
			configPath:          path.Join(tfTestdataPath, "versions-bad-module.tf"),
			wantRequiredVersion: ">= 0.13.0",
			wantModuleVersion:   "",
		},
		{
			name:                "bad core version good module version",
			configPath:          path.Join(tfTestdataPath, "versions-bad-core.tf"),
			wantRequiredVersion: "",
			wantModuleVersion:   "23.1.0",
		},
		{
			name:                "all bad",
			configPath:          path.Join(tfTestdataPath, "versions-bad-all.tf"),
			wantRequiredVersion: "",
			wantModuleVersion:   "",
		},
		{
			name:                "both versions",
			configPath:          path.Join(tfTestdataPath, "versions.tf"),
			wantRequiredVersion: ">= 0.13.0",
			wantModuleVersion:   "23.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBlueprintVersion(tt.configPath)

			if got != nil {
				if got.requiredVersion != tt.wantRequiredVersion {
					t.Errorf("getBlueprintVersion() = %v, want %v", got.requiredVersion, tt.wantRequiredVersion)
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
		configPath   string
		wantServices []string
	}{
		{
			name:       "simple list of apis",
			configPath: path.Join(tfTestdataPath, "main.tf"),
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
			content, _ := p.ParseHCLFile(tt.configPath)
			got := parseBlueprintServices(content)

			if len(got) != len(tt.wantServices) {
				t.Errorf("parseBlueprintServices() | no of services = %v, want %v", len(got), len(tt.wantServices))
			}

			for i := 0; i < len(got); i++ {
				if got[i] != tt.wantServices[i] {
					t.Errorf("parseBlueprintServices() | service = %s, want %s", got[i], tt.wantServices[i])
				}
			}
		})
	}
}
