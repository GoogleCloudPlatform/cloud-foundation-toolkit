package deployment

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	res, err := unmarshal(GetTestData("deployment", "describe-deployment.yaml", t))
	if err != nil {
		t.Errorf("Error %v", err)
	}
	deployment := res["deployment"].(map[interface{}]interface{})
	expect := "gl-akopachevsky-test"
	if deployment["name"] != expect {
		t.Errorf("got: %s, want: %s", deployment["name"], expect)
	}
}

var pathtests = []struct {
	parent string
	child  string
	out    string
}{
	{"/base/folder", "script.py", "/base/folder/script.py"},
	{"/base/folder", "./script.py", "/base/folder/script.py"},
	{"/base/folder", "../script.py", "/base/script.py"},
	{"/base/folder", "/other/script.py", "/other/script.py"},
	{"/base/folder", "templates/script.py", "/base/folder/templates/script.py"},
	{"/base/folder", "../templates/script.py", "/base/templates/script.py"},
}

func TestAbsolutePath(t *testing.T) {
	for _, tt := range pathtests {
		t.Run(tt.parent+"  "+tt.child, func(t *testing.T) {
			actual := AbsolutePath(tt.parent, tt.child)
			if actual != tt.out {
				t.Errorf("got: %s, want: %s", actual, tt.out)
			}
		})
	}
}

var fileNameTests = []struct {
	file string
	name string
}{
	{"name", "name"},
	{"name.txt", "name"},
	{"../name.txt", "name"},
	{"./name.txt", "name"},
	{"/test/name.txt", "name"},
	{"UPPERCASE.yaml", "uppercase"},
	{"under_scores_.yaml", "under-scores"},
	{"last-dash-.yaml", "last-dash"},
	{"-8first-dash-and-number.yaml", "first-dash-and-number"},
	{"8-number-followed-by-dash.yaml", "number-followed-by-dash"},
	{"--double-dash.yaml", "double-dash"},
	{"last-dash-.yaml", "last-dash"},
	{"last-dobule-dash--.yaml", "last-dobule-dash"},
	{"more-than-63-chars--------------------------------63-ends-here-----.yaml", "more-than-63-chars--------------------------------63-ends-here"},
}

func TestDeploymentNameFromFile(t *testing.T) {
	for _, tt := range fileNameTests {
		t.Run(tt.file, func(t *testing.T) {
			actual := DeploymentNameFromFile(tt.file)
			if actual != tt.name {
				t.Errorf("got: %s, want: %s", actual, tt.name)
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
