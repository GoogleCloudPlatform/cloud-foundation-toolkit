// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

// Generates a random ID containing [2-7][a-z] (base32 alphabets) of length 4.
func GenerateRandomizedSuffix() (string, error) {
	// 3 bytes will generate base32 encoded string of length 5.
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error generating random bytes: %v", err)
	}

	return strings.ToLower(base32.StdEncoding.EncodeToString(b)[0:4]), nil
}