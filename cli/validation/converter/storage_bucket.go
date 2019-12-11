package converter

import (
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

// Converts DM storage bucket resource into CAI object
func GetStorageBucketCaiObject(resType string, res cai.Resource) (cai.Asset, error) {
	version := "v1"

	return cai.Asset{
		Name: fmt.Sprintf(
			"//storage.googleapis.com/%s",
			res.Properties["name"],
		),
		Type: "storage.googleapis.com/Bucket",
		Resource: &cai.AssetResource{
			Version:              version,
			DiscoveryDocumentURI: fmt.Sprintf("https://www.googleapis.com/discovery/%s/apis/storage/%s/rest", version, version),
			DiscoveryName:        "Bucket",
			Data:                 res.Properties,
		},
	}, nil
}
