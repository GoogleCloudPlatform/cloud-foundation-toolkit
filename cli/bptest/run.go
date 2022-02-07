package bptest

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sync"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/spf13/viper"
)

const (
	allTests           = "all"
	testStageEnvVarKey = "RUN_STAGE"
	gotestBin          = "gotest"
	goBin              = "go"
)

var allTestArgs = []string{"-p", "1", "-count", "1", "-timeout", "0"}

// isValidTestName validates a given test or test regex is part of the blueprint test set
func isValidTestName(intTestDir string, name string) error {
	// user wants to run all tests
	if name == allTests {
		return nil
	}

	tests, err := getTests(intTestDir)
	if err != nil {
		return err
	}
	testNames := []string{}
	for _, test := range tests {
		if test.bptestCfg.Spec.Skip {
			Log.Info(fmt.Sprintf("skipping %s due to BlueprintTest config %s", test.name, test.bptestCfg.Name))
			continue
		}
		matched, _ := regexp.Match(name, []byte(test.name))
		if test.name == name || matched {
			return nil
		}
		testNames = append(testNames, test.name)
	}
	return fmt.Errorf("unable to find %s- one of %+q expected", name, append(testNames, allTests))
}

// streamExec runs a given cmd while streaming logs
func streamExec(cmd *exec.Cmd) error {
	op, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout
	Log.Debug(fmt.Sprintf("running %s with args %v in %s", cmd.Path, cmd.Args, cmd.Dir))

	// waitgroup to block while processing exec op
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(op)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			Log.Error(fmt.Sprintf("error reading output: %v", err))
		}
	}()

	// run command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command: %v", err)
	}
	return nil
}

// getTestCmd returns a prepared cmd for running the specified tests(s)
func getTestCmd(intTestDir string, testStage string, testName string) (*exec.Cmd, error) {
	intTestDir, err := getIntTestDir(intTestDir)
	if err != nil {
		return nil, err
	}

	// pass all current env vars to test command
	env := os.Environ()
	// set test stage env var if specified
	if testStage != "" {
		env = append(env, fmt.Sprintf("%s=%s", testStageEnvVarKey, testStage))
	}

	// determine binary and args used for test execution
	testArgs := append([]string{"./..."}, allTestArgs...)
	if testName != allTests {
		testArgs = append([]string{"./...", "-run", testName}, allTestArgs...)
	}
	cmdBin := goBin
	if utils.BinaryInPath(gotestBin) != nil {
		testArgs = append([]string{"test"}, testArgs...)
	} else {
		cmdBin = gotestBin
		// CI=true enables color op for non tty exec output
		env = append(env, "CI=true")
	}
	// verbose test output if global verbose flag is passed
	if viper.GetBool("verbose") {
		testArgs = append(testArgs, "-v")
	}
	// prepare cmd
	cmd := exec.Command(cmdBin, testArgs...)
	cmd.Env = env
	cmd.Dir = intTestDir
	return cmd, nil
}
