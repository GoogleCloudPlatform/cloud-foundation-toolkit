package validation

import (
	"context"
	"errors"
	"log"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

var forsetiValidator = validateAssets
var getDeploymentDescription = deployment.GetDeploymentDescription

func ValidateDeployment(name string, policyPath string, project string) (validated bool, err error) {
	mappers := mappers()

	if policyPath == "" {
		return false, errors.New("Policy path must be specified")
	}

	log.Printf("validating deployment %v, project %s\n", name, project)

	desc, err := getDeploymentDescription(name, project)
	if err != nil {
		log.Printf("error fetching deployment: %s, error: %v", name, err)
		return false, err
	}

	var assets []cai.Asset
	for _, data := range desc.Resources {
		res, err := parseResourceProperties(project, data)
		if err != nil {
			//noinspection GoNilness
			log.Printf("error parsing resource: %s, deployment: %s, error: %v", res.Name, name, err)
			return false, err
		}

		if mapper, ok := mappers[res.Type]; ok {
			asset, err := mapper(res.Type, res)
			if err != nil {
				log.Printf("error getting resource: %s, deployment: %s, error: %v", res.Name, name, err)
				return false, err
			}
			asset.Ancestry, err = getAncestry(project)
			if err != nil {
				log.Printf("error getting resource ancestry: %s, deployment: %s, error: %v", res.Name, name, err)
				return false, err
			}

			assets = append(assets, asset)
		} else {
			log.Printf("resource type not supported: %s, name %s, deployment: %s, project %s\n", res.Type, res.Name, name, project)
		}
	}

	ctx := context.Background()
	auditResult, err := forsetiValidator(ctx, assets, policyPath)

	if err != nil {
		log.Printf("error validating deployment: %s, error: %v", name, err)
		return false, err
	}

	if len(auditResult.Violations) > 0 {
		log.Print("Found Violations:\n\n")
		for _, v := range auditResult.Violations {
			log.Printf("Constraint %v on resource %v: %v\n\n",
				v.Constraint,
				v.Resource,
				v.Message,
			)
		}
		return false, nil
	}

	log.Println("No violations found.")
	return true, nil
}
