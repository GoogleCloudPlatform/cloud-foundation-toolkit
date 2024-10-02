package bpmetadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestUIInputFromVariables(t *testing.T) {
	tests := []struct {
		name     string
		coreVars []*BlueprintVariable
		UIinput  *BlueprintUIInput
	}{
		{
			name: "display metadata does not exist",
			coreVars: []*BlueprintVariable{
				{
					Name: "test_var_1",
				},
				{
					Name: "test_var_2",
				},
				{
					Name: "test_var_3",
				},
			},
			UIinput: &BlueprintUIInput{},
		},
		{
			name: "display metadata exists and is in line with core metadata",
			coreVars: []*BlueprintVariable{
				{
					Name: "test_var_1",
				},
				{
					Name: "test_var_2",
				},
				{
					Name: "test_var_3",
				},
			},
			UIinput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{
					"test_var_1": {
						Name: "test_var_1",
					},
					"test_var_2": {
						Name: "test_var_2",
					},
					"test_var_3": {
						Name: "test_var_3",
					},
				},
			},
		},
		{
			name: "display metadata exists and is not in line with core metadata",
			coreVars: []*BlueprintVariable{
				{
					Name: "test_var_1",
				},
				{
					Name: "test_var_2",
				},
				{
					Name: "test_var_4",
				},
			},
			UIinput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{
					"test_var_1": {
						Name: "test_var_1",
					},
					"test_var_2": {
						Name: "test_var_2",
					},
					"test_var_3": {
						Name: "test_var_3",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildUIInputFromVariables(tt.coreVars, tt.UIinput)
			for _, v := range tt.coreVars {
				dispVar := tt.UIinput.Variables[v.Name]
				assert.NotNil(t, dispVar)
				assert.Equal(t, v.Name, dispVar.Name)
			}

			assert.GreaterOrEqual(t, len(tt.UIinput.Variables), len(tt.coreVars))
		})
	}
}

func TestCreateTitleFromName(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
		wantTitle string
	}{
		{
			name:      "name with underscores",
			inputName: "foo_bar_baz",
			wantTitle: "Foo Bar Baz",
		},
		{
			name:      "name with underscores w/ numbers",
			inputName: "foo_bar_baz_01",
			wantTitle: "Foo Bar Baz 01",
		},
		{
			name:      "name w/o underscores",
			inputName: "FooBarBaz",
			wantTitle: "FooBarBaz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createTitleFromName(tt.inputName)
			if got != tt.wantTitle {
				t.Errorf("createTitleFromName() = %v, want %v", got, tt.wantTitle)
			}
		})
	}
}

func TestMergeExistingAltDefaults(t *testing.T) {
	tests := []struct {
		name          string
		newInput      *BlueprintUIInput
		existingInput *BlueprintUIInput
		expectInput   *BlueprintUIInput
	}{
		{
			name: "Merge alt default into UI input",
			newInput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{
					"test_var_1": {
						Name: "test_var_1",
					},
				},
			},
			existingInput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{
					"test_var_1": {
						Name: "test_var_1",
						AltDefaults: []*DisplayVariable_AlternateDefault{
							{
								Type:  0,
								Value: structpb.NewStringValue("alt_default_value"),
							},
						},
					},
				},
			},
			expectInput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{
					"test_var_1": {
						Name: "test_var_1",
						AltDefaults: []*DisplayVariable_AlternateDefault{
							{
								Type:  0,
								Value: structpb.NewStringValue("alt_default_value"),
							},
						},
					},
				},
			},
		},
		{
			name: "No existing input",
			newInput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{
					"test_var_1": {
						Name: "test_var_1",
					},
				},
			},
			existingInput: nil,
			expectInput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{
					"test_var_1": {
						Name: "test_var_1",
					},
				},
			},
		},
		{
			name: "Empty new input",
			newInput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{},
			},
			existingInput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{
					"test_var_1": {
						Name: "test_var_1",
						AltDefaults: []*DisplayVariable_AlternateDefault{
							{
								Type:  0,
								Value: structpb.NewStringValue("alt_default_value"),
							},
						},
					},
				},
			},
			expectInput: &BlueprintUIInput{
				Variables: map[string]*DisplayVariable{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeExistingAltDefaults(tt.newInput, tt.existingInput)
			assert.Equal(t, tt.newInput, tt.expectInput)
		})
	}
}
