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
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const (
	gfDir          = "testdata"
	gfPerms        = 0755
	gfUpdateEnvVar = "UPDATE_GOLDEN"
)

type GoldenFile struct {
	dir      string
	fileName string
	t        testing.TB
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

func NewOrUpdate(t testing.TB, data string, opts ...goldenFileOption) *GoldenFile {
	g := &GoldenFile{
		dir:      gfDir,
		fileName: fmt.Sprintf("%s.json", strings.ReplaceAll(t.Name(), "/", "-")),
		t:        t,
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
	err = ioutil.WriteFile(fp, []byte(data), gfPerms)
	if err != nil {
		g.t.Fatalf("error updating result: %v", err)
	}
}

// GetName return path of the goldenfile
func (g *GoldenFile) GetName() string {
	return path.Join(g.dir, g.fileName)
}

// GetJSON returns goldenfile as parsed json
func (g *GoldenFile) GetJSON() gjson.Result {
	return utils.LoadJSON(g.t, g.GetName())
}

// JSONEq asserts that json content in jsonPath for got and goldenfile is the same
func (g *GoldenFile) JSONEq(a *assert.Assertions, got gjson.Result, jsonPath string) {
	gf := g.GetJSON()
	a.Equal(gf.Get(jsonPath).String(), got.Get(jsonPath).String(), fmt.Sprintf("expected %s to match fixture", jsonPath))
}
