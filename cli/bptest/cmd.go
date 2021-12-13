package bptest

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	testDir string
}

func init() {
	viper.AutomaticEnv()
	Cmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&flags.testDir, "test-dir", "test/integration", "Path to directory containing integration tests")
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
