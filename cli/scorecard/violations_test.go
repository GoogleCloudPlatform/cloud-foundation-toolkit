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
	"os"
	"testing"
)

const (
	testRoot       = "../testdata/scorecard"
	localPolicyDir = testRoot + "/policy-library"
	localCaiDir    = testRoot + "/cai-dir"
)

type getAssetFromJSONTestcase struct {
	name          string
	assetJSONFile string
	ancestryPath  string
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
			assetJSONFile: "/shared/resource.json",
			ancestryPath:  "organizations/56789/projects/1234",
			isResource:    true,
			isIamPolicy:   false,
		},
		{
			name:          "iam policy",
			assetJSONFile: "/shared/iam_policy.json",
			ancestryPath:  "organizations/56789/folders/2345/projects/1234",
			isResource:    false,
			isIamPolicy:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileContent, err := os.ReadFile(testRoot + tc.assetJSONFile)
			if err != nil {
				t.Fatal("unexpected error", err)
			}

			pbAsset, err := getAssetFromJSON(fileContent)
			if err != nil {
				t.Fatal("unexpected error", err)
			}
			gotAncestryPath := pbAsset.GetAncestryPath()
			if gotAncestryPath != tc.ancestryPath {
				t.Errorf("wanted %s ancestry_path, got %s", tc.ancestryPath, gotAncestryPath)
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
			resource:   "//storage.googleapis.com/test-bucket-public",
			constraint: "GCPStorageBucketWorldReadableConstraintV1.iam-gcs-blacklist-public-users",
		},
		{
			resource:   "//cloudresourcemanager.googleapis.com/organizations/567890",
			constraint: "GCPOrgPolicySkipDefaultNetworkConstraintV1.org-policy-skip-default-network",
		},
		{
			resource:   "//cloudresourcemanager.googleapis.com/organizations/56789",
			constraint: "GCPVPCSCEnsureServicesConstraintV1.vpc-sc-ensure-services",
		},
	}
	inventory, err := NewInventory("", localCaiDir, false, false, WorkerSize(1), TargetProject("1234"), TargetFolder("2345"), TargetOrg("56789"))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	ctx := context.Background()
	config, err := NewScoringConfig(ctx, localPolicyDir)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	violations, err := getViolations(inventory, config)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	violationMap := make(map[string]int)
	for _, v := range violations {
		violationMap[v.Violation.Constraint+"-"+v.Resource] = 1
		Log.Debug("Found violation", "constraint", v.Violation.Constraint, "resource", v.Resource)
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
