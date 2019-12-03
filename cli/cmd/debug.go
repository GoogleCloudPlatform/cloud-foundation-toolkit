package cmd

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(debugCmd)
}

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Print debug information for developers.",
	Run: func(cmd *cobra.Command, args []string) {
		glog.Infof("glog info level")
	},
}
