package bpmetadata

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCmdExecution(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		expectErr   bool
	}{
		{
			description: "execute metadata command with valid inputs",
			args:        []string{"metadata", "--help"},
			expectErr:   false,
		},
		{
			description: "execute metadata command with invalid inputs",
			args:        []string{"metadata", "--invalid-flag"},
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{Use: "cft"}
			cmd.SetArgs(tt.args)

			_, err := cmd.ExecuteC()

			if tt.expectErr {
				assert.Error(t, err, "Command should return an error")
			} else {
				assert.NoError(t, err, "Command should execute without error")
			}
		})
	}
}

func TestCreateBlueprintDisplayMetadata(t *testing.T) {
	tests := []struct {
		description string
		bpPath      string
		bpDisp      *BlueprintMetadata
		bpCore      *BlueprintMetadata
		expectErr   bool
	}{
		{
			description: "create metadata with nil display metadata Spec.UI.Input",
			bpPath:      "/path/to/blueprint",
			bpDisp:      &BlueprintMetadata{},
			bpCore: &BlueprintMetadata{
				ApiVersion: "v1",
				Kind:       "Blueprint",
				Metadata: &ResourceTypeMeta{
					Name: "core-blueprint",
					Labels: map[string]string{
						"env": "core",
					},
				},
				Spec: &BlueprintMetadataSpec{
					Info: &BlueprintInfo{
						Title:            "Core Blueprint",
						Version:          "1.0.0",
						Icon:             "assets/core_icon.png",
						SingleDeployment: false,
					},
					Interfaces: &BlueprintInterface{
						Variables: []*BlueprintVariable{
							{
								Name: "test_var_1",
							},
						},
					},
					Ui: &BlueprintUI{
						Input: nil,
					},
				},
			},
			expectErr: false,
		},
		{
			description: "create metadata with valid input",
			bpPath:      "/path/to/blueprint",
			bpDisp: &BlueprintMetadata{
				Spec: &BlueprintMetadataSpec{
					Ui: &BlueprintUI{
						Input: &BlueprintUIInput{
							Variables: map[string]*DisplayVariable{
								"test_var_1": {
									Name:  "test var 1",
									Title: "This is a test input",
								},
							},
						},
					},
				},
			},
			bpCore: &BlueprintMetadata{
				ApiVersion: "v1",
				Kind:       "Blueprint",
				Metadata: &ResourceTypeMeta{
					Name: "core-blueprint",
					Labels: map[string]string{
						"env": "core",
					},
				},
				Spec: &BlueprintMetadataSpec{
					Info: &BlueprintInfo{
						Title:            "Core Blueprint",
						Version:          "1.0.0",
						Icon:             "assets/core_icon.png",
						SingleDeployment: false,
					},
					Interfaces: &BlueprintInterface{
						Variables: []*BlueprintVariable{
							{
								Name: "test_var_1",
							},
						},
					},
					Ui: &BlueprintUI{
						Input: nil,
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			metadata, err := CreateBlueprintDisplayMetadata(tt.bpPath, tt.bpDisp, tt.bpCore)
			if tt.expectErr {
				assert.Error(t, err, "Function should return an error")
				assert.Nil(t, metadata, "Metadata should be nil when there is an error")
			} else {
				assert.NoError(t, err, "Function should not return an error")
				assert.NotNil(t, metadata, "Metadata should not be nil")
				if tt.bpDisp != nil {
					assert.Equal(t, tt.bpDisp.Metadata.Name, metadata.Metadata.Name, "Metadata name should match the input")
					assert.Equal(t, tt.bpDisp.Spec.Info.Title, metadata.Spec.Info.Title, "Metadata title should match the input")
					assert.Equal(t, tt.bpDisp.Spec.Info.Version, metadata.Spec.Info.Version, "Metadata version should match the input")
				}
			}
		})
	}
}
