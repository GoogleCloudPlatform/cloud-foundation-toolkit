package validation

import (
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// protoViaJSON uses JSON as an intermediary serialization to convert a value into
// a protobuf message.
func protoViaJSON(from interface{}, to proto.Message) error {
	jsn, err := json.Marshal(from)
	if err != nil {
		return errors.Wrap(err, "marshaling to json")
	}

	if err := jsonpb.UnmarshalString(string(jsn), to); err != nil {
		return errors.Wrap(err, "unmarshaling to proto")
	}

	return nil
}
