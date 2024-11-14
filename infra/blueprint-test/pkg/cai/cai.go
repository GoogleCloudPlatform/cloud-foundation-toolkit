/**
 * Copyright 2024 Google LLC
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

// Package cai provides a set of helpers to interact with Cloud Asset Inventory
package cai

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/tidwall/gjson"
)

type CmdCfg struct {
	sleep      int      // minutes to sleep prior to CAI retreval. default: 2
	assetTypes []string // asset types to retrieve. empty: all
	args       []string // arguments to pass to call
}

type cmdOption func(*CmdCfg)

// newCmdConfig sets defaults and options
func newCmdConfig(opts ...cmdOption) (*CmdCfg) {
	caiOpts := &CmdCfg{
		sleep:      2,
		assetTypes: nil,
		args:       nil,
	}

	for _, opt := range opts {
		opt(caiOpts)
	}

	if caiOpts.assetTypes != nil {
		caiOpts.args = []string{"--asset-types", strings.Join(caiOpts.assetTypes, ",")}
	}
	caiOpts.args = append(caiOpts.args, "--content-type", "resource")

	return caiOpts
}

// Set custom sleep minutes
func WithSleep(sleep int) cmdOption {
	return func(f *CmdCfg) {
		f.sleep = sleep
	}
}

// Set asset types
func WithAssetTypes(assetTypes []string) cmdOption {
	return func(f *CmdCfg) {
		f.assetTypes = assetTypes
	}
}

// GetProjectResources returns the cloud asset inventory resources for a project as a gjson.Result
func GetProjectResources(t testing.TB, project string, opts ...cmdOption) gjson.Result {
	caiOpts := newCmdConfig(opts...)

	// Cloud Asset Inventory offers best-effort data freshness.
	t.Logf("Sleeping for %d minutes before retrieving Cloud Asset Inventory...", caiOpts.sleep)
	time.Sleep(time.Duration(caiOpts.sleep) * time.Minute)

	cmd := fmt.Sprintf("asset list --project %s", project)
	return gcloud.Runf(t, strings.Join(append([]string{cmd}, caiOpts.args...), " "))
}
