package deployment

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"hash/fnv"
	"os"
	"path/filepath"
)

func hash64(s string) int64 {
	hash := fnv.New32()
	hash.Write([]byte(s))
	return int64(hash.Sum32())
}

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

func AbsolutePath(parent string, child string) string {
	// check if file already has absolute path
	if child[0] == os.PathSeparator {
		return child
	}
	dir := filepath.Dir(parent)
	return filepath.Clean(filepath.Join(dir, child))
}
