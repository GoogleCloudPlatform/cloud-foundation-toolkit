package google

import converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"

// TODO: Lift merge into magic modules (third_party/validator/project.go).
func mergeProject(existing, incoming converter.Asset) converter.Asset {
	existing.Resource = incoming.Resource
	return existing
}
