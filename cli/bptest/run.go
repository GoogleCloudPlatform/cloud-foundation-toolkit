package bptest

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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

	// The tfplan.json files that are being used as input for the terraform validation tests
	// through the gcloud beta terraform vet are higher than the buffer default value (64*1024),
	// after some tests we had evidences that the value were around from 3MB to 5MB, so
	// we choosed a value that is at least 2x higher than the original one to avoid errors.
	// maxScanTokenSize is the maximum size used to buffer a token
	// startBufSize is the initial of the buffer token
	maxScanTokenSize = 10 * 1024 * 1024
	startBufSize     = 4096
	// This must be kept in sync with what github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft parses.
	setupEnvVarPrefix = "CFT_SETUP_"
)

var allTestArgs = []string{"-p", "1", "-count", "1", "-timeout", "0"}

// validateAndGetRelativeTestPkg validates a given test or test regex is part of the blueprint test set and returns location of test relative to intTestDir
func validateAndGetRelativeTestPkg(intTestDir string, name string) (string, error) {
	// user wants to run all tests
	if name == allTests {
		return "./...", nil
	}

	tests, err := getTests(intTestDir)
	if err != nil {
		return "", err
	}
	testNames := []string{}
	for _, test := range tests {
		if test.bptestCfg.Spec.Skip {
			Log.Info(fmt.Sprintf("skipping %s due to BlueprintTest config %s", test.name, test.bptestCfg.Name))
			continue
		}
		matched, _ := regexp.Match(name, []byte(test.name))
		if test.name == name {
			//exact match, return test relative test pkg
			relPkg, err := filepath.Rel(intTestDir, path.Dir(test.location))
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("./%s", relPkg), nil
		} else if matched {
			// loose match, more than one test could be specified
			return "./...", nil
		}
		testNames = append(testNames, test.name)
	}
	return "", fmt.Errorf("unable to find %s- one of %+q expected", name, append(testNames, allTests))
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
		scanner.Buffer(make([]byte, startBufSize), maxScanTokenSize)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			Log.Error(fmt.Sprintf("error reading output: %s", err))
		}
	}()

	// run command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}
	return nil
}

// getTestCmd returns a prepared cmd for running the specified tests(s)
func getTestCmd(intTestDir string, testStage string, testName string, relTestPkg string, setupVars map[string]string) (*exec.Cmd, error) {

	// pass all current env vars to test command
	env := os.Environ()
	// set test stage env var if specified
	if testStage != "" {
		env = append(env, fmt.Sprintf("%s=%s", testStageEnvVarKey, testStage))
	}
	// Load the env with any setup-vars specified
	for k, v := range setupVars {
		env = append(env, fmt.Sprintf("%s%s=%s", setupEnvVarPrefix, k, v))
	}

	// determine binary and args used for test execution
	testArgs := append([]string{relTestPkg}, allTestArgs...)
	if testName != allTests {
		testArgs = append([]string{relTestPkg, "-run", testName}, allTestArgs...)
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
