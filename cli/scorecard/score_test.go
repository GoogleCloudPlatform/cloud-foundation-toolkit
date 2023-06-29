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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeLineList(output []byte) []interface{} {
	outputLines := make([]interface{}, 1)
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		outputLines = append(outputLines, scanner.Text())
	}
	return outputLines
}

func TestWriteViolations(t *testing.T) {
	// Prepare violations
	inventory, err := NewInventory("", localCaiDir, false, false, WorkerSize(1), TargetOrg("56789"))
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
			listMaker: makeLineList,
		},
		{
			format: "csv", filename: "violations.csv", message: "The csv output should be equivalent.",
			listMaker: makeLineList,
		},
	}

	for _, tc := range tests {
		output := new(bytes.Buffer)
		fileContent, err := os.ReadFile(testRoot + "/output/" + tc.filename)
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		expected := tc.listMaker(fileContent)

		err = writeResults(config, output, tc.format, nil)
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		actual := tc.listMaker(output.Bytes())

		assert.ElementsMatch(t, expected, actual, tc.message)
	}
}
