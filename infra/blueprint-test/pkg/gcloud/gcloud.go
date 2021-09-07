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

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/binary"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/mitchellh/go-testing-interface"
	"github.com/tidwall/gjson"
)

type CmdCfg struct {
	*binary.BinaryCfg                       // binary config
	bOpts             []binary.BinaryOption // binary options
	commonArgs        []string              // common arguments to pass to gcloud calls
}

type cmdOption func(*CmdCfg)

func WithBinaryOptions(bOpts ...binary.BinaryOption) cmdOption {
	return func(f *CmdCfg) {
		f.bOpts = append(f.bOpts, bOpts...)
	}
}

func WithCommonArgs(commonArgs []string) cmdOption {
	return func(f *CmdCfg) {
		f.commonArgs = commonArgs
	}
}

// newCmdConfig sets defaults and validates values for gcloud Options.
func newCmdConfig(t testing.TB, opts ...cmdOption) (*CmdCfg, error) {
	gOpts := &CmdCfg{}
	// apply options
	for _, opt := range opts {
		opt(gOpts)
	}
	gOpts.BinaryCfg = binary.NewBinaryConfig(t, "gcloud", gOpts.bOpts...)
	if gOpts.commonArgs == nil {
		gOpts.commonArgs = []string{"--format", "json"}
	}
	return gOpts, nil
}

// RunCmd executes a gcloud command and fails test if there are any errors.
func RunCmd(t testing.TB, cmd string, opts ...cmdOption) string {
	gOpts, err := newCmdConfig(t, opts...)
	if err != nil {
		t.Fatal(err)
	}
	// split command into args
	args := strings.Fields(cmd)
	gcloudCmd := shell.Command{
		Command: gOpts.GetBinary(),
		Args:    append(args, gOpts.commonArgs...),
		Logger:  gOpts.GetLogger(),
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
