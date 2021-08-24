package git

import (
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/mitchellh/go-testing-interface"
)

type CmdCfg struct {
	gitBinary string         // git binary
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

func WithLogger(logger *logger.Logger) cmdOption {
	return func(f *CmdCfg) {
		f.logger = logger
	}
}

// NewCmdConfig sets defaults and validates values for git Options.
func NewCmdConfig(t testing.TB, opts ...cmdOption) *CmdCfg {
	gitOpts := &CmdCfg{
		logger: utils.GetLoggerFromT(),
		t:      t,
	}
	// apply options
	for _, opt := range opts {
		opt(gitOpts)
	}
	if gitOpts.gitBinary == "" {
		err := utils.BinaryInPath("git")
		if err != nil {
			t.Fatalf("unable to find git in path: %v", err)
		}
		gitOpts.gitBinary = "git"
	}
	return gitOpts
}

// RunCmd executes a git command
func (g *CmdCfg) RunCmdE(args ...string) (string, error) {
	kptCmd := shell.Command{
		Command:    g.gitBinary,
		Args:       args,
		Logger:     g.logger,
		WorkingDir: g.dir,
	}
	return shell.RunCommandAndGetStdOutE(g.t, kptCmd)
}

// GetLatestCommit returns latest commit
func (g *CmdCfg) GetLatestCommit() string {
	commit, err := g.RunCmdE("rev-parse", "HEAD")
	if err != nil {
		g.t.Fatalf("error getting latest commit: %v", err)
	}
	return commit
}

// Init run git init
func (g *CmdCfg) Init() {
	_, err := g.RunCmdE("init")
	if err != nil {
		g.t.Fatalf("error running git init: %v", err)
	}
}

// AddAll stages all changes.
func (g *CmdCfg) AddAll() {
	_, err := g.RunCmdE("add", "-A")
	if err != nil {
		g.t.Fatalf("error running git add: %v", err)
	}
}

// CommitWithMsg commits changes with commit msg.
func (g *CmdCfg) CommitWithMsg(msg string, commitFlags []string) {
	_, err := g.RunCmdE(append([]string{"commit", "-m", fmt.Sprintf("%q", msg)}, commitFlags...)...)
	if err != nil {
		g.t.Fatalf("error running git commit: %v", err)
	}
}

// CommitWithMsg commits changes with a generated commit msg.
func (g *CmdCfg) Commit() {
	currentTime := time.Now()
	g.CommitWithMsg(fmt.Sprintf("commit %s", currentTime.Format(time.RFC1123)), []string{"--author", "BlueprintsTest <blueprints-ci-test@google.com>"})
}
