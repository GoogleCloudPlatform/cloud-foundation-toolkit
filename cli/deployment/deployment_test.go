package deployment

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

var d = &Deployment{
	config: Config{
		Name:    "dpl",
		Project: "prj",
	},
	configFile: "/tmp/myconfig.yaml",
}

var executeTests = []struct {
	action   string
	commands []string
}{
	{
		action: ActionCreate,
		commands: []string{
			"deployment-manager deployments create dpl --config /tmp/myconfig.yaml --project prj",
			"deployment-manager manifests describe --deployment dpl --project prj --format yaml",
		},
	},
	{
		action: ActionUpdate,
		commands: []string{
			"deployment-manager deployments update dpl --config /tmp/myconfig.yaml --project prj",
			"deployment-manager manifests describe --deployment dpl --project prj --format yaml",
		},
	},
	{
		action: ActionDelete,
		commands: []string{
			"deployment-manager deployments delete dpl --project prj -q",
		},
	},
}

func TestDeploymentExecuteCreateUpdateDelete(t *testing.T) {
	for _, tt := range executeTests {
		t.Run(tt.action, func(t *testing.T) {
			var actual []string
			RunGCloud = func(args ...string) (result string, err error) {
				actual = append(actual, strings.Join(args, " "))
				return GetTestData("deployment", "describe-manifest.yaml", t), nil
			}
			_, err := d.Execute(tt.action, false)
			if err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if !reflect.DeepEqual(actual, tt.commands) {
				t.Errorf("got: %s,\nwant: %s", actual, tt.commands)
			}
		})
	}
}

var applyTests = []struct {
	name      string
	status    Status
	statusStr string
	commands  []string
	err       bool
}{
	{
		name:      "status Done",
		status:    Done,
		statusStr: "DONE",
		commands: []string{
			"deployment-manager deployments describe dpl --project prj --format yaml",
			"deployment-manager deployments update dpl --config /tmp/myconfig.yaml --project prj",
			"deployment-manager manifests describe --deployment dpl --project prj --format yaml",
		},
	},
	{
		name:      "status Pending",
		status:    Pending,
		statusStr: "PENDING",
		err:       true,
		commands: []string{
			"deployment-manager deployments describe dpl --project prj --format yaml",
		},
	},
	{
		name:      "status NotFound",
		status:    NotFound,
		statusStr: "",
		err:       false,
		commands: []string{
			"deployment-manager deployments describe dpl --project prj --format yaml",
			"deployment-manager deployments create dpl --config /tmp/myconfig.yaml --project prj",
			"deployment-manager manifests describe --deployment dpl --project prj --format yaml",
		},
	},
}

func TestDeploymentExecuteApply(t *testing.T) {
	for _, tt := range applyTests {
		t.Run(tt.name, func(t *testing.T) {
			var actual []string
			RunGCloud = func(args ...string) (result string, err error) {
				actual = append(actual, strings.Join(args, " "))
				if strings.HasPrefix(actual[len(actual)-1], "deployment-manager deployments describe") {
					if tt.status == NotFound {
						return "ResponseError: code=404, message=", errors.New("some text")
					} else {
						data := GetTestData("deployment", "describe-deployment-template.yaml", t)
						return strings.Replace(data, "%STATUS%", tt.statusStr, 1), nil
					}
				}
				return GetTestData("deployment", "describe-manifest.yaml", t), nil
			}
			_, err := d.Execute(ActionApply, false)
			if err != nil && !tt.err {
				t.Errorf("expected no error, got: %v", err)
			}
			if err == nil && tt.err {
				t.Errorf("expected to have error, got: nil")
			}
			if !reflect.DeepEqual(actual, tt.commands) {
				t.Errorf("got: %s,\nwant: %s", actual, tt.commands)
			}
		})
	}
}
