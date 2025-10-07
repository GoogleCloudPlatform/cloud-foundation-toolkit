/**
 * Copyright 2025 Google LLC
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

package tft

import (
	"fmt"
	"maps"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

const (
	RootModuleName = "root"
)

var (
	rootModuleRegexp = regexp.MustCompile("^[./]*$")
	modulesDirRegexp = regexp.MustCompile("^[./]*/modules/(.+)$")
)

// findModulesUnderTest does a graph search starting from tfDir, looking
// for any transitively referenced modules that are considered "modules under test",
// which are modules in the "modules/" dir or the root module
// (the exact definition is in isModuleUnderTest() above).
//
// Any external modules encountered are ignored.
//
// The search doesn't continue to unpack any module under test, so for example, if
// modules/aaa sources modules/bbb, then only modules/aaa will be in the returned
// set (unless modules/bbb is reachable from tfDir without going through modules/aaa).
func findModulesUnderTest(tfDir string) (stringSet, error) {
	tfDir = filepath.Clean(tfDir)

	modulesUnderTest := make(stringSet)
	pathsToVisit := stringSet{tfDir: struct{}{}}
	seen := make(stringSet)

	maxIters := 5
	for iterations := 0; iterations < maxIters; iterations++ {
		// moduleRefs is a map from filesystem path to a set of modules sourced from there.
		moduleRefs, err := findAllReferencedModules(pathsToVisit)
		if err != nil {
			return nil, err
		}

		maps.Insert(seen, maps.All(pathsToVisit))

		for _, refs := range moduleRefs {
			for ref := range refs {
				if isModuleUnderTest(ref) {
					name, err := localModuleName(ref)
					if err != nil {
						return nil, err
					}
					modulesUnderTest[name] = struct{}{}
				}
			}
		}

		pathsToVisit = stripAlreadySeen(nextPathsToVisit(moduleRefs), seen)

		if len(pathsToVisit) == 0 {
			return modulesUnderTest, nil
		}
	}

	return nil, fmt.Errorf("exceeded %v iterations when searching for referenced modules starting from %q, pathsToVisit is currently %v", maxIters, tfDir, pathsToVisit)
}

// localModuleName takes a local module source string (e.g. ../../modules/my_module)
// and extracts the short module name used to refer to the module in test/setup/*.tf
// when per_module_{roles,services} are defined (e.g. my_module).
//
// This function returns an error if the input is not a local module source
// (e.g. "terraform-google-modules/network/google").
func localModuleName(moduleRef string) (string, error) {
	if rootModuleRegexp.MatchString(moduleRef) {
		return RootModuleName, nil
	}

	matches := modulesDirRegexp.FindStringSubmatch(moduleRef)
	if len(matches) < 2 {
		// modulesDirRegexp couldn't find a match even though this function is supposed to
		// be called with local module references only.
		return "", fmt.Errorf("couldn't extract module name from source %q using regexp %v", moduleRef, modulesDirRegexp)
	}
	return matches[1], nil
}

// isLocalModule uses terraform's definition of what a local module source looks like.
func isLocalModule(moduleRef string) bool {
	return strings.HasPrefix(moduleRef, "../") ||
		strings.HasPrefix(moduleRef, "./") ||
		strings.HasPrefix(moduleRef, "/")
}

// isModuleUnderTest looks at the given module source string and returns true
// if it is the root module of this repo, or one of the modules inside the "modules/" dir.
// This is done by just checking whether the moduleRef is a local filesystem path leading
// up to the root and then potentially into the "modules/" dir.
func isModuleUnderTest(moduleRef string) bool {
	return isLocalModule(moduleRef) && (rootModuleRegexp.MatchString(moduleRef) || modulesDirRegexp.MatchString(moduleRef))
}

type stringSet = map[string]struct{}

// nextPathsToVisit looks at the given mapping of
// (module path) => (modules referenced from that path) and returns a set of
// paths to visit next. External modules and "modules under test" are excluded
// as they do not need further examination.
func nextPathsToVisit(moduleRefs map[string]stringSet) stringSet {
	nextPaths := make(stringSet)
	for modulePath, refs := range moduleRefs {
		for ref := range refs {
			if isLocalModule(ref) && !isModuleUnderTest(ref) {
				nextPaths[filepath.Clean(filepath.Join(modulePath, ref))] = struct{}{}
			}
		}
	}
	return nextPaths
}

// stripAlreadySeen returns a new set that includes everything in "modulePaths"
// that is not in "seen".
func stripAlreadySeen(modulePaths stringSet, seen stringSet) stringSet {
	newPaths := make(stringSet)
	for path := range modulePaths {
		if _, ok := seen[path]; !ok {
			newPaths[path] = struct{}{}
		}
	}
	return newPaths
}

// findAllReferencedModules takes a set of filesystem paths for terraform
// modules and returns a map from the path to a set of all modules referenced
// from the module at that path.
func findAllReferencedModules(modulePaths stringSet) (map[string]stringSet, error) {
	moduleRefs := make(map[string]stringSet)
	for path := range modulePaths {
		modules, err := findReferencedModules(path)
		if err != nil {
			return nil, err
		}
		moduleRefs[path] = modules
	}
	return moduleRefs, nil
}

// findReferencedModules looks in tfDir and extracts the sources for all module
// blocks in that directory.
// The returned value is a set of module sources. Some possible examples:
//
//	"../.."
//	"../../modules/bar"
//	"terraform-google-modules/kubernetes-engine/google"
//	"terraform-google-modules/kubernetes-engine/google//modules/workload-identity"
//
// Overridden by tests.
var findReferencedModules = func(tfDir string) (stringSet, error) {
	mod, diags := tfconfig.LoadModule(tfDir)
	err := diags.Err()
	if err != nil {
		return nil, err
	}

	sources := make(stringSet)
	for _, moduleBlock := range mod.ModuleCalls {
		sources[moduleBlock.Source] = struct{}{}
	}
	return sources, nil
}
