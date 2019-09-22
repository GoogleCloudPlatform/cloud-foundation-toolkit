package scorecard

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	policyPath       string
	targetProjectID  string
	controlProjectID string
	inputPath        string
	inputLocal       bool
}

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVar(&flags.policyPath, "policy-path", "", "Path to directory containing validation policies")
	Cmd.MarkFlagRequired("policy-path")

	Cmd.Flags().StringVar(&flags.targetProjectID, "project", "", "Project to analyze (conflicts with --organization)")

	Cmd.Flags().StringVar(&flags.inputPath, "input-path", "", "GCS bucket name (by default) OR local directory path (with --input-local option), for storing inventory")
	Cmd.MarkFlagRequired("input-path")

	Cmd.Flags().StringVar(&flags.controlProjectID, "control-project", "", "Control project to use for API calls")
	viper.BindPFlag("google_project", Cmd.Flags().Lookup("control-project"))

	Cmd.Flags().BoolVar(&flags.inputLocal, "input-local", false, "Takes inventory input from local directory")

}

// Cmd represents the base scorecard command
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
			flags.inputPath, flags.inputLocal,
			TargetProject(flags.targetProjectID))
		if err != nil {
			return err
		}

		config, err := NewScoringConfig(flags.policyPath)
		if err != nil {
			return err
		}
		err = inventory.Score(config)
		if err != nil {
			return err
		}

		return nil
	},
}
