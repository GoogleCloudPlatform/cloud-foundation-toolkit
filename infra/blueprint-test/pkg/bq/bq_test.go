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

// Package bq provides a set of helpers to interact with bq tool (part of CloudSDK)
package bq

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunf(t *testing.T) {
	tests := []struct {
		name            string
		cmd             string
		projectIdEnvVar string
	}{
		{
			name:            "Runf",
			cmd:             "query --nouse_legacy_sql 'select * FROM %s.samples.INFORMATION_SCHEMA.TABLES limit 1;'",
			projectIdEnvVar: "bigquery-public-data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if projectName, present := os.LookupEnv(tt.projectIdEnvVar); present {
				op := Runf(t, tt.cmd, projectName)
				assert := assert.New(t)
				assert.Contains(op.Array()[0], "creation_time")
			} else {
				t.Logf("Skipping test, %s envvar not set", tt.projectIdEnvVar)
				t.Skip()
			}
		})
	}
}
