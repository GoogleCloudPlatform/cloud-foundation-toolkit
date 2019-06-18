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
	expect := "gl-akopachevsky-test"
	if deployment["name"] != expect {
		t.Errorf("got: %s, want: %s", deployment["name"], expect)
	}
}

func TestReparentPath(t *testing.T) {
	var pathtests = []struct {
		parent string
		child  string
		out    string
	}{
		{"/base/folder/config.yaml", "script.py", "/base/folder/script.py"},
		{"/base/folder/config.yaml", "./script.py", "/base/folder/script.py"},
		{"/base/folder/config.yaml", "../script.py", "/base/script.py"},
		{"/base/folder/config.yaml", "/base/folder/script.py", "/base/folder/script.py"},
		{"/base/folder/config.yaml", "templates/script.py", "/base/folder/templates/script.py"},
		{"/base/folder/config.yaml", "../templates/script.py", "/base/templates/script.py"},
	}

	for _, tt := range pathtests {
		t.Run(tt.parent+"  "+tt.child, func(t *testing.T) {
			actual := ReparentPath(tt.parent, tt.child)
			if actual != tt.out {
				t.Errorf("got: %s, want: %s", actual, tt.out)
			}
		})
	}
}

// GetTestData returns file content of PROJECT_ROOT/testdata sub folder or file.
// Function suppose to call from test method, might not work otherwise
// example GetTestData("myfolder", "file.yalm") will return content of testdata/myfolder/file.yaml
func GetTestData(folder string, name string, t *testing.T) string {
	path := filepath.Join("../testdata", folder, name)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("unexpected error reading test file: %v", err)
	}
	return string(buf)
}
