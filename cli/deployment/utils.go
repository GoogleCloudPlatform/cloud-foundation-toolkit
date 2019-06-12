package deployment

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// unmarshal arbitrary yaml to map
func unmarshal(data string) (map[string]interface{}, error) {
	my := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(data), my)
	if err != nil {
		fmt.Printf("error %+v", err)
		return nil, err
	}
	return my, nil
}

func AbsolutePath(baseDir string, file string) string {
	// check if import statement path absolute
	if file[0] == os.PathSeparator {
		return file
	}
	return filepath.Clean(filepath.Join(baseDir, file))
}

// check if string is yaml
func IsYAML(text string) bool {
	obj := struct{}{}
	err := yaml.Unmarshal([]byte(text), obj)
	return err == nil
}

// append string map B to map A, returns A
func AppendMap(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}

/*
creates valid deployment name from file path satisfied Deployment resource "name" field requirements:
Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])?
which means the first character must be a lowercase letter,
and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
see more https://cloud.google.com/deployment-manager/docs/reference/latest/deployments#resource
*/
func DeploymentNameFromFile(path string) string {
	_, file := filepath.Split(path)
	name := strings.TrimSuffix(file, filepath.Ext(file))
	name = strings.ToLower(name)
	if len(name) > 63 {
		name = name[0:63]
	}
	name = strings.ReplaceAll(name, "_", "-")
	firstChar := regexp.MustCompile("^[-0-9]*")
	lastChar := regexp.MustCompile("-*$")
	name = firstChar.ReplaceAllString(name, "")
	name = lastChar.ReplaceAllString(name, "")
	return name
}
