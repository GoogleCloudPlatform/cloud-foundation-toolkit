/**
 * Copyright 2023 Google LLC
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

// Package bq provides a set of helpers to interact with bq tool (part of CloudSDK)
package bq

import (
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/mitchellh/go-testing-interface"
	"github.com/tidwall/gjson"
)

type CmdCfg struct {
	bqBinary   string         // path to bq binary
	commonArgs []string       // common arguments to pass to bq calls
	logger     *logger.Logger // custom logger
}

type cmdOption func(*CmdCfg)

func WithBinary(bqBinary string) cmdOption {
	return func(f *CmdCfg) {
		f.bqBinary = bqBinary
	}
}

func WithCommonArgs(commonArgs []string) cmdOption {
	return func(f *CmdCfg) {
		f.commonArgs = commonArgs
	}
}

func WithLogger(logger *logger.Logger) cmdOption {
	return func(f *CmdCfg) {
		f.logger = logger
	}
}

// newCmdConfig sets defaults and validates values for bq Options.
func newCmdConfig(opts ...cmdOption) (*CmdCfg, error) {
	gOpts := &CmdCfg{}
	// apply options
	for _, opt := range opts {
		opt(gOpts)
	}
	if gOpts.bqBinary == "" {
		err := utils.BinaryInPath("bq")
		if err != nil {
			return nil, err
		}
		gOpts.bqBinary = "bq"
	}
	if gOpts.commonArgs == nil {
		gOpts.commonArgs = []string{"--format", "json"}
	}
	if gOpts.logger == nil {
		gOpts.logger = utils.GetLoggerFromT()
	}
	return gOpts, nil
}

// initBq checks for a local .bigqueryrc file and creates an empty one if not to avoid forced bigquery initialization, which doesn't output valid json.
func initBq(t testing.TB) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	fileName := homeDir + "/.bigqueryrc"
	 _ , err = os.Stat(fileName)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
}

// RunCmd executes a bq command and fails test if there are any errors.
func RunCmd(t testing.TB, cmd string, opts ...cmdOption) string {
	op, err := RunCmdE(t, cmd, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return op
}

// RunCmdE executes a bq command and return output.
func RunCmdE(t testing.TB, cmd string, opts ...cmdOption) (string, error) {
	gOpts, err := newCmdConfig(opts...)
	if err != nil {
		t.Fatal(err)
	}
        initBq(t)
	// split command into args
	args := strings.Fields(cmd)
	bqCmd := shell.Command{
		Command: "bq",
		Args:    append(gOpts.commonArgs, args...),
		Logger:  gOpts.logger,
	}
	return shell.RunCommandAndGetStdOutE(t, bqCmd)
}

// Run executes a bq command and returns value as gjson.Result.
// It fails the test if there are any errors executing the bq command or parsing the output value.
func Run(t testing.TB, cmd string, opts ...cmdOption) gjson.Result {
	op := RunCmd(t, cmd, opts...)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	return gjson.Parse(op)
}

// RunWithCmdOptsf executes a bq command and returns value as gjson.Result.
//
// RunWithCmdOptsf(t, ops.., "ls --datasets --project_id=%s", "projectId")
//
// It fails the test if there are any errors executing the bq command or parsing the output value.
func RunWithCmdOptsf(t testing.TB, opts []cmdOption, cmd string, args ...interface{}) gjson.Result {
	return Run(t, utils.StringFromTextAndArgs(append([]interface{}{cmd}, args...)...), opts...)
}

// Runf executes a bq command and returns value as gjson.Result.
//
// Runf(t, "ls --datasets --project_id=%s", "projectId")
//
// It fails the test if there are any errors executing the bq command or parsing the output value.
func Runf(t testing.TB, cmd string, args ...interface{}) gjson.Result {
	return Run(t, utils.StringFromTextAndArgs(append([]interface{}{cmd}, args...)...))
}
