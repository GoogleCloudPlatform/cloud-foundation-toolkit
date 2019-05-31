package deployment

import (
	"reflect"
	"testing"
)

func TestNewDependencyGraph(t *testing.T) {

	configA := &Config{
		data:    "network: my-net",
		Project: "projectA",
		Name:    "deploymentA",
	}
	configB := &Config{
		data:    "network: $(out.projectA.deploymentA.resourceA.nameA)",
		Project: "projectB",
		Name:    "deploymentB",
	}

	graph := NewDependencyGraph([]Config{*configB, *configA})

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

func TestOrder(t *testing.T) {

	configA := &Config{
		data:    "network: my-net",
		Project: "projectA",
		Name:    "deploymentA",
	}
	configB := &Config{
		data:    "network: $(out.projectA.deploymentA.resourceA.nameA)",
		Project: "projectB",
		Name:    "deploymentB",
	}

	graph := NewDependencyGraph([]Config{*configB, *configA})

	actual, err := graph.Order()
	if err != nil {
		t.Error(err)
	}
	expected := []Config{*configB, *configA}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	graph1 := NewDependencyGraph([]Config{*configA, *configB})

	actual1, err1 := graph1.Order()
	if err1 != nil {
		t.Error(err)
	}
	expected1 := []Config{*configB, *configA}

	if !reflect.DeepEqual(actual1, expected1) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
