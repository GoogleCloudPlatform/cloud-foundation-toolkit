package scorecard

import (
	"fmt"

	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	policyPath       string
	targetProjectID  string
	controlProjectID string
	bucketName       string
	dirPath          string
	blockProfiler    bool
	cpuProfiler      bool
}

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVar(&flags.policyPath, "policy-path", "", "Path to directory containing validation policies")
	Cmd.MarkFlagRequired("policy-path")

	Cmd.Flags().StringVar(&flags.targetProjectID, "project", "", "Project to analyze (conflicts with --organization)")
	Cmd.Flags().StringVar(&flags.bucketName, "bucket", "", "GCS bucket name for storing inventory (conflicts with --dir-path)")
	Cmd.Flags().StringVar(&flags.dirPath, "dir-path", "", "Local directory path for storing inventory (conflicts with --bucket)")
	Cmd.Flags().StringVar(&flags.controlProjectID, "control-project", "", "Control project to use for API calls")
	viper.BindPFlag("google_project", Cmd.Flags().Lookup("control-project"))

	Cmd.Flags().BoolVar(&flags.blockProfiler, "blockProfiler", false, "run the block profiler.")
	Cmd.Flags().BoolVar(&flags.cpuProfiler, "cpuProfiler", false, "run the cpu profiler.")
}

// Cmd represents the base scorecard command
var Cmd = &cobra.Command{
	Use:   "scorecard",
	Short: "Print a scorecard of your GCP environment",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if (flags.bucketName == "" && flags.dirPath == "") ||
			(flags.bucketName != "" && flags.dirPath != "") {
			return fmt.Errorf("Either bucket or dir-path should be set")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var profilerOpts []func(*profile.Profile)
		if flags.blockProfiler {
			profilerOpts = append(profilerOpts, profile.BlockProfile)
		}
		if flags.cpuProfiler {
			profilerOpts = append(profilerOpts, profile.CPUProfile)
		}
		if len(profilerOpts) != 0 {
			profilerOpts = append(profilerOpts, profile.ProfilePath("."))
			defer profile.Start(profilerOpts...).Stop()
		}

		cmd.Println("Generating CFT scorecard")
		var err error

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

		stopCh := make(chan struct{})
		defer close(stopCh)
		config, err := NewScoringConfig(stopCh, flags.policyPath)
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
