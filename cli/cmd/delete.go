package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
)

func init() {
	initCommon(deleteCdm)
}

var deleteCdm = &cobra.Command{
	Use:   "delete",
	Short: "Create deployment(s)",
	Long:  `Create deployment(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		execute(deployment.ActionDelete, cmd, args)
	},
}
