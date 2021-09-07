package kpt

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/binary"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/mitchellh/go-testing-interface"
)

type CmdCfg struct {
	*binary.BinaryCfg                       // binary config
	bOpts             []binary.BinaryOption // binary options
	t                 testing.TB            // TestingT or TestingB
}

type cmdOption func(*CmdCfg)

func WithBinaryOptions(bOpts ...binary.BinaryOption) cmdOption {
	return func(f *CmdCfg) {
		f.bOpts = append(f.bOpts, bOpts...)
	}
}

// NewCmdConfig sets defaults and validates values for kpt Options.
func NewCmdConfig(t testing.TB, opts ...cmdOption) *CmdCfg {
	kOpts := &CmdCfg{
		t: t,
	}
	// apply options
	for _, opt := range opts {
		opt(kOpts)
	}
	kOpts.BinaryCfg = binary.NewBinaryConfig(t, "kpt", kOpts.bOpts...)
	return kOpts
}

func (k *CmdCfg) RunCmd(args ...string) string {
	kptCmd := shell.Command{
		Command:    k.GetBinary(),
		Args:       args,
		Logger:     k.GetLogger(),
		WorkingDir: k.GetDir(),
	}
	op, err := shell.RunCommandAndGetStdOutE(k.t, kptCmd)
	if err != nil {
		k.t.Fatal(err)
	}
	return op
}
