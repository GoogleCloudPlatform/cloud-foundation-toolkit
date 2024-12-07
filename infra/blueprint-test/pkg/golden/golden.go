/**
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless assertd by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package golden helps manage goldenfiles.
package golden

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
)

const (
	gfDir          = "testdata"
	gfPerms        = 0755
	gfUpdateEnvVar = "UPDATE_GOLDEN"
	gGoroutinesMax = 24
)

type GoldenFile struct {
	dir        string
	fileName   string
	sanitizers []Sanitizer
	t          testing.TB
}

type Sanitizer func(string) string

// StringSanitizer replaces all occurrences of old string with new string
func StringSanitizer(old, new string) Sanitizer {
	return func(s string) string {
		return strings.ReplaceAll(s, old, new)
	}
}

// ProjectIDSanitizer replaces all occurrences of current gcloud project ID with PROJECT_ID string
func ProjectIDSanitizer(t testing.TB) Sanitizer {
	return func(s string) string {
		projectID := gcloud.Run(t, "config get-value project")
		if projectID.String() == "[]" {
			t.Logf("no project ID currently set, skipping ProjectIDSanitizer: %s", projectID.String())
			return s
		}
		return strings.ReplaceAll(s, projectID.String(), "PROJECT_ID")
	}
}

type goldenFileOption func(*GoldenFile)

func WithDir(dir string) goldenFileOption {
	return func(g *GoldenFile) {
		g.dir = dir
	}
}

func WithFileName(fn string) goldenFileOption {
	return func(g *GoldenFile) {
		g.fileName = fn
	}
}

func WithSanitizer(s Sanitizer) goldenFileOption {
	return func(g *GoldenFile) {
		g.sanitizers = append(g.sanitizers, s)
	}
}

func WithStringSanitizer(old, new string) goldenFileOption {
	return func(g *GoldenFile) {
		g.sanitizers = append(g.sanitizers, StringSanitizer(old, new))
	}
}

func NewOrUpdate(t testing.TB, data string, opts ...goldenFileOption) *GoldenFile {
	g := &GoldenFile{
		dir:        gfDir,
		fileName:   fmt.Sprintf("%s.json", strings.ReplaceAll(t.Name(), "/", "-")),
		sanitizers: []Sanitizer{ProjectIDSanitizer(t)},
		t:          t,
	}
	for _, opt := range opts {
		opt(g)
	}
	g.update(data)
	return g
}

// update updates goldenfile data iff gfUpdateEnvVar is true
func (g *GoldenFile) update(data string) {
	// exit early if gfUpdateEnvVar is not set or true
	if strings.ToLower(os.Getenv(gfUpdateEnvVar)) != "true" {
		return
	}
	fp := g.GetName()
	err := os.MkdirAll(path.Dir(fp), gfPerms)
	if err != nil {
		g.t.Fatalf("error updating result: %v", err)
	}
	// apply sanitizers on data
	data = g.ApplySanitizers(data)

	err = os.WriteFile(fp, []byte(data), gfPerms)
	if err != nil {
		g.t.Fatalf("error updating result: %v", err)
	}
}

// GetName return path of the goldenfile
func (g *GoldenFile) GetName() string {
	return path.Join(g.dir, g.fileName)
}

// ApplySanitizers returns sanitized string
func (g *GoldenFile) ApplySanitizers(s string) string {
	for _, sanitizer := range g.sanitizers {
		s = sanitizer(s)
	}
	return s
}

// GetSanitizedJSON returns sanitizes and returns JSON result
func (g *GoldenFile) GetSanitizedJSON(s gjson.Result) gjson.Result {
	resultStr := s.String()
	resultStr = g.ApplySanitizers(resultStr)
	return utils.ParseJSONResult(g.t, resultStr)
}

// GetJSON returns goldenfile as parsed json
func (g *GoldenFile) GetJSON() gjson.Result {
	return utils.LoadJSON(g.t, g.GetName())
}

// JSONEq asserts that json content in jsonPath for got and goldenfile is the same
func (g *GoldenFile) JSONEq(a *assert.Assertions, got gjson.Result, jsonPath string) {
	gf := g.GetJSON()
	getPath := fmt.Sprintf("%s|@ugly", jsonPath)
	gotData := g.ApplySanitizers(got.Get(getPath).String())
	gfData := gf.Get(getPath).String()
	a.Equalf(gfData, gotData, "For path %q expected %q to match fixture %q", jsonPath, gotData, gfData)
}

// JSONPathEqs asserts that json content in jsonPaths for got and goldenfile are the same
func (g *GoldenFile) JSONPathEqs(a *assert.Assertions, got gjson.Result, jsonPaths []string) {
	syncGroup := new(errgroup.Group)
	syncGroup.SetLimit(gGoroutinesMax)
	g.t.Logf("Checking %d JSON paths with max %d goroutines", len(jsonPaths), gGoroutinesMax)
	for _, jsonPath := range jsonPaths {
		jsonPath := jsonPath
		syncGroup.Go(func() error {
			g.JSONEq(a, got, jsonPath)
			return nil
		})
	}
	if err := syncGroup.Wait(); err != nil {
		g.t.Fatal(err)
	}
}
