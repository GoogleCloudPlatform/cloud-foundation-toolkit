package converter

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

func TestGetComputeFirewallCaiObject(t *testing.T) {
	cases := []struct {
		name   string
		input  cai.Resource
		output cai.Asset
	}{
		{
			name: "v1 compute firewall",
			input: cai.Resource{
				Name:    "resource",
				Type:    "gcp-types/compute-v1:firewalls",
				Project: "project",
				Properties: map[string]interface{}{
					"name": "resource",
				},
			},
			output: cai.Asset{
				Name: "//compute.googleapis.com/projects/project/global/firewalls/resource",
				Type: "compute.googleapis.com/Firewall",
				Resource: &cai.AssetResource{
					Version:              "v1",
					DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
					DiscoveryName:        "Firewall",
					Data: map[string]interface{}{
						"name": "resource",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := GetComputeFirewallCaiObject(c.input.Type, c.input)

			if err != nil {
				t.Errorf("got error: %t", err)
			}
			if !reflect.DeepEqual(res, c.output) {
				t.Errorf("got %v, expected %v", res, c.output)
			}
		})
	}
}
