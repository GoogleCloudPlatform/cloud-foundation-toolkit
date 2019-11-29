package converter

import (
	"fmt"
	"sort"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

// mergeBindings adds members to bindings with the same roles
// and adds new bindings for roles that dont exist.
func mergeBindings(existing []cai.IAMBinding, incoming []cai.IAMBinding) []cai.IAMBinding {
	existingIdxs := make(map[string]int)
	for i, binding := range existing {
		existingIdxs[binding.Role] = i
	}

	for _, binding := range incoming {
		if ei, ok := existingIdxs[binding.Role]; ok {
			memberExists := make(map[string]bool)
			for _, m := range existing[ei].Members {
				memberExists[m] = true
			}
			for _, m := range binding.Members {
				// Only add members that don't exist.
				if !memberExists[m] {
					existing[ei].Members = append(existing[ei].Members, m)
				}
			}
		} else {
			existingIdxs[binding.Role] = len(existingIdxs)
			existing = append(existing, binding)
		}
	}

	// Sort members
	for i := range existing {
		sort.Strings(existing[i].Members)
	}

	return existing
}

// Converts DM iam resource into CAI object
func GetIamCaiObject(resType string, res cai.Resource) (cai.Asset, error) {
	assetType := ""
	assetName := ""
	switch resType {
	case "gcp-types/cloudresourcemanager-v1:virtual.organizations.iamMemberBinding":
		assetType = "cloudresourcemanager.googleapis.com/Organization"
		assetName = fmt.Sprintf(
			"//cloudresourcemanager.googleapis.com/%s",
			res.Properties["resource"],
		)
	case "gcp-types/cloudresourcemanager-v2:virtual.folders.iamMemberBinding":
		assetType = "cloudresourcemanager.googleapis.com/Folder"
		assetName = fmt.Sprintf(
			"//cloudresourcemanager.googleapis.com/%s",
			res.Properties["resource"],
		)
	case "gcp-types/cloudresourcemanager-v1:virtual.projects.iamMemberBinding":
		assetType = "cloudresourcemanager.googleapis.com/Project"
		assetName = fmt.Sprintf(
			"//cloudresourcemanager.googleapis.com/projects/%s",
			res.Properties["resource"],
		)
	}

	return cai.Asset{
		Name: assetName,
		Type: assetType,
		IAMPolicy: &cai.IAMPolicy{
			Bindings: []cai.IAMBinding{
				{
					Role:    res.Properties["role"].(string),
					Members: []string{res.Properties["member"].(string)},
				},
			},
		},
	}, nil
}

func MergeIamCaiObject(existing cai.Asset, incoming cai.Asset) (cai.Asset, error) {
	if existing.IAMPolicy != nil {
		existing.IAMPolicy.Bindings = mergeBindings(existing.IAMPolicy.Bindings, incoming.IAMPolicy.Bindings)
	} else {
		existing.IAMPolicy = incoming.IAMPolicy
	}
	return existing, nil
}
