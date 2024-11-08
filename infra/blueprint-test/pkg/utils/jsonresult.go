/**
 * Copyright 2021 Google LLC
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
	"os"
	"regexp"

	"github.com/mitchellh/go-testing-interface"
	"github.com/tidwall/gjson"
)

// LoadJSON reads and parses a json file into a gjson.Result.
// It fails test if not unable to parse.
func LoadJSON(t testing.TB, path string) gjson.Result {
	j, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error reading json file %s", path)
	}
	if !gjson.ValidBytes(j) {
		t.Fatalf("Error parsing output, invalid json: %s", path)
	}
	return gjson.ParseBytes(j)
}

// ParseJSONResult converts a JSON string into gjson result
func ParseJSONResult(t testing.TB, j string) gjson.Result {
	if !gjson.Valid(j) {
		t.Fatalf("Error parsing output, invalid json: %s", j)
	}
	return gjson.Parse(j)
}

// Kubectl transient errors
var (
	KubectlTransientErrors = []string{
		"E022[23] .* the server is currently unable to handle the request",
	}
)

// Filter transient errors from kubectl output
func ParseKubectlJSONResult(t testing.TB, str string) gjson.Result {
	for _, error := range KubectlTransientErrors {
		str = regexp.MustCompile(error).ReplaceAllString(str, "")
	}
	return ParseJSONResult(t, str)
}
