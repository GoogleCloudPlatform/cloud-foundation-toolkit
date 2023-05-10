package krmt

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	gotest "testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/discovery"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/git"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/kpt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/mitchellh/go-testing-interface"
	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

const (
	tmpBuildDir = ".build"
	prowPullSha = "PULL_PULL_SHA"
)

var CommonSetters = []string{"PROJECT_ID", "BILLING_ACCOUNT_ID", "ORG_ID"}

// KRMBlueprintTest implements bpt.Blueprint and stores information associated with a KRM blueprint test.
type KRMBlueprintTest struct {
	discovery.BlueprintTestConfig                          // additional blueprint test configs
	name                          string                   // descriptive name for the test
	exampleDir                    string                   // directory containing KRM blueprint example
	additionalResources           []string                 // paths to directories or files containing additional resources to be applied
	buildDir                      string                   // directory to hydrated blueprint configs pre apply
	kpt                           *kpt.CmdCfg              // kpt cmd config
	timeout                       string                   // timeout for KRM resource status
	setters                       map[string]string        // additional setters to populate
	updatePkgs                    bool                     // whether to update packages in exampleDir
	updateCommit                  string                   // specific commit to update to
	logger                        *logger.Logger           // custom logger
	t                             testing.TB               // TestingT or TestingB
	init                          func(*assert.Assertions) // init function
	apply                         func(*assert.Assertions) // apply function
	verify                        func(*assert.Assertions) // verify function
	teardown                      func(*assert.Assertions) // teardown function
}

type krmtOption func(*KRMBlueprintTest)

func WithName(name string) krmtOption {
	return func(f *KRMBlueprintTest) {
		f.name = name
	}
}

func WithDir(dir string) krmtOption {
	return func(f *KRMBlueprintTest) {
		f.exampleDir = dir
	}
}

func WithAdditionalResources(rscs ...string) krmtOption {
	return func(f *KRMBlueprintTest) {
		f.additionalResources = append(f.additionalResources, rscs...)
	}
}

func WithBuildDir(buildDir string) krmtOption {
	return func(f *KRMBlueprintTest) {
		f.buildDir = buildDir
	}
}

func WithUpdatePkgs(update bool) krmtOption {
	return func(f *KRMBlueprintTest) {
		f.updatePkgs = update
	}
}

func WithUpdateCommit(commit string) krmtOption {
	return func(f *KRMBlueprintTest) {
		f.updateCommit = commit
	}
}

func WithTimeout(timeout string) krmtOption {
	return func(f *KRMBlueprintTest) {
		f.timeout = timeout
	}
}

func WithLogger(logger *logger.Logger) krmtOption {
	return func(f *KRMBlueprintTest) {
		f.logger = logger
	}
}

func WithSetters(setters map[string]string) krmtOption {
	return func(f *KRMBlueprintTest) {
		f.setters = setters
	}
}

// NewKRMBlueprintTest sets defaults, validates and returns a KRMBlueprintTest.
func NewKRMBlueprintTest(t testing.TB, opts ...krmtOption) *KRMBlueprintTest {
	krmt := &KRMBlueprintTest{
		name:       fmt.Sprintf("%s KRM Blueprint", t.Name()),
		t:          t,
		setters:    make(map[string]string),
		updatePkgs: true,
		timeout:    "10m",
	}
	// default KRM blueprint methods
	krmt.init = krmt.DefaultInit
	krmt.apply = krmt.DefaultApply
	krmt.verify = krmt.DefaultVerify
	krmt.teardown = krmt.DefaultTeardown
	// apply options
	for _, opt := range opts {
		opt(krmt)
	}
	// if no custom logger, set default based on test verbosity
	if krmt.logger == nil {
		krmt.logger = utils.GetLoggerFromT()
	}
	// if explicit exampleDir is provided, validate it else try auto discovery
	if krmt.exampleDir != "" {
		_, err := os.Stat(krmt.exampleDir)
		if os.IsNotExist(err) {
			t.Fatalf("Dir path %s does not exist", krmt.exampleDir)
		}
	} else {
		exampleDir, err := discovery.GetConfigDirFromTestDir(utils.GetWD(t))
		if err != nil {
			t.Fatalf("unable to detect KRM dir :%v", err)
		}
		krmt.exampleDir = exampleDir
	}
	// if explicit resourcesDir is provided, validate it.
	if len(krmt.additionalResources) != 0 {
		for _, path := range krmt.additionalResources {
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				t.Fatalf("Path for additional resources %s does not exist", path)
			}
		}
	}

	// discover test config
	var err error
	krmt.BlueprintTestConfig, err = discovery.GetTestConfig(path.Join(krmt.exampleDir, discovery.DefaultTestConfigFilename))
	if err != nil {
		t.Fatal(err)
	}
	// if no explicit build directory is provided, setup build directory
	if krmt.buildDir == "" {
		krmt.buildDir = krmt.getDefaultBuildDir()
	}
	// configure kpt to run in buildDir
	krmt.kpt = kpt.NewCmdConfig(t, kpt.WithDir(krmt.buildDir))
	// get well known setters from env vars
	krmt.getKnownSettersFromEnv()

	krmt.logger.Logf(krmt.t, "Running tests KRM configs in %s", krmt.exampleDir)
	return krmt
}

