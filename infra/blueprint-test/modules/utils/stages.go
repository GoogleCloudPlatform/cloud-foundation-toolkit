/**
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"log"
	"os"
)

const RUN_STAGE_ENV_VAR = "RUN_STAGE"

// RunStage runs stage if stageName matches RUN_STAGE env var or RUN_STAGE is unset
// Similar to terratest RunStage but instead of skipping using env var, we match using envvar
func RunStage(stageName string, stage func()) {
	if shouldRunStage(stageName) {
		log.Printf("Running stage %s", stageName)
		stage()
	} else {
		log.Printf("Skipping stage %s", stageName)
	}

}

// shouldRunStage returns true if no expilcit stage set via RUN_STAGE env var or if stageName matches value in RUN_STAGE
func shouldRunStage(stageName string) bool {
	// no env var set, run all
	if os.Getenv(RUN_STAGE_ENV_VAR) == "" {
		log.Printf("No RUN_STAGE env var set, running stage %s", stageName)
		return true
	}
	envStage := os.Getenv(RUN_STAGE_ENV_VAR)
	log.Printf("RUN_STAGE env var set to %s", envStage)
	// if envvar matches current stage, run it
	return envStage == stageName

}
