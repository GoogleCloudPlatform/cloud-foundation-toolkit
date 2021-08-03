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

// BlueprintTest implements a generic blueprint
type BlueprintTest struct {
	init     func(*assert.Assertions) // custom init function
	apply    func(*assert.Assertions) // custom apply function
	verify   func(*assert.Assertions) // custom verify function
	teardown func(*assert.Assertions) // custom teardown function
	bp       Blueprint                // blueprint to be tested
}

type bptOption func(*BlueprintTest)

func DefineInit(init func(*assert.Assertions)) bptOption {
	return func(bt *BlueprintTest) {
		bt.init = init
	}
}

func DefineApply(apply func(*assert.Assertions)) bptOption {
	return func(bt *BlueprintTest) {
		bt.apply = apply
	}
}

func DefineTeardown(teardown func(*assert.Assertions)) bptOption {
	return func(bt *BlueprintTest) {
		bt.teardown = teardown
	}
}

func DefineVerify(verify func(*assert.Assertions)) bptOption {
	return func(bt *BlueprintTest) {
		bt.verify = verify
	}
}

func (bt *BlueprintTest) Init(a *assert.Assertions) {
	bt.init(a)
}

func (bt *BlueprintTest) Apply(a *assert.Assertions) {
	bt.apply(a)
}

func (bt *BlueprintTest) Teardown(a *assert.Assertions) {
	bt.teardown(a)
}

func (bt *BlueprintTest) Verify(a *assert.Assertions) {
	bt.verify(a)
}

// TestBlueprint runs init, apply, verify, teardown in order for a given blueprint
func TestBlueprint(t testing.TB, bp Blueprint, opts ...bptOption) {
	bpt := &BlueprintTest{
		bp:       bp,
		init:     bp.Init,
		apply:    bp.Apply,
		verify:   bp.Verify,
		teardown: bp.Teardown,
	}
	// apply any overrides to default bp methods
	for _, opt := range opts {
		opt(bpt)
	}
	a := assert.New(t)
	// run stages
	utils.RunStage("init", func() { bpt.Init(a) })
	defer utils.RunStage("teardown", func() { bpt.Teardown(a) })
	utils.RunStage("apply", func() { bpt.Apply(a) })
	utils.RunStage("verify", func() { bpt.Verify(a) })
}
