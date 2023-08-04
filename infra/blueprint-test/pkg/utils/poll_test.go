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

package utils_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
)

func TestPoll(t *testing.T) {
	testcases := []struct {
		label     string
		condition func() (bool, error)
		want      string
	}{
		{
			label: "error",
			condition: func() (bool, error) {
				return true, errors.New("condition failure")
			},
			want: "condition failure",
		},
		{
			label: "timeout",
			condition: func() (bool, error) {
				return true, nil
			},
			want: "polling timed out",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.label, func(t *testing.T) {
			it := &inspectableT{t, nil}
			utils.Poll(it, tc.condition, 3, time.Millisecond)
			if !strings.Contains(it.err.Error(), tc.want) {
				t.Errorf("got %v, want %v", it.err, tc.want)
			}
		})
	}
}

func TestPollE(t *testing.T) {
	testcases := []struct {
		label     string
		condition func() (bool, error)
		want      string
	}{
		{
			label: "error",
			condition: func() (bool, error) {
				return true, errors.New("condition failure")
			},
			want: "condition failure",
		},
		{
			label: "timeout",
			condition: func() (bool, error) {
				return true, nil
			},
			want: "polling timed out",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.label, func(t *testing.T) {
			it := &inspectableT{t, nil}
			err := utils.PollE(it, tc.condition, 3, time.Millisecond)
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("got %v, want %v", it.err, tc.want)
			}
		})
	}
}
