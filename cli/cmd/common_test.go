package cmd

import (
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
	errs  bool
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
	{name: "yaml file", input: []string{"../testdata/cmd/common/dir/*.yaml"},
		files: map[string]string{
			"../testdata/cmd/common/dir/config1.yaml": "config1",
		},
	},
	{name: "yml file", input: []string{"../testdata/cmd/common/dir/*.yml"},
		files: map[string]string{"../testdata/cmd/common/dir/config2.yml": "config2"}},
	{name: "jinja file", input: []string{"../testdata/cmd/common/dir/*.jinja"},
		files: map[string]string{"../testdata/cmd/common/dir/config3.jinja": "config3"}},
	{name: "file", input: []string{"../testdata/cmd/common/dir/config1.yaml"},
		files: map[string]string{"../testdata/cmd/common/dir/config1.yaml": "config1"}},
	{name: "two files", input: []string{"../testdata/cmd/common/dir/config1.yaml", "../testdata/cmd/common/dir/config2.yml"},
		files: map[string]string{
			"../testdata/cmd/common/dir/config1.yaml": "config1",
			"../testdata/cmd/common/dir/config2.yml":  "config2",
		},
	},
	{name: "file not exists", input: []string{"../testdata/cmd/common/dir/1.yaml"}, files: emptyMap, errs: true},
	{name: "empty dir", input: []string{"../testdata/cmd/common/dir/empty_dir"}, files: emptyMap, errs: true},
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

			hasErrors := len(errs) > 0
			if tt.errs != hasErrors {
				t.Errorf("got errors: %t, want %t", hasErrors, tt.errs)
			}
		})
	}
}
