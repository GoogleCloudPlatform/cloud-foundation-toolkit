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
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestFileLogger(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "simple",
			content: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			testFile := path.Join(t.TempDir(), fmt.Sprintf("%s.log", tt.name))
			fl, flClose := NewTestFileLogger(t, testFile)
			fl.Logf(t, tt.content)
			flClose(t)
			gotContent, err := os.ReadFile(testFile)
			assert.NoError(err)
			assert.Contains(string(gotContent), "foo")
			// assert we are wrapping logger.DoLog which prints stack/test info
			assert.Contains(string(gotContent), fmt.Sprintf("TestNewTestFileLogger/%s", tt.name))
			assert.Contains(string(gotContent), "logger_test.go")
		})
	}
}
