package bptest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/discovery"
	testing "github.com/mitchellh/go-testing-interface"
)

const (
	discoverTestFilename = "discover_test.go"
)

type bpTest struct {
	name     string
	config   string
	location string
}

// getTests returns slice of all blueprint tests
func getTests(intTestDir string) ([]bpTest, error) {
	// discover intTestDir if not provided
	if intTestDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		discoveredIntTestDir, err := discoverIntTestDir(cwd)
		if err != nil {
			return nil, err
		}
		intTestDir = discoveredIntTestDir
	}
	Log.Info(fmt.Sprintf("using test-dir: %s", intTestDir))

	tests := []bpTest{}
	discoveredTests, err := getDiscoveredTests(intTestDir)
	if err != nil {
		return nil, err
	}
	tests = append(tests, discoveredTests...)

	explicitTests, err := getExplicitTests(intTestDir)
	if err != nil {
		return nil, err
	}
	tests = append(tests, explicitTests...)

	return tests, nil
}

// getDiscoveredTests returns slice of discovered blueprint tests
func getDiscoveredTests(intTestDir string) ([]bpTest, error) {
	discoverTestFile := path.Join(intTestDir, discoverTestFilename)
	// skip discovering tests if no discoverTestFile
	_, err := os.Stat(discoverTestFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		Log.Warn(fmt.Sprintf("Skipping discovered test. %s not found.", discoverTestFilename))
		return nil, nil
	}

	// if discoverTestFile is present, find auto discovered tests
	tests := []bpTest{}
	if discoverTestFile != "" {
		discoverTestName, err := getDiscoverTestName(discoverTestFile)
		if err != nil {
			return nil, err
		}
		discoveredSubTests := discovery.FindTestConfigs(&testing.RuntimeT{}, intTestDir)
		for testName, fileName := range discoveredSubTests {
			tests = append(tests, bpTest{name: fmt.Sprintf("%s/%s", discoverTestName, testName), config: fileName, location: discoverTestFile})
		}
	}
	sort.SliceStable(tests, func(i, j int) bool { return tests[i].name < tests[j].name })
	return tests, nil
}

func getExplicitTests(intTestDir string) ([]bpTest, error) {
	// find all explicit test files ending with *_test.go excluding discover_test.go within intTestDir
	testFiles := findFiles(intTestDir,
		func(d fs.DirEntry) bool {
			return strings.HasSuffix(d.Name(), "_test.go") && d.Name() != discoverTestFilename
		},
	)

	eTests := []bpTest{}
	for _, testFile := range testFiles {
		// testDir name maps to a matching example/fixture
		testDir := path.Dir(testFile)
		testCfg, err := discovery.GetConfigDirFromTestDir(testDir)
		if err != nil {
			Log.Warn(fmt.Sprintf("unable to discover configs for %s: %v", testDir, err))
		}

		testFns, err := getTestFuncsFromFile(testFile)
		if err != nil {
			return nil, err
		}
		for _, fnName := range testFns {
			eTests = append(eTests, bpTest{name: fnName, location: testFile, config: testCfg})
		}

	}
	sort.SliceStable(eTests, func(i, j int) bool { return eTests[i].name < eTests[j].name })
	return eTests, nil
}

// getDiscoverTestName returns test name used for auto discovered tests
func getDiscoverTestName(dFileName string) (string, error) {
	fn, err := getTestFuncsFromFile(dFileName)
	if err != nil {
		return "", err
	}
	// enforce only one main test func decl for discovered tests
	if len(fn) != 1 {
		return "", fmt.Errorf("only one function should be defined in %s. Found %+q", dFileName, fn)
	}
	return fn[0], nil
}

// discoverIntTestDir attempts to discover the integration test directory
// by searching for discover_test.go in the current working directory.
// If not found, it returns current working directory.
func discoverIntTestDir(cwd string) (string, error) {
	// search for discover_test.go
	discoverTestFiles := findFiles(cwd,
		func(d fs.DirEntry) bool {
			return d.Name() == discoverTestFilename
		},
	)
	if len(discoverTestFiles) > 1 {
		return "", fmt.Errorf("found multiple %s files: %+q. Exactly one file was expected", discoverTestFilename, discoverTestFiles)
	}
	if len(discoverTestFiles) == 1 {
		relIntTestDir, err := filepath.Rel(cwd, path.Dir(discoverTestFiles[0]))
		if err != nil {
			return "", err
		}
		return relIntTestDir, nil
	}
	// no discover_test.go file discovered
	return ".", nil
}

// findFiles returns a slice of file paths matching matchFn
func findFiles(dir string, matchFn func(d fs.DirEntry) bool) []string {
	files := []string{}
	filepath.WalkDir(dir, func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && matchFn(d) {
			files = append(files, fpath)
			return nil
		}
		return nil
	})
	return files
}
