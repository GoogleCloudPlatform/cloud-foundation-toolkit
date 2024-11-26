package bpmetadata

import (
	"fmt"
	"os"
	"path"
	"slices"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	tfTestdataPath       = "../testdata/bpmetadata/tf"
	metadataTestdataPath = "../testdata/bpmetadata/metadata"
	interfaces           = "sample-module"
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
		},
		{
			name:            "with required as false",
			varName:         "regional",
			wantDescription: "Whether is a regional cluster",
			wantVarType:     "bool",
			wantDefault:     true,
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
			i := slices.IndexFunc(got.Variables, func(v *BlueprintVariable) bool { return v.Name == tt.varName })
			if got.Variables[i].Name != tt.varName {
				t.Errorf("getBlueprintInterfaces() - Variable.Name = %v, want %v", got.Variables[i].Name, tt.varName)
				return
			}

			if got.Variables[i].Description != tt.wantDescription {
				t.Errorf("getBlueprintInterfaces() - Variable.Description = %v, want %v", got.Variables[i].Description, tt.wantDescription)
				return
			}

			if got.Variables[i].DefaultValue.AsInterface() != tt.wantDefault {
				t.Errorf("getBlueprintInterfaces() - Variable.DefaultValue = %v, want %v", got.Variables[i].DefaultValue.AsInterface(), tt.wantDefault)
				return
			}

			if got.Variables[i].Required != tt.wantRequired {
				t.Errorf("getBlueprintInterfaces() - Variable.Required = %v, want %v", got.Variables[i].Required, tt.wantRequired)
				return
			}

			if got.Variables[i].VarType != tt.wantVarType {
				t.Errorf("getBlueprintInterfaces() - Variable.VarType = %v, want %v", got.Variables[i].VarType, tt.wantVarType)
				return
			}
		})
	}

	for _, tt := range outTests {
		t.Run(tt.name, func(t *testing.T) {
			i := slices.IndexFunc(got.Outputs, func(o *BlueprintOutput) bool { return o.Name == tt.outName })
			if got.Outputs[i].Name != tt.outName {
				t.Errorf("getBlueprintInterfaces() - Output.Name = %v, want %v", got.Outputs[i].Name, tt.outName)
				return
			}

			if got.Outputs[i].Description != tt.wantDescription {
				t.Errorf("getBlueprintInterfaces() - Output.Description = %v, want %v", got.Outputs[i].Description, tt.wantDescription)
				return
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
		wantRoles  []*BlueprintRoles
	}{
		{
			name:       "simple list of roles",
			configName: "iam.tf",
			wantRoles: []*BlueprintRoles{
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
		{
			name:       "simple list of roles in order for multiple level",
			configName: "iam-multi-level.tf",
			wantRoles: []*BlueprintRoles{
				{
					Level: "Project",
					Roles: []string{
						"roles/owner",
						"roles/storage.admin",
					},
				},
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

func TestSortBlueprintRoles(t *testing.T) {
	tests := []struct {
		name string
		in   []*BlueprintRoles
		want []*BlueprintRoles
	}{
		{
			name: "sort by level",
			in: []*BlueprintRoles{
				{
					Level: "Project",
					Roles: []string{
						"roles/cloudsql.admin",
					},
				},
				{
					Level: "Folder",
					Roles: []string{
						"roles/storage.admin",
					},
				},
			},
			want: []*BlueprintRoles{
				{
					Level: "Folder",
					Roles: []string{
						"roles/storage.admin",
					},
				},
				{
					Level: "Project",
					Roles: []string{
						"roles/cloudsql.admin",
					},
				},
			},
		},
		{
			name: "sort by length of roles",
			in: []*BlueprintRoles{
				{
					Level: "Project",
					Roles: []string{
						"roles/storage.admin",
					},
				},
				{
					Level: "Project",
					Roles: []string{
						"roles/cloudsql.admin",
						"roles/owner",
					},
				},
			},
			want: []*BlueprintRoles{
				{
					Level: "Project",
					Roles: []string{
						"roles/storage.admin",
					},
				},
				{
					Level: "Project",
					Roles: []string{
						"roles/cloudsql.admin",
						"roles/owner",
					},
				},
			},
		},
		{
			name: "sort by first role",
			in: []*BlueprintRoles{
				{
					Level: "Project",
					Roles: []string{
						"roles/storage.admin",
					},
				},
				{
					Level: "Project",
					Roles: []string{
						"roles/cloudsql.admin",
					},
				},
			},
			want: []*BlueprintRoles{
				{
					Level: "Project",
					Roles: []string{
						"roles/cloudsql.admin",
					},
				},
				{
					Level: "Project",
					Roles: []string{
						"roles/storage.admin",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortBlueprintRoles(tt.in)
			assert.Equal(t, tt.in, tt.want)
		})
	}
}

func TestTFProviderVersions(t *testing.T) {
	tests := []struct {
		name                 string
		configName           string
		wantProviderVersions []*ProviderVersion
	}{
		{
			name:       "Simple list of provider versions",
			configName: "versions-beta.tf",
			wantProviderVersions: []*ProviderVersion{
				{
					Source:  "hashicorp/google",
					Version: ">= 4.4.0, < 7",
				},
				{
					Source:  "hashicorp/google-beta",
					Version: ">= 4.4.0, < 7",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := hclparse.NewParser()
			content, _ := p.ParseHCLFile(path.Join(tfTestdataPath, tt.configName))
			got, err := parseBlueprintProviderVersions(content)
			require.NoError(t, err)
			assert.Equal(t, got, tt.wantProviderVersions)
		})
	}
}

func TestMergeExistingConnections(t *testing.T) {
	tests := []struct {
		name                   string
		newInterfacesFile      string
		existingInterfacesFile string
	}{
		{
			name:                   "No existing connections",
			newInterfacesFile:      "new_interfaces_no_connections_metadata.yaml",
			existingInterfacesFile: "existing_interfaces_without_connections_metadata.yaml",
		},
		{
			name:                   "One existing connection is preserved",
			newInterfacesFile:      "new_interfaces_no_connections_metadata.yaml",
			existingInterfacesFile: "existing_interfaces_with_one_connection_metadata.yaml",
		},
		{
			name:                   "Multiple existing connections are preserved",
			newInterfacesFile:      "new_interfaces_no_connections_metadata.yaml",
			existingInterfacesFile: "existing_interfaces_with_some_connections_metadata.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load new interfaces from file
			newInterfaces, err := UnmarshalMetadata(metadataTestdataPath, tt.newInterfacesFile)
			require.NoError(t, err)

			// Load existing interfaces from file
			existingInterfaces, err := UnmarshalMetadata(metadataTestdataPath, tt.existingInterfacesFile)
			require.NoError(t, err)

			// Perform the merge
			mergeExistingConnections(newInterfaces.Spec.Interfaces, existingInterfaces.Spec.Interfaces)

			// Assert that the merged interfaces match the existing ones
			assert.Equal(t, existingInterfaces.Spec.Interfaces, newInterfaces.Spec.Interfaces)
		})
	}
}

func TestMergeExistingOutputTypes(t *testing.T) {
	tests := []struct {
		name                   string
		newInterfacesFile      string
		existingInterfacesFile string
		expectedInterfacesFile string
	}{
		{
			name:                   "No existing types",
			newInterfacesFile:      "interfaces_without_output_types_metadata.yaml",
			existingInterfacesFile: "interfaces_without_output_types_metadata.yaml",
			expectedInterfacesFile: "interfaces_without_output_types_metadata.yaml",
		},
		{
			name:                   "One complex existing type is preserved",
			newInterfacesFile:      "interfaces_without_output_types_metadata.yaml",
			existingInterfacesFile: "interfaces_with_partial_output_types_metadata.yaml",
			expectedInterfacesFile: "interfaces_with_partial_output_types_metadata.yaml",
		},
		{
			name:                   "All existing types (both simple and complex) are preserved",
			newInterfacesFile:      "interfaces_without_output_types_metadata.yaml",
			existingInterfacesFile: "interfaces_with_full_output_types_metadata.yaml",
			expectedInterfacesFile: "interfaces_with_full_output_types_metadata.yaml",
		},
		{
			name:                   "Previous types are not overwriting newly generated types",
			newInterfacesFile:      "interfaces_with_new_output_types_metadata.yaml",
			existingInterfacesFile: "interfaces_with_partial_output_types_metadata.yaml",
			expectedInterfacesFile: "interfaces_with_new_output_types_metadata.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load new interfaces from file
			newInterfaces, err := UnmarshalMetadata(metadataTestdataPath, tt.newInterfacesFile)
			require.NoError(t, err)

			// Load existing interfaces from file
			existingInterfaces, err := UnmarshalMetadata(metadataTestdataPath, tt.existingInterfacesFile)
			require.NoError(t, err)

			// Perform the merge
			mergeExistingOutputTypes(newInterfaces.Spec.Interfaces, existingInterfaces.Spec.Interfaces)

			// Load expected interfaces from file
			expectedInterfaces, err := UnmarshalMetadata(metadataTestdataPath, tt.expectedInterfacesFile)
			require.NoError(t, err)

			// Assert that the merged interfaces match the expected outcome
			assert.Equal(t, expectedInterfaces.Spec.Interfaces, newInterfaces.Spec.Interfaces)
		})
	}
}

func TestTFIncompleteProviderVersions(t *testing.T) {
	tests := []struct {
		name       string
		configName string
	}{
		{
			name:       "Empty list of provider versions",
			configName: "provider-versions-empty.tf",
		},
		{
			name:       "Missing ProviderVersion field",
			configName: "provider-versions-bad.tf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := hclparse.NewParser()
			content, _ := p.ParseHCLFile(path.Join(tfTestdataPath, tt.configName))
			got, err := parseBlueprintProviderVersions(content)
			require.NoError(t, err)
			assert.Nil(t, got)
		})
	}
}

func TestTFVariableSortOrder(t *testing.T) {
	tests := []struct {
		name         string
		configPath   string
		expectOrders map[string]int
		expectError  bool
	}{
		{
			name:       "Variable order should match tf input",
			configPath: "sample-module",
			expectOrders: map[string]int{
				"description": 1,
				"project_id":  0,
				"regional":    2,
			},
			expectError: false,
		},
		{
			name:         "Empty variable name should create nil order",
			configPath:   "empty-module",
			expectOrders: map[string]int{},
			expectError:  true,
		},
		{
			name:         "No variable name should create nil order",
			configPath:   "invalid-module",
			expectOrders: map[string]int{},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBlueprintVariableOrders(path.Join(tfTestdataPath, tt.configPath))
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, got, tt.expectOrders)
			}
		})
	}
}

func TestUpdateOutputTypes(t *testing.T) {
	tests := []struct {
		name            string
		bpPath          string
		interfacesFile  string
		stateFile       string
		expectedOutputs []*BlueprintOutput
		expectError     bool
	}{
		{
			name:           "Update output types from state",
			bpPath:         "sample-module",
			interfacesFile: "interfaces_without_output_types_metadata.yaml",
			stateFile:      "terraform.tfstate",
			expectedOutputs: []*BlueprintOutput{
				{
					Name:        "cluster_id",
					Description: "Cluster ID",
					Type:        structpb.NewStringValue("string"),
				},
				{
					Name:        "endpoint",
					Description: "Cluster endpoint",
					Type: &structpb.Value{
						Kind: &structpb.Value_ListValue{
							ListValue: &structpb.ListValue{
								Values: []*structpb.Value{
									{
										Kind: &structpb.Value_StringValue{
											StringValue: "object",
										},
									},
									{
										Kind: &structpb.Value_StructValue{
											StructValue: &structpb.Struct{
												Fields: map[string]*structpb.Value{
													"host": {
														Kind: &structpb.Value_StringValue{
															StringValue: "string",
														},
													},
													"port": {
														Kind: &structpb.Value_StringValue{
															StringValue: "number",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load interfaces from file
			bpInterfaces, err := UnmarshalMetadata(metadataTestdataPath, tt.interfacesFile)
			require.NoError(t, err)

			// Override with a function that reads a hard-coded tfstate file.
			tfState = func(_ string) ([]byte, error) {
				if tt.expectError {
					return nil, fmt.Errorf("simulated error generating state file")
				}
				// Copy the test state file to the bpPath
				stateFilePath := path.Join(tfTestdataPath, tt.bpPath, tt.stateFile)
				stateData, err := os.ReadFile(stateFilePath)
				if err != nil {
					return nil, fmt.Errorf("error reading state file: %w", err)
				}
				return stateData, nil
			}

			// Update output types
			err = updateOutputTypes(path.Join(tfTestdataPath, tt.bpPath), bpInterfaces.Spec.Interfaces)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Assert that the output types are updated correctly
				expectedOutputsStr := fmt.Sprintf("%v", tt.expectedOutputs)
				actualOutputsStr := fmt.Sprintf("%v", bpInterfaces.Spec.Interfaces.Outputs)
				assert.Equal(t, expectedOutputsStr, actualOutputsStr)
			}
		})
	}
}
