package bptest

import (
	"fmt"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidTestName(t *testing.T) {
	tests := []struct {
		name       string
		intTestDir string
		testName   string
		errMsg     string
	}{
		{
			name:     "valid explicit",
			testName: "TestBar",
		},
		{
			name:     "valid discovered",
			testName: "TestAll/examples/baz",
		},
		{
			name:     "valid all regex",
			testName: "Test.*",
		},
		{
			name:     "all",
			testName: "all",
		},
		{
			name:     "invalid",
			testName: "TestBaz",
			errMsg:   "unable to find TestBaz- one of [\"TestAll/examples/baz\" \"TestAll/fixtures/qux\" \"TestBar\" \"TestFoo\" \"all\"]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			if tt.intTestDir == "" {
				tt.intTestDir = path.Join(testDirWithDiscovery, intTestDir)
			}
			err := isValidTestName(tt.intTestDir, tt.testName)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestGetTestCmd(t *testing.T) {

	tests := []struct {
		name       string
		intTestDir string
		testStage  string
		testName   string
		wantArgs   []string
		errMsg     string
	}{
		{
			name:     "single test",
			testName: "TestFoo",
			wantArgs: []string{"./...", "-run", "TestFoo", "-p", "1", "-count", "1", "-timeout", "0"},
		},
		{
			name:     "all tests",
			testName: "all",
			wantArgs: []string{"./...", "-p", "1", "-count", "1", "-timeout", "0"},
		},
		{
			name:      "custom stage",
			testName:  "TestFoo",
			testStage: "init",
			wantArgs:  []string{"./...", "-run", "TestFoo", "-p", "1", "-count", "1", "-timeout", "0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			if tt.intTestDir == "" {
				tt.intTestDir = path.Join(testDirWithDiscovery, intTestDir)
			}
			gotCmd, err := getTestCmd(tt.intTestDir, tt.testStage, tt.testName)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.Subset(gotCmd.Args, tt.wantArgs)
				if tt.testStage != "" {
					assert.Contains(gotCmd.Env, fmt.Sprintf("RUN_STAGE=%s", tt.testStage))
				}
			}
		})
	}
}
