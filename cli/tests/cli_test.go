// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	cwd       string
	cftBinary string
)

// TestCLI does integration testing.
func TestCLI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	// Test cases for each command.
	cases := []struct {
		name    string
		testers []struct {
			name   string
			tester func(t *testing.T, dirTestdata string)
		}
		setup   func(t *testing.T, dirTestdata string)
		destroy func(t *testing.T, dirTestdata string)
	}{
		{
			name: "validate",
			testers: []struct {
				name   string
				tester func(t *testing.T, dirTestdata string)
			}{
				{
					name:   "passing",
					tester: testValidatePassing,
				},
				{
					name:   "failing",
					tester: testValidateFailing,
				},
			},
			setup:   setupValidate,
			destroy: destroyValidate,
		},
	}

	var err error
	cwd, err = os.Getwd()
	if err != nil {
		t.Fatalf("cannot get current directory: %v", err)
	}
	cftBinary = filepath.Join(cwd, "..", "bin", "cft")

	for i := range cases {
		// Allocate a variable to make sure test can run in parallel.
		c := cases[i]

		dirTestdata := filepath.Join(cwd, c.name)

		// setup environment
		if c.setup != nil {
			t.Logf("Setting up environment for %s", c.name)
			c.setup(t, dirTestdata)
		}
		// destroy environment
		if c.destroy != nil {
			defer func() {
				t.Logf("Destroying environment for %s", c.name)
				c.destroy(t, dirTestdata)
			}()
		}

		for _, testcase := range c.testers {
			t.Run(fmt.Sprintf("command=%s,test=%s", c.name, testcase.name), func(t *testing.T) {
				testcase.tester(t, dirTestdata)
			})
		}
	}
}

func testValidatePassing(t *testing.T, dirTestdata string) {
	stdOut, stdErr := runCft(t, false, cwd, []string{
		"validate",
		"my-networks",
		"--policy-path",
		filepath.Join(dirTestdata, "policies", "test_all"),
	})
	checkOutput(t, append(stdOut, stdErr...), []string{"No violations found."})

	stdOut, stdErr = runCft(t, false, cwd, []string{
		"validate",
		"my-firewalls",
		"--policy-path",
		filepath.Join(dirTestdata, "policies", "test_all"),
	})
	checkOutput(t, append(stdOut, stdErr...), []string{"No violations found."})

	stdOut, stdErr = runCft(t, false, cwd, []string{
		"validate",
		"my-instance-prod-1",
		"--policy-path",
		filepath.Join(dirTestdata, "policies", "test_all"),
	})
	checkOutput(t, append(stdOut, stdErr...), []string{"No violations found."})

	stdOut, stdErr = runCft(t, false, cwd, []string{
		"validate",
		"project-iam",
		"--policy-path",
		filepath.Join(dirTestdata, "policies", "test_all"),
	})
	checkOutput(t, append(stdOut, stdErr...), []string{"No violations found."})
}

func testValidateFailing(t *testing.T, dirTestdata string) {
	stdOut, stdErr := runCft(t, false, cwd, []string{
		"validate",
		"my-networks",
		"--policy-path",
		filepath.Join(dirTestdata, "policies", "test_none"),
	})
	checkOutput(t, append(stdOut, stdErr...), []string{"No violations found."})

	stdOut, stdErr = runCft(t, false, cwd, []string{
		"validate",
		"my-firewalls",
		"--policy-path",
		filepath.Join(dirTestdata, "policies", "test_none"),
	})
	checkOutput(t, append(stdOut, stdErr...), []string{
		"Found Violations",
		"Constraint restrict-firewall-rule-world_open",
		"//compute\\.googleapis\\.com/projects/.+/global/firewalls/allow-proxy-from-inside-dev",
		"//compute\\.googleapis\\.com/projects/.+/global/firewalls/allow-proxy-from-inside-prod",
	})

	stdOut, stdErr = runCft(t, false, cwd, []string{
		"validate",
		"my-instance-prod-1",
		"--policy-path",
		filepath.Join(dirTestdata, "policies", "test_none"),
	})
	checkOutput(t, append(stdOut, stdErr...), []string{
		"Found Violations",
		"Constraint gcp-compute-zone",
		"//compute\\.googleapis\\.com/projects/.+/zones/us-central1-a/instances/my-instance-prod-1",
	})

	stdOut, stdErr = runCft(t, false, cwd, []string{
		"validate",
		"project-iam",
		"--policy-path",
		filepath.Join(dirTestdata, "policies", "test_none"),
	})
	checkOutput(t, append(stdOut, stdErr...), []string{
		"Found Violations",
		"Constraint iam_ban_roles",
		"//cloudresourcemanager\\.googleapis\\.com/projects/12345",
	})
}

func setupValidate(t *testing.T, dirTestdata string) {
	runCft(t, false, cwd, []string{"apply", filepath.Join(dirTestdata, "deployment")})
	runCft(t, false, cwd, []string{"apply", "--preview", filepath.Join(dirTestdata, "iam.yml")})
}

func destroyValidate(t *testing.T, dirTestdata string) {
	runCft(t, false, cwd, []string{"delete", filepath.Join(dirTestdata, "deployment")})
	runCft(t, false, cwd, []string{"delete", filepath.Join(dirTestdata, "iam.yml")})
}

func checkOutput(t *testing.T, data []byte, regex []string) {
	for _, reg := range regex {
		wantRe := regexp.MustCompile(reg)
		if !wantRe.Match(data) {
			t.Fatalf("Wrong output output, \ngot=%s \nwant (regex)=%s", string(data), reg)
		}
	}
}

func runCft(t *testing.T, wantError bool, dir string, args []string) ([]byte, []byte) {
	executable := cftBinary
	args = append(args, "--non-interactive")
	cmd := exec.Command(executable, args...)
	cmd.Dir = dir
	return run(t, cmd, wantError)
}

// run a command and call t.Fatal on non-zero exit.
func run(t *testing.T, cmd *exec.Cmd, wantError bool) ([]byte, []byte) {
	var stderr, stdout bytes.Buffer
	cmd.Stderr, cmd.Stdout = &stderr, &stdout
	err := cmd.Run()
	if gotError := (err != nil); gotError != wantError {
		t.Fatalf("running %s: \nerror=%v \nstderr=%s \nstdout=%s", cmdToString(cmd), err, stderr.String(), stdout.String())
	}
	// Print env, stdout and stderr if verbose flag is used.
	if len(cmd.Env) != 0 {
		t.Logf("=== Environment Variable of %s ===", cmdToString(cmd))
		t.Log(strings.Join(cmd.Env, "\n"))
	}
	if stdout.String() != "" {
		t.Logf("=== STDOUT of %s ===", cmdToString(cmd))
		t.Log(stdout.String())
	}
	if stderr.String() != "" {
		t.Logf("=== STDERR of %s ===", cmdToString(cmd))
		t.Log(stderr.String())
	}
	return stdout.Bytes(), stderr.Bytes()
}

// cmdToString clones the logic of https://golang.org/pkg/os/exec/#Cmd.String.
func cmdToString(c *exec.Cmd) string {
	// report the exact executable path (plus args)
	b := new(strings.Builder)
	b.WriteString(c.Path)
	for _, a := range c.Args[1:] {
		b.WriteByte(' ')
		b.WriteString(a)
	}
	return b.String()
}
