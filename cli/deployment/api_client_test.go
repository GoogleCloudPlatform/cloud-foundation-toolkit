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

// do actuall deployment get status call
func testGetStatus(t *testing.T) {
	dep := Deployment{
		configFile: "/1-prj/projects/google/CFT/fork/cloud-foundation-toolkit/dm/network.yaml",
		config: Config{
			Project: "gl-akopachevskyy-dm-seed",
			Name:    "my-networks-main",
		},
	}
	status, err := GetStatus(dep)
	if err != nil {
		t.Error(err)
	}
	if status != NotFound {
		t.Errorf("Expected status %d, actual %d", NotFound, status)
	}
}
