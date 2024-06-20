/**
 * Copyright 2021-2024 Google LLC
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
	"path/filepath"
	"strings"
	gotest "testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/discovery"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/alexflint/go-filemutex"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const (
	setupKeyOutputName    = "sa_key"
	tftCacheMutexFilename = "bpt-tft-cache.lock"
	planFilename          = "plan.tfplan"
)

var (
	CommonRetryableErrors = map[string]string{
		// Project deletion is eventually consistent. Even if google_project resources inside the folder are deleted there maybe a deletion error.
		".*FOLDER_TO_DELETE_NON_EMPTY_VIOLATION.*": "Failed to delete non empty folder.",

		// API activation is eventually consistent. Even if the google_project_service resource is reconciled there maybe an activation error.
		".*SERVICE_DISABLED.*": "Required API not enabled.",
	}
)

// TFBlueprintTest implements bpt.Blueprint and stores information associated with a Terraform blueprint test.
type TFBlueprintTest struct {
	discovery.BlueprintTestConfig                                                 // additional blueprint test configs
	name                          string                                          // descriptive name for the test
	saKey                         string                                          // optional setup sa key
	tfDir                         string                                          // directory containing Terraform configs
	tfEnvVars                     map[string]string                               // variables to pass to Terraform as environment variables prefixed with TF_VAR_
	backendConfig                 map[string]interface{}                          // backend configuration for terraform init
	retryableTerraformErrors      map[string]string                               // If Terraform apply fails with one of these (transient) errors, retry. The keys are a regexp to match against the error and the message is what to display to a user if that error is matched.
	maxRetries                    int                                             // Maximum number of times to retry errors matching RetryableTerraformErrors
	timeBetweenRetries            time.Duration                                   // The amount of time to wait between retries
	migrateState                  bool                                            // suppress user confirmation in a migration in terraform init
	setupDir                      string                                          // optional directory containing applied TF configs to import outputs as variables for the test
	policyLibraryPath             string                                          // optional absolute path to directory containing policy library constraints
	terraformVetProject           string                                          // optional a valid existing project that will be used when a plan has resources in a project that still does not exist.
	vars                          map[string]interface{}                          // variables to pass to Terraform as flags
	logger                        *logger.Logger                                  // custom logger
	sensitiveLogger               *logger.Logger                                  // custom logger for sensitive logging
	t                             testing.TB                                      // TestingT or TestingB
	init                          func(*assert.Assertions)                        // init function
	plan                          func(*terraform.PlanStruct, *assert.Assertions) // plan function
	apply                         func(*assert.Assertions)                        // apply function
	verify                        func(*assert.Assertions)                        // verify function
	teardown                      func(*assert.Assertions)                        // teardown function
	setupOutputOverrides          map[string]interface{}                          // override outputs from the Setup phase
	tftCacheMutex                 *filemutex.FileMutex                            // Mutex to protect Terraform plugin cache
	parallelism                   int                                             // Set the parallelism setting for Terraform
}

type tftOption func(*TFBlueprintTest)

func WithName(name string) tftOption {
	return func(f *TFBlueprintTest) {
		f.name = name
	}
}

func WithSetupSaKey(saKey string) tftOption {
	return func(f *TFBlueprintTest) {
		f.saKey = saKey
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

func WithBackendConfig(backendConfig map[string]interface{}) tftOption {
	return func(f *TFBlueprintTest) {
		f.backendConfig = backendConfig
		f.migrateState = true
	}
}

func WithRetryableTerraformErrors(retryableTerraformErrors map[string]string, maxRetries int, timeBetweenRetries time.Duration) tftOption {
	return func(f *TFBlueprintTest) {
		f.retryableTerraformErrors = retryableTerraformErrors
		f.maxRetries = maxRetries
		f.timeBetweenRetries = timeBetweenRetries
	}
}

func WithSetupPath(setupPath string) tftOption {
	return func(f *TFBlueprintTest) {
		f.setupDir = setupPath
	}
}

func WithPolicyLibraryPath(policyLibraryPath, terraformVetProject string) tftOption {
	return func(f *TFBlueprintTest) {
		f.policyLibraryPath = policyLibraryPath
		f.terraformVetProject = terraformVetProject
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

func WithSensitiveLogger(logger *logger.Logger) tftOption {
	return func(f *TFBlueprintTest) {
		f.sensitiveLogger = logger
	}
}

// WithSetupOutputs overrides output values from the setup stage
func WithSetupOutputs(vars map[string]interface{}) tftOption {
	return func(f *TFBlueprintTest) {
		f.setupOutputOverrides = vars
	}
}

func WithParallelism(p int) tftOption {
	return func(f *TFBlueprintTest) {
		f.parallelism = p
	}
}

// NewTFBlueprintTest sets defaults, validates and returns a TFBlueprintTest.
func NewTFBlueprintTest(t testing.TB, opts ...tftOption) *TFBlueprintTest {
	var err error
	tft := &TFBlueprintTest{
		name:      fmt.Sprintf("%s TF Blueprint", t.Name()),
		tfEnvVars: make(map[string]string),
		t:         t,
	}
	// initiate tft cache file mutex
	tft.tftCacheMutex, err = filemutex.New(filepath.Join(os.TempDir(), tftCacheMutexFilename))
	if err != nil {
		t.Fatalf("tft lock file <%s> could not created: %v", filepath.Join(os.TempDir(), tftCacheMutexFilename), err)
	}
	// default TF blueprint methods
	tft.init = tft.DefaultInit
	// No default plan function, plan is skipped if no custom func provided.
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
	// If no custom sensitive logger, use discard logger.
	if tft.sensitiveLogger == nil {
		tft.sensitiveLogger = logger.Discard
	}
	// if explicit tfDir is provided, validate it else try auto discovery
	if tft.tfDir != "" {
		_, err := os.Stat(tft.tfDir)
		if os.IsNotExist(err) {
			t.Fatalf("TFDir path %s does not exist", tft.tfDir)
		}
	} else {
		tfdir, err := discovery.GetConfigDirFromTestDir(utils.GetWD(t))
		if err != nil {
			t.Fatalf("unable to detect TFDir :%v", err)
		}
		tft.tfDir = tfdir
	}

	// discover test config
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
	// load setup sa Key
	if tft.saKey != "" {
		gcloud.ActivateCredsAndEnvVars(tft.t, tft.saKey)
	}
	// load TFEnvVars from setup outputs
	if tft.setupDir != "" {
		tft.logger.Logf(tft.t, "Loading env vars from setup %s", tft.setupDir)
		outputs := tft.getOutputs(tft.sensitiveOutputs(tft.setupDir))
		loadTFEnvVar(tft.tfEnvVars, tft.getTFOutputsAsInputs(outputs))
		if credsEnc, exists := tft.tfEnvVars[fmt.Sprintf("TF_VAR_%s", setupKeyOutputName)]; tft.saKey == "" && exists {
			if credDec, err := b64.StdEncoding.DecodeString(credsEnc); err == nil {
				gcloud.ActivateCredsAndEnvVars(tft.t, string(credDec))
			} else {
				tft.t.Fatalf("Unable to decode setup sa key: %v", err)
			}
		} else {
			tft.logger.Logf(tft.t, "Skipping credential activation %s output from setup", setupKeyOutputName)
		}
	}
	// Load env vars to supplement/override setup
	tft.logger.Logf(tft.t, "Loading setup from environment")
	if tft.setupOutputOverrides == nil {
		tft.setupOutputOverrides = make(map[string]interface{})
	}
	for k, v := range extractFromEnv("CFT_SETUP_") {
		tft.setupOutputOverrides[k] = v
	}

	tftVersion := gjson.Get(terraform.RunTerraformCommand(tft.t, tft.GetTFOptions(), "version", "-json"), "terraform_version")
	tft.logger.Logf(tft.t, "Running tests TF configs in %s with version %s", tft.tfDir, tftVersion)
	return tft
}

// sensitiveOutputs returns a map of sensitive output keys for module in dir.
func (b *TFBlueprintTest) sensitiveOutputs(dir string) map[string]bool {
	mod, err := tfconfig.LoadModule(dir)
	if err != nil {
		b.t.Fatalf("error loading module in %s: %v", dir, err)
	}
	sensitiveOP := map[string]bool{}
	for _, op := range mod.Outputs {
		if op.Sensitive {
			sensitiveOP[op.Name] = true
		}
	}
	return sensitiveOP
}

// getOutputs returns all output values.
func (b *TFBlueprintTest) getOutputs(sensitive map[string]bool) map[string]interface{} {
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()
	outputs := terraform.OutputAll(b.t, &terraform.Options{TerraformDir: b.setupDir, Logger: b.sensitiveLogger, NoColor: true})
	for k, v := range outputs {
		_, s := sensitive[k]
		if s {
			b.sensitiveLogger.Logf(b.t, "output key %q: %v", k, v)
		} else {
			b.logger.Logf(b.t, "output key %q: %v", k, v)
		}
	}
	return outputs
}

// GetTFOptions generates terraform.Options used by Terratest.
func (b *TFBlueprintTest) GetTFOptions() *terraform.Options {
	newOptions := terraform.WithDefaultRetryableErrors(b.t, &terraform.Options{
		TerraformDir:             b.tfDir,
		EnvVars:                  b.tfEnvVars,
		Vars:                     b.vars,
		Logger:                   b.logger,
		BackendConfig:            b.backendConfig,
		MigrateState:             b.migrateState,
		RetryableTerraformErrors: b.retryableTerraformErrors,
		NoColor:                  true,
		Parallelism:              b.parallelism,
	})
	if b.maxRetries > 0 {
		newOptions.MaxRetries = b.maxRetries
	}
	if b.timeBetweenRetries > 0 {
		newOptions.TimeBetweenRetries = b.timeBetweenRetries
	}
	return newOptions
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
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()
	return terraform.Output(b.t, b.GetTFOptions(), name)
}

// GetStringOutputList returns TF output for a given key as list.
// It fails test if given key does not output a primitive.
//
// Deprecated: Use GetJsonOutput instead.
func (b *TFBlueprintTest) GetStringOutputList(name string) []string {
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()
	return terraform.OutputList(b.t, b.GetTFOptions(), name)
}

// GetJsonOutput returns TF output for key as gjson.Result.
// An empty string for key can be used to return all values.
// It fails test on invalid JSON.
func (b *TFBlueprintTest) GetJsonOutput(key string) gjson.Result {
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()

	jsonString := terraform.OutputJson(b.t, b.GetTFOptions(), key)
	if !gjson.Valid(jsonString) {
		b.t.Fatalf("Invalid JSON: %s", jsonString)
	}

	return gjson.Parse(jsonString)
}

// GetTFSetupOutputListVal returns TF output from setup for a given key as list.
// It fails test if given key does not output a list type.
func (b *TFBlueprintTest) GetTFSetupOutputListVal(key string) []string {
	if v, ok := b.setupOutputOverrides[key]; ok {
		if listval, ok := v.([]string); ok {
			return listval
		} else {
			b.t.Fatalf("Setup Override %s is not a list value", key)
		}
	}
	if b.setupDir == "" {
		b.t.Fatal("Setup path not set")
	}
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()
	return terraform.OutputList(b.t, &terraform.Options{TerraformDir: b.setupDir, Logger: b.logger, NoColor: true}, key)
}

// GetTFSetupStringOutput returns TF setup output for a given key as string.
// It fails test if given key does not output a primitive or if setupDir is not configured.
func (b *TFBlueprintTest) GetTFSetupStringOutput(key string) string {
	if v, ok := b.setupOutputOverrides[key]; ok {
		return v.(string)
	}
	if b.setupDir == "" {
		b.t.Fatal("Setup path not set")
	}
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()
	return terraform.Output(b.t, &terraform.Options{TerraformDir: b.setupDir, Logger: b.logger, NoColor: true}, key)
}

// GetTFSetupJsonOutput returns TF setup output for a given key as gjson.Result.
// An empty string for key can be used to return all values.
// It fails test if given key does not output valid JSON or if setupDir is not configured.
func (b *TFBlueprintTest) GetTFSetupJsonOutput(key string) gjson.Result {
	if v, ok := b.setupOutputOverrides[key]; ok {
		if !gjson.Valid(v.(string)) {
			b.t.Fatalf("Invalid JSON in setup output override: %s", v)
		}
		return gjson.Parse(v.(string))
	}
	if b.setupDir == "" {
		b.t.Fatal("Setup path not set")
	}
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()

	jsonString := terraform.OutputJson(b.t, &terraform.Options{TerraformDir: b.setupDir, Logger: b.logger, NoColor: true}, key)
	if !gjson.Valid(jsonString) {
		b.t.Fatalf("Invalid JSON: %s", jsonString)
	}

	return gjson.Parse(jsonString)
}

// loadTFEnvVar adds new env variables prefixed with TF_VAR_ to an existing map of variables.
func loadTFEnvVar(m map[string]string, new map[string]string) {
	for k, v := range new {
		m[fmt.Sprintf("TF_VAR_%s", k)] = v
	}
}

// extractFromEnv parses environment variables with the given prefix, and returns a key-value map.
// e.g. CFT_SETUP_key=value returns map[string]string{"key": "value"}
func extractFromEnv(prefix string) map[string]interface{} {
	r := make(map[string]interface{})
	for _, s := range os.Environ() {
		k, v, ok := strings.Cut(s, "=")
		if !ok {
			// skip malformed entries in os.Environ
			continue
		}
		// For env vars with the prefix, extract the key and value
		if setupvar, ok := strings.CutPrefix(k, prefix); ok {
			r[setupvar] = v
		}
	}
	return r
}

// ShouldSkip checks if a test should be skipped
func (b *TFBlueprintTest) ShouldSkip() bool {
	return b.BlueprintTestConfig.Spec.Skip
}

// shouldRunTerraformVet checks if terraform vet should be executed
func (b *TFBlueprintTest) shouldRunTerraformVet() bool {
	return b.policyLibraryPath != ""
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

// DefinePlan defines a custom plan function for the blueprint.
func (b *TFBlueprintTest) DefinePlan(plan func(*terraform.PlanStruct, *assert.Assertions)) {
	b.plan = plan
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
	assert.NotEqual(2, e, "plan after apply should have no diff")
}

// DefaultInit runs TF init and validate on a blueprint.
func (b *TFBlueprintTest) DefaultInit(assert *assert.Assertions) {
	terraform.Init(b.t, b.GetTFOptions())
	// if vars are set for common options, this seems to trigger -var flag when calling validate
	// using custom tfOptions as a workaround
	terraform.Validate(b.t, terraform.WithDefaultRetryableErrors(b.t, &terraform.Options{
		TerraformDir: b.tfDir,
		Logger:       b.logger,
		NoColor:      true,
	}))
}

// Vet runs TF plan, TF show, and gcloud terraform vet on a blueprint.
func (b *TFBlueprintTest) Vet(assert *assert.Assertions) {
	jsonPlan, _ := b.PlanAndShow()
	filepath, err := utils.WriteTmpFileWithExtension(jsonPlan, "json")
	assert.NoError(err)
	defer func() {
		if err := os.Remove(filepath); err != nil {
			b.t.Fatalf("Could not remove plan json: %v", err)
		}
	}()
	results := gcloud.TFVet(b.t, filepath, b.policyLibraryPath, b.terraformVetProject).Array()
	assert.Empty(results, "Should have no Terraform Vet violations")
}

// DefaultApply runs TF apply on a blueprint.
func (b *TFBlueprintTest) DefaultApply(assert *assert.Assertions) {
	if b.shouldRunTerraformVet() {
		b.Vet(assert)
	}
	terraform.Apply(b.t, b.GetTFOptions())
}

// Init runs the default or custom init function for the blueprint.
func (b *TFBlueprintTest) Init(assert *assert.Assertions) {
	// allow only single write as Terraform plugin cache isn't concurrent safe
	if err := b.tftCacheMutex.Lock(); err != nil {
		b.t.Fatalf("Could not acquire lock: %v", err)
	}
	defer func() {
		if err := b.tftCacheMutex.Unlock(); err != nil {
			b.t.Fatalf("Could not release lock: %v", err)
		}
	}()
	b.init(assert)
}

// PlanAndShow performs a Terraform plan, show and returns the parsed plan output.
func (b *TFBlueprintTest) PlanAndShow() (string, *terraform.PlanStruct) {
	tDir, err := os.MkdirTemp(os.TempDir(), "btp")
	if err != nil {
		b.t.Fatalf("Temp directory %q could not created: %v", tDir, err)
	}
	defer func() {
		if err := os.RemoveAll(tDir); err != nil {
			b.t.Fatalf("Could not remove plan temp dir: %v", err)
		}
	}()

	planOpts := b.GetTFOptions()
	planOpts.PlanFilePath = filepath.Join(tDir, planFilename)
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()
	terraform.Plan(b.t, planOpts)
	// Logging show output is not useful since we log plan output above
	// and show output is parsed and retured.
	planOpts.Logger = logger.Discard
	planJSON := terraform.Show(b.t, planOpts)
	ps, err := terraform.ParsePlanJSON(planJSON)
	assert.NoError(b.t, err)
	return planJSON, ps
}

// Plan runs the custom plan function for the blueprint.
// If not custom plan function is defined, this stage is skipped.
func (b *TFBlueprintTest) Plan(assert *assert.Assertions) {
	if b.plan == nil {
		b.logger.Logf(b.t, "skipping plan as no function defined")
		return
	}
	_, ps := b.PlanAndShow()
	b.plan(ps, assert)
}

// Apply runs the default or custom apply function for the blueprint.
func (b *TFBlueprintTest) Apply(assert *assert.Assertions) {
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()
	b.apply(assert)
}

// Verify runs the default or custom verify function for the blueprint.
func (b *TFBlueprintTest) Verify(assert *assert.Assertions) {
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()
	b.verify(assert)
}

// Teardown runs the default or custom teardown function for the blueprint.
func (b *TFBlueprintTest) Teardown(assert *assert.Assertions) {
	// allow only parallel reads as Terraform plugin cache isn't concurrent safe
	rUnlockFn := b.rLockFn()
	defer rUnlockFn()
	b.teardown(assert)
}

const (
	initStage     = "init"
	planStage     = "plan"
	applyStage    = "apply"
	verifyStage   = "verify"
	teardownStage = "teardown"
)

// Test runs init, apply, verify, teardown in order for the blueprint.
func (b *TFBlueprintTest) Test() {
	if b.ShouldSkip() {
		b.logger.Logf(b.t, "Skipping test due to config %s", b.BlueprintTestConfig.Path)
		b.t.SkipNow()
		return
	}
	a := assert.New(b.t)
	// run stages
	utils.RunStage(initStage, func() { b.Init(a) })
	defer utils.RunStage(teardownStage, func() { b.Teardown(a) })
	utils.RunStage(planStage, func() { b.Plan(a) })
	utils.RunStage(applyStage, func() { b.Apply(a) })
	utils.RunStage(verifyStage, func() { b.Verify(a) })
}

// RedeployTest deploys the test n times in separate workspaces before teardown.
func (b *TFBlueprintTest) RedeployTest(n int, nVars map[int]map[string]interface{}) {
	if n < 2 {
		b.t.Fatalf("n should be 2 or greater but got: %d", n)
	}
	if b.ShouldSkip() {
		b.logger.Logf(b.t, "Skipping test due to config %s", b.BlueprintTestConfig.Path)
		b.t.SkipNow()
		return
	}
	a := assert.New(b.t)
	// capture currently set vars as default if no override
	defaultVars := b.vars
	overrideVars := func(i int) {
		custom, exists := nVars[i]
		if exists {
			b.vars = custom
		} else {
			b.vars = defaultVars
		}
	}
	for i := 1; i <= n; i++ {
		ws := terraform.WorkspaceSelectOrNew(b.t, b.GetTFOptions(), fmt.Sprintf("test-%d", i))
		overrideVars(i)
		utils.RunStage(initStage, func() { b.Init(a) })
		defer func(i int) {
			overrideVars(i)
			terraform.WorkspaceSelectOrNew(b.t, b.GetTFOptions(), ws)
			utils.RunStage(teardownStage, func() { b.Teardown(a) })
		}(i)
		utils.RunStage(planStage, func() { b.Plan(a) })
		utils.RunStage(applyStage, func() { b.Apply(a) })
		utils.RunStage(verifyStage, func() { b.Verify(a) })
	}
}

// rLockFn sets a read mutex lock, and returns the corresponding unlock function.
func (b *TFBlueprintTest) rLockFn() func() {
	if err := b.tftCacheMutex.RLock(); err != nil {
		b.t.Fatalf("Could not acquire read lock:%v", err)
	}

	return func() {
		if err := b.tftCacheMutex.RUnlock(); err != nil {
			b.t.Fatalf("Could not release read lock: %v", err)
		}
	}
}
