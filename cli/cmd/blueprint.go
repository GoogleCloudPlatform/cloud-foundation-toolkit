package cmd

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpbuild"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpcatalog"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bptest"
	"github.com/spf13/cobra"
)

func init() {
	blueprintCmd.AddCommand(bpmetadata.Cmd)
	blueprintCmd.AddCommand(bpbuild.Cmd)
	blueprintCmd.AddCommand(bptest.Cmd)
	blueprintCmd.AddCommand(bpcatalog.Cmd)

	rootCmd.AddCommand(blueprintCmd)
}

var blueprintCmd = &cobra.Command{
	Use:   "blueprint",
	Short: "Blueprint CLI",
	Long:  `The CFT blueprint CLI is used to execute commands specific to blueprints such as test, builds & metadata`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.HelpFunc()(cmd, args)
		}
	},
}
