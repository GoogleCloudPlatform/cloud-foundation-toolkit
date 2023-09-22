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
	"os"
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	testingiface "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
)

func TestGetTFOutputsAsInputs(t *testing.T) {
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

// newTestDir creates a new directory suitable for use as TFDir
func newTestDir(t *testing.T, pattern string, input string) string {
	t.Helper()
	assert := assert.New(t)

	// setup tf file
	tfDir, err := os.MkdirTemp("", pattern)
	assert.NoError(err)
	tfFilePath := path.Join(tfDir, "test.tf")
	err = os.WriteFile(tfFilePath, []byte(input), 0644)
	assert.NoError(err)
	return tfDir
}

func getTFOutputMap(t *testing.T, tf string) map[string]interface{} {
	t.Helper()

	tfDir := newTestDir(t, "", tf)
	defer os.RemoveAll(tfDir)

	// apply tf and get outputs
	tOpts := &terraform.Options{TerraformDir: tfDir, Logger: logger.Discard}
	terraform.Init(t, tOpts)
	terraform.Apply(t, tOpts)
	return terraform.OutputAll(t, tOpts)
}

func TestGetKVFromOutputString(t *testing.T) {
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

func TestSetupOverrideString(t *testing.T) {
	tests := []struct {
		name      string
		tfOutputs string
		overrides map[string]interface{}
		want      map[string]string
	}{
		{name: "no overrides",
			tfOutputs: `
			output "simple_string" {
			  value = "foo"
			}

			output "simple_num" {
			  value = 1
			}

			output "simple_bool" {
			  value = true
			}
			`,
			overrides: map[string]interface{}{},
			want: map[string]string{
				"simple_string": "foo",
				"simple_num":    "1",
				"simple_bool":   "true",
			},
		},
		{name: "all overrides",
			tfOutputs: `
			output "simple_string" {
			  value = "foo"
			}

			output "simple_num" {
			  value = 1
			}

			output "simple_bool" {
			  value = true
			}
			`,
			overrides: map[string]interface{}{
				"simple_string": "bar",
				"simple_num":    "2",
				"simple_bool":   "false",
			},
			want: map[string]string{
				"simple_string": "bar",
				"simple_num":    "2",
				"simple_bool":   "false",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emptyDir := newTestDir(t, "empty*", "")
			setupDir := newTestDir(t, "setup-*", tt.tfOutputs)
			defer os.RemoveAll(emptyDir)
			defer os.RemoveAll(setupDir)
			b := NewTFBlueprintTest(&testingiface.RuntimeT{},
				WithSetupOutputs(tt.overrides),
				WithTFDir(emptyDir),
				WithSetupPath(setupDir))
			// create outputs from setup
			_, err := terraform.ApplyE(t, &terraform.Options{TerraformDir: setupDir})
			if err != nil {
				t.Fatalf("Failed to apply setup: %v", err)
			}
			for k, want := range tt.want {
				if b.GetTFSetupStringOutput(k) != want {
					t.Errorf("unexpected string output for %s: want %s got %s", k, want, b.GetStringOutput(k))
				}
			}
		})
	}
}
func TestSetupOverrideList(t *testing.T) {
	tests := []struct {
		name      string
		tfOutputs string
		overrides map[string]interface{}
		want      map[string][]string
	}{
		{name: "no overrides",
			tfOutputs: `
				output "simple_list" {
					value = ["foo","bar"]
				}
			`,
			overrides: map[string]interface{}{},
			want: map[string][]string{
				"simple_list": {"foo", "bar"},
			},
		},
		{name: "all overrides",
			tfOutputs: `
				output "simple_list" {
					value = ["foo","bar"]
				}
			`,
			overrides: map[string]interface{}{
				"simple_list": []string{"apple", "orange"},
			},
			want: map[string][]string{
				"simple_list": {"apple", "orange"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emptyDir := newTestDir(t, "empty*", "")
			setupDir := newTestDir(t, "setup-*", tt.tfOutputs)
			defer os.RemoveAll(emptyDir)
			defer os.RemoveAll(setupDir)
			b := NewTFBlueprintTest(&testingiface.RuntimeT{},
				WithSetupOutputs(tt.overrides),
				WithTFDir(emptyDir),
				WithSetupPath(setupDir))
			// create outputs from setup
			_, err := terraform.ApplyE(t, &terraform.Options{TerraformDir: setupDir})
			if err != nil {
				t.Fatalf("Failed to apply setup: %v", err)
			}
			for k, want := range tt.want {
				got := b.GetTFSetupOutputListVal(k)
				assert.ElementsMatchf(t, got, want, "list mismatch: want %s got %s", want)
			}
		})
	}

}

func TestSetupOverrideFromEnv(t *testing.T) {
	t.Setenv("CFT_SETUP_my-key", "my-value")
	emptyDir := newTestDir(t, "empty*", "")
	defer os.RemoveAll(emptyDir)
	b := NewTFBlueprintTest(&testingiface.RuntimeT{},
		WithTFDir(emptyDir))
	got := b.GetTFSetupStringOutput("my-key")
	assert.Equal(t, got, "my-value")
}
