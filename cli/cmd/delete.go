package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
)

func init() {
	deleteCdm.PersistentFlags().StringVarP(&projectFlag, "project", "p", "", "project name")
	rootCmd.AddCommand(deleteCdm)
}

var deleteCdm = &cobra.Command{
	Use:   "create",
	Short: "Create deployment(s)",
	Long:  `Create deployment(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		execute(deployment.ActionDelete, cmd, args)
	},
}
