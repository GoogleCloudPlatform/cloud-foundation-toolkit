package google

import (
	"encoding/json"
	"fmt"
)

// jsonMap converts a given value to a map[string]interface{} that
// matches its JSON format.
func jsonMap(x interface{}) (map[string]interface{}, error) {
	jsn, err := json.Marshal(x)
	if err != nil {
		return nil, fmt.Errorf("marshalling: %v", err)
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(jsn, &m); err != nil {
		return nil, fmt.Errorf("unmarshalling: %v", err)
	}

	return m, nil
}
