// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package report

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	queryPath    string
	reportPath   string
	reportFormat string
	bucketName   string
	dirName      string
	listReports  bool
}

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVar(&flags.queryPath, "query-path", "", "Path to directory containing inventory queries")
	Cmd.MarkFlagRequired("query-path")

	Cmd.Flags().StringVar(&flags.reportPath, "report-path", "", "Path to directory to contain report outputs")

	Cmd.Flags().StringVar(&flags.bucketName, "bucket", "", "GCS bucket name for storing inventory (conflicts with --dir-path)")
	Cmd.Flags().StringVar(&flags.dirName, "dir-path", "", "Local directory path for storing inventory (conflicts with --bucket)")

	Cmd.Flags().StringVar(&flags.reportFormat, "report-format", "", "Format of inventory report outputs, can be json or csv, default is csv")
	viper.SetDefault("report-format", "csv")
	viper.BindPFlag("report-format", Cmd.Flags().Lookup("report-format"))

	Cmd.Flags().BoolVar(&flags.listReports, "list-available-reports", false, "List available inventory report queries")
}

// Cmd represents the base scorecard command
var Cmd = &cobra.Command{
	Use:   "report",
	Short: "Generate inventory reports based on CAI outputs in a directory.",
	Long: `Generate inventory reports for resources in Cloud Asset Inventory (CAI) output files, with reports defined in rego (in 'samplereports' folder).
	
	Example:
	  cft report --query-path ./path/to/cloud-foundation-toolkit/cli/samplereports \
		--dir-path ./path/to/cai-export-directory \
		--report-path ./path/to/report-output-directory \
	`,

	Args: cobra.NoArgs,
	PreRunE: func(c *cobra.Command, args []string) error {
		if !flags.listReports {
			if flags.reportPath == "" {
				return errors.New("missing required argument --report-path")
			}
			if (flags.bucketName == "" && flags.dirName == "") ||
				(flags.bucketName != "" && flags.dirName != "") {
				return errors.New("Either bucket or dir-path should be set")
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !flags.listReports {
			err := GenerateReports(flags.dirName, flags.queryPath, flags.reportPath, viper.GetString("report-format"))
			if err != nil {
				return err
			}
		} else {
			err := ListAvailableReports(flags.queryPath)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
