package google

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// expandIamPolicyBindings is used in google_<type>_iam_policy resources.
func expandIamPolicyBindings(d TerraformResourceData) ([]IAMBinding, error) {
	ps := d.Get("policy_data").(string)
	// The policy string is just a marshaled cloudresourcemanager.Policy.
	policy := &cloudresourcemanager.Policy{}
	if err := json.Unmarshal([]byte(ps), policy); err != nil {
		return nil, fmt.Errorf("Could not unmarshal %s:\n: %v", ps, err)
	}

	var bindings []IAMBinding
	for _, b := range policy.Bindings {
		bindings = append(bindings, IAMBinding{
			Role:    b.Role,
			Members: b.Members,
		})
	}

	return bindings, nil
}

// expandIamRoleBindings is used in google_<type>_iam_binding resources.
func expandIamRoleBindings(d TerraformResourceData) ([]IAMBinding, error) {
	var members []string
	for _, m := range d.Get("members").(*schema.Set).List() {
		members = append(members, m.(string))
	}
	return []IAMBinding{
		{
			Role:    d.Get("role").(string),
			Members: members,
		},
	}, nil
}

// expandIamMemberBindings is used in google_<type>_iam_member resources.
func expandIamMemberBindings(d TerraformResourceData) ([]IAMBinding, error) {
	return []IAMBinding{
		{
			Role:    d.Get("role").(string),
			Members: []string{d.Get("member").(string)},
		},
	}, nil
}

// mergeIamAssets merges an existing asset with the IAM bindings of an incoming
// Asset.
func mergeIamAssets(
	existing, incoming Asset,
	mergeBindings func(existing, incoming []IAMBinding) []IAMBinding,
) Asset {
	if existing.IAMPolicy != nil {
		existing.IAMPolicy.Bindings = mergeBindings(existing.IAMPolicy.Bindings, incoming.IAMPolicy.Bindings)
	} else {
		existing.IAMPolicy = incoming.IAMPolicy
	}
	return existing
}

// mergeAdditiveBindings adds members to bindings with the same roles and adds new
// bindings for roles that dont exist.
func mergeAdditiveBindings(existing, incoming []IAMBinding) []IAMBinding {
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
			existing = append(existing, binding)
		}
	}

	// Sort members
	for i := range existing {
		sort.Strings(existing[i].Members)
	}

	return existing
}

// mergeAuthoritativeBindings clobbers members to bindings with the same roles
// and adds new bindings for roles that dont exist.
func mergeAuthoritativeBindings(existing, incoming []IAMBinding) []IAMBinding {
	existingIdxs := make(map[string]int)
	for i, binding := range existing {
		existingIdxs[binding.Role] = i
	}

	for _, binding := range incoming {
		if ei, ok := existingIdxs[binding.Role]; ok {
			existing[ei].Members = binding.Members
		} else {
			existing = append(existing, binding)
		}
	}

	// Sort members
	for i := range existing {
		sort.Strings(existing[i].Members)
	}

	return existing
}
