package jumpstartsolutions

import (
	"fmt"
	"os"
	"path"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	metadataFileName = "metadata.yaml"
)

var jssConsumptionFlags struct {
	bpPath string
}

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVarP(&jssConsumptionFlags.bpPath, "path", "p", ".", "path to blueprint for metadata consumption")
}

var Cmd = &cobra.Command{
	Use:   "jump-start-solutions",
	Short: "Generates blueprint metadata for jump start solutions",
	Long:  `Generates blueprint metadata for jump start solutions`,
	Args:  cobra.NoArgs,
	RunE:  generate,
}

// The top-level command function that consumes blueprint metadata and
// generates metadata for jump start solutions.
func generate(cmd *cobra.Command, args []string) error {
	wdPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working dir: %w", err)
	}

	bpPath := jssConsumptionFlags.bpPath
	if !path.IsAbs(bpPath) {
		bpPath = path.Join(wdPath, bpPath)
	}

	err = consumeMetadata(bpPath)
	if err != nil {
		return err
	}

	return nil
}

// consumeMetadata reads the metadata.yaml from the provided path and
// generates textproto and soy files.
func consumeMetadata(bpPath string) error {
	bpObj, err := bpmetadata.UnmarshalMetadata(bpPath, metadataFileName)
	if err != nil {
		return err
	}

	err = generateSoyFile(bpObj)
	if err != nil {
		return err
	}

	err = generateTextprotoFile(bpObj)
	if err != nil {
		return err
	}

	return nil
}

// generateSoyFile consumes the blueprint metadata object to generate soy file.
func generateSoyFile(bpObj *bpmetadata.BlueprintMetadata) error {
	// TODO
	return nil
}

// generateTextprotoFile consumes the blueprint metadata object to
// generate the textproto file.
func generateTextprotoFile(bpObj *bpmetadata.BlueprintMetadata) error {
	// TODO
	return nil
}
