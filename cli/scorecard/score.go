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
	"fmt"
)

// ScoringConfig holds settings for generating a score
type ScoringConfig struct {
	PolicyPath string
}

// ScoreInventory creates a Scorecard for an inventory
func ScoreInventory(inventory *Inventory, config *ScoringConfig) error {
	auditResult, err := GetViolations(inventory, config)
	if err != nil {
		return err
	}

	fmt.Println(auditResult.Violations)

	if len(auditResult.Violations) > 0 {
		fmt.Print("Found Violations:\n\n")
		for _, v := range auditResult.Violations {
			fmt.Printf("Constraint %v on resource %v: %v\n\n",
				v.Constraint,
				v.Resource,
				v.Message,
			)
		}
	} else {
		fmt.Println("No violations found.")
	}

	return nil
}
