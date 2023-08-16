package kpt

import (
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	kptfilev1 "github.com/GoogleContainerTools/kpt-functions-sdk/go/api/kptfile/v1"
	kptutil "github.com/GoogleContainerTools/kpt-functions-sdk/go/api/util"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/mitchellh/go-testing-interface"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// MIN_KPT_VERSION format: vMAJOR[.MINOR[.PATCH[-PRERELEASE]]]
const MIN_KPT_VERSION = "v1.0.0-beta.16"

type CmdCfg struct {
	kptBinary string         // kpt binary
	dir       string         // dir to execute commands in
	logger    *logger.Logger // custom logger
	t         testing.TB     // TestingT or TestingB
	tries     int            // qty to try kpt command, default: 3
}

type cmdOption func(*CmdCfg)

func WithDir(dir string) cmdOption {
	return func(f *CmdCfg) {
		f.dir = dir
	}
}

func WithBinary(kptBinary string) cmdOption {
	return func(f *CmdCfg) {
		f.kptBinary = kptBinary
	}
}

func WithLogger(logger *logger.Logger) cmdOption {
	return func(f *CmdCfg) {
		f.logger = logger
	}
}

// NewCmdConfig sets defaults and validates values for kpt Options.
func NewCmdConfig(t testing.TB, opts ...cmdOption) *CmdCfg {
	kOpts := &CmdCfg{
		logger: utils.GetLoggerFromT(),
		t:      t,
		tries:  3,
	}
	// apply options
	for _, opt := range opts {
		opt(kOpts)
	}
	if kOpts.kptBinary == "" {
		err := utils.BinaryInPath("kpt")
		if err != nil {
			t.Fatalf("unable to find kpt in path: %v", err)
		}
		kOpts.kptBinary = "kpt"
	}
	// Validate required KPT version
	if err := utils.MinSemver("v"+GetKptVersion(t, kOpts.kptBinary), MIN_KPT_VERSION); err != nil {
		t.Fatalf("unable to validate minimum required kpt version: %v", err)
	}

	return kOpts
}

func (k *CmdCfg) RunCmd(args ...string) string {
	kptCmd := shell.Command{
		Command:    "kpt",
		Args:       args,
		Logger:     k.logger,
		WorkingDir: k.dir,
	}
	command := func() (string, error) {
		return shell.RunCommandAndGetStdOutE(k.t, kptCmd)
	}
	op, err := retry.DoWithRetryE(k.t, fmt.Sprintf("kpt %v",  kptCmd.Args), k.tries, 15*time.Second, command)
	if err != nil {
		k.t.Fatal(err)
	}
	return op
}

// findKptfile discovers Kptfile of the root package from slice of nodes
func findKptfile(nodes []*yaml.RNode) (*kptfilev1.KptFile, error) {
	for _, node := range nodes {
		if node.GetAnnotations()[kioutil.PathAnnotation] == kptfilev1.KptFileName {
			s, err := node.String()
			if err != nil {
				return nil, fmt.Errorf("unable to read Kptfile: %v", err)
			}
			kf, err := kptutil.DecodeKptfile(s)
			if err != nil {
				return nil, fmt.Errorf("unable to decode Kptfile: %v", err)
			}
			return kf, nil
		}
	}
	return nil, fmt.Errorf("unable to find Kptfile, please include --include-meta-resources flag if a Kptfile is present")
}

// GetKptVersion gets the version of kptBinary
func GetKptVersion(t testing.TB, kptBinary string) string {
	kVersionOpts := &CmdCfg{
		kptBinary: kptBinary,
		dir:       utils.GetWD(t),
		logger:    utils.GetLoggerFromT(),
		t:         t,
	}
	return kVersionOpts.RunCmd("version")
}
