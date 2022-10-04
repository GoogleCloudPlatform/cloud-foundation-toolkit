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
	readmeFileName     = "README.md"
	tfVersionsFileName = "versions.tf"
	tfRolesFileName    = "test/setup/iam.tf"
	tfServicesFileName = "test/setup/main.tf"
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
	readmeContent, err := os.ReadFile(path.Join(bpPath, readmeFileName))
	if err != nil {
		return fmt.Errorf("error reading blueprint readme markdown: %w", err)
	}

	info, err := createInfo(bpPath, readmeContent)
	if err != nil {
		return fmt.Errorf("error creating blueprint info: %w", err)
	}

	interfaces, err := getBlueprintInterfaces(bpPath)
	if err != nil {
		return fmt.Errorf("error creating blueprint interfaces: %w", err)
	}

	rolesCfgPath := path.Join(repoDetails.Source.RootPath, tfRolesFileName)
	svcsCfgPath := path.Join(repoDetails.Source.RootPath, tfServicesFileName)
	requirements, err := getBlueprintRequirements(rolesCfgPath, svcsCfgPath)
	if err != nil {
		return fmt.Errorf("error creating blueprint requirements: %w", err)
	}

	content := createContent(bpPath, readmeContent)

	bpMetadataObj.Spec = &BlueprintMetadataSpec{
		Info:         *info,
		Content:      content,
		Interfaces:   *interfaces,
		Requirements: *requirements,
	}

	return nil
}

func createInfo(bpPath string, readmeContent []byte) (*BlueprintInfo, error) {
	i := &BlueprintInfo{}
	title, err := getMdContent(readmeContent, 1, 1, "", false)
	if err != nil {
		return nil, err
	}

	i.Title = title.literal

	repoDetails, err := getRepoDetailsByPath(bpPath)
	if err != nil {
		return nil, err
	}

	i.Source = &BlueprintRepoDetail{
		Repo:       repoDetails.Source.Path,
		SourceType: "git",
	}

	versionInfo, err := getBlueprintVersion(path.Join(bpPath, tfVersionsFileName))
	if err != nil {
		return nil, err
	}

	i.Version = versionInfo.moduleVersion

	// create descriptions
	i.Description = &BlueprintDescription{}
	tagline, err := getMdContent(readmeContent, -1, -1, "Tagline", true)
	if err == nil {
		i.Description.Tagline = tagline.literal
	}

	detailed, err := getMdContent(readmeContent, -1, -1, "Detailed", true)
	if err == nil {
		i.Description.Detailed = detailed.literal
	}

	preDeploy, err := getMdContent(readmeContent, -1, -1, "PreDeploy", true)
	if err == nil {
		i.Description.PreDeploy = preDeploy.literal
	}

	// TODO:(b/246603410) create icon

	return i, nil
}

func createContent(bpPath string, readmeContent []byte) BlueprintContent {
	var docListToSet []BlueprintListContent
	documentation, err := getMdContent(readmeContent, -1, -1, "Documentation", true)
	if err == nil {
		for _, li := range documentation.listItems {
			doc := BlueprintListContent{
				Title: li.text,
				Url:   li.url,
			}

			docListToSet = append(docListToSet, doc)
		}
	}

	// TODO:(b/246603410) create sub-blueprints

	// TODO:(b/246603410) create examples

	return BlueprintContent{
		Documentation: docListToSet,
	}
}
