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

func TestConfigurationWithViolations(t *testing.T) {
	cwd, _ := os.Getwd()
	libraryPath := path.Join(cwd, "./policy-library")
	fakeT := &testingiface.RuntimeT{}
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	bpt := tft.NewTFBlueprintTest(fakeT,
		tft.WithTFDir("../examples/tf_vet_with_violations"),
		tft.WithSetupPath("./setup/simple_tf_module"),
		tft.WithPolicyLibraryPath(libraryPath))
	bpt.DefineVerify(
		func(assert *assert.Assertions) {
			bpt.DefaultVerify(assert)
		})
	bpt.Test()
	assert.True(t, fakeT.Failed(), "test should have failed")
	assert.Contains(t, logOutput.String(), "GCPServiceUsageConstraintV1")
}
