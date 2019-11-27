package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
)

func init() {
	initValidateFlags(createCmd)
	initCommon(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create deployment(s)",
	Long:  `Create deployment(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		execute(deployment.ActionCreate, cmd, args)
	},
}
