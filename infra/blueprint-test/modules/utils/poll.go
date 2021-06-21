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
	"time"

	"github.com/mitchellh/go-testing-interface"
)

// Polls on a particular condition function while the returns true.
// It fails the test if the condition is not met within numRetries.
func Poll(t testing.TB, condition func() (bool, error), numRetries int, interval time.Duration) {
	if numRetries < 0 {
		t.Fatal("invalid value for numRetries. Must be >= 0")
	}

	if interval <= 0 {
		t.Fatal("invalid value for numRetries. Must be > 0")
	}

	retry, err := condition()

	for count := 0; retry && count <= numRetries; count++ {
		time.Sleep(interval)
		if err != nil {
			t.Logf("Received error while polling: %v", err)
		}
		t.Logf("Retrying... %d", count+1)
		retry, err = condition()
	}

	if err != nil {
		t.Fatalf("failed to pull provided condition after %d retries, last error: %v", numRetries, err)
	}

	if retry {
		t.Fatalf("polling timed out after %d retries with %d second intervals", numRetries, interval/time.Second)
	}
}