// getDefaultBuildDir returns a temporary build directory for hydrated configs.
func (b *KRMBlueprintTest) getDefaultBuildDir() string {
	buildDir := path.Join(utils.GetWD(b.t), tmpBuildDir)
	err := os.MkdirAll(buildDir, 0755)
	if err != nil {
		b.t.Fatalf("unable to create %s :%v", buildDir, err)
	}
	exampleBuildDir := path.Join(buildDir, b.t.Name())
	abs, err := filepath.Abs(exampleBuildDir)
	if err != nil {
		b.t.Fatalf("unable to get absolute path for %s :%v", exampleBuildDir, err)
	}
	return abs
}

// getKnownSettersFromEnv creates setters from known CommonSetters env vars.
func (b *KRMBlueprintTest) getKnownSettersFromEnv() {
	setters := make(map[string]string)
	for _, s := range CommonSetters {
		sKey, sVal, err := kpt.GenerateSetterKVFromEnvVar(s)
		// if a common env var is not set, log and continue
		if err != nil {
			b.logger.Logf(b.t, "Skipping env var %s: %v", s, err)
		} else {
			b.logger.Logf(b.t, "Setting env var %s as setter %s: %s", s, sKey, sVal)
			setters[sKey] = sVal
		}
	}
	// merge user provided setters with env discovered setters
	// user provided setters override env discovered setters
	b.setters = kpt.MergeSetters(setters, b.setters)
}

// setupBuildDir prepares build dir with configs from exampleDir.
func (b *KRMBlueprintTest) setupBuildDir() {
	// remove buildDir if exists
	err := os.RemoveAll(b.buildDir)
	if err != nil {
		b.t.Fatalf("unable to remove %s :%v", b.buildDir, err)
	}
	// copy over configs into build dir
	err = copy.Copy(b.exampleDir, b.buildDir)
	if err != nil {
		b.t.Fatalf("unable to copy %s to %s :%v", b.exampleDir, b.buildDir, err)
	}
	// copy over additional resources into build dir, if present
	if len(b.additionalResources) != 0 {
		for _, path := range b.additionalResources {
			err = copy.Copy(path, b.buildDir)
			if err != nil {
				b.t.Fatalf("unable to copy %s to %s :%v", path, b.buildDir, err)
			}
		}
	}
	// subsequent kpt pkg update requires a clean git repo without uncommitted changes
	// init a new git repo in build dir and commit changes
	git := git.NewCmdConfig(b.t, git.WithDir(b.buildDir))
	git.Init()
	git.AddAll()
	git.Commit()
}

// updateSetters updates existing setters with user provided setters.
func (b *KRMBlueprintTest) updateSetters() {
	rs, err := kpt.ReadPkgResources(b.buildDir)
	if err != nil {
		b.t.Fatalf("unable to read resources in %s :%v", b.buildDir, err)
	}
	if err := kpt.UpsertSetters(rs, b.setters); err != nil {
		b.t.Fatalf("unable to upsert setters in %s :%v", b.buildDir, err)
	}
	err = kpt.WritePkgResources(b.buildDir, rs)
	if err != nil {
		b.t.Fatalf("unable to write resources in %s :%v", b.buildDir, err)
	}
}

// updatePkg updates a kpt pkg to a specified commit, prow PR commit or the latest commit.
func (b *KRMBlueprintTest) updatePkg() {
	g := git.NewCmdConfig(b.t, git.WithDir(b.exampleDir))
	commit := b.updateCommit
	// no explicit commit specified
	if commit == "" {
		// check if prow PR commit exists
		prowCommit, found := os.LookupEnv(prowPullSha)
		if found {
			commit = prowCommit
		} else {
			commit = g.GetLatestCommit()
		}
	}
	b.kpt.RunCmd("pkg", "update", fmt.Sprintf(".@%s", commit))

}

