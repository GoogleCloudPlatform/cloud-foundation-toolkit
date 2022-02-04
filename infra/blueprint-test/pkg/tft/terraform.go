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

// Package tft provides a set of helpers to test Terraform modules/blueprints.
package tft

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"
	gotest "testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/discovery"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
)

const setupKeyOutputName = "sa_key"

// TFBlueprintTest implements bpt.Blueprint and stores information associated with a Terraform blueprint test.
type TFBlueprintTest struct {
	discovery.BlueprintTestConfig                          // additional blueprint test configs
	name                          string                   // descriptive name for the test
	tfDir                         string                   // directory containing Terraform configs
	tfEnvVars                     map[string]string        // variables to pass to Terraform as environment variables prefixed with TF_VAR_
	setupDir                      string                   // optional directory containing applied TF configs to import outputs as variables for the test
	vars                          map[string]interface{}   // variables to pass to Terraform as flags
	logger                        *logger.Logger           // custom logger
	t                             testing.TB               // TestingT or TestingB
	init                          func(*assert.Assertions) // init function
	apply                         func(*assert.Assertions) // apply function
	verify                        func(*assert.Assertions) // verify function
	teardown                      func(*assert.Assertions) // teardown function
}

type tftOption func(*TFBlueprintTest)

func WithName(name string) tftOption {
	return func(f *TFBlueprintTest) {
		f.name = name
	}
}

func WithFixtureName(fixtureName string) tftOption {
	return func(f *TFBlueprintTest) {
		// when a test is invoked for an explicit blueprint fixture
		// expect fixture path to be ../../fixtures/fixtureName
		tfModFixtureDir := path.Join("..", "..", discovery.FixtureDir, fixtureName)
		f.tfDir = tfModFixtureDir
	}
}

func WithTFDir(tfDir string) tftOption {
	return func(f *TFBlueprintTest) {
		f.tfDir = tfDir
	}
}

func WithEnvVars(envVars map[string]string) tftOption {
	return func(f *TFBlueprintTest) {
		tfEnvVars := make(map[string]string)
		loadTFEnvVar(tfEnvVars, envVars)
		f.tfEnvVars = tfEnvVars
	}
}

func WithSetupPath(setupPath string) tftOption {
	return func(f *TFBlueprintTest) {
		f.setupDir = setupPath
	}
}

func WithVars(vars map[string]interface{}) tftOption {
	return func(f *TFBlueprintTest) {
		f.vars = vars
	}
}

func WithLogger(logger *logger.Logger) tftOption {
	return func(f *TFBlueprintTest) {
		f.logger = logger
	}
}

