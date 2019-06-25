package google

import "fmt"

func GetProjectIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newProjectIamAsset(d, config, expandIamPolicyBindings)
}

func GetProjectIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newProjectIamAsset(d, config, expandIamRoleBindings)
}

func GetProjectIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newProjectIamAsset(d, config, expandIamMemberBindings)
}

func MergeProjectIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeProjectIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeProjectIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newProjectIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use project_number, but since that is generated server-side,
	// we substitute project_id.
	name, err := assetName(d, config, "//cloudresourcemanager.googleapis.com/projects/{{project}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "cloudresourcemanager.googleapis.com/Project",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
