package cmd

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/scorecard"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scorecardCmd)

	scorecardCmd.PersistentFlags().BoolVar(&flags.scorecard.refresh, "refresh", false, "Refresh CAI inventory")

	scorecardCmd.Flags().StringVar(&flags.scorecard.policyPath, "policy-path", "", "Path to directory containing validation policies")
	scorecardCmd.MarkFlagRequired("policy-path")

	scorecardCmd.Flags().StringVar(&flags.scorecard.targetProjectID, "project", "", "Project to analyze (conflicts with --organization)")

}

var scorecardCmd = &cobra.Command{
	Use:   "scorecard",
	Short: "Print a scorecard of your GCP environment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Generating CFT scorecard")
		var err error

		inventory, err := scorecard.NewInventory("gcp-foundation-shared-devops", scorecard.TargetProject(flags.scorecard.targetProjectID))
		if err != nil {
			return err
		}

		if flags.scorecard.refresh {
			err = scorecard.ExportInventory(inventory)
			if err != nil {
				return errors.Wrap(err, "Error exporting asset inventory")
			}
		}

		config := &scorecard.ScoringConfig{
			PolicyPath: flags.scorecard.policyPath,
		}
		err = scorecard.ScoreInventory(inventory, config)
		if err != nil {
			return err
		}

		return nil
	},
}
