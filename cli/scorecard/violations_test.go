// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scorecard

import (
	"context"
	"testing"
)

const (
	testRoot       = "../testdata/scorecard"
	localPolicyDir = testRoot + "/policy-library"
	localCaiDir    = testRoot + "/cai-dir"
)

type getAssetFromJSONTestcase struct {
	name          string
	assetJson     string
	ancestry_path string
	isResource    bool
	isIamPolicy   bool
}

type getViolationsTestcase struct {
	resource   string
	constraint string
}

func TestGetAssetFromJSON(t *testing.T) {
	var testCases = []getAssetFromJSONTestcase{
		{
			name:          "resource",
			assetJson:     testResourceJSON,
			ancestry_path: "organizations/56789/projects/1234",
			isResource:    true,
			isIamPolicy:   false,
		},
		{
			name:          "iam policy",
			assetJson:     testIamPolicyJSON,
			ancestry_path: "organizations/56789/folders/2345/projects/1234",
			isResource:    false,
			isIamPolicy:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pbAsset, err := getAssetFromJSON([]byte(tc.assetJson))
			if err != nil {
				t.Fatal("unexpected error", err)
			}
			gotAncestryPath := pbAsset.GetAncestryPath()
			if gotAncestryPath != tc.ancestry_path {
				t.Errorf("wanted %s ancestry_path, got %s", tc.ancestry_path, gotAncestryPath)
			}

			if tc.isResource && pbAsset.Resource == nil {
				t.Errorf("wanted resource, got %s", pbAsset)
			}
			if tc.isIamPolicy && pbAsset.IamPolicy == nil {
				t.Errorf("wanted IAM Policy bindings, got %s", pbAsset)
			}

		})
	}
}

func TestGetViolations(t *testing.T) {
	var testCases = []getViolationsTestcase{
		{
			resource:   "//storage.googleapis.com/test-project",
			constraint: "iam-gcs-blacklist-public-users",
		},
	}
	inventory, err := NewInventory("", localCaiDir, false, false, TargetProject("1234"), TargetFolder("2345"), TargetOrg("56789"))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	ctx := context.Background()
	config, err := NewScoringConfig(ctx, localPolicyDir)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	auditResult, err := getViolations(inventory, config)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	violationMap := make(map[string]int)
	for _, v := range auditResult.Violations {
		violationMap[v.Constraint+"-"+v.Resource] = 1
	}

	for _, tc := range testCases {
		key := tc.constraint + "-" + tc.resource
		t.Run(key, func(t *testing.T) {
			if violationMap[key] != 1 {
				t.Errorf("wanted violation for resource %s & constraint %s, got none", tc.resource, tc.constraint)
			}
		})
	}
}

var defaultGetAssetFromJSONs = map[string]string{
	"testResourceJSON":  testResourceJSON,
	"testIamPolicyJSON": testIamPolicyJSON,
}

var testResourceJSON = `
{
	"name": "//compute.googleapis.com/projects/test-project",
	"asset_type": "compute.googleapis.com/Project",
	"resource": {
	  "version": "v1",
	  "discovery_document_uri": "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
	  "discovery_name": "Project",
	  "parent": "//cloudresourcemanager.googleapis.com/projects/1234",
	  "data": {
		"creationTimestamp": "2019-04-08T21:19:06.581-07:00",
		"defaultNetworkTier": "PREMIUM",
		"defaultServiceAccount": "1234-compute@developer.gserviceaccount.com",
		"id": "4321",
		"kind": "compute#project",
		"name": "test-project",
		"selfLink": "https://www.googleapis.com/compute/v1/projects/test-project",
		"xpnProjectStatus": "UNSPECIFIED_XPN_PROJECT_STATUS"
	  }
	},
	"ancestors": [
	  "projects/1234",
	  "organizations/56789"
	]
}`

var testIamPolicyJSON = `{
	"name": "//storage.googleapis.com/test-public-bucket-1",
	"asset_type": "storage.googleapis.com/Bucket",
	"iam_policy": {
	  "etag": "WaAAAaAaaaa=",
	  "bindings": [
		{
		  "role": "roles/storage.legacyBucketOwner",
		  "members": [
			"projectEditor:test-project",
			"projectOwner:test-project"
		  ]
		},
		{
		  "role": "roles/storage.objectViewer",
		  "members": [
			"allAuthenticatedUsers"
		  ]
		}
	  ]
	},
	"ancestors": [
	  "projects/1234",
	  "folders/2345",
	  "organizations/56789"
	]
}`
