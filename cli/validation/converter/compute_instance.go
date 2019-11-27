package converter

import (
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

// Converts DM compute resource into CAI object
func GetComputeInstanceCaiObject(resType string, res cai.Resource) (cai.Asset, error) {
	version := ""
	switch resType {
	case "gcp-types/compute-v1:instances":
		version = "v1"
	case "compute.v1.instance":
		version = "v1"
	case "gcp-types/compute-beta:instances":
		version = "beta"
	case "compute.beta.instance":
		version = "beta"
	}

	return cai.Asset{
		Name: fmt.Sprintf(
			"//compute.googleapis.com/projects/%s/zones/%s/instances/%s",
			res.Project,
			res.Properties["zone"],
			res.Properties["name"],
		),
		Type: "compute.googleapis.com/Instance",
		Resource: &cai.AssetResource{
			Version:              version,
			DiscoveryDocumentURI: fmt.Sprintf("https://www.googleapis.com/discovery/%s/apis/compute/%s/rest", version, version),
			DiscoveryName:        "Instance",
			Data:                 res.Properties,
		},
	}, nil
}
