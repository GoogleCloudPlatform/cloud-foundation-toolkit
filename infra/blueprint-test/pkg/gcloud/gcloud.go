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

type CmdCfg struct {
	gcloudBinary string         // path to gcloud binary
	commonArgs   []string       // common arguments to pass to gcloud calls
	logger       *logger.Logger // custom logger
}

type cmdOption func(*CmdCfg)

func WithBinary(gcloudBinary string) cmdOption {
	return func(f *CmdCfg) {
		f.gcloudBinary = gcloudBinary
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

// newCmdConfig sets defaults and validates values for gcloud Options.
func newCmdConfig(opts ...cmdOption) (*CmdCfg, error) {
	gOpts := &CmdCfg{}
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

// RunCmd executes a gcloud command and fails test if there are any errors.
func RunCmd(t testing.TB, cmd string, opts ...cmdOption) string {
	gOpts, err := newCmdConfig(opts...)
	if err != nil {
		t.Fatal(err)
	}
	// split command into args
	args := strings.Fields(cmd)
	gcloudCmd := shell.Command{
		Command: "gcloud",
		Args:    append(args, gOpts.commonArgs...),
		Logger:  gOpts.logger,
	}
	op, err := shell.RunCommandAndGetStdOutE(t, gcloudCmd)
	if err != nil {
		t.Fatal(err)
	}
	return op
}

// Run executes a gcloud command and returns value as gjson.Result.
// It fails the test if there are any errors executing the gcloud command or parsing the output value.
func Run(t testing.TB, cmd string, opts ...cmdOption) gjson.Result {
	op := RunCmd(t, cmd, opts...)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	return gjson.Parse(op)
}

// ActivateCredsAndEnvVars activates credentials and exports auth related envvars.
func ActivateCredsAndEnvVars(t testing.TB, creds string) {
	credsPath, err := utils.WriteTmpFile(creds)
	if err != nil {
		t.Fatal(err)
	}
	RunCmd(t, "auth activate-service-account", WithCommonArgs([]string{"--key-file", credsPath}))
	// set auth related env vars
	// TF provider auth
	utils.SetEnv(t, "GOOGLE_CREDENTIALS", creds)
	// gcloud SDK override
	utils.SetEnv(t, "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE", credsPath)
	// ADC
	utils.SetEnv(t, "GOOGLE_APPLICATION_CREDENTIALS", credsPath)
}
