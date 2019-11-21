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

	configC := Config{
		data:    "network: $(out.projectA.deploymentA.resourceA.nameA)",
		Project: "projectC",
		Name:    "deploymentC",
	}

	configD := Config{
		data:    "network: $(out.projectC.deploymentC.resourceC.nameC)",
		Project: "projectD",
		Name:    "deploymentD",
	}

	input := map[string]Config{
		configB.FullName(): configB,
		configA.FullName(): configA,
		configD.FullName(): configD,
		configC.FullName(): configC,
	}

	actual, err := Order(input)
	if err != nil {
		t.Errorf("unexpected error on order: %v", err)
	}

	expected1 := [][]Config{{configA}, {configB, configC}, {configD}}
	expected2 := [][]Config{{configA}, {configC, configB}, {configD}}

	if !reflect.DeepEqual(actual, expected1) && !reflect.DeepEqual(actual, expected2) {
		t.Errorf("got %v, expected %v OR %v", actual, expected1, expected2)
	}

}
