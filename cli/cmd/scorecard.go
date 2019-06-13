package cmd

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/scorecard"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scorecardCmd)
}

var scorecardCmd = &cobra.Command{
	Use:   "scorecard",
	Short: "Print a scorecard of your GCP envirment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("CFT scorecard")
		inventory := scorecard.NewInventory("gcp-foundation-shared-devops")
		err := scorecard.ExportInventory(inventory)
		if err != nil {
			return errors.Wrap(err, "Error exporting asset inventory")
		}
		return nil
	},
}
