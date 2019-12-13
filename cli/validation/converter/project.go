package converter

import (
	"fmt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

func mergeProject(existing cai.Asset, incoming cai.Asset) cai.Asset {
	existing.Resource = incoming.Resource
	return existing
}

// Converts DM iam resource into CAI object
func GetProjectCaiObject(resType string, res cai.Resource) (cai.Asset, error) {
	version := "v1"

	return cai.Asset{
		Name: fmt.Sprintf(
			"//cloudresourcemanager.googleapis.com/projects/%s",
			res.Properties["project_id"],
		),
		Type: "cloudresourcemanager.googleapis.com/Project",
		Resource: &cai.AssetResource{
			Version:              version,
			DiscoveryDocumentURI: fmt.Sprintf("https://www.googleapis.com/discovery/%s/apis/compute/%s/rest", version, version),
			DiscoveryName:        "Project",
			Data:                 res.Properties,
		},
	}, nil
}

func GetProjectBillingInfoCaiObject(resType string, res cai.Resource) (cai.Asset, error) {
	version := "v1"

	return cai.Asset{
		Name: fmt.Sprintf(
			"//cloudbilling.googleapis.com/projects/%s/billingInfo",
			res.Properties["project_id"],
		),
		Type: "cloudbilling.googleapis.com/ProjectBillingInfo",
		Resource: &cai.AssetResource{
			Version:              version,
			DiscoveryDocumentURI: fmt.Sprintf("https://www.googleapis.com/discovery/%s/apis/cloudbilling/%s/rest", version, version),
			DiscoveryName:        "ProjectBillingInfo",
			Data:                 res.Properties,
		},
	}, nil
}

func MapProjectCaiObject(existing cai.Asset, incoming cai.Asset) (cai.Asset, error) {
	existing.Resource = incoming.Resource
	return existing, nil
}
