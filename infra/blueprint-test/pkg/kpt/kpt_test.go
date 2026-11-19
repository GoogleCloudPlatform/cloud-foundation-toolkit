/**
 * Copyright 2026 Google LLC
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

package kpt

import (
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseKptVersion(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "legacy single-line beta version",
			raw:  "1.0.0-beta.67\n",
			want: "1.0.0-beta.67",
		},
		{
			name: "kpt 1.0.1 multiline output",
			raw:  "Version: 1.0.1\nGit commit: 6c8a2e16f1f51605d6db9404ddea4f05db0e2aa2 (dirty)\n",
			want: "1.0.1",
		},
		{
			name: "v-prefixed version",
			raw:  "Version: v1.0.2\n",
			want: "1.0.2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := parseKptVersion(tc.raw)
			assert.Equal(t, tc.want, got)
			assert.NoError(t, utils.MinSemver("v"+got, MIN_KPT_VERSION))
		})
	}
}
