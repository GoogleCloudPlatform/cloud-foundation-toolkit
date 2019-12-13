package converter

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

func TestGetProjectCaiObject(t *testing.T) {
	cases := []struct {
		name   string
		inputs []cai.Resource
		output cai.Asset
	}{
		{
			name: "project test",
			inputs: []cai.Resource{
				{
					Name:    "resource",
					Type:    "gcp-types/cloudresourcemanager-v1:projects",
					Project: "project",
					Properties: map[string]interface{}{
						"project_id": "test_project",
					},
				},
				{
					Name:    "resource",
					Type:    "cloudresourcemanager.v1.project",
					Project: "project",
					Properties: map[string]interface{}{
						"project_id": "test_project",
					},
				},
			},
			output: cai.Asset{
				Name: "//cloudresourcemanager.googleapis.com/projects/test_project",
				Type: "cloudresourcemanager.googleapis.com/Project",
				Resource: &cai.AssetResource{
					Version:              "v1",
					DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
					DiscoveryName:        "Project",
					Data: map[string]interface{}{
						"project_id": "test_project",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var asset cai.Asset
			for _, res := range c.inputs {
				newAsset, err := GetProjectCaiObject(res.Type, res)
				if err != nil {
					t.Errorf("got error: %t", err)
				}
				if asset.Type != "" {
					asset, err = MapProjectCaiObject(asset, newAsset)
					if err != nil {
						t.Errorf("got error: %t", err)
					}
				} else {
					asset = newAsset
				}
			}

			if !reflect.DeepEqual(asset, c.output) {
				t.Errorf("got %v, expected %v", asset, c.output)
			}
		})
	}
}

func TestGetProjectBillingInfoCaiObject(t *testing.T) {
	cases := []struct {
		name   string
		inputs []cai.Resource
		output cai.Asset
	}{
		{
			name: "project test",
			inputs: []cai.Resource{
				{
					Name:    "resource",
					Type:    "gcp-types/cloudresourcemanager-v1:projects",
					Project: "project",
					Properties: map[string]interface{}{
						"project_id": "test_project",
					},
				},
			},
			output: cai.Asset{
				Name: "//cloudbilling.googleapis.com/projects/test_project/billingInfo",
				Type: "cloudbilling.googleapis.com/ProjectBillingInfo",
				Resource: &cai.AssetResource{
					Version:              "v1",
					DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/cloudbilling/v1/rest",
					DiscoveryName:        "ProjectBillingInfo",
					Data: map[string]interface{}{
						"project_id": "test_project",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var asset cai.Asset
			for _, res := range c.inputs {
				newAsset, err := GetProjectBillingInfoCaiObject(res.Type, res)
				if err != nil {
					t.Errorf("got error: %t", err)
				}
				if asset.Type != "" {
					t.Errorf("Merging not supported")
				} else {
					asset = newAsset
				}
			}

			if !reflect.DeepEqual(asset, c.output) {
				t.Errorf("got %v, expected %v", asset, c.output)
			}
		})
	}
}
