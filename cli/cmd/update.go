package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
)

func init() {
	updateCmd.PersistentFlags().StringVarP(&projectFlag, "project", "p", "", "project name")
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update deployment(s)",
	Long:  `Update deployment(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		execute(deployment.ActionUpdate, cmd, args)
	},
}
