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

package test

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestCFTSimpleModule(t *testing.T) {
	networkBlueprint := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("../examples/simple_tf_module"),
		tft.WithSetupPath("setup/simple_tf_module"),
	)
	networkBlueprint.DefineVerify(
		func(assert *assert.Assertions) {
			networkBlueprint.DefaultVerify(assert)
			op := gcloud.Run(t, fmt.Sprintf("compute networks subnets describe subnet-01 --project %s --region us-west1", networkBlueprint.GetStringOutput("project_id")))
			assert.Equal("10.10.10.0/24", op.Get("ipCidrRange").String(), "should have the right CIDR")
			assert.Equal("false", op.Get("logConfig.enable").String(), "logConfig should not be enabled")
		})
	networkBlueprint.Test()
}
