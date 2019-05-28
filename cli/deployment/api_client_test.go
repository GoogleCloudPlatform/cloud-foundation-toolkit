package deployment

import (
	"strings"
	"testing"
)

/*
func testGCloud(t *testing.T) {
	res, err := GCloud("deployment-manager", "deployments", "describe", "gl-akopachevsky-test")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}*/

func TestGetOutputs(t *testing.T) {
	execFunction = func(args ...string) (result string, err error) {
		expected := "deployment-manager manifests describe --deployment myproject --project mydeployment"
		actual := strings.Join(args, " ")
		if expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		return GetTestData("deployment", "api-describe-manifest.yaml", t), nil
	}

	outputs, err := GetOutputs("myproject", "mydeployment")

	if err != nil {
		t.Error(err)
	}
	if "my-network-prod" != outputs["my-network-prod.name"] {
		t.Errorf("Expected \"my-network-prod\" got \"%s\"", outputs["my-network-prod.name"])
	}
}
