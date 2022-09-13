package bpmetadata

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	RunE:  generateBlueprintMetadata,
}

// The top-level command function that generates metadata based on the provided flags
func generateBlueprintMetadata(cmd *cobra.Command, args []string) error {
	var bpMetadataDetailObj = &BpMetadataDetail{}

	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working dir: %w", err)
	}

	if mdFlags.path != "." {
		path = path + "/" + mdFlags.path
	}

	// verfiy that the blueprint path is valid & get repo details
	if repoDetails, err := getRepoDetailsByPath(path); err != nil {
		return err
	} else {
		bpMetadataDetailObj.Name = repoDetails.Name
		bpMetadataDetailObj.Source.Path = repoDetails.Source.Path
		bpMetadataDetailObj.Source.SourceType = repoDetails.Source.SourceType
	}

	// TODO: generate metadata details

	// TODO: write metadata to metadata.yaml

	return nil
}
