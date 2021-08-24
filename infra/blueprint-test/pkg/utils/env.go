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

	"github.com/mitchellh/go-testing-interface"
)

// ValFromEnv returns value for a given env var.
// It fails test if not set.
func ValFromEnv(t testing.TB, k string) string {
	v, found := os.LookupEnv(k)
	if !found {
		t.Fatalf("envvar %s not set", k)
	}
	return v
}
