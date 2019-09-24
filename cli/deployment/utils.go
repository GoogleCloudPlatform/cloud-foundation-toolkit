package deployment

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// forbiddenCharRegexp searches all chars except allowed in deployment name.
var forbiddenCharRegexp = regexp.MustCompile("[^-a-z0-9]")

// firstCharRegexp is regexp that used to find number or dash at the beginning of the string
var firstCharRegexp = regexp.MustCompile("^[-0-9]*")

// lastCharRegexp is regexp that used to find trailing dash
var lastCharRegexp = regexp.MustCompile("-*$")

// Function unmarshal arbitrary YAML to map.
func unmarshal(data string) (map[string]interface{}, error) {
	my := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(data), my)
	if err != nil {
		fmt.Printf("error %+v", err)
		return nil, err
	}
	return my, nil
}

// Function ReparentPath used to create absolute path for config "import" entries.
// Absolute path composed by ReparentPath function from config file path base folder concatenated with
// import statement value. This transformation needed to make deployment YAML config file location independent,
// after import "absolutisation" deployment config might be copied to any location (as current CFT cli will copy its
// copy to tmp folder.
// Examples:
//    /base/folder/config.yaml and ../script.py will concatenate to /base/script.py
//    /base/folder/config.yaml and /base/folder/script.py will concatenate to /base/folder/script.py as long path already absolute
func ReparentPath(baseDir string, file string) string {
	// check if import statement path already absolute
	if file[0] == os.PathSeparator {
		return file
	}
	baseDir, _ = filepath.Abs(baseDir)
	baseDirStat, err := os.Stat(baseDir)
	if err != nil {
		log.Fatalf("error set stat for path: %s, error: %v", baseDir, err)
	}

	if !baseDirStat.Mode().IsDir() {
		baseDir = filepath.Dir(baseDir)
	}

	relative := filepath.Clean(filepath.Join(baseDir, file))
	result, err := filepath.Abs(relative)
	if err != nil {
		log.Fatalf("error creating absolute path, for file: %s, error: %v", relative, err)
	}
	return result
}

// IsYAML checks if text parameter passed if YAML string.
func IsYAML(text string) bool {
	obj := struct{}{}
	err := yaml.Unmarshal([]byte(text), obj)
	return err == nil
}

// AppendMap appends string map B to map A, returns A
func AppendMap(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}

// DeploymentNameFromFile creates valid deployment name from file path satisfied Deployment resource "name" field requirements:
// Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])?
// which means the first character must be a lowercase letter,
// and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
// Function cuts forbidden characters at the end and start of a name also it truncates name up to 63 characters.
// see more https://cloud.google.com/deployment-manager/docs/reference/latest/deployments#resource
func DeploymentNameFromFile(path string) string {
	_, file := filepath.Split(path)
	ext := filepath.Ext(file)
	name := strings.TrimSuffix(file, ext)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", "-")
	name = firstCharRegexp.ReplaceAllString(name, "")
	name = lastCharRegexp.ReplaceAllString(name, "")
	name = forbiddenCharRegexp.ReplaceAllString(name, "")
	if len(name) > 63 {
		name = name[0:63]
	}
	return name
}

// GetUserInput asks for user input and validate entered value is equal to one of the provided options.
// Returns validated option string
func GetUserInput(message string, options []string, rd io.Reader) string {
	reader := bufio.NewReader(rd)
	var input string
	var err error
	for !stringInSlice(input, options) {
		log.Print(message + " ")
		input, err = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			log.Fatalf("failed to get user input, error: %v", err)
		}
	}
	return input
}

// Checks if string a in slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
