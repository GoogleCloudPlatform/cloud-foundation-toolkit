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

	// TODO: write metadata to metadata.yaml

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

	return nil
}
