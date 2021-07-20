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
	"fmt"
	"os"
	"path"

	gotest "testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/modules/bpt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/modules/discovery"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/modules/utils"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
)

// TFBlueprintTest implements bpt.Blueprint and stores information associated with a Terraform blueprint test.
type TFBlueprintTest struct {
	name        string                 // descriptive name for the test
	fixtureName string                 // fixture name which will be used to compute TFDir
	tfDir       string                 // directory containing Terraform configs
	envVars     map[string]string      // variables to pass to Terraform as environment variables
	setupPath   string                 // optional directory containing applied TF configs to import outputs as variables for the test
	vars        map[string]interface{} // variables to pass to Terraform as flags
	logger      *logger.Logger         // custom logger
	t           testing.TB             // TestingT or TestingB
}

type option func(*TFBlueprintTest)

func WithName(name string) option {
	return func(f *TFBlueprintTest) {
		f.name = name
	}
}

func WithFixtureName(fixtureName string) option {
	return func(f *TFBlueprintTest) {
		f.fixtureName = fixtureName
	}
}

func WithTFDir(tfDir string) option {
	return func(f *TFBlueprintTest) {
		f.tfDir = tfDir
	}
}

func WithEnvVars(envVars map[string]string) option {
	return func(f *TFBlueprintTest) {
		f.envVars = envVars
	}
}

func WithSetupPath(setupPath string) option {
	return func(f *TFBlueprintTest) {
		f.setupPath = setupPath
	}
}

func WithVars(vars map[string]interface{}) option {
	return func(f *TFBlueprintTest) {
		f.vars = vars
	}
}

func WithLogger(logger *logger.Logger) option {
	return func(f *TFBlueprintTest) {
		f.logger = logger
	}
}

// Init sets defaults, validates and returns a TFBlueprintTest.
func Init(t testing.TB, opts ...option) *TFBlueprintTest {
	tft := &TFBlueprintTest{
		name: fmt.Sprintf("%s TF Blueprint", t.Name()),
		t:    t,
	}
	// apply options
	for _, opt := range opts {
		opt(tft)
	}
	// if no custom logger, set default based on test verbosity
	if tft.logger == nil {
		tft.logger = utils.GetLoggerFromT()
	}
	// one of fixture name or tfDir should be provided
	if tft.fixtureName != "" && tft.tfDir != "" {
		t.Fatalf("Only one of FixtureName or TFDir must be provided, found FixtureName=%s, TFDir=%s", tft.fixtureName, tft.tfDir)
	}
	// if both fixture name and tfDir are empty, auto discover tfDir based on cwd
	if tft.fixtureName == "" && tft.tfDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("unable to get wd :%v", err)
		}
		tfdir, err := discovery.ConfigDirFromCWD(cwd)
		if err != nil {
			t.Fatalf("unable to detect TFDir :%v", err)
		}
		tft.tfDir = tfdir
	}
	// if fixture name compute tfDir from given fixture name
	if tft.fixtureName != "" {
		tfModFixtureDir := getTFModuleFixtureDir(tft.fixtureName)
		if _, err := os.Stat(tfModFixtureDir); os.IsNotExist(err) {
			t.Fatalf("TFDir path derived from %s as %s does not exist", tft.fixtureName, tfModFixtureDir)
		}
		tft.tfDir = tfModFixtureDir
	}
	// setupPath is empty, try known setupPath
	if tft.setupPath == "" {
		tfModSetupDir := getTFModuleSetupDir()
		if _, err := os.Stat(tfModSetupDir); os.IsNotExist(err) {
			t.Logf("Setup dir %s not found, skipping loading setup outputs as fixture inputs", tfModSetupDir)
		}
		tft.setupPath = tfModSetupDir
	}
	// load TFEnvVars
	tfEnvVars := make(map[string]string)
	//load from setup outputs
	if tft.setupPath != "" {
		tft.logger.Logf(tft.t, "Loading env vars from setup %s", tft.setupPath)
		loadTFEnvVar(tfEnvVars, tft.getTFSetupOPMap())
	}
	// load additional env variables
	if tft.envVars != nil {
		loadTFEnvVar(tfEnvVars, tft.envVars)
	}
	tft.envVars = tfEnvVars

	tft.logger.Logf(tft.t, "Running tests TF configs in %s", tft.tfDir)
	return tft
}

