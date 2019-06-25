package google

import (
	"fmt"
	"strings"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func GetProjectCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	// NOTE: asset.name should use the project number, but we use project_id b/c
	// the number is computed server-side.
	name, err := assetName(d, config, "//cloudresourcemanager.googleapis.com/projects/{{project}}")
	if err != nil {
		return Asset{}, err
	}
	if obj, err := GetProjectApiObject(d, config); err == nil {
		return Asset{
			Name: name,
			Type: "cloudresourcemanager.googleapis.com/Project",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Project",
				Data:                 obj,
			},
		}, nil
	} else {
		return Asset{}, err
	}
}

func GetProjectApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	pid := d.Get("project_id").(string)

	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      d.Get("name").(string),
	}

	if err := getParentResourceId(d, project); err != nil {
		return nil, err
	}

	if _, ok := d.GetOk("labels"); ok {
		project.Labels = expandLabels(d)
	}

	return jsonMap(project)
}

func getParentResourceId(d TerraformResourceData, p *cloudresourcemanager.Project) error {
	orgId := d.Get("org_id").(string)
	folderId := d.Get("folder_id").(string)

	if orgId != "" && folderId != "" {
		return fmt.Errorf("'org_id' and 'folder_id' cannot be both set.")
	}

	if orgId != "" {
		p.Parent = &cloudresourcemanager.ResourceId{
			Id:   orgId,
			Type: "organization",
		}
	}

	if folderId != "" {
		p.Parent = &cloudresourcemanager.ResourceId{
			Id:   strings.TrimPrefix(folderId, "folders/"),
			Type: "folder",
		}
	}

	return nil
}

func GetProjectBillingInfoCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	name, err := assetName(d, config, "//cloudbilling.googleapis.com/projects/{{project}}/billingInfo")
	if err != nil {
		return Asset{}, err
	}
	if obj, err := GetProjectBillingInfoApiObject(d, config); err == nil {
		return Asset{
			Name: name,
			Type: "cloudbilling.googleapis.com/ProjectBillingInfo",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/cloudbilling/v1/rest",
				DiscoveryName:        "ProjectBillingInfo",
				Data:                 obj,
			},
		}, nil
	} else {
		return Asset{}, err
	}
}

func GetProjectBillingInfoApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	if _, ok := d.GetOk("billing_account"); !ok {
		// TODO: If the project already exists, we could ask the API about it's
		// billing info here.
		return nil, ErrNoConversion
	}

	ba := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: fmt.Sprintf("billingAccounts/%s", d.Get("billing_account")),
		Name:               fmt.Sprintf("projects/%s/billingInfo", d.Get("project_id")),
		ProjectId:          d.Get("project_id").(string),
	}

	return jsonMap(ba)
}
