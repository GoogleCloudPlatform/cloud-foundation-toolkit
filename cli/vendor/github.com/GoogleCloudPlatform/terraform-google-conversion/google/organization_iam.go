package google

import "fmt"

func GetOrganizationIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newOrganizationIamAsset(d, config, expandIamPolicyBindings)
}

func GetOrganizationIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newOrganizationIamAsset(d, config, expandIamRoleBindings)
}

func GetOrganizationIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newOrganizationIamAsset(d, config, expandIamMemberBindings)
}

func MergeOrganizationIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeOrganizationIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeOrganizationIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newOrganizationIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//cloudresourcemanager.googleapis.com/organizations/{{org_id}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "cloudresourcemanager.googleapis.com/Organization",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
