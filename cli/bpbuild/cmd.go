package bpbuild

import (
	"fmt"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var avgTimeFlags struct {
	projectId       string
	repoName        string
	buildFilePath   string
	buildStepID     string
	lookUpStart     string
	lookUpStartTime time.Time
	lookUpEnd       string
	lookUpEndTime   time.Time
}

const defaultBuildFilePath = "build/int.cloudbuild.yaml"

func init() {
	viper.AutomaticEnv()
	Cmd.AddCommand(avgTimeCmd)

	avgTimeCmd.Flags().StringVar(&avgTimeFlags.buildFilePath, "build-file", defaultBuildFilePath, "Path to file containing CloudBuild configs.")
	avgTimeCmd.Flags().StringVar(&avgTimeFlags.buildStepID, "step", "", "ID of build step to compute avg.")
	avgTimeCmd.Flags().StringVar(&avgTimeFlags.lookUpStart, "start-time", "", "Time to start computing build step avg in form MM-DD-YYYY. Defaults to one month ago.")
	avgTimeCmd.Flags().StringVar(&avgTimeFlags.lookUpEnd, "end-time", "", "Time to stop computing build step avg in form MM-DD-YYYY. Defaults to current date.")
	avgTimeCmd.Flags().StringVar(&avgTimeFlags.projectId, "project-id", "cloud-foundation-cicd", "Project ID where builds are executed.")
	avgTimeCmd.Flags().StringVar(&avgTimeFlags.repoName, "repo", "", "Name of repo that triggered the builds. Defaults to extracting from git config.")
}

var Cmd = &cobra.Command{
	Use:   "builds",
	Short: "Blueprint builds",
	Long:  `Blueprint builds CLI is used to get information about blueprint builds.`,
	Args:  cobra.NoArgs,
}

var avgTimeCmd = &cobra.Command{
	Use:   "avgtime",
	Short: "average time for build step",
	Long:  `Compute average time for a given build step across build executions from a given start-time to end-time.`,
	Args:  cobra.NoArgs,
	RunE:  calcAvgTime,
}

func calcAvgTime(cmd *cobra.Command, args []string) error {
	// set any computed defaults
	if err := setAvgTimeFlagDefaults(); err != nil {
		return err
	}

	// build filters
	filterExpr := successBuildsBtwFilterExpr(avgTimeFlags.lookUpStartTime, avgTimeFlags.lookUpEndTime)
	cFilters := []clientBuildFilter{
		filterRealBuilds,
		filterGHRepoBuilds(avgTimeFlags.repoName),
	}

	// get builds and compute avg
	builds, err := getCBBuildsWithFilter(avgTimeFlags.projectId, filterExpr, cFilters)
	if err != nil {
		return fmt.Errorf("error retrieving builds: %w", err)
	}
	durations, err := findBuildStageDurations(avgTimeFlags.buildStepID, builds)
	if err != nil {
		return err
	}
	if len(durations) < 1 {
		return fmt.Errorf("error no successful build stage %s found", avgTimeFlags.buildStepID)
	}
	avgTime := durationAvg(durations)

	// todo(bharathkkb): Add JSON output
	fmt.Printf("Discovered %d samples for %s stage between %s and %s\n", len(durations), avgTimeFlags.buildStepID, avgTimeFlags.lookUpStart, avgTimeFlags.lookUpEnd)
	color.Green("Computed average time: %s", avgTime)
	return nil
}

// setAvgTimeFlagDefaults sets computed defaults for any missing flags.
// An error is thrown if a default cannot be computed.
func setAvgTimeFlagDefaults() error {
	// if no explicit repo name specified via flag, try to auto discover
	if avgTimeFlags.repoName == "" {
		Log.Info("No repo specified, attempting to detect repo name from current dir")
		path, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working dir: %w", err)
		}
		r, err := getRepoName(path)
		if err != nil {
			return fmt.Errorf("error finding repo name: %w", err)
		}
		if r == "" {
			return fmt.Errorf("unable to detect repo name, please specify a name using --repo")
		}
		avgTimeFlags.repoName = r
		Log.Info("Found repo", "default", avgTimeFlags.repoName)
	}

	// if no explicit build step specified via flag, prompt user with possible options from CloudBuild configs.
	if avgTimeFlags.buildStepID == "" {
		Log.Info("No build ID specified, attempting to find and prompt for build step ID from build file.")
		buildFile, err := getBuildFromFile(avgTimeFlags.buildFilePath)
		if err != nil {
			return fmt.Errorf("error finding build file: %w", err)
		}
		steps := getBuildStepIDs(buildFile)
		avgTimeFlags.buildStepID = util.PromptSelect("Select build step to compute average", steps)
	}

	// if no explicit start time, default to starting computation from one month ago.
	if avgTimeFlags.lookUpStart == "" {
		avgTimeFlags.lookUpStart = time.Now().AddDate(0, -1, 0).Format("01-02-2006")
		Log.Info("No start time specified.", "default", avgTimeFlags.lookUpStart)
	}

	startTime, err := getTimeFromStr(avgTimeFlags.lookUpStart)
	if err != nil {
		return fmt.Errorf("error converting %s to time: %w", avgTimeFlags.lookUpStart, err)
	}
	avgTimeFlags.lookUpStartTime = startTime

	// if no explicit end time, default to ending computation to now.
	if avgTimeFlags.lookUpEnd == "" {
		avgTimeFlags.lookUpEnd = time.Now().Format("01-02-2006")
		Log.Info("No end time specified.", "default", avgTimeFlags.lookUpEnd)
	}

	endTime, err := getTimeFromStr(avgTimeFlags.lookUpEnd)
	if err != nil {
		return fmt.Errorf("error converting %s to time: %w", avgTimeFlags.lookUpEnd, err)
	}
	avgTimeFlags.lookUpEndTime = endTime
	return nil
}
