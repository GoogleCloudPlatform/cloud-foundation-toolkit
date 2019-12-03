package cmd

import (
	"os"
	"flag"
	
	log "github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/policies"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/report"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/scorecard"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cft",
	Short: "Google Cloud Foundation Toolkit CLI",
	Long:  "Google Cloud Foundation Toolkit CLI",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// no params means same as -h flag
		if args == nil || len(args) == 0 {
			cmd.HelpFunc()(cmd, args)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !flags.verbose {
			// discard logs
			scorecard.Log.SetHandler(log.DiscardHandler())
		}
	},
}

var flags struct {
	// Common flags
	verbose bool
}

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// For https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})

	rootCmd.SetUsageTemplate(`Usage:
  {{if .Runnable}}{{.UseLine}}{{end}}
  {{if .HasAvailableSubCommands}}{{.CommandPath}} [command] [flags]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
	if os.Args == nil {
		rootCmd.SetArgs([]string{"-h"})
	}

	rootCmd.PersistentFlags().BoolVar(&flags.verbose, "verbose", false, "Log output to stdout")

	rootCmd.AddCommand(scorecard.Cmd)
	rootCmd.AddCommand(report.Cmd)
	rootCmd.AddCommand(policies.Cmd)


	// flag.Set("logtostderr", "true")
	// flag.Set("stderrthreshold", "INFO")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