// getTFModuleFixtureDir gets TF fixture config path from a fixture name.
func getTFModuleFixtureDir(fixture string) string {
	return fmt.Sprintf("../../test/fixtures/%s", fixture)
}

// getTFModuleSetupDir returns well known setup path.
func getTFModuleSetupDir() string {
	return "../../setup"
}

// getTFOptions generates terraform.Options used by Terratest.
func (b *TFBlueprintTest) getTFOptions() *terraform.Options {
	return terraform.WithDefaultRetryableErrors(b.t, &terraform.Options{
		TerraformDir: b.tfDir,
		EnvVars:      b.envVars,
		Vars:         b.vars,
		Logger:       b.logger,
	})
}

// getTFSetupOPMap computes a map of TF outputs from setup.
// Currently only returns string outputs.
func (b *TFBlueprintTest) getTFSetupOPMap() map[string]string {
	o := terraform.OutputAll(b.t, &terraform.Options{TerraformDir: b.setupPath, Logger: b.logger})
	n := make(map[string]string)
	for k, v := range o {
		s, ok := v.(string)
		if !ok {
			b.logger.Logf(b.t, "Unable to convert output %s value in setup dir to string", k)
			continue
		}
		n[k] = s
	}
	return n
}

// GetStringOutput returns TF output for a given key as string.
// It fails test if given key does not output a primitive.
func (b *TFBlueprintTest) GetStringOutput(name string) string {
	return terraform.Output(b.t, b.getTFOptions(), name)
}

// GetTFSetupOPListVal returns TF output from setup for a given key as list.
// It fails test if given key does not output a list type.
func (b *TFBlueprintTest) GetTFSetupOPListVal(key string) []string {
	if b.setupPath == "" {
		b.t.Fatal("Setup path not set")
	}
	return terraform.OutputList(b.t, &terraform.Options{TerraformDir: b.setupPath, Logger: b.logger}, key)
}

// loadTFEnvVar adds new env variables prefixed with TF_VAR_ to an existing map of variables.
func loadTFEnvVar(m map[string]string, new map[string]string) {
	for k, v := range new {
		m[fmt.Sprintf("TF_VAR_%s", k)] = v
	}
}

// AutoDiscoverAndTest discovers TF config from examples/fixtures and runs tests.
func AutoDiscoverAndTest(t *gotest.T) {
	configs := discovery.FindTestConfigs(t, "./")
	for _, dir := range configs {
		// dir must be of the form ../fixture/name or ../examples/name
		testName := fmt.Sprintf("test-%s-%s", path.Base(path.Dir(dir)), path.Base(dir))
		t.Run(testName, func(t *gotest.T) {
			nt := Init(t, WithTFDir(dir), WithSetupPath("../setup"))
			bpt.TestBlueprint(t, nt, nil)
		})
	}
}

// Teardown runs TF destroy on a blueprint.
func (b *TFBlueprintTest) Teardown() {
	terraform.Destroy(b.t, b.getTFOptions())
}

// Verify asserts no resource changes exist after apply.
func (b *TFBlueprintTest) Verify(assert *assert.Assertions) {
	e := terraform.PlanExitCode(b.t, b.getTFOptions())
	// exit code 0 is success with no diffs, 2 is success with non-empty diff
	// https://www.terraform.io/docs/cli/commands/plan.html#detailed-exitcode
	assert.Equal(e, 0, "plan after apply should have exit code 0")
}

// Setup runs TF init and validate on a blueprint.
func (b *TFBlueprintTest) Setup() {
	terraform.Init(b.t, b.getTFOptions())
	terraform.Validate(b.t, b.getTFOptions())
}

// Apply runs TF apply on a blueprint.
func (b *TFBlueprintTest) Apply() {
	terraform.Apply(b.t, b.getTFOptions())
}
