package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
)

func init() {
	initValidateFlags(applyCmd)
	initCommon(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply deployment(s)",
	Long:  `Apply deployment(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		execute(deployment.ActionApply, cmd, args)
	},
}
