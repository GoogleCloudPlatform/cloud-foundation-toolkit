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
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/cai"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectResources(t *testing.T) {
	tests := []struct {
		name        string
		assetTypes  []string
		wantKeyPath string
		wantVal     string
	}{
		{name: "all", assetTypes: nil, wantKeyPath: "resource.data.nodeConfig.imageType", wantVal: "COS_CONTAINERD"},
		{name: "cluster", assetTypes: []string{"container.googleapis.com/Cluster", "compute.googleapis.com/Project"}, wantKeyPath: "resource.data.nodeConfig.imageType", wantVal: "COS_CONTAINERD"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			tfBlueprint := tft.NewTFBlueprintTest(t,
				tft.WithTFDir("setup"),
			)

			clusterResourceName := fmt.Sprintf("//container.googleapis.com/projects/%s/locations/%s/clusters/%s",
				tfBlueprint.GetStringOutput("project_id"),
				tfBlueprint.GetStringOutput("cluster_region"),
				tfBlueprint.GetStringOutput("cluster_name"),
			)

			projectResourceName := fmt.Sprintf("//compute.googleapis.com/projects/%s",
				tfBlueprint.GetStringOutput("project_id"),
			)

			// Use the test SA for cai call
			credDec, _ := base64.StdEncoding.DecodeString(tfBlueprint.GetStringOutput("sa_key"))
			gcloud.ActivateCredsAndEnvVars(t, string(credDec))

			cai := cai.GetProjectResources(t, tfBlueprint.GetStringOutput("project_id"), cai.WithAssetTypes(tt.assetTypes))
			assert.Equal(tfBlueprint.GetStringOutput("project_id"), cai.Get("#(name=\"" + projectResourceName + "\").resource.data.name").String(), "project_id exists in cai")
			assert.Equal(tt.wantVal, cai.Get("#(name=\"" + clusterResourceName + "\")." + tt.wantKeyPath).String(), "correct cluster image type")
		})
	}
}
