/**
 * Copyright 2022 Google LLC
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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringFromTextAndArgs(t *testing.T) {
	tests := []struct {
		name   string
		cmd    string
		args   []interface{}
		output string
	}{
		{
			name:   "one arg",
			cmd:    "project list --filter=%s",
			args:   []interface{}{"TEST_PROJECT"},
			output: "project list --filter=TEST_PROJECT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			funcOut := StringFromTextAndArgs(append([]interface{}{tt.cmd}, tt.args...)...)
			assert.Equal(tt.output, funcOut)
		})
	}
}
