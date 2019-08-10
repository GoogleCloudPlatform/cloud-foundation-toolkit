package cmd

import (
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/scorecard"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scorecardCmd)

	scorecardCmd.Flags().StringVar(&flags.scorecard.policyPath, "policy-path", "", "Path to directory containing validation policies")
	scorecardCmd.MarkFlagRequired("policy-path")

	scorecardCmd.Flags().StringVar(&flags.scorecard.targetProjectID, "project", "", "Project to analyze (conflicts with --organization)")

	scorecardCmd.Flags().StringVar(&flags.scorecard.bucketName, "bucket", "", "GCS bucket name for storing inventory")
	scorecardCmd.MarkFlagRequired("bucket")
}

// getEnvProjectID finds the implict environment project
func getEnvProjectID() (string, error) {
	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		return project, fmt.Errorf("Please set $GOOGLE_PROJECT environment variable")
	}
	return project, nil
}

var scorecardCmd = &cobra.Command{
	Use:   "scorecard",
	Short: "Print a scorecard of your GCP environment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Generating CFT scorecard")
		var err error

		controlProjectID, err := getEnvProjectID()
		if err != nil {
			return err
		}

		inventory, err := scorecard.NewInventory(controlProjectID,
			flags.scorecard.bucketName,
			scorecard.TargetProject(flags.scorecard.targetProjectID))
		if err != nil {
			return err
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
