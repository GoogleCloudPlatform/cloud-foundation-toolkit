/**
 * Copyright 2022 Google LLC
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

package golden

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/stretchr/testify/assert"

	gotest "github.com/mitchellh/go-testing-interface"
)

const testProjectID = "foo"

func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		skipUpdate bool
		want       string
	}{
		{
			name: "simple",
			data: "foo",
			want: "foo",
		},
		{
			name:       "with-prev-data",
			data:       "{\"baz\":\"qux\"}",
			skipUpdate: true,
			want:       "{\"foo\":\"bar\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			if !tt.skipUpdate {
				os.Setenv(gfUpdateEnvVar, "true")
				defer os.Unsetenv(gfUpdateEnvVar)
			}

			got := NewOrUpdate(t, tt.data)

			if !tt.skipUpdate {
				defer os.Remove(got.GetName())
			}
			j, err := os.ReadFile(got.GetName())
			assert.NoError(err)
			assert.Equal(tt.want, string(j))

		})
	}
}

func TestJSONEq(t *testing.T) {
	tests := []struct {
		name         string
		data         string
		eqPath       string
		opts         []goldenFileOption
		want         string
		setProjectID bool
	}{
		{
			name:   "nested",
			data:   "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"}}",
			eqPath: "baz",
			want:   "{\"qux\":\"quz\"}",
		},
		{
			name:   "empty path",
			data:   "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"}}",
			eqPath: "cookie",
			want:   "",
		},
		{
			name:   "sanitize quz",
			data:   "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"}}",
			opts:   []goldenFileOption{WithSanitizer(StringSanitizer("quz", "REPLACED"))},
			eqPath: "baz",
			want:   "{\"qux\":\"REPLACED\"}",
		},
		{
			name:         "sanitize projectID",
			data:         fmt.Sprintf("{\"foo\":\"bar\",\"baz\":{\"qux\":\"%s\"}}", testProjectID),
			opts:         []goldenFileOption{WithSanitizer(ProjectIDSanitizer(t))},
			setProjectID: true,
			eqPath:       "baz",
			want:         "{\"qux\":\"PROJECT_ID\"}",
		},
		{
			name:   "no gcloud projectID set",
			data:   fmt.Sprintf("{\"foo\":\"bar\",\"baz\":{\"qux\":\"%s\"}}", testProjectID),
			opts:   []goldenFileOption{WithSanitizer(ProjectIDSanitizer(t))},
			eqPath: "baz",
			want:   fmt.Sprintf("{\"qux\":\"%s\"}", testProjectID),
		},
		{
			name: "multiple sanitizers quz",
			data: "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\",\"quux\":\"quuz\"}}",
			opts: []goldenFileOption{
				WithSanitizer(StringSanitizer("quz", "REPLACED")),
				WithSanitizer(func(s string) string { return strings.ReplaceAll(s, "quuz", "NEW") }),
			},
			eqPath: "baz",
			want:   "{\"qux\":\"REPLACED\",\"quux\":\"NEW\"}",
		},
		{
			name:   "diff_whitespace",
			data:   "{\"list\":[\n  \"SYSTEM_COMPONENTS\",\n              \"POD\",\n              \"DAEMONSET\",\n              \"DEPLOYMENT\",\n              \"STATEFULSET\",\n              \"STORAGE\",\n              \"HPA\",\n              \"CADVISOR\",\n              \"KUBELET\"\n            ]}",
			eqPath: "list",
			want:   "[\"SYSTEM_COMPONENTS\",\"POD\",\"DAEMONSET\",\"DEPLOYMENT\",\"STATEFULSET\",\"STORAGE\",\"HPA\",\"CADVISOR\",\"KUBELET\"]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			if tt.setProjectID {
				gcloud.Runf(t, "config set project %s", testProjectID)
				defer gcloud.Run(t, "config unset project")
			}
			os.Setenv(gfUpdateEnvVar, "true")
			defer os.Unsetenv(gfUpdateEnvVar)
			got := NewOrUpdate(t, tt.data, tt.opts...)
			defer os.Remove(got.GetName())
			got.JSONEq(assert, utils.ParseJSONResult(t, tt.data), tt.eqPath)
			assert.Equal(tt.want, got.GetJSON().Get(tt.eqPath).Get("@ugly").String())
		})
	}
}

func TestJSONEqs(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		eqPaths  []string
		opts     []goldenFileOption
		want     string
		hasError bool
	}{
		{
			name:     "simple",
			data:     "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"},\"fizz\":\"pop\"}",
			eqPaths:  []string{"foo","baz"},
			want:     "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"}}",
			hasError: false,
		},
		{
			name:     "simple space diff",
			data:     "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"},\"fizz\":\"pop\"}",
			eqPaths:  []string{"foo","baz"},
			want:     "{          \"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"}}",
			hasError: false,
		},
		{
			name:     "simple order diff",
			data:     "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"},\"fizz\":\"pop\"}",
			eqPaths:  []string{"foo","baz"},
			want:     "{\"baz\":{\"qux\":\"quz\"},\"foo\":\"bar\",\"foo\":\"bar\"}",
			hasError: false,
		},
		{
			name:     "false",
			data:     "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"},\"fizz\":\"pop\"}",
			eqPaths:  []string{"foo","baz"},
			want:     "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz1\"}}",
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			innerT := &gotest.RuntimeT{}
			innerAssert := assert.New(innerT)
			os.Setenv(gfUpdateEnvVar, "true")
			defer os.Unsetenv(gfUpdateEnvVar)
			got := NewOrUpdate(t, tt.data, tt.opts...)
			defer os.Remove(got.GetName())
			got.JSONPathEqs(innerAssert, utils.ParseJSONResult(t, tt.want), tt.eqPaths)

			assert := assert.New(t)
			assert.True(innerT.Failed() == tt.hasError)
		})
	}
}
