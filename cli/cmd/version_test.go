package cmd

import (
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	c, output, err := ExecuteCommandC(rootCmd, "version")

	if c.Name() != "version" {
		t.Errorf(`invalid command returned from ExecuteC: expected "version"', got %q`, c.Name())
	}

	if output == "" {
		t.Errorf("Unexpected output: %v", output)
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestVersionCommandHelp(t *testing.T) {
	output, err := ExecuteCommand(rootCmd, "version", "-h")
	if !strings.HasPrefix(output, versionCmd.Long) {
		t.Errorf("Unexpected output: %v", output)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