// NewTFBlueprintTest sets defaults, validates and returns a TFBlueprintTest.
func NewTFBlueprintTest(t testing.TB, opts ...tftOption) *TFBlueprintTest {
	tft := &TFBlueprintTest{
		name:      fmt.Sprintf("%s TF Blueprint", t.Name()),
		tfEnvVars: make(map[string]string),
		t:         t,
	}
	// default TF blueprint methods
	tft.init = tft.DefaultInit
	tft.apply = tft.DefaultApply
	tft.verify = tft.DefaultVerify
	tft.teardown = tft.DefaultTeardown
	// apply options
	for _, opt := range opts {
		opt(tft)
	}
	// if no custom logger, set default based on test verbosity
	if tft.logger == nil {
		tft.logger = utils.GetLoggerFromT()
	}
	// if explicit tfDir is provided, validate it else try auto discovery
	if tft.tfDir != "" {
		_, err := os.Stat(tft.tfDir)
		if os.IsNotExist(err) {
			t.Fatalf("TFDir path %s does not exist", tft.tfDir)
		}
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("unable to get wd :%v", err)
		}
		tfdir, err := discovery.GetConfigDirFromTestDir(cwd)
		if err != nil {
			t.Fatalf("unable to detect TFDir :%v", err)
		}
		tft.tfDir = tfdir
	}
	// discover test config
	var err error
	tft.BlueprintTestConfig, err = discovery.GetTestConfig(path.Join(tft.tfDir, discovery.DefaultTestConfigFilename))
	if err != nil {
		t.Fatal(err)
	}
	// setupDir is empty, try known setupDir paths
	if tft.setupDir == "" {
		setupDir, err := discovery.GetKnownDirInParents(discovery.SetupDir, 2)
		if err != nil {
			t.Logf("Setup dir not found, skipping loading setup outputs as fixture inputs: %v", err)
		} else {
			tft.setupDir = setupDir
		}
	}
	//load TFEnvVars from setup outputs
	if tft.setupDir != "" {
		tft.logger.Logf(tft.t, "Loading env vars from setup %s", tft.setupDir)
		loadTFEnvVar(tft.tfEnvVars, tft.getTFOutputsAsInputs(terraform.OutputAll(tft.t, &terraform.Options{TerraformDir: tft.setupDir, Logger: tft.logger})))
		// setup credentials if explicit sa_key output from setup
		credsEnc, exists := tft.tfEnvVars[fmt.Sprintf("TF_VAR_%s", setupKeyOutputName)]
		if !exists {
			tft.logger.Logf(tft.t, "Unable to find %s output from setup, skipping credential activation", setupKeyOutputName)
		} else {
			credDec, err := b64.StdEncoding.DecodeString(credsEnc)
			if err != nil {
				t.Fatalf("Unable to decode %s output from setup: %v", setupKeyOutputName, err)
			}
			gcloud.ActivateCredsAndEnvVars(tft.t, string(credDec))
		}
	}

	tft.logger.Logf(tft.t, "Running tests TF configs in %s", tft.tfDir)
	return tft
}

// GetTFOptions generates terraform.Options used by Terratest.
func (b *TFBlueprintTest) GetTFOptions() *terraform.Options {
	return terraform.WithDefaultRetryableErrors(b.t, &terraform.Options{
		TerraformDir: b.tfDir,
		EnvVars:      b.tfEnvVars,
		Vars:         b.vars,
		Logger:       b.logger,
	})
}

// getTFOutputsAsInputs computes a map of TF inputs from outputs map.
func (b *TFBlueprintTest) getTFOutputsAsInputs(o map[string]interface{}) map[string]string {
	n := make(map[string]string)
	// TF requires complex values to be an HCL expression passed literally.
	// However, Terratest only exposes a way to format strings as HCL expressions to be used as var flags.
	// Var flags requires the root module to declare a variable of that name.
	// Hence, we extract the HCL formated string from the var arg slice of form [-var, key1=value1, -var, key2={"complex"="data"}...]
	for _, v := range terraform.FormatTerraformVarsAsArgs(o) {
		if v == "-var" {
			continue
		}
		parsedKey, parsedVal, err := getKVFromOutputString(v)
		if err != nil {
			b.t.Logf("Unable to parse output from setup: %v", err)
			continue
		}
		n[parsedKey] = parsedVal
	}
	return n
}

// getKVFromOutputString parses string kv pairs of form k=v
func getKVFromOutputString(v string) (string, string, error) {
	// v of format key1=value1
	kv := strings.SplitN(v, "=", 2)
	if len(kv) < 2 {
		return "", "", fmt.Errorf("error parsing %s", v)
	}
	return kv[0], kv[1], nil
}

// GetStringOutput returns TF output for a given key as string.
// It fails test if given key does not output a primitive.
func (b *TFBlueprintTest) GetStringOutput(name string) string {
	return terraform.Output(b.t, b.GetTFOptions(), name)
}

// GetTFSetupOutputListVal returns TF output from setup for a given key as list.
// It fails test if given key does not output a list type.
func (b *TFBlueprintTest) GetTFSetupOutputListVal(key string) []string {
	if b.setupDir == "" {
		b.t.Fatal("Setup path not set")
	}
	return terraform.OutputList(b.t, &terraform.Options{TerraformDir: b.setupDir, Logger: b.logger}, key)
}

// loadTFEnvVar adds new env variables prefixed with TF_VAR_ to an existing map of variables.
func loadTFEnvVar(m map[string]string, new map[string]string) {
	for k, v := range new {
		m[fmt.Sprintf("TF_VAR_%s", k)] = v
	}
}

