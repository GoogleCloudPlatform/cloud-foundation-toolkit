package cf

import (
	"encoding/json"
	"reflect"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	pb "github.com/golang/protobuf/ptypes/struct"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
)

// expressionVal is patterned off of the response object provided by the audit script.
// The purpose of this object is to be able to to parse the generic result provided by Rego using json parsing.
type expressionVal struct {
	Asset            string                 `json:"asset"`
	Constraint       string                 `json:"constraint"`
	ConstraintConfig map[string]interface{} `json:"constraint_config"`
	Violation        *struct {
		Msg      string                 `json:"msg"`
		Metadata map[string]interface{} `json:"details"`
	} `json:"violation"`
}

func validateExpression(expr *expressionVal) error {
	switch {
	case expr.Asset == "":
		return errors.New("No asset field found")
	case expr.Constraint == "":
		return errors.New("No constraint field found")
	case expr.Violation == nil:
		return errors.New("No violation field found")
	case expr.Violation.Msg == "":
		return errors.New("No violation.msg field found")
	default:
		return nil
	}
}

func parseExpression(expression *rego.ExpressionValue) ([]*expressionVal, error) {
	jsonBytes, err := json.Marshal(expression.Value)
	if err != nil {
		return nil, err
	}
	var ret []*expressionVal
	if err := json.Unmarshal(jsonBytes, &ret); err != nil {
		return nil, err
	}
	// Validate fields
	for _, expr := range ret {
		if validationError := validateExpression(expr); validationError != nil {
			return nil, validationError
		}
	}
	return ret, nil
}

func convertToViolations(expression *rego.ExpressionValue) ([]*validator.Violation, error) {
	parsedExpression, err := parseExpression(expression)
	if err != nil {
		return nil, err
	}
	var violations []*validator.Violation
	for i := 0; i < len(parsedExpression); i++ {
		violationToAdd := &validator.Violation{
			Constraint: parsedExpression[i].Constraint,
			Resource:   parsedExpression[i].Asset,
			Message:    parsedExpression[i].Violation.Msg,
		}
		if parsedExpression[i].Violation.Metadata != nil {
			convertedMetadata, err := convertToProtoVal(parsedExpression[i].Violation.Metadata)
			if err != nil {
				return nil, err
			}
			violationToAdd.Metadata = convertedMetadata
		}
		if parsedExpression[i].ConstraintConfig != nil {
			constraintMetadata, err := convertToProtoVal(parsedExpression[i].ConstraintConfig["metadata"])
			if err != nil {
				return nil, err
			}
			violationToAdd.ConstraintConfig = &validator.Constraint{
				Metadata: constraintMetadata,
			}
		}

		violations = append(violations, violationToAdd)
	}
	return violations, nil
}

type convertFailed struct {
	err error
}

func convertToProtoVal(from interface{}) (val *pb.Value, err error) {
	defer func() {
		if x := recover(); x != nil {
			convFail, ok := x.(*convertFailed)
			if !ok {
				panic(x)
			}
			val = nil
			err = errors.Errorf("failed to convert proto val: %s", convFail.err)
		}
	}()
	val = convertToProtoValInternal(from)
	return
}

func convertToProtoValInternal(from interface{}) *pb.Value {
	if from == nil {
		return nil
	}
	switch val := from.(type) {
	case map[string]interface{}:
		fields := map[string]*pb.Value{}
		for k, v := range val {
			fields[k] = convertToProtoValInternal(v)
		}
		return &pb.Value{
			Kind: &pb.Value_StructValue{
				StructValue: &pb.Struct{
					Fields: fields,
				},
			}}

	case []interface{}:
		vals := make([]*pb.Value, len(val))
		for idx, v := range val {
			vals[idx] = convertToProtoValInternal(v)
		}
		return &pb.Value{
			Kind: &pb.Value_ListValue{
				ListValue: &pb.ListValue{Values: vals},
			},
		}

	case string:
		return &pb.Value{Kind: &pb.Value_StringValue{StringValue: val}}
	case int:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: float64(val)}}
	case int64:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: float64(val)}}
	case float64:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: val}}
	case float32:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: float64(val)}}
	case bool:
		return &pb.Value{Kind: &pb.Value_BoolValue{BoolValue: val}}

	default:
		panic(&convertFailed{errors.Errorf("Unhandled type %v (%s)", from, reflect.TypeOf(from).String())})
	}
}
