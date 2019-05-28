package deployment

import (
	"fmt"
	"strings"
	"testing"
)

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

func TestCreate(t *testing.T) {
	dep := Deployment{
		configFile: "/1-prj/projects/google/CFT/fork/cloud-foundation-toolkit/dm/network.yaml",
		config: Config{
			Project: "gl-akopachevskyy-dm-seed",
			Name:    "my-first-auto-deployment",
		},
	}
	deployment, err := Create(&dep)
	if err != nil {
		t.Error(err)
	}

	if len(deployment.outputs) == 0 {
		t.Errorf("Should be more outputs")
	}

	fmt.Println(deployment)
}
