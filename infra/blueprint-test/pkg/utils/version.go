/**
 * Copyright 2023 Google LLC
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
	"fmt"

	"golang.org/x/mod/semver"
)

// MinSemver validates gotSemver is not less than minSemver
func MinSemver(gotSemver string, minSemver string) error {
	if !semver.IsValid(gotSemver) {
		return fmt.Errorf("unable to parse got version %q", gotSemver)
	} else if !semver.IsValid(minSemver) {
		return fmt.Errorf("unable to parse minimum version %q", minSemver)
	}
	if semver.Compare(gotSemver, minSemver) == -1 {
		return fmt.Errorf("got version %q is less than minimum version %q", gotSemver, minSemver)
	}

	return nil
}
