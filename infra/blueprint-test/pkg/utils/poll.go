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
	"time"

	"github.com/mitchellh/go-testing-interface"
)

// Polls on a particular condition function while the returns true.
// It fails the test if the condition is not met within numRetries.
func Poll(t testing.TB, condition func() (bool, error), numRetries int, interval time.Duration) {
	err := PollE(t, condition, numRetries, interval)
	if err != nil {
		t.Fatalf("failed to pull provided condition after %d retries, last error: %v", numRetries, err)
	}
}

// Polls on a particular condition function while the returns true.
// Returns an error if the condition is not met within numRetries.
func PollE(t testing.TB, condition func() (bool, error), numRetries int, interval time.Duration) error {
	if numRetries < 0 {
		return &PollParameterError{"invalid value for numRetries. Must be >= 0"}
	}

	if interval <= 0 {
		return &PollParameterError{"invalid value for interval. Must be > 0"}
	}

	retry, err := condition()

	for count := 0; retry && count <= numRetries; count++ {
		time.Sleep(interval)
		if err != nil {
			GetLoggerFromT().Logf(t, "Received error while polling: %v", err)
		}
		GetLoggerFromT().Logf(t, "Retrying... %d", count+1)
		retry, err = condition()
	}

	if err != nil {
		return &PollConditionError{err: err, numRetries: numRetries}
	}

	if retry {
		return &PollRetryLimitExceededError{interval: interval, numRetries: numRetries}
	}

	return nil
}

// PollParameterError is returend by PollE when input parameters are invalid.
type PollParameterError struct {
	msg string
}

func (e *PollParameterError) Error() string {
	return e.msg
}

// PollRetryLimitExceededError is returned by PollE when retries exceed numRetries.
type PollRetryLimitExceededError struct {
	numRetries int
	interval   time.Duration
}

func (e *PollRetryLimitExceededError) Error() string {
	return fmt.Sprintf("polling timed out after %d retries with %.2f second intervals", e.numRetries, e.interval.Seconds())
}

// PollConditionError is an error returned on the final PollE attempt.
type PollConditionError struct {
	err        error
	numRetries int
}

func (e *PollConditionError) Error() string {
	return fmt.Sprintf("failed to pull provided condition after %d retries, last error: %v", e.numRetries, e.err)
}

func (e *PollConditionError) Unwrap() error {
	return e.err
}
