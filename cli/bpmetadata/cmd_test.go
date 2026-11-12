package bpmetadata

import (
	"os"
	"path/filepath"
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

func TestGenerateMetadataForBpPathExistingDisplay(t *testing.T) {
	oldDisplay := mdFlags.display
	defer func() { mdFlags.display = oldDisplay }()
	mdFlags.display = false

	oldQuiet := mdFlags.quiet
	defer func() { mdFlags.quiet = oldQuiet }()
	mdFlags.quiet = true

	tmpDir := t.TempDir()

	// Create README.md
	readmeContent := []byte("# Test Module\n\nModule description.\n")
	err := os.WriteFile(filepath.Join(tmpDir, readmeFileName), readmeContent, 0644)
	assert.NoError(t, err)

	// Create variables.tf
	variablesTf := []byte(`variable "existing_var" {
  type        = string
  description = "Existing variable"
}

variable "new_var" {
  type        = string
  description = "New variable"
}
`)
	err = os.WriteFile(filepath.Join(tmpDir, "variables.tf"), variablesTf, 0644)
	assert.NoError(t, err)

	// Create metadata.display.yaml with existing_var having a custom title
	initialDisplayYaml := []byte(`apiVersion: blueprints.cloud.google.com/v1alpha1
kind: BlueprintMetadata
metadata:
  name: test-module-display
spec:
  ui:
    input:
      variables:
        existing_var:
          name: existing_var
          title: Custom Existing Title
`)
	err = os.WriteFile(filepath.Join(tmpDir, metadataDisplayFileName), initialDisplayYaml, 0644)
	assert.NoError(t, err)

	// Call generateMetadataForBpPath with mdFlags.display = false
	err = generateMetadataForBpPath(tmpDir)
	assert.NoError(t, err)

	// Read updated metadata.display.yaml
	updatedDisplay, err := UnmarshalMetadata(tmpDir, metadataDisplayFileName)
	assert.NoError(t, err)
	assert.NotNil(t, updatedDisplay)
	assert.NotNil(t, updatedDisplay.Spec)
	assert.NotNil(t, updatedDisplay.Spec.Ui)
	assert.NotNil(t, updatedDisplay.Spec.Ui.Input)
	assert.NotNil(t, updatedDisplay.Spec.Ui.Input.Variables)

	// Verify existing_var preserved its custom title
	existingVar, ok := updatedDisplay.Spec.Ui.Input.Variables["existing_var"]
	assert.True(t, ok, "existing_var should be present in display metadata")
	assert.Equal(t, "Custom Existing Title", existingVar.Title)

	// Verify new_var was added
	newVar, ok := updatedDisplay.Spec.Ui.Input.Variables["new_var"]
	assert.True(t, ok, "new_var should be automatically added to existing display metadata")
	assert.Equal(t, "New Var", newVar.Title)
}

func TestGenerateMetadataForBpPathNoDisplayWhenFlagFalse(t *testing.T) {
	oldDisplay := mdFlags.display
	defer func() { mdFlags.display = oldDisplay }()
	mdFlags.display = false

	oldQuiet := mdFlags.quiet
	defer func() { mdFlags.quiet = oldQuiet }()
	mdFlags.quiet = true

	tmpDir := t.TempDir()

	readmeContent := []byte("# Test Module\n\nModule description.\n")
	err := os.WriteFile(filepath.Join(tmpDir, readmeFileName), readmeContent, 0644)
	assert.NoError(t, err)

	err = generateMetadataForBpPath(tmpDir)
	assert.NoError(t, err)

	// Verify metadata.display.yaml was NOT created
	_, err = os.Stat(filepath.Join(tmpDir, metadataDisplayFileName))
	assert.True(t, os.IsNotExist(err), "metadata.display.yaml should not be created when flag is false and it did not exist")
}
