// Copyright 2021 Google LLC
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
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteViolations(t *testing.T) {
	// Prepare violations
	inventory, err := NewInventory("", localCaiDir, false, false, TargetOrg("56789"))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	ctx := context.Background()
	config, err := NewScoringConfig(ctx, localPolicyDir)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	err = inventory.findViolations(config)
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	tests := []struct {
		format    string
		filename  string
		message   string
		listMaker func([]byte) []interface{}
	}{
		{
			format: "json", filename: "violations.json", message: "The JSON output should be equivalent.",
			listMaker: func(output []byte) []interface{} {
				var outputJSON []interface{}
				if err = json.Unmarshal(output, &outputJSON); err != nil {
					t.Fatal("unexpected error", err)
				}
				return outputJSON
			},
		},
		{
			format: "txt", filename: "violations.txt", message: "The text output should be equivalent.",
			listMaker: func(output []byte) []interface{} {
				outputLines := make([]interface{}, 1)
				scanner := bufio.NewScanner(bytes.NewReader(output))
				for scanner.Scan() {
					outputLines = append(outputLines, scanner.Text())
				}
				return outputLines
			},
		},
	}

	for _, tc := range tests {
		output := new(bytes.Buffer)
		fileContent, err := ioutil.ReadFile(testRoot + "/output/" + tc.filename)
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		expected := tc.listMaker(fileContent)

		writeResults(config, output, tc.format, nil)
		actual := tc.listMaker(output.Bytes())

		assert.ElementsMatch(t, expected, actual, tc.message)
	}

	// Test CSV output
	csvOutput := new(bytes.Buffer)
	expectedLines := []string{
		"Category,Constraint,Resource,Message,Parent",
		"Other,forbid-subnets,//compute.googleapis.com/projects/my-cai-project/regions/europe-north1/subnetworks/default,//compute.googleapis.com/projects/my-cai-project/regions/europe-north1/subnetworks/default is in violation.,",
		"Other,org-policy-skip-default-network,//cloudresourcemanager.googleapis.com/organizations/567890,Required enforcement of skipDefaultNetworkCreation at org level,",
		"Other,vpc-sc-ensure-services,//cloudresourcemanager.googleapis.com/organizations/56789,Required services compute.googleapis.com missing from service perimeter: accessPolicies/12345/servicePerimeters/perimeter_gcs.,",
		"Security,iam-gcs-blacklist-public-users,//storage.googleapis.com/test-bucket-public,//storage.googleapis.com/test-bucket-public is publicly accessable,",
	}
	writeResults(config, csvOutput, "csv", nil)

	// assert equality without caring about order
	scanner := bufio.NewScanner(csvOutput)
	var outputLines []string
	for scanner.Scan() {
		outputLines = append(outputLines, scanner.Text())
	}
	assert.ElementsMatch(t, expectedLines, outputLines, "The CSV output should contain the same values.")
}
