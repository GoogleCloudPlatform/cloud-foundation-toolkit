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
	"github.com/mitchellh/go-testing-interface"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// Blueprint represents a config that can be setup, applied, verified and torndown.
type Blueprint interface {
	Setup()
	Apply()
	Verify(*assert.Assertions)
	Teardown()
}

// BlueprintTest implements a generic blueprint
type BlueprintTest struct {
	Setup    func()
	Apply    func()
	Verify   func(*assert.Assertions)
	Teardown func()
}

func (b *BlueprintTest) DefineSetup(setup func()) {
	b.Setup = setup

}
func (b *BlueprintTest) DefineApply(apply func()) {
	b.Apply = apply

}
func (b *BlueprintTest) DefineTeardown(teardown func()) {
	b.Teardown = teardown

}
func (b *BlueprintTest) DefineVerify(verify func(*assert.Assertions)) {
	b.Verify = verify
}

// TestBlueprint runs setup, apply, verify, teardown in order for a given blueprint
func TestBlueprint(t testing.TB, bp Blueprint, bptf func(*BlueprintTest)) {
	bpt := &BlueprintTest{}
	// apply any overrides to default bp methods
	if bptf != nil {
		bptf(bpt)
	}
	a := assert.New(t)
	// set default blueprint methods if not overriden by blueprint test
	if bpt.Setup == nil {
		bpt.Setup = func() { bp.Setup() }
	}
	if bpt.Apply == nil {
		bpt.Apply = func() { bp.Apply() }
	}
	if bpt.Teardown == nil {
		bpt.Teardown = func() { bp.Teardown() }
	}
	if bpt.Verify == nil {
		bpt.Verify = func(a *assert.Assertions) { bp.Verify(a) }
	}
	// run stages
	utils.RunStage("setup", func() { bpt.Setup() })
	defer utils.RunStage("teardown", func() { bpt.Teardown() })
	utils.RunStage("apply", func() { bpt.Apply() })
	utils.RunStage("verify", func() { bpt.Verify(a) })
}
