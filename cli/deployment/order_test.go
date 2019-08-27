package deployment

import (
	"reflect"
	"testing"
)

func TestOrder(t *testing.T) {
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

	input := map[string]Config{
		configB.FullName(): configB,
		configA.FullName(): configA,
	}

	actual, err := Order(input)
	if err != nil {
		t.Errorf("unexpected error on order: %v", err)
	}
	expected := []Config{configA, configB}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

}
