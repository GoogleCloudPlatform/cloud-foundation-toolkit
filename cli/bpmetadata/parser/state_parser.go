package parser

import (
	"bytes"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func ParseOutputTypesFromState(stateData []byte) (map[string]*structpb.Value, error) {
	// Unmarshal the state data into a map[string]interface{} first
	var rawState map[string]interface{}
	err := json.Unmarshal(stateData, &rawState)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state data: %w", err)
	}

	// Check if "format_version" key exists
	if _, ok := rawState["format_version"]; !ok {
		// If not present, add it with a default value
		rawState["format_version"] = "1.0"
	}

	// Create a "values" field with "outputs" inside
	rawState["values"] = map[string]interface{}{
		"outputs": rawState["outputs"],
	}

	// Remove the top-level "outputs" field
	delete(rawState, "outputs")

	// Now marshal the updated map back to JSON
	updatedStateData, err := json.MarshalIndent(rawState, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated state data: %w", err)
	}

	var state tfjson.State

	// Unmarshal the updated JSON into tfjson.State
	err = json.Unmarshal(updatedStateData, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state data: %w", err)
	}

	outputTypeMap := make(map[string]*structpb.Value)
	for name, output := range state.Values.Outputs {
		pbValue, err := convertOutputTypeToStructpb(output)
		if err != nil {
			return nil, fmt.Errorf("failed to convert output %q to structpb.Value: %w", name, err)
		}
		outputTypeMap[name] = pbValue
	}

	return outputTypeMap, nil
}

func convertOutputTypeToStructpb(output *tfjson.StateOutput) (*structpb.Value, error) {
	// Handle nil values explicitly
	if output.Value == nil {
		return structpb.NewNullValue(), nil
	}

	// Handle cases where output.Type is NilType
	if output.Type == cty.NilType {
		return structpb.NewNullValue(), nil
	}

	// Marshal the output value to JSON
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(output.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal output type to JSON: %w", err)
	}

	// Unmarshal the JSON into a structpb.Value
	pbValue := &structpb.Value{}
	err = pbValue.UnmarshalJSON(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON into structpb.Value: %w", err)
	}

	return pbValue, nil
}
