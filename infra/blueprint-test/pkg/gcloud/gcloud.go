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

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/mitchellh/go-testing-interface"
	"github.com/tidwall/gjson"
)

type Options struct {
	gcloudBinary string         // path to gcloud binary
	commonArgs   []string       // common arguments to pass to gcloud calls
	logger       *logger.Logger // custom logger
}

type option func(*Options)

func WithBinary(gcloudBinary string) option {
	return func(f *Options) {
		f.gcloudBinary = gcloudBinary
	}
}

func WithCommonArgs(commonArgs []string) option {
	return func(f *Options) {
		f.commonArgs = commonArgs
	}
}

func WithLogger(logger *logger.Logger) option {
	return func(f *Options) {
		f.logger = logger
	}
}

// getCommonOptions sets defaults and validates values for gcloud Options.
func GetCommonOptions(opts ...option) (*Options, error) {
	gOpts := &Options{}
	// apply options
	for _, opt := range opts {
		opt(gOpts)
	}
	if gOpts.gcloudBinary == "" {
		err := utils.BinaryInPath("gcloud")
		if err != nil {
			return nil, err
		}
		gOpts.gcloudBinary = "gcloud"
	}
	if gOpts.commonArgs == nil {
		gOpts.commonArgs = []string{"--format", "json"}
	}
	if gOpts.logger == nil {
		gOpts.logger = utils.GetLoggerFromT()
	}
	return gOpts, nil
}

// generateCmd prepares gcloud command to be executed.
func generateCmd(opts *Options, args ...string) shell.Command {
	cmd := shell.Command{
		Command: "gcloud",
		Args:    append(args, opts.commonArgs...),
		Logger:  opts.logger,
	}
	return cmd
}

// RunCmd executes a gcloud command and fails test if there are any errors.
func RunCmd(t testing.TB, options *Options, additionalArgs ...string) string {
	cmd := generateCmd(options, additionalArgs...)
	op, err := shell.RunCommandAndGetStdOutE(t, cmd)
	if err != nil {
		t.Fatal(err)
	}
	return op
}

// Run executes a gcloud command and returns value as gjson.Result.
// It fails the test if there are any errors executing the gcloud command or parsing the output value.
func Run(t testing.TB, cmd string, opts ...option) gjson.Result {
	args := strings.Fields(cmd)
	gOpts, err := GetCommonOptions(opts...)
	if err != nil {
		t.Fatal(err)
	}
	op := RunCmd(t, gOpts, args...)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	return gjson.Parse(op)
}
