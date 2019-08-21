package scorecard

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var flags struct {
	policyPath      string
	targetProjectID string
	bucketName      string
}

func init() {
	Cmd.Flags().StringVar(&flags.policyPath, "policy-path", "", "Path to directory containing validation policies")
	Cmd.MarkFlagRequired("policy-path")

	Cmd.Flags().StringVar(&flags.targetProjectID, "project", "", "Project to analyze (conflicts with --organization)")

	Cmd.Flags().StringVar(&flags.bucketName, "bucket", "", "GCS bucket name for storing inventory")
	Cmd.MarkFlagRequired("bucket")
}

// getEnvProjectID finds the implict environment project
func getEnvProjectID() (string, error) {
	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		return project, fmt.Errorf("Please set $GOOGLE_PROJECT environment variable")
	}
	return project, nil
}

var Cmd = &cobra.Command{
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

		inventory, err := NewInventory(controlProjectID,
			flags.bucketName,
			TargetProject(flags.targetProjectID))
		if err != nil {
			return err
		}

		config := &ScoringConfig{
			PolicyPath: flags.policyPath,
		}
		err = ScoreInventory(inventory, config)
		if err != nil {
			return err
		}

		return nil
	},
}
