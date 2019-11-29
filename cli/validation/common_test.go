package validation

import (
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

func TestParseResourceProperties(t *testing.T) {
	cases := []struct {
		name    string
		project string
		input   deployment.DeploymentDescriptionResource
		output  cai.Resource
	}{
		{
			name:    "properties",
			project: "project",
			input: deployment.DeploymentDescriptionResource{
				Name:       "properties",
				Type:       "test",
				Properties: "foo: bar",
			},
			output: cai.Resource{
				Name:    "properties",
				Type:    "test",
				Project: "project",
				Properties: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
		{
			name:    "final properties",
			project: "project",
			input: deployment.DeploymentDescriptionResource{
				Name:            "final properties",
				Type:            "test",
				Properties:      "foo: bar",
				FinalProperties: "foo2: bar2",
			},
			output: cai.Resource{
				Name:    "final properties",
				Type:    "test",
				Project: "project",
				Properties: map[string]interface{}{
					"foo2": "bar2",
				},
			},
		},
		{
			name:    "update properties",
			project: "project",
			input: deployment.DeploymentDescriptionResource{
				Name:            "update properties",
				Type:            "test",
				Properties:      "foo: bar",
				FinalProperties: "foo2: bar2",
				Update: struct {
					Properties      string
					FinalProperties string `yaml:",omitempty"`
					State           string
				}{
					Properties:      "foo3: bar3",
					FinalProperties: "",
					State:           "",
				},
			},
			output: cai.Resource{
				Name:    "update properties",
				Type:    "test",
				Project: "project",
				Properties: map[string]interface{}{
					"foo3": "bar3",
				},
			},
		},
		{
			name:    "final update properties",
			project: "project",
			input: deployment.DeploymentDescriptionResource{
				Name:            "final update properties",
				Type:            "test",
				Properties:      "foo: bar",
				FinalProperties: "foo2: bar2",
				Update: struct {
					Properties      string
					FinalProperties string `yaml:",omitempty"`
					State           string
				}{
					Properties:      "foo3: bar3",
					FinalProperties: "foo4: bar4",
					State:           "",
				},
			},
			output: cai.Resource{
				Name:    "final update properties",
				Type:    "test",
				Project: "project",
				Properties: map[string]interface{}{
					"foo4": "bar4",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := parseResourceProperties(c.project, c.input)

			if err != nil {
				t.Errorf("got error: %t", err)
			}
			if !reflect.DeepEqual(res, c.output) {
				t.Errorf("got %v, expected %v", res, c.output)
			}
		})
	}
}

func TestGetAncestry(t *testing.T) {
	cmd := ""

	runGCloud = func(args ...string) (result string, err error) {
		cmd = strings.Join(args, " ")
		return "[{\"Id\":\"test1\",\"Type\":\"project\"},{\"Id\":\"test2\",\"Type\":\"organization\"}]", nil
	}

	anc, err := getAncestry("foobar")
	if err != nil {
		t.Errorf("got error: %t", err)
	}

	testAnc := "organization/test2/project/test1"
	if anc != testAnc {
		t.Errorf("got %v, expected %v", anc, testAnc)
	}

	testCmd := "projects get-ancestors foobar --format json"
	if cmd != testCmd {
		t.Errorf("got %v, expected %v", cmd, testCmd)
	}
}
