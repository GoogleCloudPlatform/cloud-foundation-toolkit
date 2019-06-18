package deployment

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	data := GetTestData("config", "simple.yaml", t)
	config := NewConfig(data, "")

	if len(config.Imports) != 1 {
		t.Errorf("expected 1 import, got: %d", len(config.Imports))
	}
	if len(config.Resources) != 2 {
		t.Errorf("want 2, got: %d", len(config.Resources))
	}
}

func TestFindAllOutRefs(t *testing.T) {
	var tests = []struct {
		name string
		in   string
		out  []string
	}{
		{
			"one refs",
			"$(out.project1.deployment1.resource1.output1)",
			[]string{
				"project1.deployment1.resource1.output1",
			},
		},
		{
			"several refs",
			`$(out.project1.deployment1.resource1.output1)
                    $(out.deployment2.resource2.output2)`,
			[]string{
				"project1.deployment1.resource1.output1",
				"deployment2.resource2.output2",
			},
		},
		{"empty file", "", nil},
		{"no refs", "name: myname", nil},
		{"invalid delimiter", "${out1.account.project.resource.output}", nil},
		{"missing closing delimiter", "$(out1.account.project.resource.output", nil},
		{"no delimiter", "out1.account.project.resource.output", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{data: tt.in}
			actual := config.findAllOutRefs()
			if !reflect.DeepEqual(actual, tt.out) {
				t.Errorf("want: %v, got: %v.", tt.out, actual)
			}
		})
	}
}

func TestFindAllDependencies(t *testing.T) {
	configA := Config{
		data:    "network: my-net",
		Project: "projectA",
		Name:    "deploymentA",
	}
	configB := Config{
		data:    "network: $(out.projectA.deploymentA.resourceA.nameA)",
		Project: "projectB",
		Name:    "deploymentB",
	}

	configs := map[string]Config{
		configB.FullName(): configB,
		configA.FullName(): configA,
	}

	actual := configA.findAllDependencies(configs)
	if len(actual) != 0 {
		t.Errorf("want %d, got %v", 0, len(actual))
	}

	actual = configB.findAllDependencies(configs)
	if !reflect.DeepEqual(actual, []Config{configA}) {
		t.Errorf("want %v, got %v", []Config{configA}, actual)
	}
}

func TestYAML(t *testing.T) {
	data, err := Config{
		data: GetTestData("config", "custom-fields.yaml", t),
	}.YAML()
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(string(data), "project:") {
		t.Errorf("Should not contain, project, name or descriptions")
	}
}
