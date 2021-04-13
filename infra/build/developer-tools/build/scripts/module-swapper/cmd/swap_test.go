package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
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

func getFileHelper(p string) []byte {
	f, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
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

func getProcessFileTestArgs(p, m string) struct {
	f       []byte
	p       string
	modules []LocalTerraformModule
} {
	f := struct {
		f       []byte
		p       string
		modules []LocalTerraformModule
	}{
		getFileHelper(p),
		p,
		append(
			findSubModules("testdata/"+m+"/modules", "terraform-google-modules/"+m+"/google"),
			LocalTerraformModule{m, getAbsPathHelper("testdata/" + m), fmt.Sprintf("%s/%s/%s", moduleRegistryPrefix, m, moduleRegistrySuffix)},
		),
	}
	return f
}

func Test_processFile(t *testing.T) {
	type args struct {
		f       []byte
		p       string
		modules []LocalTerraformModule
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"simple", getProcessFileTestArgs("testdata/example-module-simple/examples/example-one/main.tf", "example-module-simple"), getFileHelper("testdata/example-module-simple/examples/example-one/main.tf.good"), false},
		{"simple-submodules-single-submod", getProcessFileTestArgs("testdata/example-module-with-submodules/examples/example-one/main.tf", "example-module-with-submodules"), getFileHelper("testdata/example-module-with-submodules/examples/example-one/main.tf.good"), false},
		{"simple-submodules-multiple-modules", getProcessFileTestArgs("testdata/example-module-with-submodules/examples/example-two/main.tf", "example-module-with-submodules"), getFileHelper("testdata/example-module-with-submodules/examples/example-two/main.tf.good"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupProcessFileTest(tt.args.modules)
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()
			got, err := replaceLocalModules(tt.args.f, tt.args.p)
			t.Log(buf.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("processFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processFile() = %v, want %v", string(got), string(tt.want))
			}
			tearDownProcessFileTest()
		})
	}
}

func getTempDir() string {
	d, err := ioutil.TempDir("", "gitrmtest")
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
		{"simple", args{tempGitRepoWithRemote("https://github.com/foo/terraform-google-bar", "origin")}, "bar", "foo", false, ""},
		{"simple-with-trailing-slash", args{tempGitRepoWithRemote("https://github.com/foo/terraform-google-bar/", "origin")}, "bar", "foo", false, ""},
		{"simple-with-trailing-git", args{tempGitRepoWithRemote("https://github.com/foo/terraform-google-bar.git", "origin")}, "bar", "foo", false, ""},
		{"err-no-remote-origin", args{tempGitRepoWithRemote("https://github.com/foo/terraform-google-bar", "foo")}, "", "", true, ""},
		{"err-not-git-repo", args{getTempDir()}, "", "", true, ""},
		{"err-not-github-repo", args{tempGitRepoWithRemote("https://gitlab.com/foo/terraform-google-bar", "origin")}, "", "", true, "Expected GitHub remote of form https://github.com/ModuleRegistry/ModuleRepo"},
		{"err-not-prefixed-repo", args{tempGitRepoWithRemote("https://github.com/foo/bar", "origin")}, "", "", true, "Expected to find repo name prefixed with terraform-google-"},
		{"err-malformed-remote", args{tempGitRepoWithRemote("https://github.com/footerraform-google-bar", "origin")}, "", "", true, "Expected GitHub org/owner of form ModuleRegistry/ModuleRepo"},
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
