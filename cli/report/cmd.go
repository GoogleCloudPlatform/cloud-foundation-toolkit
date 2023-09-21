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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	queryPath    string
	outputPath   string
	reportFormat string
	bucketName   string
	dirName      string
}

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVar(&flags.queryPath, "query-path", "", "Path to directory containing inventory queries")
	err := Cmd.MarkFlagRequired("query-path")
	if err != nil {
		panic(err)
	}

	Cmd.Flags().StringVar(&flags.outputPath, "output-path", "", "Path to directory to contain report outputs")
	err = Cmd.MarkFlagRequired("output-path")
	if err != nil {
		panic(err)
	}

	//Cmd.Flags().StringVar(&flags.bucketName, "bucket", "", "GCS bucket name for storing inventory (conflicts with --dir-path)")
	Cmd.Flags().StringVar(&flags.dirName, "dir-path", "", "Local directory path for storing inventory ")
	err = Cmd.MarkFlagRequired("dir-path")
	if err != nil {
		panic(err)
	}

	Cmd.Flags().StringVar(&flags.reportFormat, "report-format", "", "Format of inventory report outputs, can be json or csv, default is csv")
	viper.SetDefault("report-format", "csv")
	err = viper.BindPFlag("report-format", Cmd.Flags().Lookup("report-format"))
	if err != nil {
		panic(err)
	}

	Cmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&flags.queryPath, "query-path", "", "Path to directory containing inventory queries")
	err = listCmd.MarkFlagRequired("query-path")
	if err != nil {
		panic(err)
	}
}

// Cmd represents the base report command
var Cmd = &cobra.Command{
	Use:   "report",
	Short: "Generate inventory reports based on CAI outputs in a directory.",
	Long: `Generate inventory reports for resources in Cloud Asset Inventory (CAI) output files, with reports defined in rego (in '<path_to_cloud-foundation-toolkit>/reports/sample' folder).

	Example:
	  cft report --query-path <path_to_cloud-foundation-toolkit>/reports/sample \
		--dir-path <path-to-directory-containing-cai-export> \
		--report-path <path-to-directory-for-report-output>
	`,

	Args: cobra.NoArgs,
	/*
		PreRunE: func(c *cobra.Command, args []string) error {
			if (flags.bucketName == "" && flags.dirName == "") ||
				(flags.bucketName != "" && flags.dirName != "") {
				return errors.New("Either bucket or dir-path should be set")
			}
			return nil
		},
	*/
	RunE: func(cmd *cobra.Command, args []string) error {
		err := GenerateReports(flags.dirName, flags.queryPath, flags.outputPath, viper.GetString("report-format"))
		if err != nil {
			return err
		}
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list-available-reports",
	Short: "List available inventory report queries.",
	Long: `List available inventory report queries for resources in Cloud Asset Inventory (CAI).

	Example:
	  cft report list-available-reports --query-path <path_to_cloud-foundation-toolkit>/reports/sample
	`,

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := ListAvailableReports(flags.queryPath)
		if err != nil {
			return err
		}
		return nil
	},
}
