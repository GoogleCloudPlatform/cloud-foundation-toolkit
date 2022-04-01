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
	"github.com/mitchellh/go-testing-interface"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// GetFirstMatchResult returns the first matching result with a given k/v
func GetFirstMatchResult(t testing.TB, rs []gjson.Result, k, v string) gjson.Result {
	for _, r := range rs {
		if r.Get(k).Exists() && r.Get(k).String() == v {
			return r
		}
	}
	t.Fatalf("unable to find key %s with value %s in %s", k, v, rs)
	return gjson.Result{}
}

// GetResultStrSlice parses results into a string slice
func GetResultStrSlice(rs []gjson.Result) []string {
	s := make([]string, 0)
	for _, r := range rs {
		s = append(s, r.String())
	}
	return s
}

// DeleteFromResult deletes given path from result and returns the modified result
func DeleteFromResult(t testing.TB, pth string, r gjson.Result) gjson.Result {
	n, err := sjson.Delete(r.String(), pth)
	if err != nil {
		t.Fatalf("unable to delete path %s in %s", pth, r.String())
	}
	return ParseJSONResult(t, n)
}
