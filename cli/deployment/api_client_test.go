package deployment

import (
	"strings"
	"testing"
)

func TestGetOutputs(t *testing.T) {
	runGCloud = func(args ...string) (result string, err error) {
		expected := "deployment-manager manifests describe --deployment myproject --project mydeployment"
		actual := strings.Join(args, " ")
		if expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		return GetTestData("deployment", "api-describe-manifest.yaml", t), nil
	}

	outputs, err := GetOutputs("myproject", "mydeployment")
	if err != nil {
		t.Errorf("erorr getting deployment outputs: %v", err)
	}
	expected := "my-network-prod"
	if expected != outputs["my-network-prod.name"] {
		t.Errorf("expected: %s, got: %s", expected, outputs["my-network-prod.name"])
	}
}
