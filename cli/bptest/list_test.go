package bptest

import (
	"os"
	"path"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/discovery"
	"github.com/stretchr/testify/assert"
)

const (
	testDirWithDiscovery = "testdata/with-discovery"
	intTestDir           = "test/integration"
)

func TestGetDiscoveredTests(t *testing.T) {
	tests := []struct {
		name    string
		testDir string
		want    []bpTest
		errMsg  string
	}{
		{
			name:    "simple",
			testDir: path.Join(testDirWithDiscovery, intTestDir),
			want: []bpTest{
				getBPTest("TestAll/examples/baz", path.Join(testDirWithDiscovery, "examples/baz"), path.Join(testDirWithDiscovery, intTestDir, discoverTestFilename), false),
				getBPTest("TestAll/fixtures/qux", path.Join(testDirWithDiscovery, "test/fixtures/qux"), path.Join(testDirWithDiscovery, intTestDir, discoverTestFilename), false),
				getBPTest("TestAll/examples/quux", path.Join(testDirWithDiscovery, "examples/quux"), path.Join(testDirWithDiscovery, intTestDir, discoverTestFilename), true),
			},
		},
		{
			name:    "no discovery",
			testDir: path.Join(testDirWithDiscovery, "doesnotexist"),
			want:    []bpTest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got, err := getDiscoveredTests(tt.testDir)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.ElementsMatch(tt.want, got)
			}
		})
	}
}

func TestGetExplicitTests(t *testing.T) {
	tests := []struct {
		name    string
		testDir string
		want    []bpTest
		errMsg  string
	}{
		{
			name:    "simple",
			testDir: path.Join(testDirWithDiscovery, intTestDir),
			want: []bpTest{
				getBPTest("TestBar", path.Join(testDirWithDiscovery, "examples/bar"), path.Join(testDirWithDiscovery, intTestDir, "bar/bar_test.go"), false),
				getBPTest("TestFoo", path.Join(testDirWithDiscovery, "test/fixtures/foo"), path.Join(testDirWithDiscovery, intTestDir, "foo/foo_test.go"), false),
				getBPTest("TestQuuz", path.Join(testDirWithDiscovery, "test/fixtures/quuz"), path.Join(testDirWithDiscovery, intTestDir, "quuz/quuz_test.go"), true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got, err := getExplicitTests(tt.testDir)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.ElementsMatch(tt.want, got)
			}
		})
	}
}

func TestGetTests(t *testing.T) {
	tests := []struct {
		name    string
		testDir string
		want    []bpTest
		errMsg  string
	}{
		{
			name:    "simple",
			testDir: path.Join(testDirWithDiscovery, intTestDir),
			want: []bpTest{
				getBPTest("TestAll/examples/baz", path.Join(testDirWithDiscovery, "examples/baz"), path.Join(testDirWithDiscovery, intTestDir, discoverTestFilename), false),
				getBPTest("TestAll/fixtures/qux", path.Join(testDirWithDiscovery, "test/fixtures/qux"), path.Join(testDirWithDiscovery, intTestDir, discoverTestFilename), false),
				getBPTest("TestAll/examples/quux", path.Join(testDirWithDiscovery, "examples/quux"), path.Join(testDirWithDiscovery, intTestDir, discoverTestFilename), true),
				getBPTest("TestBar", path.Join(testDirWithDiscovery, "examples/bar"), path.Join(testDirWithDiscovery, intTestDir, "bar/bar_test.go"), false),
				getBPTest("TestFoo", path.Join(testDirWithDiscovery, "test/fixtures/foo"), path.Join(testDirWithDiscovery, intTestDir, "foo/foo_test.go"), false),
				getBPTest("TestQuuz", path.Join(testDirWithDiscovery, "test/fixtures/quuz"), path.Join(testDirWithDiscovery, intTestDir, "quuz/quuz_test.go"), true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got, err := getTests(tt.testDir)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.ElementsMatch(tt.want, got)
			}
		})
	}
}

func getBPTest(n string, c string, l string, s bool) bpTest {
	b := discovery.BlueprintTestConfig{}
	b.Spec.Skip = s
	if s {
		b.APIVersion = "blueprints.cloud.google.com/v1alpha1"
		b.Kind = "BlueprintTest"
		b.Name = path.Base(c)
		b.Path = path.Join(c, discovery.DefaultTestConfigFilename)
	}
	return bpTest{name: n, config: c, location: l, bptestCfg: b}
}

func TestGetDiscoverTestName(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   string
		errMsg string
	}{
		{
			name: "simple",
			data: `package test

import "testing"

func TestAll(t *testing.T) {
}
`,
			want: "TestAll",
		},
		{
			name: "multiple",
			data: `package test

import "testing"

const ShouldNotErr = "foo"

func TestA(t *testing.T) {
}

func TestB(t *testing.T) {
}

func OtherHelper(t *testing.T) {
}
`,
			errMsg: "only one function should be defined",
		},
		{
			name: "empty",
			data: `package test
`,
			errMsg: "only one function should be defined",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			filePath, cleanup := writeTmpFile(t, tt.data)
			defer cleanup()
			got, err := getDiscoverTestName(filePath)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}

func Test_discoverIntTestDir(t *testing.T) {
	tests := []struct {
		name   string
		files  []string
		want   string
		errMsg string
	}{
		{
			name:  "with single discover_test.go",
			files: []string{discoverTestFilename},
			want:  ".",
		},
		{
			name:  "with single discover_test.go in a dir",
			files: []string{path.Join("test/integration", discoverTestFilename)},
			want:  "test/integration",
		},
		{
			name:  "with single discover_test.go in a dir and other files",
			files: []string{path.Join("foo/bar/baz", discoverTestFilename), "foo.go", "test.tf", "other/test/bar_test.go"},
			want:  "foo/bar/baz",
		},
		{
			name:  "with single discover_test.go and multiple hidden discover_test.go",
			files: []string{path.Join("foo/bar/baz", discoverTestFilename), path.Join("foo/bar/baz/.terraform", discoverTestFilename), "foo.go", "test.tf", "other/test/bar_test.go"},
			want:  "foo/bar/baz",
		},
		{
			name:   "with multiple discover_test.go",
			files:  []string{path.Join("mod1/test/integration", discoverTestFilename), path.Join("mod2/test/integration", discoverTestFilename)},
			errMsg: "found multiple discover_test.go files:",
		},
		{
			name:  "no discover_test.go files",
			files: []string{},
			want:  ".",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			dir, cleanup := createFilesInTmpDir(t, tt.files)
			defer cleanup()
			got, err := discoverIntTestDir(dir)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}

func createFilesInTmpDir(t *testing.T, files []string) (string, func()) {
	assert := assert.New(t)
	tempDir, err := os.MkdirTemp("", "bpt-")
	assert.NoError(err)
	cleanup := func() { os.RemoveAll(tempDir) }

	//create files in tmpdir
	for _, f := range files {
		p := path.Join(tempDir, path.Dir(f))
		err = os.MkdirAll(p, 0755)
		assert.NoError(err)
		_, err = os.Create(path.Join(p, path.Base(f)))
		assert.NoError(err)
	}
	return tempDir, cleanup
}
