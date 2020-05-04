// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scorecard

import (
	"encoding/json"
	"strconv"
	"strings"

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
	umar := &jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err := umar.Unmarshal(strings.NewReader(string(jsn)), to); err != nil {
		return errors.Wrap(err, "unmarshaling to proto")
	}

	return nil
}

// interfaceViaJSON uses JSON as an intermediary serialization to convert a protobuf message
// into an interface value
func interfaceViaJSON(from proto.Message) (interface{}, error) {
	marshaler := &jsonpb.Marshaler{}
	jsn, err := marshaler.MarshalToString(from)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling to json")
	}

	var to interface{}
	if err := json.Unmarshal([]byte(jsn), &to); err != nil {
		return nil, errors.Wrap(err, "unmarshaling to interface")
	}

	return to, nil
}

// stringViaJSON uses JSON as an intermediary serialization to convert a protobuf message
// into an string value
func stringViaJSON(from proto.Message) (string, error) {
	marshaler := &jsonpb.Marshaler{}
	jsn, err := marshaler.MarshalToString(from)
	if err != nil {
		return "", errors.Wrap(err, "marshaling to json")
	}
	str, err := strconv.Unquote(jsn)
	if err != nil {
		// return original json string if it's not a quoted string
		return jsn, nil
	}
	return str, nil
}
