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
	"bytes"
	"log"
	"os"
	"path"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	testingiface "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
)

func TestTerraformVet(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to get wd :%v", err)
	}
	libraryPath := path.Join(cwd, "../examples/policy-library")

	for _, tc := range []struct {
		name           string
		service        string
		errMsgContains string
	}{
		{
			name:    "Valid configuration",
			service: "cloudresourcemanager.googleapis.com",
		},
		{
			name:           "Configuration with violations",
			service:        "oslogin.googleapis.com",
			errMsgContains: "GCPServiceUsageConstraintV1",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fakeT := &testingiface.RuntimeT{}
			var logOutput bytes.Buffer
			log.SetOutput(&logOutput)
			vars := map[string]interface{}{"service": tc.service}

			temp := tft.NewTFBlueprintTest(fakeT,
				tft.WithVars(vars),
				tft.WithTFDir("../examples/tf_vet"),
				tft.WithSetupPath("./setup"),
			)

			bpt := tft.NewTFBlueprintTest(fakeT,
				tft.WithVars(vars),
				tft.WithTFDir("../examples/tf_vet"),
				tft.WithSetupPath("./setup"),
				tft.WithPolicyLibraryPath(libraryPath, temp.GetTFSetupStringOutput("project_id")))
			bpt.DefineVerify(
				func(assert *assert.Assertions) {
					bpt.DefaultVerify(assert)
				})
			bpt.Test()

			if tc.errMsgContains == "" {
				assert.False(t, fakeT.Failed(), "test should be sucessful")
			} else {
				assert.True(t, fakeT.Failed(), "test should have failed")
				assert.Contains(t, logOutput.String(), tc.errMsgContains)
			}

		})
	}
}
