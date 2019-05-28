package deployment

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"hash/fnv"
)

func hash64(s string) int64 {
	hash := fnv.New32()
	hash.Write([]byte(s))
	return int64(hash.Sum32())
}

// unmarshal arbitraty yaml to map
func unmarshal(data string) (map[string]interface{}, error) {
	my := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(data), my)
	if err != nil {
		fmt.Printf("error %+v", err)
		return nil, err
	}
	return my, nil
}
