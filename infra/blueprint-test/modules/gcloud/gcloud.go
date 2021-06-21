/**
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package gcloud provides a set of helpers to interact with gcloud(Cloud SDK) binary
package gcloud

import (
	"strings"

	gotest "testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/modules/utils"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/mitchellh/go-testing-interface"
	"github.com/tidwall/gjson"
)

type Options struct {
	GCloudBinary string         // path to gcloud binary
	CommonArgs   []string       // common arguments to pass to gcloud calls
	Logger       *logger.Logger // custom logger
}

// getCommonOptions sets defaults and validates values for gcloud Options
func getCommonOptions(options *Options, args ...string) (*Options, []string, error) {
	if options.GCloudBinary == "" {
		err := utils.BinaryInPath("gcloud")
		if err != nil {
			return nil, nil, err
		}
		options.GCloudBinary = "gcloud"
	}
	if options.CommonArgs == nil {
		options.CommonArgs = []string{"--format", "json"}
	}
	if options.Logger == nil {
		if gotest.Verbose() {
			options.Logger = logger.Default
		} else {
			options.Logger = logger.Discard
		}

	}
	return options, args, nil
}

// generateCmd prepares gcloud command to be executed
func generateCmd(opts *Options, args ...string) shell.Command {
	cmd := shell.Command{
		Command: "gcloud",
		Args:    append(args, opts.CommonArgs...),
		Logger:  opts.Logger,
	}
	return cmd
}

// RunCmd executes a gcloud command and fails test if there are any errors
func RunCmd(t testing.TB, additionalOptions *Options, additionalArgs ...string) string {
	options, args, err := getCommonOptions(additionalOptions, additionalArgs...)
	if err != nil {
		t.Fatal(err)
	}
	cmd := generateCmd(options, args...)
	op, err := shell.RunCommandAndGetStdOutE(t, cmd)
	if err != nil {
		t.Fatal(err)
	}
	return op
}

// RunWithOptsAndOutput executes a gcloud command with custom options and returns value as gjson.Result
// It fails the test if there are any errors executing the gcloud command or parsing the output value
func RunWithOptsAndOutput(t testing.TB, options *Options, cmd string) gjson.Result {
	args := strings.Fields(cmd)
	op := RunCmd(t, options, args...)

	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	return gjson.Parse(op)
}

// Run executes a gcloud command with default options and returns value as gjson.Result
// It fails the test if there are any errors executing the gcloud command or parsing the output value
func Run(t testing.TB, cmd string) gjson.Result {
	return RunWithOptsAndOutput(t, &Options{}, cmd)
}
