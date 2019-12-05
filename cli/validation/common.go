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

func mappers() map[string]convertFunc {
	return map[string]convertFunc{
		"gcp-types/compute-v1:firewalls":   converter.GetComputeFirewallCaiObject,
		"compute.v1.firewall":              converter.GetComputeFirewallCaiObject,
		"gcp-types/compute-beta:firewalls": converter.GetComputeFirewallCaiObject,
		"compute.beta.firewall":            converter.GetComputeFirewallCaiObject,
		"gcp-types/compute-v1:instances":   converter.GetComputeInstanceCaiObject,
		"compute.v1.instance":              converter.GetComputeInstanceCaiObject,
		"gcp-types/compute-beta:instances": converter.GetComputeInstanceCaiObject,
		"compute.beta.instance":            converter.GetComputeInstanceCaiObject,
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
