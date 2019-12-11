package converter

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

func TestGetStorageBucketCaiObject(t *testing.T) {
	cases := []struct {
		name   string
		input  cai.Resource
		output cai.Asset
	}{
		{
			name: "v1 storage bucket",
			input: cai.Resource{
				Name:    "resource",
				Type:    "gcp-types/storage-v1:buckets",
				Project: "project",
				Properties: map[string]interface{}{
					"name": "resource",
					"foo":  "bar",
				},
			},
			output: cai.Asset{
				Name: "//storage.googleapis.com/resource",
				Type: "storage.googleapis.com/Bucket",
				Resource: &cai.AssetResource{
					Version:              "v1",
					DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/storage/v1/rest",
					DiscoveryName:        "Bucket",
					Data: map[string]interface{}{
						"name": "resource",
						"foo":  "bar",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := GetStorageBucketCaiObject(c.input.Type, c.input)

			if err != nil {
				t.Errorf("got error: %t", err)
			}
			if !reflect.DeepEqual(res, c.output) {
				t.Errorf("got %v, expected %v", res, c.output)
			}
		})
	}
}
