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
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/stretchr/testify/assert"
)

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
			j, err := ioutil.ReadFile(got.GetName())
			assert.NoError(err)
			assert.Equal(tt.want, string(j))

		})
	}
}

func TestJSONEq(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		eqPath string
		opts   []goldenFileOption
		want   string
	}{
		{
			name:   "nested",
			data:   "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"}}",
			eqPath: "baz",
			want:   "{\"qux\":\"quz\"}",
		},
		{
			name:   "sanitize quz",
			data:   "{\"foo\":\"bar\",\"baz\":{\"qux\":\"quz\"}}",
			opts:   []goldenFileOption{WithSanitizer(StringSanitizer("quz", "REPLACED"))},
			eqPath: "baz",
			want:   "{\"qux\":\"REPLACED\"}",
		},
		{
			name:   "sanitize projectID",
			data:   fmt.Sprintf("{\"foo\":\"bar\",\"baz\":{\"qux\":\"%s\"}}", os.Getenv("TEST_GCLOUD_PROJECT")),
			opts:   []goldenFileOption{WithSanitizer(ProjectIDSanitizer(t))},
			eqPath: "baz",
			want:   "{\"qux\":\"PROJECT_ID\"}",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			os.Setenv(gfUpdateEnvVar, "true")
			defer os.Unsetenv(gfUpdateEnvVar)
			got := NewOrUpdate(t, tt.data, tt.opts...)
			defer os.Remove(got.GetName())
			got.JSONEq(assert, utils.ParseJSONResult(t, tt.data), tt.eqPath)
			assert.JSONEq(tt.want, got.GetJSON().Get(tt.eqPath).String())
		})
	}
}
