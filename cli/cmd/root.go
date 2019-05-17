package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cft",
	Short: "Google Cloud Formation Toolkit CLI",
	Long:  "Google Cloud Formation Toolkit CLI",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// no params means same as -h flag
		if args == nil || len(args) == 0 {
			cmd.HelpFunc()(cmd, args)
		}
	},
}

func init() {
	if os.Args == nil {
		rootCmd.SetArgs([]string{"-h"})
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
