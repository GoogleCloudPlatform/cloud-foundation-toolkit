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
	configC := Config{
		data:    "network: $(out.projectA.deploymentA.resourceA.nameA)-$(out.projectB.deploymentB.resourceB.nameB)",
		Project: "projectC",
		Name:    "deploymentC",
	}

	configs := map[string]Config{
		configC.FullName(): configC,
		configB.FullName(): configB,
		configA.FullName(): configA,
	}

	actual, err := configA.findAllDependencies(configs)
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}
	if len(actual) != 0 {
		t.Errorf("want %d, got %v", 0, len(actual))
	}

	actual, err = configB.findAllDependencies(configs)
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}
	if !reflect.DeepEqual(actual, []Config{configA}) {
		t.Errorf("want %v, got %v", []Config{configA}, actual)
	}

	actual, err = configC.findAllDependencies(configs)
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}
	if !reflect.DeepEqual(actual, []Config{configA, configB}) {
		t.Errorf("want %v, got %v", []Config{configA, configB}, actual)
	}
}

func TestFindAllDependencies_corner_cases(t *testing.T) {
	var tests = []struct {
		name          string
		describeFile  string
		errorExpected bool
	}{
		{"done", "done.yaml", false},
		{"running", "running.yaml", true},
		{"pending", "pending.yaml", true},
		{"error", "error.yaml", true},
		{"not_found", "not-found.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunGCloud = func(args ...string) (result string, err error) {
				expected := "deployment-manager deployments describe deploymentA --project projectA --format yaml"
				actual := strings.Join(args, " ")
				if expected != actual {
					t.Errorf("expected %v, got %v", expected, actual)
				}
				return GetTestData("deployment/describe", tt.describeFile, t), nil
			}

			configB := Config{
				data:    "network: $(out.projectA.deploymentA.resourceA.nameA)",
				Project: "projectB",
				Name:    "deploymentB",
			}

			configs := map[string]Config{
				configB.FullName(): configB,
			}

			res, err := configB.findAllDependencies(configs)

			if tt.errorExpected {
				if res != nil {
					t.Errorf("expected %v, got %v", nil, res)
				}
				if err == nil {
					t.Errorf("expected error")
				}
			}
		})
	}
}

func TestYAMLStripCustomFields(t *testing.T) {
	data, err := Config{
		data: GetTestData("config", "custom-fields.yaml", t),
	}.YAML(map[string]map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(string(data), "project:") {
		t.Errorf("Should not contain, project, name or descriptions")
	}
}

func TestYamlReplaceOutRefs(t *testing.T) {
	data := GetTestData("cross-ref", "main-manifest.yaml", t)
	output, err := parseOutputs(data)
	if err != nil {
		t.Errorf("failed to parse outputs: %v", err)
	}
	DefaultProjectID = "my-project"
	config := NewConfig(GetTestData("cross-ref", "dependent-with-refs.yaml", t), "/home/test.yaml")
	outputs := map[string]map[string]interface{}{}
	outputs["prj1.name1"] = output
	actual, err := config.YAML(outputs)
	if err != nil {
		t.Errorf("failed to export config YAML: %v", err)
	}
	expected := GetTestData("cross-ref", "dependent-final-expected.yaml", t)
	if strings.TrimSpace(string(actual)) != strings.TrimSpace(expected) {
		t.Errorf("got: \n%s, want: \n%s", actual, expected)
	}
}
