package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestParseOutputTypesFromState_WithSimpleTypes(t *testing.T) {
	t.Parallel()
	stateData := []byte(`
{
  "version": 4,
  "terraform_version": "0.12.0",
  "serial": 8,
  "lineage": "491fd8f4-81b5-9890-520c-8a173c36e483",
  "outputs": {
    "boolean_output": {
      "value": true,
      "type": "bool"
    },
    "number_output": {
      "value": 42,
      "type": "number"
    },
    "string_output": {
      "value": "foo",
      "type": "string"
    }
  },
  "resources": []
}
`)
	want := map[string]*structpb.Value{
		"boolean_output": structpb.NewStringValue("bool"),
		"number_output":  structpb.NewStringValue("number"),
		"string_output":  structpb.NewStringValue("string"),
	}
	got, err := ParseOutputTypesFromState(stateData)
	if err != nil {
		t.Errorf("ParseOutputTypesFromState() error = %v", err)
		return
	}
	if diff := cmp.Diff(got, want, cmp.Comparer(compareStructpbValues)); diff != "" {
		t.Errorf("ParseOutputTypesFromState() mismatch (-got +want):\n%s", diff)
	}
}

func TestParseOutputTypesFromState_WithComplexTypes(t *testing.T) {
	t.Parallel()
	stateData := []byte(`
{
  "version": 4,
  "terraform_version": "0.12.0",
  "serial": 8,
  "lineage": "491fd8f4-81b5-9890-520c-8a173c36e483",
  "outputs": {
    "interpolated_deep": {
      "value": {
        "foo": "bar",
        "map": {
          "bar": "baz",
          "id": "424881806176056736"
        },
        "number": 42
      },
      "type": [
        "object",
        {
          "foo": "string",
          "map": [
            "object",
            {
              "bar": "string",
              "id": "string"
            }
          ],
          "number": "number"
        }
      ]
    },
    "list": {
      "value": [
        "foo",
        "bar"
      ],
      "type": [
        "tuple",
        [
          "string",
          "string"
        ]
      ]
    },
    "map": {
      "value": {
        "foo": "bar",
        "number": 42
      },
      "type": [
        "object",
        {
          "foo": "string",
          "number": "number"
        }
      ]
    }
  },
  "resources": []
}
`)
	want := map[string]*structpb.Value{
		"interpolated_deep": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStringValue("object"),
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"foo":    structpb.NewStringValue("string"),
				"map":    structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{structpb.NewStringValue("object"), structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{"bar": structpb.NewStringValue("string"), "id": structpb.NewStringValue("string")}})}}),
				"number": structpb.NewStringValue("number"),
			}}),
		}}),
		"list": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStringValue("tuple"),
			structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{structpb.NewStringValue("string"), structpb.NewStringValue("string")}}),
		}}),
		"map": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStringValue("object"),
			structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"foo":    structpb.NewStringValue("string"),
				"number": structpb.NewStringValue("number"),
			}}),
		}}),
	}
	got, err := ParseOutputTypesFromState(stateData)
	if err != nil {
		t.Errorf("ParseOutputTypesFromState() error = %v", err)
		return
	}
	if diff := cmp.Diff(got, want, cmp.Comparer(compareStructpbValues)); diff != "" {
		t.Errorf("ParseOutputTypesFromState() mismatch (-got +want):\n%s", diff)
	}
}

func TestParseOutputTypesFromState_WithoutTypes(t *testing.T) {
	t.Parallel()
	stateData := []byte(`
{
  "version": 4,
  "terraform_version": "0.12.0",
  "serial": 8,
  "lineage": "491fd8f4-81b5-9890-520c-8a173c36e483",
  "outputs": {
    "no_type_output": {
      "value": "some_value"
    }
  },
  "resources": []
}
`)
	want := map[string]*structpb.Value{
		"no_type_output": structpb.NewNullValue(), // Expecting null value when type is missing
	}

	got, err := ParseOutputTypesFromState(stateData)
	if err != nil {
		t.Errorf("ParseOutputTypesFromState() error = %v", err)
		return
	}
	if diff := cmp.Diff(got, want, cmp.Comparer(compareStructpbValues)); diff != "" {
		t.Errorf("ParseOutputTypesFromState() mismatch (-got +want):\n%s", diff)
	}
}

// compareStructpbValues is a custom comparer for structpb.Value
func compareStructpbValues(x, y *structpb.Value) bool {
	// Marshal to JSON and compare the JSON strings
	xJSON, _ := x.MarshalJSON()
	yJSON, _ := y.MarshalJSON()
	return string(xJSON) == string(yJSON)
}
