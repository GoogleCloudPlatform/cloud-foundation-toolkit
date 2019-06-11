package deployment

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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
//    /base/folder/config.YAML and ../script.py will concatenate to /base/script.py
//    /base/folder/config.YAML and /base/folder/script.py will concatenate to /base/folder/script.py as long path already absolute
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
