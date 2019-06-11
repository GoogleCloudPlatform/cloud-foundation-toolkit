package cmd

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

var emptyMap = map[string]string{}

var tests = []struct {
	name  string
	input []string
	files map[string]string
	yamls string
	errs  map[string]error
}{
	{name: "yaml", input: []string{"name: name1\nproject: prj"}, files: emptyMap, yamls: "name: name1\nproject: prj"},
	{name: "dir",
		input: []string{"../testdata/cmd/common/dir"},
		files: map[string]string{
			"../testdata/cmd/common/dir/config1.yaml":  "config1",
			"../testdata/cmd/common/dir/config2.yml":   "config2",
			"../testdata/cmd/common/dir/config3.jinja": "config3"},
	},
	{name: "dir wildcard",
		input: []string{"../testdata/cmd/common/dir/*"},
		files: map[string]string{
			"../testdata/cmd/common/dir/config1.yaml":  "config1",
			"../testdata/cmd/common/dir/config2.yml":   "config2",
			"../testdata/cmd/common/dir/config3.jinja": "config3",
			"../testdata/cmd/common/dir/config4.txt":   "config4"},
	},
	{name: "yaml file", input: []string{"../testdata/cmd/common/dir/*.yaml"}, files: map[string]string{"../testdata/cmd/common/dir/config1.yaml": "config1"}},
	{name: "yml file", input: []string{"../testdata/cmd/common/dir/*.yml"}, files: map[string]string{"../testdata/cmd/common/dir/config2.yml": "config2"}},
	{name: "jinja file", input: []string{"../testdata/cmd/common/dir/*.jinja"}, files: map[string]string{"../testdata/cmd/common/dir/config3.jinja": "config3"}},
	{name: "file", input: []string{"../testdata/cmd/common/dir/config1.yaml"}, files: map[string]string{"../testdata/cmd/common/dir/config1.yaml": "config1"}},
	{name: "two files", input: []string{"../testdata/cmd/common/dir/config1.yaml", "../testdata/cmd/common/dir/config2.yml"},
		files: map[string]string{"../testdata/cmd/common/dir/config1.yaml": "config1", "../testdata/cmd/common/dir/config2.yml": "config2"}},

	{name: "file not exists", input: []string{"../testdata/cmd/common/dir/1.yaml"}, files: emptyMap,
		errs: map[string]error{"../testdata/cmd/common/dir/1.yaml": errors.New("no file(s) exists or valid yaml for config param: ../testdata/cmd/common/dir/1.yaml")}},
	{name: "empty dir", input: []string{"../testdata/cmd/common/emptydir"}, files: emptyMap,
		errs: map[string]error{"../testdata/cmd/common/emptydir": errors.New("no *.yaml, *.yml, *.jinja files found in directory: ../testdata/cmd/common/emptydir")}},
}

func TestListConfigs(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := map[string]error{}
			files, yamlsSlice := listConfigs(tt.input, errs)
			yamls := strings.Join(yamlsSlice, " ")

			if !reflect.DeepEqual(files, tt.files) {
				t.Errorf("for files got %v, want %v", files, tt.files)
			}

			if yamls != tt.yamls {
				t.Errorf("got %v, want %v", yamls, tt.yamls)
			}

			for key, value := range tt.errs {
				if actualValue, ok := tt.errs[key]; !ok {
					t.Errorf("errors map should contain error for file: %s", key)
				} else {
					if !reflect.DeepEqual(actualValue, value) {
						t.Errorf("error for file: %s, got: %v, expected %v", key, actualValue, value)
					}
				}
			}
		})
	}
}
