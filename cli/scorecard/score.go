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
	"context"
	"crypto/md5"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/pkg/errors"
)

// ScoringConfig holds settings for generating a score
type ScoringConfig struct {
	categories  map[string]*constraintCategory   // available constraint categories
	constraints map[string]*constraintViolations // a map of constraints violated and their violations
	validator   *gcv.Validator                   // the validator instance used for scoring
}

// NewScoringConfigFromValidator creates a scoring engine with a given validator.
func NewScoringConfigFromValidator(v *gcv.Validator) *ScoringConfig {
	config := &ScoringConfig{}
	config.validator = v
	return config
}

// NewScoringConfig creates a scoring engine for the given policy library
func NewScoringConfig(ctx context.Context, policyPath string) (*ScoringConfig, error) {
	flag.Parse()
	v, err := gcv.NewValidator(
		[]string{filepath.Join(policyPath, "policies")},
		filepath.Join(policyPath, "lib"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "initializing gcv validator")
	}
	config := NewScoringConfigFromValidator(v)
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
	constraint string
	Violations []*validator.Violation `protobuf:"bytes,1,rep,name=violations,proto3" json:"violations,omitempty"`
}

func (cv constraintViolations) Count() int {
	return len(cv.Violations)
}

func getConstraintShortName(constraintName string) string {
	return strings.Split(constraintName, ".")[1]
}

// RichViolation holds a violation with its category
type RichViolation struct {
	Category string // category of violation
	Resource string
	Message  string
	Metadata *_struct.Value `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

// NewRichViolation creates a new RichViolation
func NewRichViolation(categoryName string, violation *validator.Violation) (*RichViolation, error) {
	richViolation := &RichViolation{}
	richViolation.Category = categoryName
	richViolation.Resource = violation.Resource
	richViolation.Message = violation.Message
	richViolation.Metadata = violation.Metadata
	return richViolation, nil
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
		constraint := key
		cv = &constraintViolations{
			constraint: constraint,
		}
		config.constraints[key] = cv

		metadata := violation.GetMetadata().GetStructValue().GetFields()["constraint"]
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
func (inventory *InventoryConfig) Score(config *ScoringConfig, outputPath string, outputFormat string, outputMetadataFields []string) error {
	auditResult, err := getViolations(inventory, config)
	if err != nil {
		return err
	}
	Log.Debug("AuditResult from Config Validator", "# of Violations", len(auditResult.Violations))
	auditResult.Violations = uniqueViolations(auditResult.Violations)
	Log.Debug("AuditResult from Config Validator", "# of Unique Violations", len(auditResult.Violations))

	err = config.attachViolations(auditResult)
	if err != nil {
		return err
	}
	var dest io.Writer
	if len(auditResult.Violations) > 0 {
		if outputPath == "" {
			dest = os.Stdout
		} else {
			outputFile := "scorecard." + outputFormat
			dest, err = os.Create(filepath.Join(outputPath, outputFile))
			if err != nil {
				return err
			}
		}
		switch outputFormat {
		case "json":
			var richViolations []*RichViolation
			for _, category := range config.categories {
				for _, cv := range category.constraints {
					for _, v := range cv.Violations {
						richViolation, err := NewRichViolation(category.Name, v)
						if err != nil {
							return err
						}
						if len(outputMetadataFields) > 0 {
							newMetadata := make(map[string]interface{})
							oldMetadata := v.Metadata.GetStructValue().Fields["details"].GetStructValue()
							for _, field := range outputMetadataFields {
								newMetadata[field], _ = interfaceViaJSON(oldMetadata.Fields[field])
							}
							err := protoViaJSON(newMetadata, richViolation.Metadata)
							if err != nil {
								return err
							}
						}
						richViolations = append(richViolations, richViolation)
						Log.Debug("Violation metadata", "metadata", v.GetMetadata())
					}
				}
			}
			byteContent, err := json.MarshalIndent(richViolations, "", "  ")
			if err != nil {
				return err
			}
			io.WriteString(dest, string(byteContent)+"\n")
		case "csv":
			w := csv.NewWriter(dest)
			header := []string{"Category", "Constraint", "Resource", "Message"}
			for _, field := range outputMetadataFields {
				header = append(header, field)
			}
			w.Write(header)
			w.Flush()
			for _, category := range config.categories {
				for _, cv := range category.constraints {
					for _, v := range cv.Violations {
						record := []string{category.Name, getConstraintShortName(v.Constraint), v.Resource, v.Message}
						for _, field := range outputMetadataFields {
							metadata := v.Metadata.GetStructValue().Fields["details"].GetStructValue().Fields[field]
							value, _ := stringViaJSON(metadata)
							record = append(record, value)
						}
						w.Write(record)
						w.Flush()
						Log.Debug("Violation metadata", "metadata", v.GetMetadata())
					}
				}
			}
		case "txt":
			io.WriteString(dest, fmt.Sprintf("\n\n%v total issues found\n", len(auditResult.Violations)))
			for _, category := range config.categories {
				io.WriteString(dest, fmt.Sprintf("\n\n%v: %v issues found\n", category.Name, category.Count()))
				io.WriteString(dest, fmt.Sprintf("----------\n"))
				for _, cv := range category.constraints {
					io.WriteString(dest, fmt.Sprintf("%v: %v issues\n", getConstraintShortName(cv.constraint), cv.Count()))
					for _, v := range cv.Violations {
						io.WriteString(dest, fmt.Sprintf("- %v\n", v.Message))
						for _, field := range outputMetadataFields {
							metadata := v.Metadata.GetStructValue().Fields["details"].GetStructValue().Fields[field]
							value, _ := stringViaJSON(metadata)
							if value != "" {
								io.WriteString(dest, fmt.Sprintf("  %v: %v\n", field, value))
							}
						}
						io.WriteString(dest, "\n")
						Log.Debug("Violation metadata", "metadata", v.GetMetadata())
					}
				}
			}
		default:
			return fmt.Errorf("Unsupported output format %v", outputFormat)
		}
	} else {
		fmt.Println("No issues found found! You have a perfect score.")
	}

	return nil
}

func uniqueViolations(violations []*validator.Violation) []*validator.Violation {
	uniqueViolationMap := make(map[string]*validator.Violation)
	for _, v := range violations {
		b, _ := json.Marshal(v)
		hash := md5.Sum(b)
		uniqueViolationMap[string(hash[:])] = v
	}
	uniqueViolations := make([]*validator.Violation, 0, len(uniqueViolationMap))
	for _, v := range uniqueViolationMap {
		uniqueViolations = append(uniqueViolations, v)
	}
	return uniqueViolations
}
