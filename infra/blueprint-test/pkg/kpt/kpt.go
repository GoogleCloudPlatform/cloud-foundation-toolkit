package kpt

import (
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	kptfilev1 "github.com/GoogleContainerTools/kpt-functions-sdk/go/api/kptfile/v1"
	kptutil "github.com/GoogleContainerTools/kpt-functions-sdk/go/api/util"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/mitchellh/go-testing-interface"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type CmdCfg struct {
	kptBinary string         // kpt binary
	dir       string         // dir to execute commands in
	logger    *logger.Logger // custom logger
	t         testing.TB     // TestingT or TestingB
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
	return kOpts
}

func (k *CmdCfg) RunCmd(args ...string) string {
	kptCmd := shell.Command{
		Command:    "kpt",
		Args:       args,
		Logger:     k.logger,
		WorkingDir: k.dir,
	}
	op, err := shell.RunCommandAndGetStdOutE(k.t, kptCmd)
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
			return kf, fmt.Errorf("unable to read Kptfile: %v", err)
		}
	}
	return nil, fmt.Errorf("unable to find Kptfile, please include --include-meta-resources flag if a Kptfile is present")
}
