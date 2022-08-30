package bpmetadata

import (
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
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("command under construction")
	},
}
