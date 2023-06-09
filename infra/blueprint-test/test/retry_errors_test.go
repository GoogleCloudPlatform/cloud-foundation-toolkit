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
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestRetryErrors(t *testing.T) {
	bpt := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("../examples/retry_errors"),
		tft.WithSetupPath("./setup"),
		tft.WithRetryableTerraformErrors(tft.CommonRetryableErrors, 2, 3*time.Second),
	)
	bpt.DefineVerify(func(assert *assert.Assertions) {})
	bpt.DefineTeardown(func(assert *assert.Assertions) {})

	// The default apply is `terraform.Apply(b.t, b.GetTFOptions())` which has a `require.NoError(t, err)`
	// calling `terraform.ApplyE(t, bpt.GetTFOptions())` we are able to process the error end check if it has the
	// "unsuccessful after X retries" message. this works for the this test because the code to sent the retry options
	// to terraform is in the `bpt.GetTFOptions()` function.
	bpt.DefineApply(
		func(assert *assert.Assertions) {
			out, err := terraform.ApplyE(t, bpt.GetTFOptions())
			assert.Contains(out, "SERVICE_DISABLED")
			errMsg := err.Error()
			assert.Equal(errMsg, "'terraform [apply -input=false -auto-approve -no-color -lock=false]' unsuccessful after 2 retries")
		})
	bpt.Test()
}
