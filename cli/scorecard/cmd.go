package scorecard

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	policyPath       string
	targetProjectID	 string
	targetFolderID	 string
	targetOrgID		 string
	bucketName       string
	dirPath          string
	stdin            bool
	refresh          bool
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

	Cmd.Flags().StringVar(&flags.bucketName, "bucket", "", "GCS bucket name for storing inventory (conflicts with --dir-path or --stdin)")
	Cmd.Flags().StringVar(&flags.dirPath, "dir-path", "", "Local directory path for storing inventory (conflicts with --bucket or --stdin)")
	Cmd.Flags().BoolVar(&flags.stdin, "stdin", false, "Passed Cloud Asset Inventory json string as standard input (conflicts with --dir-path or --bucket)")
	Cmd.Flags().BoolVar(&flags.refresh, "refresh", false, "Refresh Cloud Asset Inventory export files in GCS bucket. If set, Application Default Credentials must be a service account (Works with --bucket)")
	Cmd.Flags().StringVar(&flags.targetProjectID, "target-project", "", "Project ID to analyze (Works with --bucket and --refresh; conflicts with --target-folder or --target--organization)")
	Cmd.Flags().StringVar(&flags.targetFolderID, "target-folder", "", "Folder ID to analyze (Works with --bucket and --refresh; conflicts with --target-project or --target--organization)")
	Cmd.Flags().StringVar(&flags.targetOrgID, "target-organization", "", "Organization ID to analyze (Works with --bucket and --refresh; conflicts with --target-project or --target--folder)")
}

// Cmd represents the base scorecard command
var Cmd = &cobra.Command{
	Use:   "scorecard",
	Short: "Print a scorecard of your GCP environment",
	Long: `Print a scorecard of your GCP environment, for resources and IAM policies in Cloud Asset Inventory (CAI) exports, and constraints and constraint templates from Config Validator policy library.

	Read from a bucket:
		  cft scorecard --policy-path <path-to>/policy-library \
			  --bucket <name-of-bucket-containing-cai-export>

	Read from a local directory:
		  cft scorecard --policy-path <path-to>/policy-library \
			  --dir-path <path-to-directory-containing-cai-export>

	Read from standard input:
		  cft scorecard --policy-path <path-to>/policy-library \
			  --stdin

	As of now, CAI export file names need to be resource_inventory.json and iam_inventory.json

	`,
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if (flags.bucketName == "" && flags.dirPath == "" && !flags.stdin) ||
			(flags.bucketName != "" && flags.stdin) ||
			(flags.bucketName != "" && flags.dirPath != "") ||
			(flags.dirPath != "" && flags.stdin) {
			return fmt.Errorf("One and only one of bucket, dir-path, or stdin should be set")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Generating CFT scorecard")
		var err error
		ctx := context.Background()

		targetProjectID := flags.targetProjectID
		if (targetProjectID == "" && flags.targetFolderID == "" && flags.targetOrgID == ""){
			targetProjectID = viper.GetString("google_project")
		}
		if (flags.bucketName != "" && flags.refresh){
			if  (targetProjectID == "" && flags.targetFolderID == "" && flags.targetOrgID == "") ||
				(targetProjectID != "" && flags.targetFolderID != "") ||
				(targetProjectID != "" && flags.targetOrgID != "") ||
				(flags.targetFolderID != "" && flags.targetOrgID != "") {
				return fmt.Errorf("When using --refresh and --bucket, one and only one of target-project, target-folder, or target-org should be set")
			}
		}
		inventory, err := NewInventory(flags.bucketName, flags.dirPath, flags.stdin, flags.refresh,
			TargetProject(targetProjectID), TargetFolder(flags.targetFolderID), TargetOrg(flags.targetOrgID))
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
