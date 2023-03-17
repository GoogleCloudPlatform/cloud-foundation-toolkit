package schema

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mdFlags struct {
	output string
}

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVarP(&mdFlags.output, "output", "o", ".", "Output path for the JSON schema.")
}

var Cmd = &cobra.Command{
	Use:   "schema",
	Short: "Generates the JSON schema for blueprint metadata.",
	Args:  cobra.NoArgs,
	RunE:  process,
}

// The top-level command function that processes metadata based on the provided flags
func process(cmd *cobra.Command, args []string) error {
	// get the working directory for the command
	wdPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working dir: %w", err)
	}

	if err := generateSchema(mdFlags.output, wdPath); err != nil {
		return err
	}

	return nil
}
