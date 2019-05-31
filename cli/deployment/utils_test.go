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

func TestAbsolutePath(t *testing.T) {
	expected := "/base/folder/script.py"
	actual := AbsolutePath("/base/folder/config.yaml", "script.py")
	if actual != expected {
		t.Errorf("Expected %s, actual %s", expected, actual)
	}

	actual = AbsolutePath("/base/folder/config.yaml", "./script.py")
	if actual != expected {
		t.Errorf("Expected %s, actual %s", expected, actual)
	}

	actual = AbsolutePath("/base/folder/config.yaml", expected)
	if actual != expected {
		t.Errorf("Expected %s, actual %s", expected, actual)
	}

	expected = "/base/script.py"
	actual = AbsolutePath("/base/folder/config.yaml", "../script.py")
	if actual != expected {
		t.Errorf("Expected %s, actual %s", expected, actual)
	}

	expected = "/base/folder/templates/script.py"
	actual = AbsolutePath("/base/folder/config.yaml", "templates/script.py")
	if actual != expected {
		t.Errorf("Expected %s, actual %s", expected, actual)
	}

	expected = "/base/templates/script.py"
	actual = AbsolutePath("/base/folder/config.yaml", "../templates/script.py")
	if actual != expected {
		t.Errorf("Expected %s, actual %s", expected, actual)
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
