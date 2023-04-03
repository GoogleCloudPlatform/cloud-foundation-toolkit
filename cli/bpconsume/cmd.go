package bpconsume

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpconsume/jumpstartsolutions"
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(jumpstartsolutions.Cmd)
}

var Cmd = &cobra.Command{
	Use:   "consume",
	Short: "Consumes blueprint metadata",
	Long:  `Consumes blueprint metadata to generate solution specific metadata`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if args == nil || len(args) == 0 {
			cmd.HelpFunc()(cmd, args)
		}
	},
}
