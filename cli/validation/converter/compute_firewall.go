package converter

import (
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

// Converts DM compute firewall resource into CAI object
func GetComputeFirewallCaiObject(resType string, res cai.Resource) (cai.Asset, error) {
	version := ""
	switch resType {
	case "gcp-types/compute-v1:firewalls":
		version = "v1"
	case "compute.v1.firewall":
		version = "v1"
	case "gcp-types/compute-beta:firewalls":
		version = "beta"
	case "compute.beta.firewall":
		version = "beta"
	}

	return cai.Asset{
		Name: fmt.Sprintf(
			"//compute.googleapis.com/projects/%s/global/firewalls/%s",
			res.Project,
			res.Properties["name"],
		),
		Type: "compute.googleapis.com/Firewall",
		Resource: &cai.AssetResource{
			Version:              version,
			DiscoveryDocumentURI: fmt.Sprintf("https://www.googleapis.com/discovery/%s/apis/compute/%s/rest", version, version),
			DiscoveryName:        "Firewall",
			Data:                 res.Properties,
		},
	}, nil
}
