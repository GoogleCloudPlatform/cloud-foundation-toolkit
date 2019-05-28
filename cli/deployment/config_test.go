package deployment

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	data := GetTestData("config", "simple.yaml", t)
	config := NewConfig(data, "")
	if config == nil {
		t.Errorf("Config is nil")
	}
}

func TestFindAllOutRefs(t *testing.T) {
	executeFindAllOutRefsAndAssert(
		`$(out.project1.deployment1.resource1.output1)
                    $(out.deployment2.resource2.output2)`,
		[]string{
			"project1.deployment1.resource1.output1",
			"deployment2.resource2.output2",
		},
		t)
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

func executeFindAllOutRefsAndAssert(yaml_string string, expected []string, t *testing.T) {
	config := &Config{data: yaml_string}
	actual := config.findAllOutRefs()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got: %v, expected: %v.", actual, expected)
	}
}
