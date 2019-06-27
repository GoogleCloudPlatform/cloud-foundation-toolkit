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
	"fmt"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"
	"github.com/pkg/errors"
)

// ScoringConfig holds settings for generating a score
type ScoringConfig struct {
	PolicyPath  string
	categories  map[string]*constraintCategory
	constraints map[string]*constraintViolations
	validator   *gcv.Validator
	ctx         context.Context
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
	return cv.constraint.GetMetadata().GetName()
}

var availableCategories = map[string]string{
	"operational-efficiency": "Operational Efficiency",
	otherCategoryKey:         "Other",
}

func getConstraintForViolation(config *ScoringConfig, violation *validator.Violation) (*constraintViolations, error) {
	key := violation.GetConstraint()
	cv, found := config.constraints[key]
	if !found {
		response, err := config.validator.GetConstraint(config.ctx, &validator.GetConstraintRequest{
			Name: key,
		})
		if err != nil {
			return nil, errors.Wrap(err, "Finding matching constraint")
		}

		constraint := response.GetConstraint()
		cv = &constraintViolations{
			constraint: response.GetConstraint(),
		}
		config.constraints[key] = cv

		annotations := constraint.GetMetadata().GetAnnotations()
		categoryKey, found := annotations["scorecard.cft.dev/category"]
		if !found {
			categoryKey = otherCategoryKey
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
func attachViolations(audit *validator.AuditResponse, config *ScoringConfig) error {
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
		cv, err := getConstraintForViolation(config, v)
		if err != nil {
			return errors.Wrap(err, "Categorizing violation")
		}

		cv.Violations = append(cv.Violations, v)
	}

	return nil
}

// ScoreInventory creates a Scorecard for an inventory
func ScoreInventory(inventory *Inventory, config *ScoringConfig) error {
	config.ctx = context.Background()

	err := attachValidator(config)
	if err != nil {
		return errors.Wrap(err, "initializing gcv validator")
	}

	auditResult, err := getViolations(inventory, config)
	if err != nil {
		return err
	}

	err = attachViolations(auditResult, config)

	if len(auditResult.Violations) > 0 {
		fmt.Printf("\n\n%v total issues found\n", len(auditResult.Violations))
		for _, category := range config.categories {
			fmt.Printf("\n\n%v: %v issues found\n", category.Name, category.Count())
			fmt.Printf("----------\n")
			for _, cv := range category.constraints {
				fmt.Printf("%v: %v issues\n", cv.GetName(), cv.Count())
				for _, v := range cv.Violations {
					fmt.Printf("- %v\n\n",
						v.Message,
					)
					Log.Debug("Violation metadata", "metadata", v.GetMetadata())
				}
			}
		}
	} else {
		fmt.Println("No issues found found! You have a perfect score.")
	}

	return nil
}
