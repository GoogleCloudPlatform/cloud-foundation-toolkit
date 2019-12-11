package validation

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/converter"
)

var runGCloud = deployment.RunGCloud
var ancestryCache = map[string]string{}

type convertFunc func(resType string, res cai.Resource) (cai.Asset, error)

// mergeFunc combines DM resources that have a many-to-one relationship
// with CAI assets, i.e:
// gcp-types/cloudresourcemanager-v1:virtual.projects.iamMemberBinding -> google.cloud.resourcemanager/Project
type mergeFunc func(existing cai.Asset, incoming cai.Asset) (cai.Asset, error)

// mapper pairs related conversion/merging functions.
type mapper struct {
	// convert must be defined.
	convert convertFunc
	// merge may be defined.
	merge mergeFunc
}

// mappers maps DM resource types (i.e. `compute.v1.firewall`) into
// a slice of mapperFuncs.
//
// Modelling of relationships:
// DM resources to CAI assets as []mapperFuncs:
// 1:1 = [mapper{convert: convertAbc}]                  (len=1)
// 1:N = [mapper{convert: convertAbc}, ...]             (len=N)
// N:1 = [mapper{convert: convertAbc, merge: mergeAbc}] (len=1)
func mappers() map[string][]mapper {
	return map[string][]mapper{
		"gcp-types/compute-v1:firewalls":   {{convert: converter.GetComputeFirewallCaiObject}},
		"compute.v1.firewall":              {{convert: converter.GetComputeFirewallCaiObject}},
		"gcp-types/compute-beta:firewalls": {{convert: converter.GetComputeFirewallCaiObject}},
		"compute.beta.firewall":            {{convert: converter.GetComputeFirewallCaiObject}},

		"gcp-types/compute-v1:instances":   {{convert: converter.GetComputeInstanceCaiObject}},
		"compute.v1.instance":              {{convert: converter.GetComputeInstanceCaiObject}},
		"gcp-types/compute-beta:instances": {{convert: converter.GetComputeInstanceCaiObject}},
		"compute.beta.instance":            {{convert: converter.GetComputeInstanceCaiObject}},

		"gcp-types/storage-v1:buckets": {{convert: converter.GetStorageBucketCaiObject}},
		"storage.v1.bucket":            {{convert: converter.GetStorageBucketCaiObject}},

		"gcp-types/cloudresourcemanager-v1:virtual.organizations.iamMemberBinding": {
			{
				convert: converter.GetIamCaiObject,
				merge:   converter.MergeIamCaiObject,
			},
		},
		"gcp-types/cloudresourcemanager-v2:virtual.folders.iamMemberBinding": {
			{
				convert: converter.GetIamCaiObject,
				merge:   converter.MergeIamCaiObject,
			},
		},
		"gcp-types/cloudresourcemanager-v1:virtual.projects.iamMemberBinding": {
			{
				convert: converter.GetIamCaiObject,
				merge:   converter.MergeIamCaiObject,
			},
		},
	}
}

// getAncestry uses the resource manager API to get ancestry paths for
// projects. It implements a cache because many resources share the same
// project.
func getAncestry(project string) (string, error) {
	if path, ok := ancestryCache[project]; ok {
		return path, nil
	}

	args := []string{
		"projects",
		"get-ancestors",
		project,
		"--format", "json",
	}
	data, err := runGCloud(args...)
	if err != nil {
		return "", err
	}

	out := []struct {
		Id   string
		Type string
	}{}
	err = json.Unmarshal([]byte(data), &out)
	if err != nil {
		return "", err
	}

	var paths []string
	for i := len(out) - 1; i >= 0; i-- {
		paths = append(paths, out[i].Type+"/"+out[i].Id)
	}
	path := strings.Join(paths, "/")
	ancestryCache[project] = path

	return path, nil
}

// Parses Deployment description into cai DM Resource. API returns properties as serialized YAML + it can be in one
// of multiple fields.
func parseResourceProperties(project string, res deployment.DeploymentDescriptionResource) (cai.Resource, error) {
	ret := cai.Resource{
		Project: project,
		Name:    res.Name,
		Type:    res.Type,
	}

	props := res.Properties
	if res.FinalProperties != "" {
		props = res.FinalProperties
	}
	if res.Update.Properties != "" {
		props = res.Update.Properties
	}
	if res.Update.FinalProperties != "" {
		props = res.Update.FinalProperties
	}

	err := yaml.Unmarshal([]byte(props), &ret.Properties)
	return ret, err
}
