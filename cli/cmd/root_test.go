package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = ExecuteCommandC(root, args...)
	return output, err
}

func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	setOutput(root, buf)

	// reset command state after prev test execution
	root.SetArgs([]string{})
	root.ResetFlags()

	// set child command and/or command line args
	if args != nil {
		root.SetArgs(args)
	}

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func setOutput(rootCommand *cobra.Command, buf *bytes.Buffer) {
	rootCommand.SetOutput(buf)
	for _, command := range rootCommand.Commands() {
		setOutput(command, buf)
	}
}

func TestRootCommand(t *testing.T) {
	rootCmd.SetArgs([]string{})
	output, err := ExecuteCommand(rootCmd)
	if output == "" {
		t.Errorf("Unexpected output: %v", output)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRootCommandHelpArgs(t *testing.T) {
	output, err := ExecuteCommand(rootCmd, "-h")
	if output == "" {
		t.Errorf("Unexpected output: %v", output)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRootCommandWithUnknownCommand(t *testing.T) {
	output, err := ExecuteCommand(rootCmd, "unknown")
	if !strings.HasPrefix(output, "Error: unknown command \"unknown\" for \"cft\"") {
		t.Errorf("Unexpected output: %v", output)
	}
	if err == nil {
		t.Errorf("Expected unkwnown command error")
	}
}
