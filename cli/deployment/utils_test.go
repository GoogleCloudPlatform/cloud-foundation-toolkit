package deployment

import (
	"io/ioutil"
	"path/filepath"
	"strings"
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
	{"../testdata/reparent_path/folder1/config.yaml", "script1.py", "testdata/reparent_path/folder1/script1.py"},
	{"../testdata/reparent_path/folder1/config.yaml", "./script1.py", "testdata/reparent_path/folder1/script1.py"},
	{"../testdata/reparent_path/folder1/config.yaml", "../script4.py", "testdata/reparent_path/script4.py"},
	{"../testdata/reparent_path/folder1/config.yaml", "templates/script3.py", "testdata/reparent_path/folder1/templates/script3.py"},
	{"../testdata/reparent_path/folder1/config.yaml", "../folder2/script2.py", "testdata/reparent_path/folder2/script2.py"},
}

func TestReparentPath(t *testing.T) {
	basePath, _ := filepath.Abs("..")
	for _, tt := range pathtests {
		t.Run(tt.parent+"  "+tt.child, func(t *testing.T) {
			path := filepath.Join(basePath, tt.out)
			actual := ReparentPath(tt.parent, tt.child)
			if actual != path {
				t.Errorf("got: %s,\n want: %s", actual, path)
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
	{"space in name.yaml", "spaceinname"},
	{"my.config.yaml", "myconfig"},
	{"forbidden chars,.?\\)(*&^%$#@@!\\!`\";<>.yaml", "forbiddenchars"},
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

var inputTests = []struct {
	input    string
	expected string
}{
	{"a\n", "a"},
	{"a \n", "a"},
	{" a \n", "a"},
	{"a\r\n", "a"},
	{"u\n", "u"},
	{"s\n", "s"},
}

func TestGetUserInput(t *testing.T) {
	for _, tt := range inputTests {
		t.Run(tt.input, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			actual := GetUserInput("Update(u), Skip (s), or Abort(a) Deployment?", []string{"u", "s", "a"}, input)
			if actual != tt.expected {
				t.Errorf("got: %s, want: %s", actual, tt.expected)
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
