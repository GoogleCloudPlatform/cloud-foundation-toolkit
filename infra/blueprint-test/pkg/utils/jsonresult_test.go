/**
 * Copyright 2024 Google LLC
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

func TestKubectlJSONResult(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output map[string]interface{}
	}{
		{
			name:   "one error",
			input:	`E0223 error the server is currently unable to handle the request
							{
								"apiVersion": "v1"
							}`,
			output: map[string]interface{}{
				"apiVersion": "v1",
			},
		},
		{
			name:   "two error",
			input:	`E0223 error the server is currently unable to handle the request
							E0222 some other error so the server is currently unable to handle the request
							{
								"apiVersion": "v1"
							}`,
			output: map[string]interface{}{
				"apiVersion": "v1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			funcOut := ParseKubectlJSONResult(t, tt.input)
			assert.Equal(tt.output, funcOut.Value().(map[string]interface{}))
		})
	}
}
