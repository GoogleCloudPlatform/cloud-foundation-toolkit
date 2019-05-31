package deployment

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	data := GetTestData("config", "simple.yaml", t)
	config := NewConfig(data, "")
	if config == nil {
		t.Errorf("Config is nil")
	}
	if len(config.Imports) != 1 {
		t.Errorf("Expected to have 1 import")
	}
	if len(config.Resources) != 2 {
		t.Errorf("Expected to have 2 resources")
	}
}

func TestFindAllOutRefs(t *testing.T) {
	executeFindAllOutRefsAndAssert(
		`$(out.project1.deployment1.resource1.output1)
                    $(out.deployment2.resource2.output2)`,
		[]string{
			"project1.deployment1.resource1.output1",
			"deployment2.resource2.output2",
		}, t)
	// empty string
	executeFindAllOutRefsAndAssert("", nil, t)
	// invalid notation
	executeFindAllOutRefsAndAssert("${out1.account.project.resource.output", nil, t)
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

	configs := []Config{configB, configA}

	if configA.findAllDependencies(configs) != nil {
		t.Errorf("ConfigA should not have any deps")
	}

	if !reflect.DeepEqual(configB.findAllDependencies(configs), []Config{configA}) {
		t.Errorf("ConfigB should have ConfigA as dep")
	}
}

func TestYaml(t *testing.T) {
	data, err := Config{
		data: GetTestData("config", "custom-fields.yaml", t),
	}.Yaml()
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(string(data), "project:") {
		t.Errorf("Should not contain, project, name or descriptions")
	}
}

func executeFindAllOutRefsAndAssert(yaml_string string, expected []string, t *testing.T) {
	config := &Config{data: yaml_string}
	actual := config.findAllOutRefs()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got: %v, expected: %v.", actual, expected)
	}
}
