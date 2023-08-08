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
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/go-testing-interface"
)

// AssertHTTP provides a collection of HTTP asserts.
type AssertHTTP struct {
	httpClient    *http.Client
	retryCount    int
	retryInterval time.Duration
}

type assertOption func(*AssertHTTP)

// WithHTTPClient specifies an HTTP client for the AssertHTTP use.
func WithHTTPClient(c *http.Client) assertOption {
	return func(ah *AssertHTTP) {
		ah.httpClient = c
	}
}

// WithHTTPRequestRetries specifies a HTTP request retry policy.
func WithHTTPRequestRetries(count int, interval time.Duration) assertOption {
	return func(ah *AssertHTTP) {
		ah.retryCount = count
		ah.retryInterval = interval
	}
}

// NewAssertHTTP creates a new AssertHTTP with option overrides.
func NewAssertHTTP(opts ...assertOption) *AssertHTTP {
	ah := &AssertHTTP{
		httpClient:    http.DefaultClient,
		retryCount:    3,
		retryInterval: 2 * time.Second,
	}
	for _, opt := range opts {
		opt(ah)
	}
	return ah
}

// AssertSuccessWithRetry runs httpRequest and retries on errors outside client control.
func (ah *AssertHTTP) AssertSuccessWithRetry(t testing.TB, r *http.Request) {
	t.Helper()
	if ah.retryCount == 0 || ah.retryInterval == 0 {
		ah.AssertSuccess(t, r)
		return
	}

	err := PollE(t, ah.httpRequest(t, r), ah.retryCount, ah.retryInterval)
	if err != nil {
		t.Error(err.Error())
	}
}

// AssertSuccess runs httpRequest without retry.
func (ah *AssertHTTP) AssertSuccess(t testing.TB, r *http.Request) {
	t.Helper()
	_, err := ah.httpRequest(t, r)()
	if err != nil {
		t.Error(err)
	}
}

// AssertResponseWithRetry runs httpResponse and retries on errors outside client control.
func (ah *AssertHTTP) AssertResponseWithRetry(t testing.TB, r *http.Request, wantCode int, want ...string) {
	t.Helper()
	if ah.retryCount == 0 || ah.retryInterval == 0 {
		ah.AssertSuccess(t, r)
		return
	}

	err := PollE(t, ah.httpResponse(t, r, wantCode, want...), ah.retryCount, ah.retryInterval)
	if err != nil {
		t.Error(err.Error())
	}
}

// AssertResponse runs httpResponse without retry.
func (ah *AssertHTTP) AssertResponse(t testing.TB, r *http.Request, wantCode int, want ...string) {
	t.Helper()
	_, err := ah.httpResponse(t, r, wantCode, want...)()
	if err != nil {
		t.Error(err)
	}
}

// httpRequest verifies the request is successful by HTTP status code.
func (ah *AssertHTTP) httpRequest(t testing.TB, r *http.Request) func() (bool, error) {
	t.Helper()
	logger := GetLoggerFromT()

	return func() (bool, error) {
		logger.Logf(t, "Sending HTTP Request %s %s", r.Method, r.URL.String())
		got, err := ah.httpClient.Do(r)
		if err != nil {
			return false, err
		}
		// Keep trying until the result is success or the request responsibility.
		ok, retry := httpRetryCondition(got.StatusCode)
		if !ok {
			return retry, fmt.Errorf("want 2xx, got %d", got.StatusCode)
		}
		logger.Logf(t, "Successful HTTP Request %s %s", r.Method, r.URL.String())

		return retry, nil
	}
}

// httpResponse verifies the requested response has the wanted status code and payload.
func (ah *AssertHTTP) httpResponse(t testing.TB, r *http.Request, wantCode int, want ...string) func() (bool, error) {
	t.Helper()
	logger := GetLoggerFromT()

	return func() (bool, error) {
		t.Logf("Sending HTTP Request %s %s", r.Method, r.URL.String())
		got, err := ah.httpClient.Do(r)
		if err != nil {
			return false, err
		}
		defer got.Body.Close()

		// Determine if the request is successful, and if the response indicates
		// we should attempt a retry.
		ok, retry := httpRetryCondition(got.StatusCode)
		if ok {
			logger.Logf(t, "Successful HTTP Request %s %s", r.Method, r.URL.String())
		}

		// e is the wrapped error for all expectation mismatches.
		var e error
		if got.StatusCode != wantCode {
			e = errors.Join(e, fmt.Errorf("response code: got %d, want %d", got.StatusCode, wantCode))
		}

		// No further processing required.
		if len(want) == 0 {
			return false, e
		}

		b, err := io.ReadAll(got.Body)
		if err != nil {
			return retry, errors.Join(e, err)
		}

		if len(b) == 0 {
			return retry, errors.Join(e, errors.New("empty response body"))
		}

		out := string(b)
		var bodyErr error
		for _, fragment := range want {
			if !strings.Contains(out, fragment) {
				bodyErr = errors.Join(bodyErr, fmt.Errorf("response body does not contain %q", fragment))
			}
		}

		// Only log errors and response body once.
		if bodyErr != nil {
			logger.Logf(t, "response output:")
			logger.Logf(t, strings.TrimSpace(out))
			return retry, errors.Join(e, bodyErr)
		}

		return retry, e
	}
}

// httpRetryCondition indicates retry should be attempted on HTTP 1xx, 401, 403, and 5xx errors.
// 401 and 403 are retried in case of lagging authorization configuration.
// First return value indicates successful response.
// Second return value, on true a retry is preferred.
func httpRetryCondition(code int) (bool, bool) {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return true, false
	case code < http.StatusOK:
		return false, false
	case code >= http.StatusInternalServerError:
		return false, true
	// IAM & network configuration propagation is a source of delayed access.
	case code == http.StatusUnauthorized || code == http.StatusForbidden:
		return false, true
	case code >= http.StatusBadRequest:
		return false, false
	}

	return false, false
}
