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
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/mod/semver"
)

// BinaryInPath checks if a given binary is in path.
func BinaryInPath(bin string) error {
	if _, err := exec.LookPath(bin); err != nil {
		return err
	}
	return nil
}

// Return validated canonical KPT version string
func KptVersion(bin string) (string, error) {
	cmd, err := exec.Command(bin, "version").Output()
	kptVersion := "v" + strings.TrimSpace(string(cmd))
	if err != nil {
		return "", err
	}
	if semver.IsValid(kptVersion) != true {
		return "", fmt.Errorf("Unable to parse kpt version")
	}

	return semver.Canonical(kptVersion), nil
}
