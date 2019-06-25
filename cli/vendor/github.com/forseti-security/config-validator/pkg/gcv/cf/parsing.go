package cf

import (
	"encoding/json"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/golang/protobuf/jsonpb"
	pb "github.com/golang/protobuf/ptypes/struct"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
)

// expressionVal is patterned off of the response object provided by the audit script.
// The purpose of this object is to be able to to parse the generic result provided by Rego using json parsing.
type expressionVal struct {
	Asset      string `json:"asset"`
	Constraint string `json:"constraint"`
	Violation  *struct {
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
	violations := []*validator.Violation{}
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
		violations = append(violations, violationToAdd)
	}
	return violations, nil
}

func convertToProtoVal(from interface{}) (*pb.Value, error) {
	to := &pb.Value{}
	jsn, err := json.Marshal(from)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling to json")
	}

	if err := jsonpb.UnmarshalString(string(jsn), to); err != nil {
		return nil, errors.Wrap(err, "unmarshalling to proto")
	}

	return to, nil
}
