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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
)

func TestAssertSuccess(t *testing.T) {
	tests := []struct {
		label           string
		serverFunc      func(t *testing.T) *httptest.Server
		requestFunc     func(t *testing.T, s *httptest.Server) *http.Request
		assertFunc      func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request)
		assertRetryFunc func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request)
	}{
		{
			label: "success",
			serverFunc: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "Hello World")
				}))
			},
			requestFunc: func(t *testing.T, s *httptest.Server) *http.Request {
				r, err := http.NewRequest(http.MethodGet, s.URL, nil)
				if err != nil {
					t.Fatal(err)
				}
				return r
			},
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertSuccess(it, r)
				if it.err != nil {
					t.Errorf("wanted success, got %v", it.err)
				}
			},
			assertRetryFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertSuccessWithRetry(it, r)
				if it.err != nil {
					t.Errorf("wanted success, got %v", it.err)
				}
			},
		},
		{
			label: "request failure",
			serverFunc: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Not Available", http.StatusServiceUnavailable)
				}))
			},
			requestFunc: func(t *testing.T, s *httptest.Server) *http.Request {
				r, err := http.NewRequest(http.MethodGet, "/nope", nil)
				if err != nil {
					t.Fatal(err)
				}
				return r
			},
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertSuccess(it, r)
				if it.err == nil {
					t.Error("wanted error, got success")
				}
			},
			assertRetryFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertSuccessWithRetry(it, r)
				if it.err == nil {
					t.Error("wanted error, got success")
				}
			},
		},
		{
			label: "response error",
			serverFunc: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Not Available", http.StatusServiceUnavailable)
				}))
			},
			requestFunc: func(t *testing.T, s *httptest.Server) *http.Request {
				r, err := http.NewRequest(http.MethodGet, s.URL, nil)
				if err != nil {
					t.Fatal(err)
				}
				return r
			},
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertSuccess(it, r)
				if it.err == nil {
					t.Errorf("wanted error, got success")
				}
			},
			assertRetryFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertSuccessWithRetry(it, r)
				if it.err == nil {
					t.Error("wanted error, got success")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			server := tc.serverFunc(t)
			defer server.Close()
			t.Run("default", func(t *testing.T) {
				it := &inspectableT{t, nil}
				ah := utils.NewAssertHTTP(utils.WithHTTPClient(server.Client()))
				tc.assertFunc(t, it, ah, tc.requestFunc(t, server))
			})
			t.Run("retry", func(t *testing.T) {
				it := &inspectableT{t, nil}
				ah := utils.NewAssertHTTP(
					utils.WithHTTPClient(server.Client()),
					utils.WithHTTPRequestRetries(3, time.Millisecond),
				)
				tc.assertRetryFunc(t, it, ah, tc.requestFunc(t, server))
			})
		})
	}
}

