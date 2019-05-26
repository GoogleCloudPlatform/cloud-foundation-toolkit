package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	path := filepath.Join("testdata", "config", "simple.yaml")
	config := NewConfig(path)
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
		yaml_string: "network: my-net",
		Project:     "projectA",
		Name:        "deploymentA",
	}
	configB := Config{
		yaml_string: "network: $(out.projectA.deploymentA.resourceA.nameA)",
		Project:     "projectB",
		Name:        "deploymentB",
	}

	configs := []Config{configB, configA}

	if configA.findAllDependencies(configs) != nil {
		t.Errorf("ConfigA should not have any deps")
	}

	if !reflect.DeepEqual(configB.findAllDependencies(configs), []Config{configA}) {
		t.Errorf("ConfigB should have ConfigA as dep")
	}
}

func TestNewConfigGraph(t *testing.T) {

	configA := &Config{
		yaml_string: "network: my-net",
		Project:     "projectA",
		Name:        "deploymentA",
	}
	configB := &Config{
		yaml_string: "network: $(out.projectA.deploymentA.resourceA.nameA)",
		Project:     "projectB",
		Name:        "deploymentB",
	}

	graph := NewConfigGraph([]Config{*configB, *configA})

	if graph.Nodes().Len() != 2 {
		t.Errorf("expected graph size not correct expected %d, got %d", 2, graph.Nodes().Len())
	}

	if graph.From(configA.ID()).Len() != 0 {
		t.Errorf("ConfigA should not have any dependencies")
	}

	if graph.To(configB.ID()).Len() != 0 {
		t.Errorf("ConfigB should not have any dependents")
	}

	if graph.To(configA.ID()).Len() != 1 {
		t.Errorf("ConfigA should have one dependent, got %d", graph.To(configA.ID()).Len())
	}

	if graph.From(configB.ID()).Len() != 1 {
		t.Errorf("From ConfigA to ConfigB should be 1 connection, got %d", graph.From(configB.ID()).Len())
	}

	if graph.To(configA.ID()).Node().ID() != configB.ID() {
		t.Errorf("ConfigB should depends on ConfigA, "+
			"expeted %v, actual %v", configB.ID(), graph.To(configA.ID()).Node().ID())
	}

	if graph.From(configB.ID()).Node().ID() != configA.ID() {
		t.Errorf("ConfigB should depends on ConfigA, "+
			"expeted %v, actual %v", configA.ID(), graph.From(configB.ID()).Node().ID())
	}
}

func TestSort(t *testing.T) {

	configA := &Config{
		yaml_string: "network: my-net",
		Project:     "projectA",
		Name:        "deploymentA",
	}
	configB := &Config{
		yaml_string: "network: $(out.projectA.deploymentA.resourceA.nameA)",
		Project:     "projectB",
		Name:        "deploymentB",
	}

	graph := NewConfigGraph([]Config{*configB, *configA})

	fmt.Println(graph.sort())

}

func executeFindAllOutRefsAndAssert(yaml_string string, expected []string, t *testing.T) {
	config := &Config{yaml_string: yaml_string}
	actual := config.findAllOutRefs()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got: %v, expected: %v.", actual, expected)
	}
}
