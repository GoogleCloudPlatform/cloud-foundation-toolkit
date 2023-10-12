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

func TestActivateCredsAndEnvVars(t *testing.T) {
	tests := []struct {
		name      string
		keyEnvVar string
		user      string
	}{
		{
			name:      "with sa key",
			keyEnvVar: "TEST_KEY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds, present := os.LookupEnv(tt.keyEnvVar)
			if !present {
				t.Logf("Skipping test, %s envvar not set", tt.keyEnvVar)
				t.Skip()
			}
			ActivateCredsAndEnvVars(t, creds)
			assert := assert.New(t)
			assert.Equal(os.Getenv("GOOGLE_CREDENTIALS"), creds)
			pathEnvVars := []string{"CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE", "GOOGLE_APPLICATION_CREDENTIALS"}
			for _, v := range pathEnvVars {
				c, err := os.ReadFile(os.Getenv(v))
				assert.NoError(err)
				assert.Equal(string(c), creds)
			}

		})
	}
}

func TestRunf(t *testing.T) {
	tests := []struct {
		name            string
		cmd             string
		projectIdEnvVar string
	}{
		{
			name:            "Runf",
			cmd:             "ls --datasets --project_id=%s",
			projectIdEnvVar: "bigquery-public-data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if projectName, present := os.LookupEnv(tt.projectIdEnvVar); present {
				op := Runf(t, tt.cmd, projectName)
				assert := assert.New(t)
				assert.Equal("bigquery#dataset", op.Array()[0].Get("kind").String())
			} else {
				t.Logf("Skipping test, %s envvar not set", tt.projectIdEnvVar)
				t.Skip()
			}
		})
	}
}
