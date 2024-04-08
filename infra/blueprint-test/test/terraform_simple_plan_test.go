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

package test

import (
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	tfjson "github.com/hashicorp/terraform-json"
)

func TestPlan(t *testing.T) {
	networkBlueprint := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("../examples/simple_tf_module"),
		tft.WithSetupPath("setup"),
	)
	networkBlueprint.DefinePlan(func(ps *terraform.PlanStruct, assert *assert.Assertions) {
		assert.Equal(4, len(ps.ResourceChangesMap), "expected 4 resources")
	})
	networkBlueprint.DefineVerify(
		func(assert *assert.Assertions) {
			_, ps := networkBlueprint.PlanAndShow()
			for _, r := range ps.ResourceChangesMap {
				assert.Equal(tfjson.Actions{tfjson.ActionNoop}, r.Change.Actions, "must be no-op")
			}
			op := gcloud.Runf(t, "compute networks subnets describe subnet-01 --project %s --region us-west1", networkBlueprint.GetStringOutput("project_id"))
			assert.Equal("10.10.10.0/24", op.Get("ipCidrRange").String(), "should have the right CIDR")
		})
	networkBlueprint.Test()
}
