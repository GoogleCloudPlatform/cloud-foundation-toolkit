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

// Package tf provides a set of helpers to test Terraform modules/blueprints
package tf

import (
	"fmt"
	"os"

	gotest "testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/go-testing-interface"
)

// TFBlueprintTest stores information associated with a Terraform blueprint test
type TFBlueprintTest struct {
	Name        string                 // descriptive name for the test
	FixtureName string                 // fixture name which will be used to compute TFDir
	TFDir       string                 // directory containing Terraform configs
	EnvVars     map[string]string      // variables to pass to Terraform as environment variables
	SetupPath   string                 // optional directory containing applied TF configs to import outputs as variables for the test
	Vars        map[string]interface{} //variables to pass to Terraform as flags
	Logger      *logger.Logger         // custom logger
}

// Init sets defaults and validates values for a TFBlueprintTest
func Init(t testing.TB, opts *TFBlueprintTest) *TFBlueprintTest {
	if opts.Name == "" {
		opts.Name = fmt.Sprintf("%s TF Blueprint", t.Name())
	}

	if opts.FixtureName == "" && opts.TFDir == "" {
		t.Fatal("One of FixtureName or TFDir must be provided")
	}

	if opts.FixtureName != "" && opts.TFDir != "" {
		t.Fatalf("Only one of FixtureName or TFDir must be provided, found FixtureName=%s, TFDir=%s", opts.FixtureName, opts.TFDir)
	}

	// compute TFDir path from given fixture name
	if opts.FixtureName != "" {
		tfModFixtureDir := getTFModuleFixtureDir(opts.FixtureName)
		if _, err := os.Stat(tfModFixtureDir); os.IsNotExist(err) {
			t.Fatalf("TFDir path derived from %s as %s does not exist", opts.FixtureName, tfModFixtureDir)
		}
		opts.TFDir = tfModFixtureDir
	}

	if opts.SetupPath == "" {
		tfModSetupDir := getTFModuleSetupDir()
		if _, err := os.Stat(tfModSetupDir); os.IsNotExist(err) {
			t.Logf("Setup dir %s not found, skipping loading setup outputs as fixture inputs", tfModSetupDir)
		}
		opts.SetupPath = tfModSetupDir
	}
	if opts.Logger == nil {
		if gotest.Verbose() {
			opts.Logger = logger.Default
		} else {
			opts.Logger = logger.Discard
		}
	}

	return opts
}

// get TF dir path from a fixture
func getTFModuleFixtureDir(fixture string) string {
	return fmt.Sprintf("../../test/fixtures/%s", fixture)
}

// get potential setup path
func getTFModuleSetupDir() string {
	return "../../test/setup"
}

// generate terraform.Options for a TFBlueprint
// this will be used by Terratest
func (b *TFBlueprintTest) getTFOptions(t testing.TB) *terraform.Options {
	tfEnvVars := make(map[string]string)
	// load TF outputs from setup as input variables
	if b.SetupPath != "" {
		t.Logf("Loading env vars from setup %s", b.SetupPath)
		loadTFEnvVar(tfEnvVars, b.getTFSetupOPMap(t))
	}
	// load additional env variables
	if b.EnvVars != nil {
		loadTFEnvVar(tfEnvVars, b.EnvVars)
	}
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: b.TFDir,
		EnvVars:      tfEnvVars,
		Vars:         b.Vars,
		Logger:       b.Logger,
	})
}

// getTFSetupOPMap computes a map of TF outputs from setup
// currently only returns string outputs
func (b *TFBlueprintTest) getTFSetupOPMap(t testing.TB) map[string]string {
	o := terraform.OutputAll(t, &terraform.Options{TerraformDir: b.SetupPath, Logger: b.Logger})
	n := make(map[string]string)
	for k, v := range o {
		s, ok := v.(string)
		if !ok {
			t.Logf("Unable to convert output %s value in setup dir to string", k)
		}
		n[k] = s
	}
	return n
}

// GetStringOutput returns TF output for a given key as string
// fails test if given key does not output a primitive
func (b *TFBlueprintTest) GetStringOutput(t testing.TB, name string) string {
	return terraform.Output(t, b.getTFOptions(t), name)
}

// GetTFSetupOPListVal returns TF output from setup for a given key as list
// fails test if given key does not output a list type
func (b *TFBlueprintTest) GetTFSetupOPListVal(t testing.TB, key string) []string {
	if b.SetupPath == "" {
		t.Fatal("Setup path not set")
	}
	return terraform.OutputList(t, &terraform.Options{TerraformDir: b.SetupPath, Logger: b.Logger}, key)
}

// loadTFEnvVar adds new env variables prefixed with TF_VAR_ to an existing map of variables
func loadTFEnvVar(m map[string]string, new map[string]string) {
	for k, v := range new {
		m[fmt.Sprintf("TF_VAR_%s", k)] = v
	}
}

// Teardown runs TF destroy on a blueprint
func (b *TFBlueprintTest) Teardown(t testing.TB) {
	t.Logf("Destroying %s", b.Name)
	terraform.Destroy(t, b.getTFOptions(t))
}

// Teardown runs TF init and apply on a blueprint
func (b *TFBlueprintTest) Setup(t testing.TB) {
	t.Logf("Initializing and applying %s", b.Name)
	terraform.InitAndApply(t, b.getTFOptions(t))
}

// TFInit runs TF init on a blueprint
func (b *TFBlueprintTest) TFInit(t testing.TB) {
	t.Logf("Initializing %s", b.Name)
	terraform.Init(t, b.getTFOptions(t))
}

// TFInit runs TF apply on a blueprint
func (b *TFBlueprintTest) TFApply(t testing.TB) {
	t.Logf("Applying %s", b.Name)
	terraform.Apply(t, b.getTFOptions(t))
}
