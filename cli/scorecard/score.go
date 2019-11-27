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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"
	"github.com/pkg/errors"
)

// ScoringConfig holds settings for generating a score
type ScoringConfig struct {
	PolicyPath  string                           // the directory path of a policy library to use
	categories  map[string]*constraintCategory   // available constraint categories
	constraints map[string]*constraintViolations // a map of constraints violated and their violations
	validator   *gcv.Validator                   // the validator instance used for scoring
}

// NewScoringConfig creates a scoring engine for the given policy library
func NewScoringConfig(policyPath string) (*ScoringConfig, error) {
	config := &ScoringConfig{}

	config.PolicyPath = policyPath

	v, err := gcv.NewValidator(
		gcv.PolicyPath(filepath.Join(config.PolicyPath, "policies")),
		gcv.PolicyLibraryDir(filepath.Join(config.PolicyPath, "lib")),
	)
	if err != nil {
		return nil, errors.Wrap(err, "initializing gcv validator")
	}
	config.validator = v

	return config, nil
}

const otherCategoryKey = "other"

// constraintCategory holds constraints by category
type constraintCategory struct {
	Name        string
	constraints []*constraintViolations
}

func (c constraintCategory) Count() int {
	sum := 0
	for _, cv := range c.constraints {
		sum += cv.Count()
	}
	return sum
}

// constraintViolations holds violations for a particular constraint
type constraintViolations struct {
	constraint *validator.Constraint
	Violations []*validator.Violation `protobuf:"bytes,1,rep,name=violations,proto3" json:"violations,omitempty"`
}

func (cv constraintViolations) Count() int {
	return len(cv.Violations)
}

func (cv constraintViolations) GetName() string {
	return cv.constraint.GetMetadata().GetStructValue().GetFields()["name"].GetStringValue()
}

var availableCategories = map[string]string{
	"operational-efficiency": "Operational Efficiency",
	"security":               "Security",
	"reliability":            "Reliability",
	otherCategoryKey:         "Other",
}

func (config *ScoringConfig) getConstraintForViolation(violation *validator.Violation) (*constraintViolations, error) {
	key := violation.GetConstraint()
	cv, found := config.constraints[key]
	if !found {
		constraint := violation.GetConstraintConfig()
		cv = &constraintViolations{
			constraint: constraint,
		}
		config.constraints[key] = cv

		metadata := constraint.GetMetadata()
		annotations := metadata.GetStructValue().GetFields()["annotations"].GetStructValue().GetFields()

		categoryKey := otherCategoryKey
		categoryValue, found := annotations["bundles.validator.forsetisecurity.org/scorecard-v1"]
		if found {
			categoryKey = categoryValue.GetStringValue()
		}

		category, found := config.categories[categoryKey]
		if !found {
			return nil, fmt.Errorf("Unknown constraint category %v for constraint %v", categoryKey, key)
		}
		category.constraints = append(category.constraints, cv)
	}
	return cv, nil
}

// attachViolations puts violations into their appropriate categories
func (config *ScoringConfig) attachViolations(audit *validator.AuditResponse) error {
	// Build map of categories
	config.categories = make(map[string]*constraintCategory)
	for k, name := range availableCategories {
		config.categories[k] = &constraintCategory{
			Name: name,
		}
	}

	// Categorize violations
	config.constraints = make(map[string]*constraintViolations)
	for _, v := range audit.Violations {
		cv, err := config.getConstraintForViolation(v)
		if err != nil {
			return errors.Wrap(err, "Categorizing violation")
		}

		cv.Violations = append(cv.Violations, v)
	}

	return nil
}

// Score creates a Scorecard for an inventory
func (inventory *InventoryConfig) Score(config *ScoringConfig, outputPath string, outputFormat string) error {
	auditResult, err := getViolations(inventory, config)
	if err != nil {
		return err
	}

	err = config.attachViolations(auditResult)

	if len(auditResult.Violations) > 0 {
		content := ""
		switch outputFormat {
		case "json":
			type violationOutput struct {
				Category   string
				Constraint string
				Resource   string
				Message    string
			}
			var vOutput violationOutput
			for _, category := range config.categories {
				for _, cv := range category.constraints {
					for _, v := range cv.Violations {
						vOutput.Category = category.Name
						vOutput.Constraint = v.Constraint
						vOutput.Resource = v.Resource
						vOutput.Message = v.Message
						byteContent, err := json.MarshalIndent(vOutput, "", "  ")
						if err != nil {
							return err
						}
						content += string(byteContent) + "\n"
						Log.Debug("Violation metadata", "metadata", v.GetMetadata())
					}
				}
			}
		case "csv":
			content = "Category,Constraint,Resource,Message\n"
			for _, category := range config.categories {
				for _, cv := range category.constraints {
					for _, v := range cv.Violations {
						content += fmt.Sprintf("%v,%v,%v,%v\n", category.Name, v.Constraint, v.Resource, v.Message)
						Log.Debug("Violation metadata", "metadata", v.GetMetadata())
					}
				}
			}
		default:
			outputFormat = "txt"
			content = fmt.Sprintf("\n\n%v total issues found\n", len(auditResult.Violations))
			for _, category := range config.categories {
				content += fmt.Sprintf("\n\n%v: %v issues found\n", category.Name, category.Count())
				content += fmt.Sprintf("----------\n")
				for _, cv := range category.constraints {
					content += fmt.Sprintf("%v: %v issues\n", cv.GetName(), cv.Count())
					for _, v := range cv.Violations {
						content += fmt.Sprintf("- %v\n\n",
							v.Message,
						)
						Log.Debug("Violation metadata", "metadata", v.GetMetadata())
					}
				}
			}
		}
		if outputPath == "" {
			fmt.Printf("%v", content)
		} else {
			outputFile := "scorecard." + outputFormat
			err = ioutil.WriteFile(filepath.Join(outputPath, outputFile), []byte(content), 0644)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println("No issues found found! You have a perfect score.")
	}

	return nil
}
