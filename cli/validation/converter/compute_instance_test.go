package converter

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

func TestGetComputeInstanceCaiObject(t *testing.T) {
	cases := []struct {
		name   string
		input  cai.Resource
		output cai.Asset
	}{
		{
			name: "v1 compute instance",
			input: cai.Resource{
				Name:    "resource",
				Type:    "gcp-types/compute-v1:instances",
				Project: "project",
				Properties: map[string]interface{}{
					"name": "resource",
					"zone": "zone",
				},
			},
			output: cai.Asset{
				Name: "//compute.googleapis.com/projects/project/zones/zone/instances/resource",
				Type: "compute.googleapis.com/Instance",
				Resource: &cai.AssetResource{
					Version:              "v1",
					DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
					DiscoveryName:        "Instance",
					Data: map[string]interface{}{
						"name": "resource",
						"zone": "zone",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := GetComputeInstanceCaiObject(c.input.Type, c.input)

			if err != nil {
				t.Errorf("got error: %t", err)
			}
			if !reflect.DeepEqual(res, c.output) {
				t.Errorf("got %v, expected %v", res, c.output)
			}
		})
	}
}
