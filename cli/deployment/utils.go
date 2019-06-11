package deployment

import (
	"fmt"
	"os"
	"path/filepath"

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

func ReparentPath(parent string, child string) string {
	// check if file already has absolute path
	if child[0] == os.PathSeparator {
		return child
	}
	dir := filepath.Dir(parent)
	return filepath.Clean(filepath.Join(dir, child))
}

// check if string is yaml
func IsYAML(text string) bool {
	obj := struct{}{}
	err := yaml.Unmarshal([]byte(text), obj)
	return err == nil
}

func AppendMap(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}
