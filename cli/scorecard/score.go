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
	PolicyPath string
	Categories map[string]*CategoryViolations
	Validator  *gcv.Validator
	ctx        context.Context
}

const otherCategoryKey = "other"

// CategoryViolations holds actual scores for a particular category
type CategoryViolations struct {
	Name       string
	Violations []*validator.Violation `protobuf:"bytes,1,rep,name=violations,proto3" json:"violations,omitempty"`
}

var availableCategories = map[string]string{
	"operational-efficiency": "Operational Efficiency",
	otherCategoryKey:         "Other",
}

// attachViolations puts violations into their appropriate categories
func attachViolations(audit *validator.AuditResponse, config *ScoringConfig) error {
	// Build map of categories
	config.Categories = make(map[string]*CategoryViolations)
	for k, name := range availableCategories {
		config.Categories[k] = &CategoryViolations{
			Name: name,
		}
	}

	// Categorize violations
	for _, v := range audit.Violations {
		c, err := config.Validator.GetConstraint(config.ctx, &validator.GetConstraintRequest{
			Name: v.GetConstraint(),
		})
		if err != nil {
			return errors.Wrap(err, "Finding matching constraint")
		}

		annotations := c.GetConstraint().GetMetadata().GetAnnotations()
		categoryKey, found := annotations["scorecard.cft.dev/category"]
		if !found {
			categoryKey = otherCategoryKey
		}

		category, found := config.Categories[categoryKey]
		if !found {
			return fmt.Errorf("Unknown constraint category %v for constraint %v", categoryKey, v.GetConstraint())
		}

		category.Violations = append(category.Violations, v)
	}

	return nil
}

// ScoreInventory creates a Scorecard for an inventory
func ScoreInventory(inventory *Inventory, config *ScoringConfig) error {
	config.ctx = context.Background()

	err := AttachValidator(config)
	if err != nil {
		return errors.Wrap(err, "initializing gcv validator")
	}

	auditResult, err := GetViolations(inventory, config)
	if err != nil {
		return err
	}

	err = attachViolations(auditResult, config)

	if len(auditResult.Violations) > 0 {
		fmt.Print("\n\nFound %v issues:\n\n")
		for _, v := range auditResult.Violations {
			fmt.Printf("Constraint %v on resource %v: %v\n\n",
				v.Constraint,
				v.Resource,
				v.Message,
			)
			Log.Debug("Violation metadata", "metadata", v.GetMetadata())
		}
	} else {
		fmt.Println("No issues found found! You have a perfect score.")
	}

	return nil
}
