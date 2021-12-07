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

// Package discovery attempts to discover test configs from well known directories.
package discovery

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	baseDir     = "testdata/tf"
	testDir     = "test"
	intDir      = "integration"
	testFixture = "foo"
	testExample = "bar"
	testInvalid = "baz"
)

func TestGetConfigDirFromTestDir(t *testing.T) {
	tests := []struct {
		name    string
		testDir string
		want    string
		errMsg  string
	}{
		{
			name:    "With example and fixture",
			testDir: path.Join(baseDir, testDir, intDir, testFixture),
			want:    path.Join(baseDir, testDir, FixtureDir, testFixture),
		},
		{
			name:    "With example",
			testDir: path.Join(baseDir, testDir, intDir, testExample),
			want:    path.Join(baseDir, ExamplesDir, testExample),
		},
		{
			name:    "Invalid no fixture/example",
			testDir: path.Join(baseDir, testDir, intDir, testInvalid),
			errMsg:  "unable to find config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got, err := GetConfigDirFromTestDir(tt.testDir)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}
