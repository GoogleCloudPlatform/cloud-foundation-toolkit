package validation

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
)

func TestValidateDeployment(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	type defI struct {
		name       string
		policyPath string
		project    string
	}
	type defO struct {
		validated bool
		err       bool
	}

	cases := []struct {
		name   string
		input  defI
		output defO
	}{
		{
			name: "no policy",
			input: defI{
				"",
				"",
				"",
			},
			output: defO{
				validated: false,
				err:       true,
			},
		},
		{
			name: "deployment error",
			input: defI{
				"test1",
				"foobar",
				"demo",
			},
			output: defO{
				validated: false,
				err:       true,
			},
		},
	}

	type descRet struct {
		desc *deployment.DeploymentDescription
		err  error
	}
	descs := map[string]descRet{
		"test1__demo": descRet{
			desc: nil,
			err:  errors.New("test"),
		},
	}
	getDeploymentDescription = func(name string, project string) (description *deployment.DeploymentDescription, e error) {
		desc := descs[name+"__"+project]
		return desc.desc, desc.err
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			validated, err := ValidateDeployment(c.input.name, c.input.policyPath, c.input.project)

			if err != nil && !c.output.err || err == nil && c.output.err {
				t.Errorf("got err=%v, expected present err=%v", err, c.output.err)
			}
			if validated != c.output.validated {
				t.Errorf("got validated=%v, expected validated=%v", validated, c.output.validated)
			}
		})
	}
}
