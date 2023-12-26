package cmd

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/google/go-cmp/cmp"
)

var (
	moduleRegistryPrefix = "terraform-google-modules"
	moduleRegistrySuffix = "google"
)

func getAbsPathHelper(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		log.Fatalf("Unable to find absolute path %s: %v", p, err)
	}
	return a
}

func getFileHelper(t *testing.T, p string) []byte {
	f, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	return f
}

func setupProcessFileTest(modules []LocalTerraformModule) {
	localModules = modules
}

func tearDownProcessFileTest() {
	localModules = []LocalTerraformModule{}
}

func Test_getTFFiles(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"simple", args{"testdata/example-module-simple"}, []string{"testdata/example-module-simple/examples/example-one/main.tf", "testdata/example-module-simple/examples/main.tf"}},
		{"simple-single-submodule", args{"testdata/example-module-with-submodules/modules/bar-module"}, []string{"testdata/example-module-with-submodules/modules/bar-module/main.tf"}},
		{"simple-single-submodule-empty", args{"testdata/example-module-with-submodules/docs"}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTFFiles(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTFFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findSubModules(t *testing.T) {
	type args struct {
		path          string
		rootModuleFQN string
	}
	tests := []struct {
		name string
		args args
		want []LocalTerraformModule
	}{
		{"simple-no-submodules", args{"testdata/example-module-simple/modules", "terraform-google-modules/example-module-simple/google"}, []LocalTerraformModule{}},
		{"simple-with-submodules", args{"testdata/example-module-with-submodules/modules", "terraform-google-modules/example-module-with-submodules/google"},
			[]LocalTerraformModule{
				{"bar-module", filepath.Join(getAbsPathHelper("testdata/example-module-with-submodules/modules"), "bar-module"), "terraform-google-modules/example-module-with-submodules/google//modules/bar-module"},
				{"foo-module", filepath.Join(getAbsPathHelper("testdata/example-module-with-submodules/modules"), "foo-module"), "terraform-google-modules/example-module-with-submodules/google//modules/foo-module"},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findSubModules(tt.args.path, tt.args.rootModuleFQN); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findSubModules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processFile(t *testing.T) {
	tests := []struct {
		name              string
		modules           []LocalTerraformModule
		exampleRemotePath string
		exampleLocalPath  string
	}{
		{
			name:              "simple",
			modules:           testModules("example-module-simple"),
			exampleRemotePath: "example-module-simple/examples/example-one/main.tf",
			exampleLocalPath:  "example-module-simple/examples/example-one/main.tf.local",
		},
		{
			name:              "simple-submodules-single-submod",
			modules:           testModules("example-module-with-submodules"),
			exampleRemotePath: "example-module-with-submodules/examples/example-one/main.tf",
			exampleLocalPath:  "example-module-with-submodules/examples/example-one/main.tf.local",
		},
		{
			name:              "simple-submodules-multiple-modules",
			modules:           testModules("example-module-with-submodules"),
			exampleRemotePath: "example-module-with-submodules/examples/example-two/main.tf",
			exampleLocalPath:  "example-module-with-submodules/examples/example-two/main.tf.local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupProcessFileTest(tt.modules)
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()
			tt.exampleRemotePath = path.Join(testDataDir, tt.exampleRemotePath)
			tt.exampleLocalPath = path.Join(testDataDir, tt.exampleLocalPath)
			remoteExample := getFileHelper(t, tt.exampleRemotePath)
			localExample := getFileHelper(t, tt.exampleLocalPath)

			// Swap remote references to local.
			got, err := remoteToLocal(remoteExample, tt.exampleRemotePath)
			if err != nil {
				t.Fatalf("remoteToLocal() error = %v", err)
				return
			}
			if diff := cmp.Diff(localExample, got); diff != "" {
				t.Errorf("remoteToLocal() returned unexpected difference (-want +got):\n%s", diff)
			}

			// Swap local references to remote.
			got, err = localToRemote(localExample, tt.exampleLocalPath)
			t.Log(buf.String())
			if err != nil {
				t.Errorf("localToRemote() error = %v", err)
				return
			}
			if diff := cmp.Diff(remoteExample, got); diff != "" {
				t.Errorf("localToRemote() returned unexpected difference (-want +got):\n%s", diff)
			}
			tearDownProcessFileTest()
		})
	}
}

const testDataDir = "testdata"

func testModules(m string) []LocalTerraformModule {
	root := LocalTerraformModule{m, getAbsPathHelper(path.Join(testDataDir, m)), path.Join(moduleRegistryPrefix, m, moduleRegistrySuffix)}
	return append(findSubModules(path.Join(testDataDir, m, "modules"), path.Join(moduleRegistryPrefix, m, moduleRegistrySuffix)), root)
}

func getTempDir() string {
	d, err := os.MkdirTemp("", "gitrmtest")
	if err != nil {
		log.Fatalf("Error creating tempdir: %v", err)
	}
	return d
}

func tempGitRepoWithRemote(repoURL, remote string) string {
	dir := getTempDir()
	r, err := git.PlainInit(dir, true)
	if err != nil {
		log.Fatalf("Error creating repo in tempdir: %v", err)
	}
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: remote,
		URLs: []string{repoURL},
	})
	if err != nil {
		log.Fatalf("Error creating remote in tempdir repo: %v", err)
	}
	return dir
}

func Test_getModuleNameRegistry(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name       string
		args       args
		want       string
		want1      string
		wantErr    bool
		wantErrStr string
	}{
		{"simple-https", args{tempGitRepoWithRemote("https://github.com/foo/terraform-google-bar", "origin")}, "bar", "foo", false, ""},
		{"simple-git", args{tempGitRepoWithRemote("git@github.com:foo/terraform-google-bar.git", "origin")}, "bar", "foo", false, ""},
		{"simple-with-trailing-slash", args{tempGitRepoWithRemote("https://github.com/foo/terraform-google-bar/", "origin")}, "bar", "foo", false, ""},
		{"simple-with-trailing-git", args{tempGitRepoWithRemote("https://github.com/foo/terraform-google-bar.git", "origin")}, "bar", "foo", false, ""},
		{"err-no-remote-origin", args{tempGitRepoWithRemote("https://github.com/foo/terraform-google-bar", "foo")}, "", "", true, ""},
		{"err-not-git-repo", args{getTempDir()}, "", "", true, ""},
		{"err-not-github-repo-https", args{tempGitRepoWithRemote("https://gitlab.com/foo/terraform-google-bar", "origin")}, "", "", true, "expected GitHub remote, got: https://gitlab.com/foo/terraform-google-bar"},
		{"err-not-github-repo-ssh", args{tempGitRepoWithRemote("git@gitlab.com:foo/terraform-google-bar.git", "origin")}, "", "", true, "expected GitHub remote, got: git@gitlab.com:foo/terraform-google-bar.git"},
		{"err-not-prefixed-repo", args{tempGitRepoWithRemote("https://github.com/foo/bar", "origin")}, "", "", true, "expected to find repo name prefixed with terraform-google-"},
		{"err-malformed-remote", args{tempGitRepoWithRemote("https://github.com/footerraform-google-bar", "origin")}, "", "", true, "expected GitHub remote of form https://github.com/ModuleRegistry/ModuleRepo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getModuleNameRegistry(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("getModuleNameRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				if tt.wantErrStr != "" {
					if !strings.Contains(err.Error(), tt.wantErrStr) {
						t.Errorf("getModuleNameRegistry() error = %v, expected to contain %v", err, tt.wantErrStr)
					}
				}
			}
			if got != tt.want {
				t.Errorf("getModuleNameRegistry() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getModuleNameRegistry() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
