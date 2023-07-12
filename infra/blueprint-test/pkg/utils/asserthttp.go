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
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/go-testing-interface"
)

// AssertHTTP provides a collection of HTTP asserts.
type AssertHTTP struct {
	httpClient *http.Client
}

type assertOption func(*AssertHTTP)

// WithHTTPClient specifies an HTTP client for the AssertHTTP use.
func WithHTTPClient(c *http.Client) assertOption {
	return func(ah *AssertHTTP) {
		ah.httpClient = c
	}
}

// NewAssertHTTP creates a new AssertHTTP with option overrides.
func NewAssertHTTP(opts ...assertOption) *AssertHTTP {
	ah := &AssertHTTP{http.DefaultClient}
	for _, opt := range opts {
		opt(ah)
	}
	return ah
}

// AssertSuccessWithRetry runs httpRequest and retries on errors outside client control.
func (ah *AssertHTTP) AssertSuccessWithRetry(t testing.TB, r *http.Request) {
	t.Helper()
	Poll(t, ah.httpRequest(t, r), 3, 2*time.Second)
}

// AssertSuccess runs httpRequest without retry.
func (ah *AssertHTTP) AssertSuccess(t testing.TB, r *http.Request) {
	t.Helper()
	_, err := ah.httpRequest(t, r)()
	if err != nil {
		t.Fatal(err)
	}
}

// AssertResponseWithRetry runs httpResponse and retries on errors outside client control.
func (ah *AssertHTTP) AssertResponseWithRetry(t testing.TB, r *http.Request, wantCode int, want ...string) {
	t.Helper()
	Poll(t, ah.httpResponse(t, r, wantCode, want...), 3, 2*time.Second)
}

// AssertResponse runs httpResponse without retry.
func (ah *AssertHTTP) AssertResponse(t testing.TB, r *http.Request, wantCode int, want ...string) {
	t.Helper()
	_, err := ah.httpResponse(t, r, wantCode, want...)()
	if err != nil {
		t.Fatal(err)
	}
}

// httpRequest verifies the request is successful by HTTP status code.
func (ah *AssertHTTP) httpRequest(t testing.TB, r *http.Request) func() (bool, error) {
	t.Helper()
	return func() (bool, error) {
		t.Logf("Sending HTTP Request %s %s", r.Method, r.URL.String())
		got, err := ah.httpClient.Do(r)
		if err != nil {
			return false, err
		}
		// Keep trying until the result is success or the request responsibility.
		if retry := httpRetryCondition(got.StatusCode); retry {
			return true, nil
		}
		// Any HTTP success will work. For a specific status use AssertResponse.
		if got.StatusCode < http.StatusOK || got.StatusCode >= http.StatusMultipleChoices {
			t.Errorf("want 2xx, got %d", got.StatusCode)
		}

		return false, nil
	}
}

// httpResponse verifies the requested response has the wanted status code and payload.
func (ah *AssertHTTP) httpResponse(t testing.TB, r *http.Request, wantCode int, want ...string) func() (bool, error) {
	t.Helper()
	return func() (bool, error) {
		t.Logf("Sending HTTP Request %s %s", r.Method, r.URL.String())
		got, err := ah.httpClient.Do(r)
		if err != nil {
			return false, err
		}
		defer got.Body.Close()

		if got.StatusCode != wantCode {
			t.Errorf("response code: got %d, want %d", got.StatusCode, wantCode)
			// Unwanted status code may be a server-side error condition that will clear.
			// Assume unwanted success is not going to change.
			return httpRetryCondition(got.StatusCode), nil
		}

		b, err := io.ReadAll(got.Body)
		if err != nil {
			return true, err
		}
		out := string(b)

		atLeastOneError := false
		for _, fragment := range want {
			if !strings.Contains(out, fragment) {
				t.Errorf("response body: want contained %q", fragment)
				atLeastOneError = true
			}
		}

		// Only output received HTTP response body once.
		if atLeastOneError {
			t.Log("response output:")
			t.Log(out)
		}

		return false, nil
	}
}

// httpRetryCondition indicates retry should be attempted on HTTP 1xx, 401, 403, and 5xx errors.
// 401 and 403 are retried in case of lagging authorization configuration.
// On true return value a retry is preferred.
func httpRetryCondition(code int) bool {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return false
	case code < http.StatusOK:
		return false
	case code >= http.StatusInternalServerError:
		return true
	case code == http.StatusUnauthorized || code == http.StatusForbidden:
		return true
	case code >= http.StatusBadRequest:
		return false
	}

	return false
}
