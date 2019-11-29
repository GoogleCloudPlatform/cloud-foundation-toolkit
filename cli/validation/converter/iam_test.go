package converter

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

func TestGetIamCaiObject(t *testing.T) {
	cases := []struct {
		name   string
		inputs []cai.Resource
		output cai.Asset
	}{
		{
			name: "organization iam binding",
			inputs: []cai.Resource{
				{
					Name:    "resource",
					Type:    "gcp-types/cloudresourcemanager-v1:virtual.organizations.iamMemberBinding",
					Project: "project",
					Properties: map[string]interface{}{
						"resource": "organizations/123456",
						"role":     "testRole1",
						"member":   "member1",
					},
				},
				{
					Name:    "resource",
					Type:    "gcp-types/cloudresourcemanager-v1:virtual.organizations.iamMemberBinding",
					Project: "project",
					Properties: map[string]interface{}{
						"resource": "organizations/123456",
						"role":     "testRole2",
						"member":   "member2",
					},
				},
				{
					Name:    "resource",
					Type:    "gcp-types/cloudresourcemanager-v1:virtual.organizations.iamMemberBinding",
					Project: "project",
					Properties: map[string]interface{}{
						"resource": "organizations/123456",
						"role":     "testRole2",
						"member":   "member1",
					},
				},
			},
			output: cai.Asset{
				Name: "//cloudresourcemanager.googleapis.com/organizations/123456",
				Type: "cloudresourcemanager.googleapis.com/Organization",
				IAMPolicy: &cai.IAMPolicy{
					Bindings: []cai.IAMBinding{
						{
							Role:    "testRole1",
							Members: []string{"member1"},
						},
						{
							Role:    "testRole2",
							Members: []string{"member1", "member2"},
						},
					},
				},
			},
		},
		{
			name: "folder iam binding",
			inputs: []cai.Resource{
				{
					Name:    "resource",
					Type:    "gcp-types/cloudresourcemanager-v2:virtual.folders.iamMemberBinding",
					Project: "project",
					Properties: map[string]interface{}{
						"resource": "folders/123456",
						"role":     "testRole",
						"member":   "member",
					},
				},
			},
			output: cai.Asset{
				Name: "//cloudresourcemanager.googleapis.com/folders/123456",
				Type: "cloudresourcemanager.googleapis.com/Folder",
				IAMPolicy: &cai.IAMPolicy{
					Bindings: []cai.IAMBinding{
						{
							Role:    "testRole",
							Members: []string{"member"},
						},
					},
				},
			},
		},
		{
			name: "project iam binding",
			inputs: []cai.Resource{
				{
					Name:    "resource",
					Type:    "gcp-types/cloudresourcemanager-v1:virtual.projects.iamMemberBinding",
					Project: "project",
					Properties: map[string]interface{}{
						"resource": "123456",
						"role":     "testRole",
						"member":   "member2",
					},
				},
				{
					Name:    "resource",
					Type:    "gcp-types/cloudresourcemanager-v1:virtual.projects.iamMemberBinding",
					Project: "project",
					Properties: map[string]interface{}{
						"resource": "123456",
						"role":     "testRole",
						"member":   "member1",
					},
				},
				{
					Name:    "resource",
					Type:    "gcp-types/cloudresourcemanager-v1:virtual.projects.iamMemberBinding",
					Project: "project",
					Properties: map[string]interface{}{
						"resource": "123456",
						"role":     "testRole",
						"member":   "member3",
					},
				},
			},
			output: cai.Asset{
				Name: "//cloudresourcemanager.googleapis.com/projects/123456",
				Type: "cloudresourcemanager.googleapis.com/Project",
				IAMPolicy: &cai.IAMPolicy{
					Bindings: []cai.IAMBinding{
						{
							Role:    "testRole",
							Members: []string{"member1", "member2", "member3"},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var asset cai.Asset
			for _, res := range c.inputs {
				newAsset, err := GetIamCaiObject(res.Type, res)
				if err != nil {
					t.Errorf("got error: %t", err)
				}
				if asset.Type != "" {
					asset, err = MergeIamCaiObject(asset, newAsset)
					if err != nil {
						t.Errorf("got error: %t", err)
					}
				} else {
					asset = newAsset
				}
			}

			if !reflect.DeepEqual(asset, c.output) {
				t.Errorf("got %v, expected %v", asset.IAMPolicy, c.output.IAMPolicy)
			}
		})
	}
}

func TestMergeBindings(t *testing.T) {
	cases := []struct {
		name   string
		inputs [][]cai.IAMBinding
		output []cai.IAMBinding
	}{
		{
			name: "primitive test",
			inputs: [][]cai.IAMBinding{
				{
					{
						Role:    "testRole",
						Members: []string{"member1"},
					},
					{
						Role:    "testRole1",
						Members: []string{"member2"},
					},
				},
			},
			output: []cai.IAMBinding{
				{
					Role:    "testRole",
					Members: []string{"member1"},
				},
				{
					Role:    "testRole1",
					Members: []string{"member2"},
				},
			},
		},
		{
			name: "multiple members",
			inputs: [][]cai.IAMBinding{
				{
					{
						Role:    "testRole",
						Members: []string{"member1"},
					},
				},
				{
					{
						Role:    "testRole",
						Members: []string{"member2"},
					},
				},
			},
			output: []cai.IAMBinding{
				{
					Role:    "testRole",
					Members: []string{"member1", "member2"},
				},
			},
		},
		{
			name: "multiple roles",
			inputs: [][]cai.IAMBinding{
				{
					{
						Role:    "testRole1",
						Members: []string{"member1"},
					},
				},
				{
					{
						Role:    "testRole1",
						Members: []string{"member2"},
					},
					{
						Role:    "testRole2",
						Members: []string{"member1"},
					},
				},
			},
			output: []cai.IAMBinding{
				{
					Role:    "testRole1",
					Members: []string{"member1", "member2"},
				},
				{
					Role:    "testRole2",
					Members: []string{"member1"},
				},
			},
		},
		{
			name: "multiple roles in single bindings list",
			inputs: [][]cai.IAMBinding{
				{
					{
						Role:    "testRole1",
						Members: []string{"member0"},
					},
					{
						Role:    "testRole2",
						Members: []string{"member0"},
					},
				},
				{
					{
						Role:    "testRole1",
						Members: []string{"member1"},
					},
					{
						Role:    "testRole1",
						Members: []string{"member1"},
					},
				},
				{
					{
						Role:    "testRole2",
						Members: []string{"member1"},
					},
					{
						Role:    "testRole2",
						Members: []string{"member1", "member2"},
					},
				},
				{
					{
						Role:    "testRole3",
						Members: []string{"member1"},
					},
					{
						Role:    "testRole3",
						Members: []string{"member1", "member2"},
					},
				},
			},
			output: []cai.IAMBinding{
				{
					Role:    "testRole1",
					Members: []string{"member0", "member1"},
				},
				{
					Role:    "testRole2",
					Members: []string{"member0", "member1", "member2"},
				},
				{
					Role:    "testRole3",
					Members: []string{"member1", "member2"},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var bindingsMerged []cai.IAMBinding
			for _, bindings := range c.inputs {
				if len(bindingsMerged) == 0 {
					bindingsMerged = bindings
				} else {
					bindingsMerged = mergeBindings(bindingsMerged, bindings)
				}
			}

			if !reflect.DeepEqual(bindingsMerged, c.output) {
				t.Errorf("got %v, expected %v", bindingsMerged, c.output)
			}
		})
	}
}
