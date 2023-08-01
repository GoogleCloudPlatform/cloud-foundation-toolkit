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

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
)

func TestAssertSuccess(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Hello World")
		}))
		defer ts.Close()

		r, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		it := &inspectableT{t, nil}
		ah := utils.NewAssertHTTP(utils.WithHTTPClient(ts.Client()))
		ah.AssertSuccess(it, r)

		if it.err != nil {
			t.Errorf("wanted success, got %v", it.err)
		}
	})
	t.Run("request error", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/nope", nil)
		if err != nil {
			t.Fatal(err)
		}
		it := &inspectableT{t, nil}
		utils.NewAssertHTTP().AssertSuccess(it, r)

		if it.err == nil {
			t.Error("wanted error, got success")
		}
	})
	t.Run("response error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not Available", http.StatusServiceUnavailable)
		}))
		defer ts.Close()

		r, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		it := &inspectableT{t, nil}
		ah := utils.NewAssertHTTP(utils.WithHTTPClient(ts.Client()))
		ah.AssertSuccess(it, r)

		if it.err != nil {
			t.Errorf("wanted error, got %v", it.err)
		}
	})
}

func TestAssertResponse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Hello World")
		}))
		defer ts.Close()

		r, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		it := &inspectableT{t, nil}
		ah := utils.NewAssertHTTP(utils.WithHTTPClient(ts.Client()))
		ah.AssertResponse(it, r, http.StatusOK)
		if it.err != nil {
			t.Errorf("wanted success, got %v", it.err)
		}
	})
	t.Run("request error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not Available", http.StatusServiceUnavailable)
		}))
		defer ts.Close()

		r, err := http.NewRequest(http.MethodGet, "/nope", nil)
		if err != nil {
			t.Fatal(err)
		}

		it := &inspectableT{t, nil}
		ah := utils.NewAssertHTTP(utils.WithHTTPClient(ts.Client()))
		ah.AssertResponse(it, r, http.StatusServiceUnavailable)

		if it.err == nil {
			t.Error("wanted error, got success")
		}
	})
	t.Run("response error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not Available", http.StatusServiceUnavailable)
		}))
		defer ts.Close()

		r, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		it := &inspectableT{t, nil}
		ah := utils.NewAssertHTTP(utils.WithHTTPClient(ts.Client()))
		ah.AssertResponse(it, r, http.StatusServiceUnavailable)

		if it.err != nil {
			t.Errorf("wanted error, got %v", it.err)
		}
	})
	t.Run("response contains", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello World")
		}))
		defer ts.Close()

		r, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		t.Run("success", func(t *testing.T) {
			it := &inspectableT{t, nil}
			ah := utils.NewAssertHTTP(utils.WithHTTPClient(ts.Client()))
			ah.AssertResponse(it, r, http.StatusOK, "Hello", "World")
			if it.err != nil {
				t.Errorf("wanted success, got %v", it.err)
			}
		})
		t.Run("error", func(t *testing.T) {
			it := &inspectableT{t, nil}
			ah := utils.NewAssertHTTP(utils.WithHTTPClient(ts.Client()))
			ah.AssertResponse(it, r, http.StatusOK, "Hello", "Moon")
			if it.err == nil {
				t.Error("wanted error, got success")
			}
		})
		t.Run("error multiple", func(t *testing.T) {
			it := &inspectableT{t, nil}
			ah := utils.NewAssertHTTP(utils.WithHTTPClient(ts.Client()))
			ah.AssertResponse(it, r, http.StatusOK, "Hello", "Moon", "People")
			if it.err == nil {
				t.Error("wanted error, got success")
				return
			}
			if e := it.err.Error(); !strings.Contains(e, "Moon") || !strings.Contains(e, "People") {
				t.Errorf("wanted multiple errors, got one: %v", it.err)
			}
		})
	})
}
