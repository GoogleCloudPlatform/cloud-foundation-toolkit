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
	"github.com/tidwall/gjson"
)

 func TestGetJSONPaths(t *testing.T) {
	 tests := []struct {
		 name  string
		 json  gjson.Result
		 paths []string
	 }{
		 {
			 name:   "one",
			 json:	gjson.Parse(`{
								 "apiVersion": "v1",
								 "autopilot": {},
								 "locations": [
									 "europe-west4-b"
								 ],
								 "metadata": {
								   "annotations": [
									   {"my-annotation": "test"}
									 ]
								 },
								"bool": true,
								"number": 3,
								"null": null,
							 }`),
			 paths: []string{"apiVersion", "autopilot", "bool", "locations", "locations.0", "metadata", "metadata.annotations", "metadata.annotations.0", "metadata.annotations.0.my-annotation", "null", "number"},
		 },
	 }
	 for _, tt := range tests {
		 t.Run(tt.name, func(t *testing.T) {
			 assert := assert.New(t)
			 assert.Equal(tt.paths, GetJSONPaths(tt.json))
		 })
	 }
 }

 func TestTerminalGetJSONPaths(t *testing.T) {
	tests := []struct {
		name  string
		json  gjson.Result
		paths []string
	}{
		{
			name:   "one",
			json:	gjson.Parse(`{
								"apiVersion": "v1",
								"autopilot": {},
								"locations": [
									"europe-west4-b"
								],
								"metadata": {
									"annotations": [
										{"my-annotation": "test"}
									]
								},
								"bool": true,
								"number": 3,
								"null": null,
							}`),
			paths: []string{"apiVersion", "bool", "locations.0", "metadata.annotations.0.my-annotation", "null", "number"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tt.paths, GetTerminalJSONPaths(tt.json))
		})
	}
}
