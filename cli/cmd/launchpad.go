package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/launchpad"
)

func init() {
	launchpadGenerateCmd.Flags().StringVarP(&launchpadOutputFlavor, "flavor", "f", "tf", "Output Flavor: [tf]/dm")
	launchpadGenerateCmd.Flags().StringVarP(&launchpadOutputDirectory, "directory", "d", "config", "Output Directory: [config]")

	rootCmd.AddCommand(launchpadCmd)
	launchpadCmd.AddCommand(launchpadGenerateCmd)
}

var launchpadOutputFlavor string
var launchpadOutputDirectory string
var launchpadCmd = &cobra.Command{
	Use:     "launchpad",
	Aliases: []string{"lp"},
	Short:   "launchpad (lp)",
	Long: `Cloud Foundation Toolkit Launchpad
	bootstraps foundational GCP infrastructure by following the
	Cloud Foundation Ecosystem Convention. Taking YAML and generate opinionated
	infrastructure resources ready to be deployed in Infrastructure as Code style`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.HelpFunc()(cmd, args)
		}
	},
}

var launchpadGenerateCmd = &cobra.Command{
	Use:     "generate [YAML files]",
	Aliases: []string{"g", "gen"},
	Short:   "generate (g)",
	Long:    `Generate infrastructure foundation via defined YAML`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.HelpFunc()(cmd, args)
		} else {
			launchpad.NewGenerate(args, launchpad.NewOutputFlavor(launchpadOutputFlavor), launchpadOutputDirectory)
		}
	},
}
