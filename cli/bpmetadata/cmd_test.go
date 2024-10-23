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

			execCmd := func() error {
				_, err := cmd.ExecuteC()
				return err
			}

			if tt.expectErr {
				assert.Error(t, execCmd(), "Command should return an error")
			} else {
				assert.NoError(t, execCmd(), "Command should execute without error")
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
		expectNil   bool
	}{
		{
			description: "create metadata with valid input",
			bpPath:      "/path/to/blueprint",
			bpDisp: &BlueprintMetadata{
				ApiVersion: "v1",
				Kind:       "Blueprint",
				Metadata: &ResourceTypeMeta{
					Name: "test-blueprint",
					Labels: map[string]string{
						"env": "test",
					},
				},
				Spec: &BlueprintMetadataSpec{
					Info: &BlueprintInfo{
						Title:            "Test Blueprint",
						Version:          "1.0.0",
						Icon:             "assets/icon.png",
						SingleDeployment: true,
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
				},
			},
			expectErr: false,
			expectNil: false,
		},
		{
			description: "create metadata with invalid input",
			bpPath:      "",
			bpDisp: &BlueprintMetadata{
				ApiVersion: "",
				Kind:       "",
				Metadata: &ResourceTypeMeta{
					Name:   "",
					Labels: map[string]string{},
				},
			},
			bpCore: &BlueprintMetadata{
				ApiVersion: "",
				Kind:       "",
				Metadata: &ResourceTypeMeta{
					Name:   "",
					Labels: map[string]string{},
				},
			},
			expectErr: true,
			expectNil: true,
		},
		{
			description: "create metadata with nil Spec.UI.Input",
			bpPath:      "/path/to/blueprint",
			bpDisp:      nil,
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
					Ui: &BlueprintUI{},
				},
			},
			expectErr: false,
			expectNil: true,
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
					assert.Equal(t, tt.bpDisp.Spec.Info.SingleDeployment, metadata.Spec.Info.SingleDeployment, "Single deployment flag should match the input")
				}
			}
		})
	}
}
