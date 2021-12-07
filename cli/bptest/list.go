package bptest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
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
	discoverTestFile, err := getDiscoverTestFile(intTestDir)
	// skip discovering tests if no discoverTestFile
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		Log.Warn(fmt.Sprintf("Skipping discovered test. %s not found.", discoverTestFilename))
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

func getExplicitTests(path string) ([]bpTest, error) {
	// find dirs within test/integration
	dirs, err := findDirContentsFilter(path, func(f fs.FileInfo) bool { return f.IsDir() })
	if err != nil {
		return nil, err
	}
	// find test files within each dir
	eTests := []bpTest{}
	filterTestFiles := func(f fs.FileInfo) bool {
		return !f.IsDir() && strings.HasSuffix(f.Name(), "_test.go")
	}
	// each dir within test/integration could contain explict tests
	for _, dir := range dirs {
		testCfg, err := discovery.GetConfigDirFromTestDir(dir)
		if err != nil {
			Log.Warn(fmt.Sprintf("unable to discover configs for %s: %v", dir, err))
		}
		// find all test files in an explicit test dir that ends with _test.go
		files, err := findDirContentsFilter(dir, filterTestFiles)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			testFns, err := getTestFuncsFromFile(f)
			if err != nil {
				return nil, err
			}
			for _, fnName := range testFns {
				eTests = append(eTests, bpTest{name: fnName, location: f, config: testCfg})
			}

		}
	}
	sort.SliceStable(eTests, func(i, j int) bool { return eTests[i].name < eTests[j].name })
	return eTests, nil
}

// getDiscoverTestFile returns test file path used for auto discovered tests if exists
func getDiscoverTestFile(intDir string) (string, error) {
	discoverTestFilePath := path.Join(intDir, discoverTestFilename)
	_, err := os.Stat(discoverTestFilePath)
	if err != nil {
		return "", err
	}
	return discoverTestFilePath, nil

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
