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
	"os"
	"path"

	"github.com/mitchellh/go-testing-interface"
)

const (
	SetupDir    = "setup"    // known setup directory
	FixtureDir  = "fixtures" // known fixtures directory
	ExamplesDir = "examples" // known fixtures directory
)

// GetConfigDirFromTestDir attempts to autodiscover config for a given explicit test based on dirpath for the test.
func GetConfigDirFromTestDir(testDir string) (string, error) {
	name := path.Base(testDir)
	// check if fixture dir exists at ../../fixture/fixtureName
	fixturePath := path.Clean(path.Join(testDir, "../../", FixtureDir, name))
	_, err := os.Stat(fixturePath)
	if err == nil {
		return fixturePath, nil
	}
	// check if example dir exists at ../../../examples/exampleName
	examplePath := path.Clean(path.Join(testDir, "../../../", ExamplesDir, name))
	_, err = os.Stat(examplePath)
	if err == nil {
		return examplePath, nil
	}
	return "", fmt.Errorf("unable to find config in %s nor %s", fixturePath, examplePath)

}

// FindTestConfigs attempts to auto discover configs to test and is expected to be executed from a directory containing explicit integration tests.
// Order of discovery is all explicit tests, followed by all fixtures that do not have explicit tests, followed by all examples that do not have fixtures nor explicit tests.
func FindTestConfigs(t testing.TB, intTestDir string) map[string]string {
	testBase := intTestDir
	examplesBase := path.Join(testBase, "../../", ExamplesDir)
	fixturesBase := path.Join(testBase, "../", FixtureDir)
	explicitTests, err := findDirs(testBase)
	if err != nil {
		t.Logf("Skipping explicit tests discovery: %v", err)
	}
	fixtures, err := findDirs(fixturesBase)
	if err != nil {
		t.Logf("Skipping fixtures discovery: %v", err)
	}
	examples, err := findDirs(examplesBase)
	if err != nil {
		t.Logf("Skipping examples discovery: %v", err)
	}
	testCases := make(map[string]string)

	// if a fixture exists but no explicit test defined
	for n := range fixtures {
		_, ok := explicitTests[n]
		if !ok {
			testDir := path.Join(fixturesBase, n)
			testName := fmt.Sprintf("%s/%s", path.Base(path.Dir(testDir)), path.Base(testDir))
			testCases[testName] = testDir
		}
	}
	// if an example exists that does not have a fixture nor explicit test defined
	for n := range examples {
		_, okTest := explicitTests[n]
		_, okFixture := fixtures[n]
		if !okTest && !okFixture {
			testDir := path.Join(examplesBase, n)
			testName := fmt.Sprintf("%s/%s", path.Base(path.Dir(testDir)), path.Base(testDir))
			testCases[testName] = testDir
		}
	}
	// explicit tests in integration/test_name are not gathered since they are invoked directly
	return testCases

}

// GetKnownDirInParents checks if a well known dir exists in parent dir upto max parents.
func GetKnownDirInParents(dir string, max int) (string, error) {
	if max <= 0 {
		return "", fmt.Errorf("unable to find %s dir, searched upto %s", path.Base(dir), dir)
	}
	dirInParent := path.Join("..", dir)
	_, err := os.Stat(dirInParent)
	if !os.IsNotExist(err) {
		return dirInParent, err
	}
	return GetKnownDirInParents(dirInParent, max-1)
}

// findDirs returns a map of directories in path
func findDirs(path string) (map[string]bool, error) {
	dirs := make(map[string]bool)
	files, err := os.ReadDir(path)
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