func TestAssertResponse(t *testing.T) {
	tests := []struct {
		label           string
		serverFunc      func(t *testing.T) *httptest.Server
		requestFunc     func(t *testing.T, s *httptest.Server) *http.Request
		assertFunc      func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request)
		assertRetryFunc func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request)
	}{
		{
			label: "success",
			serverFunc: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "Hello World")
				}))
			},
			requestFunc: func(t *testing.T, s *httptest.Server) *http.Request {
				r, err := http.NewRequest(http.MethodGet, s.URL, nil)
				if err != nil {
					t.Fatal(err)
				}
				return r
			},
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertResponse(it, r, http.StatusOK)
				if it.err != nil {
					t.Errorf("wanted success, got %v", it.err)
				}
			},
			assertRetryFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertResponseWithRetry(it, r, http.StatusOK)
				if it.err != nil {
					t.Errorf("wanted success, got %v", it.err)
				}
			},
		},
		{
			label: "request failure",
			serverFunc: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Not Available", http.StatusServiceUnavailable)
				}))
			},
			requestFunc: func(t *testing.T, s *httptest.Server) *http.Request {
				r, err := http.NewRequest(http.MethodGet, "/nope", nil)
				if err != nil {
					t.Fatal(err)
				}
				return r
			},
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertResponse(it, r, http.StatusServiceUnavailable)
				if it.err == nil {
					t.Error("got success, wanted error")
				}
			},
			assertRetryFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertResponseWithRetry(it, r, http.StatusServiceUnavailable)
				if it.err == nil {
					t.Error("got success, wanted error")
				}
			},
		},
		{
			label: "assert HTTP 503",
			serverFunc: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Not Available", http.StatusServiceUnavailable)
				}))
			},
			requestFunc: func(t *testing.T, s *httptest.Server) *http.Request {
				r, err := http.NewRequest(http.MethodGet, s.URL, nil)
				if err != nil {
					t.Fatal(err)
				}
				return r
			},
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertResponse(it, r, http.StatusServiceUnavailable)
				if it.err != nil {
					t.Errorf("got %v, wanted success", it.err)
				}
			},
			assertRetryFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertResponseWithRetry(it, r, http.StatusServiceUnavailable)
				if it.err != nil {
					t.Errorf("got %v, wanted success", it.err)
				}
			},
		},
		{
			label: "response error",
			serverFunc: func(t *testing.T) *httptest.Server {
				var n int = 0
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					n++
					http.Error(w, fmt.Sprintf("Not Available #%d", n), http.StatusServiceUnavailable)
				}))
			},
			requestFunc: func(t *testing.T, s *httptest.Server) *http.Request {
				r, err := http.NewRequest(http.MethodGet, s.URL, nil)
				if err != nil {
					t.Fatal(err)
				}
				return r
			},
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				// The single assert is the first request to the test server.
				ah.AssertResponse(it, r, http.StatusOK, "#1")
				if it.err != nil && !strings.Contains(it.err.Error(), "got 503, want 200") {
					t.Error(it.err.Error())
				} else if it.err == nil {
					t.Error("got success, want error")
				}
			},
			assertRetryFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				// This function is called given the AssertHTTP object is configured for 3 retries.
				// The final request count for three retries is 5:
				// - the first request is not a retry and counts as #1
				// - utils.Poll retries n+1 times
				// The number is included in this assertion to confirm the error
				// is associated with the last retry attempt.
				ah.AssertResponseWithRetry(it, r, http.StatusOK, "#5")
				if it.err != nil && !strings.Contains(it.err.Error(), "got 503, want 200") {
					t.Error(it.err.Error())
				} else if it.err == nil {
					t.Error("got success, want error")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			t.Run("default", func(t *testing.T) {
				// Unlike other test implementations in this file, server is
				// instantiated per sub-case. This ensure the specific count of
				// requests matches configured retry expectations.
				server := tc.serverFunc(t)
				defer server.Close()
				it := &inspectableT{t, nil}
				ah := utils.NewAssertHTTP(utils.WithHTTPClient(server.Client()))
				tc.assertFunc(t, it, ah, tc.requestFunc(t, server))
			})
			t.Run("retry", func(t *testing.T) {
				server := tc.serverFunc(t)
				defer server.Close()
				it := &inspectableT{t, nil}
				ah := utils.NewAssertHTTP(
					utils.WithHTTPClient(server.Client()),
					utils.WithHTTPRequestRetries(3, time.Millisecond),
				)
				tc.assertRetryFunc(t, it, ah, tc.requestFunc(t, server))
			})
		})
	}
}

func TestAssertResponse_contains(t *testing.T) {
	tests := []struct {
		label      string
		assertFunc func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request)
	}{
		{
			label: "success",
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertResponse(it, r, http.StatusOK, "Hello", "World")
				if it.err != nil {
					t.Errorf("wanted success, got %v", it.err)
				}
			},
		},
		{
			label: "error",
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertResponse(it, r, http.StatusOK, "Hello", "Moon")
				if it.err == nil {
					t.Error("wanted error, got success")
				}
			},
		},
		{
			label: "multi error",
			assertFunc: func(t *testing.T, it *inspectableT, ah *utils.AssertHTTP, r *http.Request) {
				ah.AssertResponse(it, r, http.StatusOK, "Hello", "Moon", "Dwellers")
				if it.err == nil {
					t.Error("wanted error, got success")
					return
				}
				if e := it.err.Error(); !strings.Contains(e, "Moon") || !strings.Contains(e, "Dwellers") {
					t.Errorf("wanted multiple errors, got one: %v", it.err)
				}
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	}))
	defer server.Close()

	r, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			it := &inspectableT{t, nil}
			ah := utils.NewAssertHTTP(utils.WithHTTPClient(server.Client()))
			tc.assertFunc(t, it, ah, r)
		})
	}
}
