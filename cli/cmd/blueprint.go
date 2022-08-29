package cmd

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpbuild"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bptest"
	"github.com/spf13/cobra"
)

func init() {
	blueprintCmd.AddCommand(bpmetadata.MdCmd)
	blueprintCmd.AddCommand(bpbuild.Cmd)
	blueprintCmd.AddCommand(bptest.Cmd)

	rootCmd.AddCommand(blueprintCmd)
}

var blueprintCmd = &cobra.Command{
	Use:   "blueprint",
	Short: "Blueprint CLI",
	Long:  `The CFT blueprint CLI is used to execute commands specific to blueprints such as test, avgtime & metadata`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if args == nil || len(args) == 0 {
			cmd.HelpFunc()(cmd, args)
		}
	},
}
