package scorecard

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	policyPath       string
	targetProjectID  string
	controlProjectID string
	bucketName       string
	dirPath          string
	outputPath       string
	outputFormat     string
}

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVar(&flags.policyPath, "policy-path", "", "Path to directory containing validation policies")
	Cmd.MarkFlagRequired("policy-path")

	Cmd.Flags().StringVar(&flags.outputPath, "output-path", "", "Path to directory to contain scorecard outputs. Output to console if not specified")

	Cmd.Flags().StringVar(&flags.outputFormat, "output-format", "", "Format of scorecard outputs, can be txt, json or csv, default is txt")
	viper.SetDefault("output-format", "txt")
	viper.BindPFlag("output-format", Cmd.Flags().Lookup("output-format"))

	//Cmd.Flags().StringVar(&flags.targetProjectID, "project", "", "Project to analyze (conflicts with --organization)")
	Cmd.Flags().StringVar(&flags.bucketName, "bucket", "", "GCS bucket name for storing inventory (conflicts with --dir-path)")
	Cmd.Flags().StringVar(&flags.dirPath, "dir-path", "", "Local directory path for storing inventory (conflicts with --bucket)")
	Cmd.Flags().StringVar(&flags.controlProjectID, "control-project", "", "Control project to use for API calls")
	viper.BindPFlag("google_project", Cmd.Flags().Lookup("control-project"))

}

// Cmd represents the base scorecard command
var Cmd = &cobra.Command{
	Use:   "scorecard",
	Short: "Print a scorecard of your GCP environment",
	Long: `Print a scorecard of your GCP environment, for resources and IAM policies in Cloud Asset Inventory (CAI) exports, and constraints and constraint templates from Config Validator policy library.

	Example:
		  cft scorecard --policy-path <path-to>/policy-library \
			  --bucket <name-of-bucket-containing-cai-export>
	Or:
		  cft scorecard --policy-path <path-to>/policy-library \
			  --dir-path <path-to-directory-containing-cai-export>

	As of now, CAI export file names need to be resource_inventory.json and/or iam_inventory.json

	`,
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if (flags.bucketName == "" && flags.dirPath == "") ||
			(flags.bucketName != "" && flags.dirPath != "") {
			return fmt.Errorf("Either bucket or dir-path should be set")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Generating CFT scorecard")
		var err error
		ctx := context.Background()

		controlProjectID := viper.GetString("google_project")
		if controlProjectID == "" {
			controlProjectID = flags.targetProjectID
			Log.Info("No control project specified, using target project", "project", controlProjectID)
		}

		inventory, err := NewInventory(controlProjectID,
			flags.bucketName, flags.dirPath,
			TargetProject(flags.targetProjectID))
		if err != nil {
			return err
		}

		config, err := NewScoringConfig(ctx, flags.policyPath)
		if err != nil {
			return err
		}
		err = inventory.Score(config, flags.outputPath, viper.GetString("output-format"))
		if err != nil {
			return err
		}

		return nil
	},
}
