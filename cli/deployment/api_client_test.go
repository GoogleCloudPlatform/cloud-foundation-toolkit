package deployment

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestUnmarshal(t *testing.T) {

	res, err := unmarshal(testData("api-client-describe.yaml", t))
	if err != nil {
		t.Errorf("Error %v", err)
	}
	deployment := res["deployment"].(map[interface{}]interface{})
	if deployment["name"] != "gl-akopachevsky-test" {
		t.Errorf("Expected \"gl-akopachevsky-test\", got \"%s\"", deployment["name"])
	}
}

func testGCloud(t *testing.T) {
	res, err := GCloud("deployment-manager", "deployments", "describe", "gl-akopachevsky-test")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestGetOutputs(t *testing.T) {
	outputs, err := GetOutputs(testData("api-describe-manifest.yaml", t))
	if err != nil {
		t.Error(err)
	}
	if "my-network-prod" != outputs["my-network-prod.name"] {
		t.Errorf("Expected \"my-network-prod\" got \"%s\"", outputs["my-network-prod.name"])
	}
}

func testData(name string, t *testing.T) string {
	path := filepath.Join("../testdata", "deployment", name)
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	return string(buff)
}
