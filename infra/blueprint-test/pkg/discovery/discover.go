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

// Package discovery attempts to discover test configs from well known directories.
package discovery

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-testing-interface"
)

// ConfigDirFromCWD attempts to autodiscover config for a given explicit test based on dirpath for the test.
func ConfigDirFromCWD(cwd string) (string, error) {
	name := path.Base(cwd)
	// check if fixture dir exists at ../../fixture/fixtureName
	fixturePath := fmt.Sprintf("../../fixture/%s", name)
	_, err := os.Stat(fixturePath)
	if err == nil {
		return fixturePath, nil
	}
	// check if example dir exists at ../../../examples/exampleName
	examplePath := fmt.Sprintf("../../../examples/%s", name)
	_, err = os.Stat(examplePath)
	if err == nil {
		return examplePath, nil
	}
	return "", fmt.Errorf("unable to find config in %s nor %s", fixturePath, examplePath)

}

// FindTestConfigs attempts to auto discover configs to test and expected to be executed from directory containing explicit integration tests.
// Order of discovery is all explicit tests, followed by all fixtures that do not have explicit tests, followed by all examples that do not have fixtures nor explicit tests.
func FindTestConfigs(t testing.TB, intTestDir string) []string {
	testBase := intTestDir
	examplesBase := path.Join(testBase, "../../examples")
	fixturesBase := path.Join(testBase, "../fixtures")
	explicitTests, err := findDirs(testBase)
	if err != nil {
		t.Logf("Error discovering explicit tests: %v", err)
	}
	fixtures, err := findDirs(fixturesBase)
	if err != nil {
		t.Logf("Error discovering fixtures: %v", err)
	}
	examples, err := findDirs(examplesBase)
	if err != nil {
		t.Logf("Error discovering examples: %v", err)
	}
	configsToRun := make([]string, 0)
	//TODO(bharathkkb): add overrides
	// if a fixture exists but no explicit test defined
	for n := range fixtures {
		_, ok := explicitTests[n]
		if !ok {
			configsToRun = append(configsToRun, path.Join(fixturesBase, n))
		}
	}
	// if an example exists that does not have a fixture nor explicit test defined
	for n := range examples {
		_, okTest := explicitTests[n]
		_, okFixture := fixtures[n]
		if !okTest && !okFixture {
			configsToRun = append(configsToRun, path.Join(examplesBase, n))
		}
	}
	return configsToRun

}

// getKnownDirInParents checks if a well known dir exists in parent or grandparent directory.
func GetKnownDirInParents(dir string) (string, error) {
	dirInParent := path.Join("..", dir)
	_, err := os.Stat(dirInParent)
	if !os.IsNotExist(err) {
		return dirInParent, err
	}
	dirInGrandparent := path.Join("..", dirInParent)
	_, err = os.Stat(dirInGrandparent)
	if !os.IsNotExist(err) {
		return dirInGrandparent, err
	}
	return "nil", fmt.Errorf("unable to find %s nor %s", dirInParent, dirInGrandparent)
}

// findDirs returns a map of directories in path
func findDirs(path string) (map[string]bool, error) {
	dirs := make(map[string]bool)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return dirs, err
	}
	for _, f := range files {
		if f.IsDir() {
			dirs[f.Name()] = true
		}
	}
	return dirs, nil
}
