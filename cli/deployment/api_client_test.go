package deployment

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetOutputs(t *testing.T) {
	RunGCloud = func(args ...string) (result string, err error) {
		expected := "deployment-manager manifests describe --deployment mydeployment --project myproject --format yaml"
		actual := strings.Join(args, " ")
		if expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		return GetTestData("deployment", "describe-manifest.yaml", t), nil
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

func TestStatus_String(t *testing.T) {
	var inputTests = []struct {
		status   Status
		expected string
	}{
		{Done, "DONE"},
		{Pending, "PENDING"},
		{Running, "RUNNING"},
		{NotFound, "NOT_FOUND"},
		{Error, "ERROR"},
	}

	for _, tt := range inputTests {
		t.Run(tt.expected, func(t *testing.T) {
			var actual string = fmt.Sprintf("%v", tt.status)
			if actual != tt.expected {
				t.Errorf("got: %s, want: %s", actual, tt.expected)
			}
		})
	}
}
