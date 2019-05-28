package deployment

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	res, err := unmarshal(GetTestData("deployment", "api-client-describe.yaml", t))
	if err != nil {
		t.Errorf("Error %v", err)
	}
	deployment := res["deployment"].(map[interface{}]interface{})
	if deployment["name"] != "gl-akopachevsky-test" {
		t.Errorf("Expected \"gl-akopachevsky-test\", got \"%s\"", deployment["name"])
	}
}

func GetTestData(folder string, name string, t *testing.T) string {
	path := filepath.Join("../testdata", folder, name)
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	return string(buff)
}
