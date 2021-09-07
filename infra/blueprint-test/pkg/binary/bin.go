package binary

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/mitchellh/go-testing-interface"
)

type BinaryCfg struct {
	binary string         // path to  binary
	logger *logger.Logger // custom logger
	dir    string         // dir to execute commands in
}

type BinaryOption func(*BinaryCfg)

func WithBinary(binary string) BinaryOption {
	return func(f *BinaryCfg) {
		f.binary = binary
	}
}

func WithLogger(logger *logger.Logger) BinaryOption {
	return func(f *BinaryCfg) {
		f.logger = logger
	}
}

func WithDir(dir string) BinaryOption {
	return func(f *BinaryCfg) {
		f.dir = dir
	}
}

// NewBinaryConfig sets defaults and validates configuration for running a binary
func NewBinaryConfig(t testing.TB, bin string, opts ...BinaryOption) *BinaryCfg {
	b := &BinaryCfg{binary: bin, logger: utils.GetLoggerFromT()}
	for _, opt := range opts {
		opt(b)
	}
	err := utils.BinaryInPath(b.binary)
	if err != nil {
		t.Fatalf("unable to find binary %s in path: %v", b.binary, err)
	}
	return b
}

func (b *BinaryCfg) GetBinary() string {
	return b.binary
}

func (b *BinaryCfg) GetLogger() *logger.Logger {
	return b.logger
}

func (b *BinaryCfg) GetDir() string {
	return b.dir
}