// DefaultInit sets up build directory, updates pkg, upserts setters and renders config.
func (b *KRMBlueprintTest) DefaultInit(assert *assert.Assertions) {
	b.setupBuildDir()
	if b.updatePkgs {
		b.updatePkg()
	}
	b.updateSetters()
	kpt.NewCmdConfig(b.t, kpt.WithDir(b.buildDir)).RunCmd("fn", "render")
	b.kpt.RunCmd("live", "install-resource-group")
	b.kpt.RunCmd("live", "init")
}

// DefaultApply installs resource-group, initializes inventory, applies pkg and polls resource statuses until current.
func (b *KRMBlueprintTest) DefaultApply(assert *assert.Assertions) {
	b.kpt.RunCmd("live", "apply")
	b.kpt.RunCmd("live", "status", "--output", "json", "--poll-until", "current", "--timeout", b.timeout)
}

// DefaultVerify asserts all resources are status successful
func (b *KRMBlueprintTest) DefaultVerify(assert *assert.Assertions) {
	jsonOp := b.kpt.RunCmd("live", "apply", "--output", "json")

	// assert each resource status is successful
	resourceStatus, err := kpt.GetPkgApplyResourcesStatus(jsonOp)
	assert.NoError(err, "Resource statuses should be parsable")
	for _, r := range resourceStatus {
		assert.Equal(kpt.ResourceOperationSuccessful, r.Status, "Status should be successful")
	}

	// assert count of resources applied equals count of resources successful
	groupStatus, err := kpt.GetPkgApplyGroupStatus(jsonOp)
	assert.NoError(err, "Group status should be parsable")
	assert.Equal(groupStatus.Count, groupStatus.Successful, "All resources should be successful")

}

// DefaultTeardown destroys resources from cluster and polls until deleted.
func (b *KRMBlueprintTest) DefaultTeardown(assert *assert.Assertions) {
	b.kpt.RunCmd("live", "destroy")
	b.kpt.RunCmd("live", "status", "--output", "json", "--poll-until", "deleted", "--timeout", b.timeout)
}

// ShouldSkip checks if a test should be skipped
func (b *KRMBlueprintTest) ShouldSkip() bool {
	return b.BlueprintTestConfig.Spec.Skip
}

// AutoDiscoverAndTest discovers KRM config from examples/fixtures and runs tests.
func AutoDiscoverAndTest(t *gotest.T) {
	configs := discovery.FindTestConfigs(t, "./")
	for testName, dir := range configs {
		t.Run(testName, func(t *gotest.T) {
			nt := NewKRMBlueprintTest(t, WithDir(dir))
			nt.Test()
		})
	}
}

// DefineInit defines a custom init function for the blueprint.
func (b *KRMBlueprintTest) DefineInit(init func(*assert.Assertions)) {
	b.init = init
}

// DefineApply defines a custom apply function for the blueprint.
func (b *KRMBlueprintTest) DefineApply(apply func(*assert.Assertions)) {
	b.apply = apply
}

// DefineVerify defines a custom verify function for the blueprint.
func (b *KRMBlueprintTest) DefineVerify(verify func(*assert.Assertions)) {
	b.verify = verify
}

// DefineTeardown defines a custom teardown function for the blueprint.
func (b *KRMBlueprintTest) DefineTeardown(teardown func(*assert.Assertions)) {
	b.teardown = teardown
}

// Init runs the default or custom init function for the blueprint.
func (b *KRMBlueprintTest) Init(assert *assert.Assertions) {
	b.init(assert)
}

// Apply runs the default or custom apply function for the blueprint.
func (b *KRMBlueprintTest) Apply(assert *assert.Assertions) {
	b.apply(assert)
}

// Verify runs the default or custom verify function for the blueprint.
func (b *KRMBlueprintTest) Verify(assert *assert.Assertions) {
	b.verify(assert)
}

// Teardown runs the default or custom teardown function for the blueprint.
func (b *KRMBlueprintTest) Teardown(assert *assert.Assertions) {
	b.teardown(assert)
}

// Test runs init, apply, verify, teardown in order for the blueprint.
func (b *KRMBlueprintTest) Test() {
	if b.ShouldSkip() {
		b.logger.Logf(b.t, "Skipping test due to config %s", b.BlueprintTestConfig.Path)
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

// GetBuildDir returns the temporary build dir created for hydrating config. Defaults to .build/test-name.
func (b *KRMBlueprintTest) GetBuildDir() string {
	if b.buildDir == "" {
		b.t.Fatalf("unable to get a valid build directory")
	}

	return b.buildDir
}
