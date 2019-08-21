package scorecard

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	policyPath       string
	targetProjectID  string
	controlProjectID string
	bucketName       string
}

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVar(&flags.policyPath, "policy-path", "", "Path to directory containing validation policies")
	Cmd.MarkFlagRequired("policy-path")

	Cmd.Flags().StringVar(&flags.targetProjectID, "project", "", "Project to analyze (conflicts with --organization)")

	Cmd.Flags().StringVar(&flags.bucketName, "bucket", "", "GCS bucket name for storing inventory")
	Cmd.MarkFlagRequired("bucket")

	Cmd.Flags().StringVar(&flags.controlProjectID, "control-project", "", "Control project to use for API calls")
	viper.BindPFlag("google_project", Cmd.Flags().Lookup("control-project"))
}

var Cmd = &cobra.Command{
	Use:   "scorecard",
	Short: "Print a scorecard of your GCP environment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Generating CFT scorecard")
		var err error

		controlProjectID := viper.GetString("google_project")
		if controlProjectID == "" {
			controlProjectID = flags.targetProjectID
			Log.Info("No control project specified, using target project", "project", controlProjectID)
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
