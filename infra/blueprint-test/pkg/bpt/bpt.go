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

// Package bpt defines a blueprint and implements default stages and execution order for a blueprint test

package bpt

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
)

// Blueprint represents a config that can be initialized, applied, verified and torndown.
type Blueprint interface {
	Init(*assert.Assertions)
	Apply(*assert.Assertions)
	Verify(*assert.Assertions)
	Teardown(*assert.Assertions)
}

// TestBlueprint runs init, apply, verify, teardown in order for a given blueprint
func TestBlueprint(t testing.TB, bp Blueprint) {
	a := assert.New(t)
	// run stages
	utils.RunStage("init", func() { bp.Init(a) })
	defer utils.RunStage("teardown", func() { bp.Teardown(a) })
	utils.RunStage("apply", func() { bp.Apply(a) })
	utils.RunStage("verify", func() { bp.Verify(a) })
}
