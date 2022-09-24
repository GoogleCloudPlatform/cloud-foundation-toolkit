package bpmetadata

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

var mdFlags struct {
	path   string
	nested bool
}

const (
	readmeFileName     string = "README.md"
	tfVersionsFileName string = "versions.tf"
	tfRolesFileName    string = "../../test/setup/iam.tf"
	tfServicesFileName string = "../../test/setup/main.tf"
)

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVar(&mdFlags.path, "path", ".", "Path to the blueprint for generating metadata.")
	Cmd.Flags().BoolVar(&mdFlags.nested, "nested", true, "Flag for generating metadata for nested blueprint, if any.")
}

var Cmd = &cobra.Command{
	Use:   "metadata",
	Short: "Generates blueprint metatda",
	Long:  `Generates metadata.yaml for specified blueprint`,
	Args:  cobra.NoArgs,
	RunE:  generate,
}

// The top-level command function that generates metadata based on the provided flags
func generate(cmd *cobra.Command, args []string) error {
	wdPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working dir: %w", err)
	}

	// create metadata details
	bpPath := path.Join(wdPath, mdFlags.path)
	err = CreateBlueprintMetadata(bpPath)
	if err != nil {
		return fmt.Errorf("error creating metadata for blueprint: %w", err)
	}

	// TODO:(b/248642744) write metadata to metadata.yaml

	return nil
}

func CreateBlueprintMetadata(bpPath string) error {
	// verfiy that the blueprint path is valid & get repo details
	repoDetails, err := getRepoDetailsByPath(bpPath)
	if err != nil {
		return err
	}

	// start creating blueprint metadata
	var bpMetadataObj = &BlueprintMetadata{}
	bpMetadataObj.Meta = yaml.ResourceMeta{
		TypeMeta: yaml.TypeMeta{
			APIVersion: "blueprints.cloud.google.com/v1alpha1",
			Kind:       "BlueprintMetadata",
		},
		ObjectMeta: yaml.ObjectMeta{
			NameMeta: yaml.NameMeta{
				Name:      repoDetails.Name,
				Namespace: "",
			},
			Annotations: map[string]string{"config.kubernetes.io/local-config": "true"},
		},
	}

	// start creating the Spec node
	bpMetadataObj.Spec = &BlueprintMetadataSpec{
		Info:         BlueprintInfo{},
		Content:      BlueprintContent{},
		Interfaces:   BlueprintInterface{},
		Requirements: BlueprintRequirements{},
	}

	//create blueprint title
	readmeContent, err := os.ReadFile(path.Join(bpPath, readmeFileName))
	if err != nil {
		return fmt.Errorf("error reading blueprint readme markdown: %w", err)
	}
	title := getMdContent(readmeContent, 1, 1, "", false)
	bpMetadataObj.Spec.Info.Title = title.literal

	// create blueprint source
	bpMetadataObj.Spec.Info.Source = &BlueprintRepoDetail{
		Repo:       repoDetails.Source.Path,
		SourceType: "git",
	}

	// create version
	versionInfo := getBlueprintVersion(path.Join(bpPath, tfVersionsFileName))
	bpMetadataObj.Spec.Info.Version = versionInfo.moduleVersion

	// create blueprint actuation tool & core version
	bpMetadataObj.Spec.Info.ActuationTool = &BlueprintActuationTool{
		Flavor:  "Terraform",
		Version: versionInfo.requiredVersion,
	}

	// create descriptions
	tagline := getMdContent(readmeContent, -1, -1, "Tagline", true)
	detailed := getMdContent(readmeContent, -1, -1, "Detailed", true)
	preDeploy := getMdContent(readmeContent, -1, -1, "PreDeploy", true)
	bpMetadataObj.Spec.Info.Description = &BlueprintDescription{
		Tagline:   tagline.literal,
		Detailed:  detailed.literal,
		PreDeploy: preDeploy.literal,
	}

	// TODO:(b/246603410) create icon

	// create documentation
	documentation := getMdContent(readmeContent, -1, -1, "Documentation", true)
	var docListToSet []BlueprintMdListContent

	for _, li := range documentation.listItems {
		doc := BlueprintMdListContent{
			Title: li.text,
			Url:   li.url,
		}

		docListToSet = append(docListToSet, doc)
	}

	bpMetadataObj.Spec.Content.Documentation = docListToSet

	// TODO:(b/246603410) create sub-blueprints

	// TODO:(b/246603410) create examples

	// create variables & outputs
	interfaces := getBlueprintInterfaces(bpPath)
	bpMetadataObj.Spec.Interfaces.Variables = interfaces.Variables
	bpMetadataObj.Spec.Interfaces.Outputs = interfaces.Outputs

	// create roles & services
	requirements := getBlueprintRequirements(path.Join(bpPath, tfRolesFileName), path.Join(bpPath, tfServicesFileName))
	bpMetadataObj.Spec.Requirements = requirements

	return nil
}
