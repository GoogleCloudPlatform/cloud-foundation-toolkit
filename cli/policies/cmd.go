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

package policies

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/forseti-security/config-validator/cmd/policy-tool/status"
)

var flags struct {
	libraryPath string
	bundle      string
}

func init() {
	viper.AutomaticEnv()

	Cmd.AddCommand(status.Cmd)

	Cmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&flags.libraryPath, "path", "", "Path to the policy library root.")

	listCmd.Flags().StringVar(&flags.bundle, "bundle", "scorecard-v1", "Policy bundle to use, default is scorecard")
	// viper.SetDefault("bundle", "txt")
	viper.BindPFlag("bundle", listCmd.Flags().Lookup("bundle"))
}

// Cmd represents the base policies command
var Cmd = &cobra.Command{
	Use:   "policies",
	Short: "Tool to manage a local policy library.",
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List available constraints and constraint templates from a library.",
	Example: `cft policies list`,
	Args:    cobra.NoArgs,
	RunE:    list,
}
