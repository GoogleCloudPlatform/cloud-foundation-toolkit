package bptest

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	testDir   string
	testStage string
}

func init() {
	viper.AutomaticEnv()
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(runCmd)

	Cmd.PersistentFlags().StringVar(&flags.testDir, "test-dir", "", "Path to directory containing integration tests (default is computed by scanning current working directory)")
	runCmd.Flags().StringVar(&flags.testStage, "stage", "", "Test stage to execute (default is running all stages in order - init, apply, verify, teardown)")
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
	Long:  "Lists both auto discovered and explicit intergration tests",

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
		if err := isValidTestName(flags.testDir, args[0]); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		intTestDir := flags.testDir
		testStage, err := validateAndGetStage(flags.testStage)
		if err != nil {
			return err
		}
		testCmd, err := getTestCmd(intTestDir, testStage, args[0])
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
