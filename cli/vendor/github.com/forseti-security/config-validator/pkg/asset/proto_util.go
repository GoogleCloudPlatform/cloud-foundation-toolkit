package asset

import (
	structpb "github.com/golang/protobuf/ptypes/struct"
)

// CleanProtoValue recursively updates proto Values that have a nil .Kind field
// to be a NullValue to avoid issues with the jsonpb.Marshaler.
// This issue arose when calling GCV from python.
func CleanProtoValue(v *structpb.Value) {
	if v == nil {
		return
	}
	switch t := v.Kind.(type) {
	case *structpb.Value_NullValue, *structpb.Value_NumberValue, *structpb.Value_StringValue, *structpb.Value_BoolValue:
	case *structpb.Value_StructValue:
		CleanStructValue(t.StructValue)
	case *structpb.Value_ListValue:
		if list := t.ListValue; list != nil {
			for i := range list.Values {
				CleanProtoValue(list.Values[i])
			}
		}
	default: // No other kinds should be allowed (including nil).
		v.Kind = &structpb.Value_NullValue{}
	}
}

func CleanStructValue(s *structpb.Struct) {
	if s == nil {
		return
	}
	for k := range s.Fields {
		CleanProtoValue(s.Fields[k])
	}
}
