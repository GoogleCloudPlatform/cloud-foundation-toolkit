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
	"encoding/json"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/config-validator/pkg/api/validator"
	"github.com/google/go-cmp/cmp"
)

func jsonToInterface(jsonStr string) (map[string]interface{}, error) {
	var interfaceVar map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &interfaceVar)
	if err != nil {
		return nil, err
	}

	return interfaceVar, nil
}

func TestDataTypeTransformation(t *testing.T) {
	fileContent, err := os.ReadFile(testRoot + "/shared/iam_policy_audit_logs.json")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	asset, err := jsonToInterface(string(fileContent))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	wantedName := "//cloudresourcemanager.googleapis.com/projects/23456"

	pbAsset := &validator.Asset{}
	err = protoViaJSON(asset, pbAsset)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	t.Run("protoViaJSON - CAI asset with unknown fieldsto Proto", func(t *testing.T) {
		if pbAsset.Name != wantedName {
			t.Errorf("wanted %s pbAsset.Name, got %s", wantedName, pbAsset.Name)
		}
	})
	t.Run("interfaceViaJSON", func(t *testing.T) {
		var gotInterface interface{}
		gotInterface, err = interfaceViaJSON(pbAsset)
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		gotName := gotInterface.(map[string]interface{})["name"]
		if gotName != wantedName {
			t.Errorf("wanted %s, got %s", wantedName, gotName)
		}
	})
	t.Run("stringViaJSON", func(t *testing.T) {
		// Compare as structured JSON objects, since eventually this
		// should use the protojson package, which does not support
		// stable serialization. See
		// https://github.com/golang/protobuf/issues/1121#issuecomment-627554847
		gotStr, err := stringViaJSON(pbAsset)
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		var gotJSON map[string]interface{}
		if err := json.Unmarshal([]byte(gotStr), &gotJSON); err != nil {
			t.Fatalf("failed to parse JSON string %v: %v", gotStr, err)
		}

		wantStr := `{"name":"//cloudresourcemanager.googleapis.com/projects/23456","assetType":"cloudresourcemanager.googleapis.com/Project","iamPolicy":{"version":1,"bindings":[{"role":"roles/owner","members":["user:user@example.com"]}],"auditConfigs":[{"service":"storage.googleapis.com","auditLogConfigs":[{"logType":"ADMIN_READ"},{"logType":"DATA_READ"},{"logType":"DATA_WRITE"}]}]},"ancestors":["projects/1234","organizations/56789"]}`
		var wantJSON map[string]interface{}
		if err := json.Unmarshal([]byte(wantStr), &wantJSON); err != nil {
			t.Fatalf("failed to parse JSON string %v: %v", wantStr, err)
		}

		if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
			t.Errorf("stringViaJSON() returned unexpected difference (-want +got):\n%s", diff)
		}
	})
}
