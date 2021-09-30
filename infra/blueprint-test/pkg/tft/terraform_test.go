/**
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless assertd by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tft

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	testingiface "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
)

func TestTFBlueprintTest_getTFOutputsAsInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name: "simple",
			input: `
output "simple_string" {
  value = "foo"
}

output "simple_num" {
  value = 1
}

output "simple_bool" {
  value = true
}

output "simple_list" {
  value = ["foo","bar"]
}

output "simple_map" {
  value = {test="hello"}
}
`,
			want: map[string]string{"simple_string": "foo", "simple_num": "1", "simple_bool": "true", "simple_list": "[\"foo\", \"bar\"]", "simple_map": "{\"test\" = \"hello\"}"},
		},
		{
			name:  "empty",
			input: "",
			want:  map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			b := &TFBlueprintTest{t: &testingiface.RuntimeT{}}
			op := getTFOutputMap(t, tt.input)
			got := b.getTFOutputsAsInputs(op)
			assert.Equal(tt.want, got, "inputs should match")
		})
	}
}

func getTFOutputMap(t *testing.T, tf string) map[string]interface{} {
	t.Helper()
	assert := assert.New(t)

	// setup tf file
	tfDir, err := ioutil.TempDir("", "")
	assert.NoError(err)
	defer os.RemoveAll(tfDir)
	tfFilePath := path.Join(tfDir, "test.tf")
	err = ioutil.WriteFile(tfFilePath, []byte(tf), 0644)
	assert.NoError(err)

	// apply tf and get outputs
	tOpts := &terraform.Options{TerraformDir: path.Dir(tfFilePath), Logger: logger.Discard}
	terraform.Init(t, tOpts)
	terraform.Apply(t, tOpts)
	return terraform.OutputAll(t, tOpts)
}

func Test_getKVFromOutputString(t *testing.T) {
	tests := []struct {
		name    string
		kv      string
		wantKey string
		wantVal string
		errMsg  string
	}{
		{name: "simple", kv: "foo=bar", wantKey: "foo", wantVal: "bar"},
		{name: "adjacent equals", kv: "foo==bar", wantKey: "foo", wantVal: "=bar"},
		{name: "no equals invalid", kv: "foobar", errMsg: "error parsing foobar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			gotKey, gotVal, err := getKVFromOutputString(tt.kv)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Equal(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.Equal(tt.wantKey, gotKey)
				assert.Equal(tt.wantVal, gotVal)
			}
		})
	}
}
