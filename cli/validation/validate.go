package validation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

var forsetiValidator = validateAssets
var getDeploymentDescription = deployment.GetDeploymentDescription

type byName []cai.Asset

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

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

	assets := make(map[string]cai.Asset, 0)
	assetList := make(byName, 0)
	for _, data := range desc.Resources {
		res, err := parseResourceProperties(project, data)
		if err != nil {
			//noinspection GoNilness
			log.Printf("error parsing resource: %s, deployment: %s, error: %v", res.Name, name, err)
			return false, err
		}

		if mapperList, ok := mappers[res.Type]; ok {
			for _, mapper := range mapperList {
				asset, err := mapper.convert(res.Type, res)
				if err != nil {
					log.Printf("error getting resource: %s, deployment: %s, error: %v", res.Name, name, err)
					return false, err
				}
				asset.Ancestry, err = getAncestry(project)
				if err != nil {
					log.Printf("error getting resource ancestry: %s, deployment: %s, error: %v", res.Name, name, err)
					return false, err
				}

				key := asset.Type + "__" + asset.Name

				if existing, exists := assets[key]; exists {
					// The existence of a merge function signals that this resource maps to a
					// patching operation on an API resource.
					if mapper.merge != nil {
						asset, err = mapper.merge(existing, asset)
						if err != nil {
							return false, err
						}
					} else {
						log.Printf("duplicate asset: %s, deployment: %s, error: %v", key, name, err)
						return false, errors.New(fmt.Sprintf("duplicate asset: %s", key))
					}
				}

				assets[key] = asset
				assetList = append(assetList, asset)
			}
		} else {
			log.Printf("resource type not supported: %s, name %s, deployment: %s, project %s\n", res.Type, res.Name, name, project)
		}
	}

	ctx := context.Background()
	sort.Sort(assetList)
	auditResult, err := forsetiValidator(ctx, assetList, policyPath)

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