// ShouldSkip checks if a test should be skipped
func (b *TFBlueprintTest) ShouldSkip() bool {
	return b.Spec.Skip
}

// AutoDiscoverAndTest discovers TF config from examples/fixtures and runs tests.
func AutoDiscoverAndTest(t *gotest.T) {
	configs := discovery.FindTestConfigs(t, "./")
	for testName, dir := range configs {
		t.Run(testName, func(t *gotest.T) {
			nt := NewTFBlueprintTest(t, WithTFDir(dir))
			nt.Test()
		})
	}
}

// DefineInit defines a custom init function for the blueprint.
func (b *TFBlueprintTest) DefineInit(init func(*assert.Assertions)) {
	b.init = init
}

// DefineApply defines a custom apply function for the blueprint.
func (b *TFBlueprintTest) DefineApply(apply func(*assert.Assertions)) {
	b.apply = apply
}

// DefineVerify defines a custom verify function for the blueprint.
func (b *TFBlueprintTest) DefineVerify(verify func(*assert.Assertions)) {
	b.verify = verify
}

// DefineTeardown defines a custom teardown function for the blueprint.
func (b *TFBlueprintTest) DefineTeardown(teardown func(*assert.Assertions)) {
	b.teardown = teardown
}

// DefaultTeardown runs TF destroy on a blueprint.
func (b *TFBlueprintTest) DefaultTeardown(assert *assert.Assertions) {
	terraform.Destroy(b.t, b.GetTFOptions())
}

// DefaultVerify asserts no resource changes exist after apply.
func (b *TFBlueprintTest) DefaultVerify(assert *assert.Assertions) {
	e := terraform.PlanExitCode(b.t, b.GetTFOptions())
	// exit code 0 is success with no diffs, 2 is success with non-empty diff
	// https://www.terraform.io/docs/cli/commands/plan.html#detailed-exitcode
	assert.NotEqual(1, e, "plan after apply should not fail")
	assert.NotEqual(2, e, "plan after apply should have non-empty diff")
}

// DefaultInit runs TF init and validate on a blueprint.
func (b *TFBlueprintTest) DefaultInit(assert *assert.Assertions) {
	terraform.Init(b.t, b.GetTFOptions())
	// if vars are set for common options, this seems to trigger -var flag when calling validate
	// using custom tfOptions as a workaround
	terraform.Validate(b.t, terraform.WithDefaultRetryableErrors(b.t, &terraform.Options{
		TerraformDir: b.tfDir,
		Logger:       b.logger,
	}))
}

// DefaultApply runs TF apply on a blueprint.
func (b *TFBlueprintTest) DefaultApply(assert *assert.Assertions) {
	terraform.Apply(b.t, b.GetTFOptions())
}

// Init runs the default or custom init function for the blueprint.
func (b *TFBlueprintTest) Init(assert *assert.Assertions) {
	b.init(assert)
}

// Apply runs the default or custom apply function for the blueprint.
func (b *TFBlueprintTest) Apply(assert *assert.Assertions) {
	b.apply(assert)
}

// Verify runs the default or custom verify function for the blueprint.
func (b *TFBlueprintTest) Verify(assert *assert.Assertions) {
	b.verify(assert)
}

// Teardown runs the default or custom teardown function for the blueprint.
func (b *TFBlueprintTest) Teardown(assert *assert.Assertions) {
	b.teardown(assert)
}

// Test runs init, apply, verify, teardown in order for the blueprint.
func (b *TFBlueprintTest) Test() {
	if b.ShouldSkip() {
		b.logger.Logf(b.t, "Skipping test due to config %s", b.Path)
		b.t.SkipNow()
		return
	}
	a := assert.New(b.t)
	// run stages
	utils.RunStage("init", func() { b.Init(a) })
	defer utils.RunStage("teardown", func() { b.Teardown(a) })
	utils.RunStage("apply", func() { b.Apply(a) })
	utils.RunStage("verify", func() { b.Verify(a) })
}
