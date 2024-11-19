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
	"slices"
	"strconv"
	"strings"

	"github.com/mitchellh/go-testing-interface"
	"github.com/tidwall/gjson"
)

// GetJSONPaths returns a []string of all possible JSON paths for a gjson.Result
 func GetJSONPaths(t testing.TB, result gjson.Result) []string {
	 return getJSONPaths(t, result.Value(), []string{})
 }

 func getJSONPaths(t testing.TB, item interface{}, crumbs []string) []string {
	var paths []string

	switch val := item.(type) {
		case []interface{}:
			for i, v := range val {
				// Add this item to paths
				paths = append(paths, strings.Join(append(crumbs, strconv.Itoa(i)), "."))
				// Search child items
				paths = append(paths,
					getJSONPaths(t, v, append(crumbs, strconv.Itoa(i)))...,
				)
			}
		case map[string]interface{}:
			for k, v := range val {
				// Add this item to paths
				paths = append(paths, strings.Join(append(crumbs, k), "."))
				// Search child items
				paths = append(paths,
					getJSONPaths(t, v, append(crumbs, k))...,
				)

			}
		}

		slices.Sort(paths)
		return paths
 }

// GetTerminalJSONPaths returns a []string of all terminal JSON paths for a gjson.Result
func GetTerminalJSONPaths(t testing.TB, result gjson.Result) []string {
	return getTerminalJSONPaths(t, result.Value(), []string{})
}

func getTerminalJSONPaths(t testing.TB, item interface{}, crumbs []string) []string {
	var paths []string

	// Only add paths for JSON bool, number, string, and null
	switch val := item.(type) {
		case bool:
			return []string{strings.Join(crumbs, ".")}
		case float64:
			return []string{strings.Join(crumbs, ".")}
		case string:
			return []string{strings.Join(crumbs, ".")}
		case nil:
			return []string{strings.Join(crumbs, ".")}
		case []interface{}:
			for i, v := range val {
				paths = append(paths,
					getTerminalJSONPaths(t, v, append(crumbs, strconv.Itoa(i)))...,
				)
			}
		case map[string]interface{}:
			for k, v := range val {
				paths = append(paths,
					getTerminalJSONPaths(t, v, append(crumbs, k))...,
				)
			}
		default:
			t.Logf("found value %s\n", item)
		}

		slices.Sort(paths)
		return paths
 }
