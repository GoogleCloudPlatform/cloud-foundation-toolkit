// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package google

import (
	"sort"

	converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"
)

type convertFunc func(d converter.TerraformResourceData, config *converter.Config) (converter.Asset, error)

// mergeFunc combines terraform resources that have a many-to-one relationship
// with CAI assets, i.e:
// google_project_iam_member -> google.cloud.resourcemanager/Project
type mergeFunc func(existing, incoming converter.Asset) converter.Asset

// mapper pairs related conversion/merging functions.
type mapper struct {
	// convert must be defined.
	convert convertFunc
	// merge may be defined.
	merge mergeFunc
}

// mappers maps terraform resource types (i.e. `google_project`) into
// a slice of mapperFuncs.
//
// Modelling of relationships:
// terraform resources to CAI assets as []mapperFuncs:
// 1:1 = [mapper{convert: convertAbc}]                  (len=1)
// 1:N = [mapper{convert: convertAbc}, ...]             (len=N)
// N:1 = [mapper{convert: convertAbc, merge: mergeAbc}] (len=1)
func mappers() map[string][]mapper {
	return map[string][]mapper{
		// TODO: Use a generated mapping once it lands in the conversion library.
		"google_compute_firewall":      {{convert: converter.GetComputeFirewallCaiObject}},
		"google_compute_disk":          {{convert: converter.GetComputeDiskCaiObject}},
		"google_compute_instance":      {{convert: converter.GetComputeInstanceCaiObject}},
		"google_storage_bucket":        {{convert: converter.GetStorageBucketCaiObject}},
		"google_sql_database_instance": {{convert: converter.GetSQLDatabaseInstanceCaiObject}},

		// Terraform resources of type "google_project" have a 1:N relationship with CAI assets.
		"google_project": {
			{
				convert: converter.GetProjectCaiObject,
				merge:   mergeProject,
			},
			{convert: converter.GetProjectBillingInfoCaiObject},
		},

		// Terraform IAM policy resources have a N:1 relationship with CAI assets.
		"google_organization_iam_policy": {
			{
				convert: converter.GetOrganizationIamPolicyCaiObject,
				merge:   converter.MergeOrganizationIamPolicy,
			},
		},
		"google_organization_iam_binding": {
			{
				convert: converter.GetOrganizationIamBindingCaiObject,
				merge:   converter.MergeOrganizationIamBinding,
			},
		},
		"google_organization_iam_member": {
			{
				convert: converter.GetOrganizationIamMemberCaiObject,
				merge:   converter.MergeOrganizationIamMember,
			},
		},
		"google_folder_iam_policy": {
			{
				convert: converter.GetFolderIamPolicyCaiObject,
				merge:   converter.MergeFolderIamPolicy,
			},
		},
		"google_folder_iam_binding": {
			{
				convert: converter.GetFolderIamBindingCaiObject,
				merge:   converter.MergeFolderIamBinding,
			},
		},
		"google_folder_iam_member": {
			{
				convert: converter.GetFolderIamMemberCaiObject,
				merge:   converter.MergeFolderIamMember,
			},
		},
		"google_project_iam_policy": {
			{
				convert: converter.GetProjectIamPolicyCaiObject,
				merge:   converter.MergeProjectIamPolicy,
			},
		},
		"google_project_iam_binding": {
			{
				convert: converter.GetProjectIamBindingCaiObject,
				merge:   converter.MergeProjectIamBinding,
			},
		},
		"google_project_iam_member": {
			{
				convert: converter.GetProjectIamMemberCaiObject,
				merge:   converter.MergeProjectIamMember,
			},
		},
	}
}

// SupportedResources returns a sorted list of terraform resource names.
func SupportedTerraformResources() []string {
	fns := mappers()
	list := make([]string, 0, len(fns))
	for k := range fns {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}
