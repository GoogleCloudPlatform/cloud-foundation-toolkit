/**
 * Copyright 2022 Google LLC
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

package test

import (
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestRedeploy(t *testing.T) {
	nt := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("../examples/simple_pet_module"),
		tft.WithSetupPath(""),
	)
	nt.DefineVerify(func(a *assert.Assertions) {
		if nt.GetStringOutput("current_ws") == "test-2" {
			a.Equal("custom", nt.GetStringOutput("test"), "should have custom var override")
		} else {
			a.Equal("", nt.GetStringOutput("test"), "should have not have custom var override")
		}
	})
	nt.RedeployTest(3, map[int]map[string]interface{}{2: {"test": "custom"}})
	expectedWorkspaces := []string{"test-1", "test-2", "test-3"}
	for _, ws := range expectedWorkspaces {
		terraform.RunTerraformCommand(t, nt.GetTFOptions(), "workspace", "select", ws)
	}
}
