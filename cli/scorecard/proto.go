// Copyright 2019-2023 Google LLC
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

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func unMarshallAsset(from []byte, to proto.Message) error {
	// CAI export returns org_policy [1] with update_time if Timestamp format in Seconds and Nanos
	// but in protojson, Timestamp is expected to be a string in the RFC 3339 format [2].
	// i.e. "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
	// Hence doing a workaround to remove the field so that protojson.Unmarshaler can handle org policy.
	// [1] https://github.com/googleapis/googleapis/blob/master/google/cloud/orgpolicy/v1/orgpolicy.proto
	// [2] https://godoc.org/google.golang.org/protobuf/types/known/timestamppb#Timestamp

	// Using json.Unmarshal will return no error
	// but this approach will lose the "oneof" proto fields in org_policy and access_policy

	var temp map[string]interface{}
	err := json.Unmarshal(from, &temp)
	if err != nil {
		return errors.Wrap(err, "marshaling to interface")
	}
	if val, ok := temp["org_policy"]; ok {
		for _, op := range val.([]interface{}) {
			orgPolicy := op.(map[string]interface{})
			delete(orgPolicy, "update_time")
		}
	}
	err = protoViaJSON(temp, to)
	if err == nil {
		return nil
	}
	return err
}

// protoViaJSON uses JSON as an intermediary serialization to convert a value into
// a protobuf message.
func protoViaJSON(from interface{}, to proto.Message) error {
	if m, ok := from.(map[string]interface{}); ok {
		if val, ok := m["iam_policy"]; ok {
			if m, ok := val.(map[string]interface{}); ok {
				delete(m, "etag")
			}
		}
	}
	jsn, err := json.Marshal(from)
	if err != nil {
		return errors.Wrap(err, "marshaling to json")
	}
	umar := &protojson.UnmarshalOptions{DiscardUnknown: true}
	if err := umar.Unmarshal(jsn, to); err != nil {
		return errors.Wrap(err, "unmarshaling to proto")
	}

	return nil
}

// interfaceViaJSON uses JSON as an intermediary serialization to convert a protobuf message
// into an interface value
func interfaceViaJSON(from proto.Message) (interface{}, error) {
	jsn, err := protojson.Marshal(from)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling to json")
	}

	var to interface{}
	if err := json.Unmarshal(jsn, &to); err != nil {
		return nil, errors.Wrap(err, "unmarshaling to interface")
	}

	return to, nil
}

// stringViaJSON uses JSON as an intermediary serialization to convert a protobuf message
// into an string value
func stringViaJSON(from proto.Message) (string, error) {
	jsn, err := protojson.Marshal(from)
	if err != nil {
		return "", errors.Wrap(err, "marshaling to json")
	}
	str, err := strconv.Unquote(string(jsn))
	if err != nil {
		// return original json string if it's not a quoted string
		return string(jsn), nil
	}
	return str, nil
}
