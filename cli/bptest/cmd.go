package bptest

import (
	"fmt"
	"os"
	"path"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	testDir   string
	testStage string
	setupVars map[string]string
}

func init() {
	viper.AutomaticEnv()
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(runCmd)
	Cmd.AddCommand(convertCmd)
	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(lintCmd)

	Cmd.PersistentFlags().StringVar(&flags.testDir, "test-dir", "", "Path to directory containing integration tests (default is computed by scanning current working directory)")
	runCmd.Flags().StringVar(&flags.testStage, "stage", "", "Test stage to execute (default is running all stages in order - init, plan, apply, verify, teardown)")
	runCmd.Flags().StringToStringVar(&flags.setupVars, "setup-var", map[string]string{}, "Specify outputs from the setup phase (useful with --stage=verify)")
}

var Cmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"bptest"},
	Short:   "Blueprint test CLI",
	Long:    `Blueprint test CLI is used to actuate the Blueprint test framework used for testing KRM and Terraform Blueprints`,
	Args:    cobra.NoArgs,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list tests",
	Long:  "Lists both auto discovered and explicit integration tests",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		intTestDir := flags.testDir
		tests, err := getTests(intTestDir)
		if err != nil {
			return err
		}
		// Warn if no tests found
		if len(tests) < 1 {
			Log.Warn("no tests discovered")
			return nil
		}
		tbl := newTable()
		tbl.AppendHeader(table.Row{"Name", "Config", "Location"})
		for _, t := range tests {
			if t.bptestCfg.Spec.Skip {
				Log.Info(fmt.Sprintf("skipping %s due to BlueprintTest config %s", t.name, t.bptestCfg.Name))
				continue
			}
			tbl.AppendRow(table.Row{t.name, t.config, t.location})
		}
		tbl.Render()
		return nil
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run tests",
	Long:  "Runs auto discovered and explicit integration tests",

	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		intTestDir, err := getIntTestDir(flags.testDir)
		if err != nil {
			return fmt.Errorf("error discovering test dir: %w", err)
		}
		testStage, err := validateAndGetStage(flags.testStage)
		if err != nil {
			return err
		}
		relTestPkg, err := validateAndGetRelativeTestPkg(intTestDir, args[0])
		if err != nil {
			return err
		}
		testCmd, err := getTestCmd(intTestDir, testStage, args[0], relTestPkg, flags.setupVars)
		if err != nil {
			return err
		}
		// if err during exec, exit instead of returning an error
		// this prevents printing usage as the args were validated above
		if err := streamExec(testCmd); err != nil {
			Log.Error(err.Error())
			os.Exit(1)
		}
		return nil
	},
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert kitchen tests (experimental)",
	Long:  "Convert all kitchen tests to blueprint tests (experimental)",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return convertKitchenTests()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize blueprint test",
	Long:  "Initialize a new blueprint test",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var initTestName string
		// if no args, prompt user to select from examples
		if len(args) < 1 {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			examplePaths, err := util.WalkTerraformDirs(path.Join(cwd, "examples"))
			if err != nil {
				return err
			}
			exampleNames := make([]string, 0, len(examplePaths))
			for _, examplePath := range examplePaths {
				exampleNames = append(exampleNames, path.Base(examplePath))
			}
			initTestName = util.PromptSelect("Select example for test", exampleNames)
		} else {
			initTestName = args[0]
		}
		return initTest(initTestName)
	},
}

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lints blueprint",
	Long:  "Lints TF blueprint",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		RunLintCommand()
		return nil
	},
}
